package tenancy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

// tenantIDPattern mirrors the JSON-Schema pattern for tenant_profile.tenant_id
// per Contextus-Spec-Addendum-Research-Aid-Tenancy §3.2. The loader applies it
// as defence-in-depth after JSON-Schema validation has already run; the two
// surfaces MUST agree on the pattern.
var tenantIDPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

// LoadResult separates Contextus-payload structs from Wyrd-envelope metadata.
//
// The split is mandated by Wyrd PR #40 §2.1: fields like tier_immune, salience,
// and weight_tier appear in the YAML but are NOT part of the Contextus
// pkg/types.Scope* structs — they thread to model.Node.* / model.Weight.Tier
// at the Wyrd side. Mixing them into the Contextus types would violate the
// Spec v1.3 §11.4 type contract and break Wyrd's compile-time guarantee.
type LoadResult struct {
	PhysicalScopes   []types.ScopePhysical
	ConceptualScopes []types.ScopeConceptual
	Memberships      []types.ScopeMembership

	// NodeOptions[scope_id] carries envelope metadata that threads to
	// model.Node.* fields at the Wyrd side; absent from the Contextus types.
	NodeOptions map[string]NodeOptions

	// EdgeOptions[i] carries weight_tier for Memberships[i]; absent from
	// pkg/types.ScopeMembership (per Wyrd PR #40 §2.1 note on weight_tier).
	EdgeOptions []EdgeOptions

	// TenantProfile carries the optional top-level tenant_profile block per
	// Contextus-Spec-Addendum-Research-Aid-Tenancy §3.3. Nil when the block
	// is absent in the YAML (v1.3 baseline backwards compatibility per §3.4);
	// populated when the tenant declares its federation identity and
	// subscriber profile. See pkg/types.TenantProfile for field semantics.
	TenantProfile *types.TenantProfile
}

// NodeOptions carries the Wyrd-envelope fields for a scope node that are NOT
// part of pkg/types.ScopePhysical or pkg/types.ScopeConceptual.
type NodeOptions struct {
	// TierImmune maps to model.Node.TierImmune. When true, the node is exempt
	// from tier-based Ebbinghaus decay. Used for foundational conceptual scopes
	// such as the hamilton-product algebraic structure (QBP tenant).
	TierImmune bool

	// Salience maps to model.Node.Salience [0.0, 1.0].
	Salience float64
}

// EdgeOptions carries the Wyrd-envelope fields for a membership edge that are
// NOT part of pkg/types.ScopeMembership.
type EdgeOptions struct {
	// WeightTier maps to model.Weight.Tier for the constructed membership
	// hyperedge. Defaults to "complex" per Wyrd PR #40 §2.1 if absent.
	WeightTier string
}

// rawConfig is the intermediate struct that yaml.v3 decodes into before the
// envelope split. It mirrors the JSON Schema shape exactly so KnownFields
// strict-mode catches unknown keys.
type rawConfig struct {
	PhysicalScopes   []rawPhysical     `yaml:"physical_scopes"`
	ConceptualScopes []rawConceptual   `yaml:"conceptual_scopes"`
	ScopeMemberships []rawMembership   `yaml:"scope_memberships"`
	TenantProfile    *rawTenantProfile `yaml:"tenant_profile"`
}

// rawTenantProfile mirrors the JSON Schema 2020-12 tenant_profile shape from
// Contextus-Spec-Addendum-Research-Aid-Tenancy §3.2. It is the intermediate
// decode target; buildResult lifts it to *types.TenantProfile, synthesising
// tenant_subgraph_ref.uri when omitted per spec §2.4 + §3.3 step 3.
type rawTenantProfile struct {
	TenantID          string                `yaml:"tenant_id"`
	SubscriberProfile rawSubscriberProfile  `yaml:"subscriber_profile"`
	TenantSubgraphRef *rawTenantSubgraphRef `yaml:"tenant_subgraph_ref"`
}

type rawSubscriberProfile struct {
	AcceptedScaffoldTypes    []string `yaml:"accepted_scaffold_types"`
	AcceptedCorpusClasses    []string `yaml:"accepted_corpus_classes"`
	IntendedConsumersDefault []string `yaml:"intended_consumers_default"`
}

type rawTenantSubgraphRef struct {
	URI string `yaml:"uri"`
}

type rawPhysical struct {
	ID            string     `yaml:"id"`
	Description   string     `yaml:"description"`
	TypeNodes     []string   `yaml:"type_nodes"`
	Bounds        *rawBounds `yaml:"bounds"`
	ParentScopeID string     `yaml:"parent_scope_id"`
	Tags          []string   `yaml:"tags"`
	TierImmune    bool       `yaml:"tier_immune"`
	Salience      float64    `yaml:"salience"`
}

type rawConceptual struct {
	ID              string   `yaml:"id"`
	Description     string   `yaml:"description"`
	TypeNodes       []string `yaml:"type_nodes"`
	OntologyURI     string   `yaml:"ontology_uri"`
	ParentScopeID   string   `yaml:"parent_scope_id"`
	RelatedScopeIDs []string `yaml:"related_scope_ids"`
	Tags            []string `yaml:"tags"`
	TierImmune      bool     `yaml:"tier_immune"`
	Salience        float64  `yaml:"salience"`
}

type rawMembership struct {
	Scope         string  `yaml:"scope"`
	Member        string  `yaml:"member"`
	Since         string  `yaml:"since"` // ISO 8601; parsed to time.Time
	Confidence    float64 `yaml:"confidence"`
	ProvenanceTag string  `yaml:"provenance_tag"`
	Method        string  `yaml:"method"`
	WeightTier    string  `yaml:"weight_tier"`
}

type rawBounds struct {
	Lat    [2]float64 `yaml:"lat"`
	Lon    [2]float64 `yaml:"lon"`
	Time   [2]string  `yaml:"time"`
	Height [2]float64 `yaml:"height"`
}

// compiledSchema is the parsed JSON Schema, compiled once at package init so
// the file-system read (embedded schema bytes) does not repeat per call.
var compiledSchema *jsonschema.Schema

func init() {
	compiledSchema = mustCompileSchema()
}

// LoadScopeConfig validates configPath against the JSON Schema 2020-12, then
// decodes the YAML into a LoadResult with the Wyrd-envelope fields split out.
//
// Error taxonomy (all wrapped so errors.Is works):
//   - ErrScopeConfigParse   — YAML decode failure
//   - ErrScopeConfigInvalid — schema validation failure or semantic constraint
//   - ErrScopeLoadConflict  — duplicate scope_id within the config
//
// Path security: configPath must not contain ".." components; the function
// returns ErrScopeConfigInvalid if a traversal attempt is detected.
//
// ctx is accepted for API consistency with future cancellation / deadline
// propagation; it is not yet inspected in v0.1.
func LoadScopeConfig(_ context.Context, configPath string) (*LoadResult, error) {
	if err := guardPath(configPath); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(configPath) //nolint:gosec // path is validated by guardPath above
	if err != nil {
		return nil, fmt.Errorf("tenancy: read %q: %w", configPath, ErrScopeConfigParse)
	}

	// YAML → JSON for schema validation (jsonschema/v6 validates JSON documents).
	jsonBytes, err := yamlToJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("tenancy: yaml-to-json %q: %w", configPath, ErrScopeConfigParse)
	}

	if err := validateAgainstSchema(jsonBytes); err != nil {
		return nil, fmt.Errorf("tenancy: schema validation %q: %w: %w", configPath, err, ErrScopeConfigInvalid)
	}

	var cfg rawConfig
	dec := yaml.NewDecoder(bytes.NewReader(raw))
	dec.KnownFields(true) // reject unknown keys; belt-and-suspenders after schema validation
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("tenancy: decode %q: %w: %w", configPath, err, ErrScopeConfigParse)
	}

	return buildResult(cfg)
}

// guardPath rejects configPath values that contain ".." components, preventing
// path-traversal attacks when configPath originates from user input.
//
// We split on the separator rather than relying solely on filepath.Clean because
// Clean resolves ".." silently, masking the attack intent. Splitting lets us
// detect the raw traversal attempt before any normalisation occurs.
func guardPath(configPath string) error {
	for _, part := range strings.Split(filepath.ToSlash(configPath), "/") {
		if part == ".." {
			return fmt.Errorf("tenancy: path traversal rejected %q: %w", configPath, ErrScopeConfigInvalid)
		}
	}
	return nil
}

// yamlToJSON converts YAML bytes to JSON bytes for schema validation.
// gopkg.in/yaml.v3 unmarshals into interface{} and encoding/json marshals
// back out; this is the minimal-dependency round-trip approach.
func yamlToJSON(yamlBytes []byte) ([]byte, error) {
	var v interface{}
	if err := yaml.Unmarshal(yamlBytes, &v); err != nil {
		return nil, err
	}
	// yaml.v3 produces map[string]interface{} which json.Marshal handles.
	return json.Marshal(normaliseYAMLValue(v))
}

// normaliseYAMLValue converts yaml.v3's map[string]interface{} recursively to
// map[string]interface{} with JSON-compatible types. yaml.v3 emits
// map[string]interface{} at top level but nested maps may be
// map[interface{}]interface{} in older decoders; yaml.v3 avoids this, but we
// normalise defensively.
func normaliseYAMLValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, vv := range val {
			out[k] = normaliseYAMLValue(vv)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, vv := range val {
			out[i] = normaliseYAMLValue(vv)
		}
		return out
	default:
		return val
	}
}

// validateAgainstSchema validates jsonBytes against the compiled scope-config
// JSON Schema 2020-12. Returns a plain error (not wrapped) so the caller can
// wrap it with ErrScopeConfigInvalid.
func validateAgainstSchema(jsonBytes []byte) error {
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("unmarshal for validation: %w", err)
	}
	return compiledSchema.Validate(inst)
}

// buildResult converts a decoded rawConfig into a LoadResult, enforcing
// duplicate-id detection and applying defaults for optional envelope fields.
func buildResult(cfg rawConfig) (*LoadResult, error) {
	seen := make(map[string]struct{})

	res := &LoadResult{
		NodeOptions: make(map[string]NodeOptions),
	}

	for _, rp := range cfg.PhysicalScopes {
		if _, dup := seen[rp.ID]; dup {
			return nil, fmt.Errorf("tenancy: duplicate scope_id %q: %w", rp.ID, ErrScopeLoadConflict)
		}
		seen[rp.ID] = struct{}{}

		res.PhysicalScopes = append(res.PhysicalScopes, types.ScopePhysical{
			ScopeID:       rp.ID,
			Name:          rp.Description, // YAML "description" is the human-readable name
			ParentScopeID: rp.ParentScopeID,
			Tags:          rp.Tags,
			// Geometry, ElevationRangeM, TemporalRange, Grain: not yet surfaced
			// in the v0.1 YAML shape; populated from bounds in a future iteration.
		})
		res.NodeOptions[rp.ID] = NodeOptions{
			TierImmune: rp.TierImmune,
			Salience:   rp.Salience,
		}
	}

	for _, rc := range cfg.ConceptualScopes {
		if _, dup := seen[rc.ID]; dup {
			return nil, fmt.Errorf("tenancy: duplicate scope_id %q: %w", rc.ID, ErrScopeLoadConflict)
		}
		seen[rc.ID] = struct{}{}

		res.ConceptualScopes = append(res.ConceptualScopes, types.ScopeConceptual{
			ScopeID:         rc.ID,
			Name:            rc.Description,
			OntologyURI:     rc.OntologyURI,
			ParentScopeID:   rc.ParentScopeID,
			RelatedScopeIDs: rc.RelatedScopeIDs,
			Tags:            rc.Tags,
		})
		res.NodeOptions[rc.ID] = NodeOptions{
			TierImmune: rc.TierImmune,
			Salience:   rc.Salience,
		}
	}

	for _, rm := range cfg.ScopeMemberships {
		since := time.Time{}
		if rm.Since != "" {
			t, err := time.Parse(time.RFC3339, rm.Since)
			if err != nil {
				return nil, fmt.Errorf("tenancy: membership scope=%q member=%q bad since %q: %w", rm.Scope, rm.Member, rm.Since, ErrScopeConfigInvalid)
			}
			since = t
		}

		wt := rm.WeightTier
		if wt == "" {
			wt = "complex" // default per Wyrd PR #40 §2.1
		}

		res.Memberships = append(res.Memberships, types.ScopeMembership{
			ScopeID:       rm.Scope,
			MemberID:      rm.Member,
			Since:         since,
			Confidence:    rm.Confidence,
			ProvenanceTag: rm.ProvenanceTag,
			Method:        rm.Method,
		})
		res.EdgeOptions = append(res.EdgeOptions, EdgeOptions{WeightTier: wt})
	}

	// tenant_profile is optional per Contextus-Spec-Addendum-Research-Aid-Tenancy
	// §3.1; nil-block path preserves the v1.3 backwards-compat contract (§3.4).
	if cfg.TenantProfile != nil {
		tp, err := buildTenantProfile(cfg.TenantProfile)
		if err != nil {
			return nil, err
		}
		res.TenantProfile = tp
	}

	return res, nil
}

// buildTenantProfile lifts a decoded rawTenantProfile into a *types.TenantProfile,
// applying defence-in-depth validation on TenantID and synthesising the
// tenant_subgraph_ref.uri convention form (cth://tenant/<tenant_id>/subgraph)
// when omitted, per Contextus-Spec-Addendum-Research-Aid-Tenancy §3.3 steps 1–4.
//
// JSON-Schema validation has already run by the time this function executes;
// the tenant_id regex check here is the spec §3.3 step 3a defence-in-depth
// requirement.
func buildTenantProfile(rt *rawTenantProfile) (*types.TenantProfile, error) {
	if !tenantIDPattern.MatchString(rt.TenantID) {
		return nil, fmt.Errorf("tenancy: tenant_profile.tenant_id %q does not match required pattern: %w", rt.TenantID, ErrScopeConfigInvalid)
	}

	uri := ""
	if rt.TenantSubgraphRef != nil {
		uri = rt.TenantSubgraphRef.URI
	}
	if uri == "" {
		// Synthesize the canonical convention form per spec §2.4 + §3.3 step 3b.
		uri = fmt.Sprintf("cth://tenant/%s/subgraph", rt.TenantID)
	}

	return &types.TenantProfile{
		TenantID: rt.TenantID,
		SubscriberProfile: types.SubscriberProfile{
			AcceptedScaffoldTypes:    rt.SubscriberProfile.AcceptedScaffoldTypes,
			AcceptedCorpusClasses:    rt.SubscriberProfile.AcceptedCorpusClasses,
			IntendedConsumersDefault: rt.SubscriberProfile.IntendedConsumersDefault,
		},
		TenantSubgraphRef: types.TenantSubgraphRef{URI: uri},
	}, nil
}

// schemaURI is the canonical $id for the embedded scope-config schema.
const schemaURI = "https://github.com/JamesPagetButler/contextus/schema/scope-config.schema.json"

// schemaJSON is the inline JSON Schema 2020-12 for scope-config. It is
// kept structurally in sync with schema/scope-config.schema.json (the
// authoritative on-disk schema). The two MUST agree under JSON-decode
// (whitespace/description differences are tolerated; structural drift is not).
// TestSchemaFileMatchesEmbedded enforces this contract.
const schemaJSON = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/JamesPagetButler/contextus/schema/scope-config.schema.json",
  "title": "ScopeConfig",
  "type": "object",
  "required": ["physical_scopes", "conceptual_scopes", "scope_memberships"],
  "additionalProperties": false,
  "properties": {
    "physical_scopes": { "type": "array", "items": { "$ref": "#/$defs/PhysicalScopeEntry" } },
    "conceptual_scopes": { "type": "array", "items": { "$ref": "#/$defs/ConceptualScopeEntry" } },
    "scope_memberships": { "type": "array", "items": { "$ref": "#/$defs/ScopeMembershipEntry" } },
    "tenant_profile": {
      "type": "object",
      "additionalProperties": false,
      "required": ["tenant_id", "subscriber_profile"],
      "properties": {
        "tenant_id": { "type": "string", "pattern": "^[a-z][a-z0-9]*(-[a-z0-9]+)*$" },
        "subscriber_profile": {
          "type": "object",
          "additionalProperties": false,
          "required": ["accepted_scaffold_types", "accepted_corpus_classes", "intended_consumers_default"],
          "properties": {
            "accepted_scaffold_types": {
              "type": "array",
              "items": { "type": "string", "enum": ["PRECEDENT_GRAPH", "EVIDENCE_LATTICE", "ALGEBRAIC_STRUCTURE_SCAFFOLD", "SOURCE_LOCATION_HYPOTHESIS"] },
              "uniqueItems": true
            },
            "accepted_corpus_classes": {
              "type": "array",
              "items": { "type": "string", "enum": ["PHYSICS_PREPRINT", "JOURNAL_ARTICLE", "DATASET_DESCRIPTOR", "CODE_REPO", "CONTRACT_PRECEDENT", "REGULATORY_TEXT", "BEEKEEPER_NOTE", "OTHER"] },
              "uniqueItems": true
            },
            "intended_consumers_default": {
              "type": "array",
              "items": { "type": "string", "pattern": "^(self|[a-z][a-z0-9]*(-[a-z0-9]+)*)$" },
              "uniqueItems": true
            }
          }
        },
        "tenant_subgraph_ref": {
          "type": "object",
          "additionalProperties": false,
          "properties": {
            "uri": { "type": "string", "pattern": "^cth://" }
          }
        }
      }
    }
  },
  "$defs": {
    "PhysicalScopeEntry": {
      "type": "object",
      "required": ["id", "description", "type_nodes"],
      "additionalProperties": false,
      "properties": {
        "id": { "type": "string", "minLength": 1 },
        "description": { "type": "string" },
        "type_nodes": { "type": "array", "items": { "type": "string", "minLength": 1 } },
        "bounds": { "$ref": "#/$defs/PhysicalBounds" },
        "parent_scope_id": { "type": "string", "minLength": 1 },
        "tags": { "type": "array", "items": { "type": "string" } },
        "tier_immune": { "type": "boolean" },
        "salience": { "type": "number", "minimum": 0, "maximum": 1 }
      }
    },
    "ConceptualScopeEntry": {
      "type": "object",
      "required": ["id", "description", "type_nodes"],
      "additionalProperties": false,
      "properties": {
        "id": { "type": "string", "minLength": 1 },
        "description": { "type": "string" },
        "type_nodes": { "type": "array", "items": { "type": "string", "minLength": 1 } },
        "ontology_uri": { "type": "string" },
        "parent_scope_id": { "type": "string", "minLength": 1 },
        "related_scope_ids": { "type": "array", "items": { "type": "string", "minLength": 1 } },
        "tags": { "type": "array", "items": { "type": "string" } },
        "tier_immune": { "type": "boolean" },
        "salience": { "type": "number", "minimum": 0, "maximum": 1 }
      }
    },
    "ScopeMembershipEntry": {
      "type": "object",
      "required": ["scope", "member", "confidence", "provenance_tag", "method"],
      "additionalProperties": false,
      "properties": {
        "scope": { "type": "string", "minLength": 1 },
        "member": { "type": "string", "minLength": 1 },
        "since": { "type": "string", "format": "date-time" },
        "confidence": { "type": "number", "minimum": 0, "maximum": 1 },
        "provenance_tag": { "type": "string", "enum": ["P", "D", "T", "I"] },
        "method": { "type": "string" },
        "weight_tier": { "type": "string" }
      }
    },
    "PhysicalBounds": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "lat": { "type": "array", "items": { "type": "number" }, "minItems": 2, "maxItems": 2 },
        "lon": { "type": "array", "items": { "type": "number" }, "minItems": 2, "maxItems": 2 },
        "time": { "type": "array", "items": { "type": "string", "format": "date-time" }, "minItems": 2, "maxItems": 2 },
        "height": { "type": "array", "items": { "type": "number" }, "minItems": 2, "maxItems": 2 }
      }
    }
  }
}`

// mustCompileSchema compiles the JSON Schema 2020-12 for scope-config. Called
// once from init(); panics on failure (programmer error — schema is embedded in
// the binary and must always be valid).
func mustCompileSchema() *jsonschema.Schema {
	c := jsonschema.NewCompiler()

	// AddResource expects an already-parsed document (any), not a raw reader.
	var schemaDoc interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schemaDoc); err != nil {
		panic(fmt.Sprintf("tenancy: parse embedded schema JSON: %v", err))
	}
	if err := c.AddResource(schemaURI, schemaDoc); err != nil {
		panic(fmt.Sprintf("tenancy: compile scope-config schema: %v", err))
	}

	s, err := c.Compile(schemaURI)
	if err != nil {
		panic(fmt.Sprintf("tenancy: compile scope-config schema: %v", err))
	}
	return s
}
