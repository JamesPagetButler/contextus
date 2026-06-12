package measurement

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

// schemaURI is the JSON Schema $id for the embedded measurement schema.
// Matches the $id in schema/nt-signal-measurement.schema.json and in schemaJSON.
const schemaURI = "https://github.com/JamesPagetButler/contextus/schema/nt-signal-measurement.schema.json"

// schemaJSON is the embedded measurement JSON Schema 2020-12. It MUST be kept
// structurally in sync with schema/nt-signal-measurement.schema.json (the
// authoritative on-disk schema). The two MUST agree under JSON-decode
// (whitespace/description differences are tolerated; structural drift is not).
// TestSchemaFileMatchesEmbedded enforces this contract.
const schemaJSON = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/JamesPagetButler/contextus/schema/nt-signal-measurement.schema.json",
  "title": "NTSignalMeasurement",
  "type": "object",
  "required": ["quantity_kind", "value", "uncertainty", "confidence", "unit_system", "registry_version"],
  "additionalProperties": false,
  "properties": {
    "quantity_kind":          { "type": "string", "minLength": 1 },
    "extended_quantity_kind": { "type": "string", "minLength": 1 },
    "value":                  { "type": "number" },
    "uncertainty": {
      "type": "object",
      "required": ["plus", "minus"],
      "additionalProperties": false,
      "properties": {
        "plus":  { "type": "number", "minimum": 0 },
        "minus": { "type": "number", "minimum": 0 }
      }
    },
    "confidence":       { "type": "number", "minimum": 0, "maximum": 1 },
    "unit_system":      { "type": "string", "minLength": 1 },
    "registry_version": { "type": "integer", "minimum": 1 }
  }
}`

var compiledSchema *jsonschema.Schema

func init() {
	compiledSchema = mustCompileSchema()
}

// mustCompileSchema compiles the embedded JSON Schema 2020-12 for measurement
// payloads. Called once from init(); panics on failure (programmer error —
// schema is embedded in the binary and must always be valid).
func mustCompileSchema() *jsonschema.Schema {
	c := jsonschema.NewCompiler()

	var schemaDoc interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schemaDoc); err != nil {
		panic(fmt.Sprintf("measurement: parse embedded schema JSON: %v", err))
	}
	if err := c.AddResource(schemaURI, schemaDoc); err != nil {
		panic(fmt.Sprintf("measurement: add embedded schema resource: %v", err))
	}

	s, err := c.Compile(schemaURI)
	if err != nil {
		panic(fmt.Sprintf("measurement: compile embedded schema: %v", err))
	}
	return s
}

// Validate validates raw JSON measurement payload bytes against the embedded
// schema. It returns ErrMeasurementInvalid (wrapping the schema error) on
// failure, nil on success.
//
// Validate does NOT check quantity_kind membership against the registry. That
// check requires the registry (see Registry.ValidateKind).
func Validate(_ context.Context, payload []byte) error {
	var doc interface{}
	d := json.NewDecoder(bytes.NewReader(payload))
	d.UseNumber()
	if err := d.Decode(&doc); err != nil {
		return fmt.Errorf("%w: %v", ErrMeasurementInvalid, err)
	}
	if err := compiledSchema.Validate(doc); err != nil {
		return fmt.Errorf("%w: %v", ErrMeasurementInvalid, err)
	}
	return nil
}

// Decode validates and decodes raw JSON measurement payload bytes into a
// types.Measurement. Returns ErrMeasurementInvalid on schema failure.
func Decode(_ context.Context, payload []byte) (types.Measurement, error) {
	var doc interface{}
	d := json.NewDecoder(bytes.NewReader(payload))
	d.UseNumber()
	if err := d.Decode(&doc); err != nil {
		return types.Measurement{}, fmt.Errorf("%w: %v", ErrMeasurementInvalid, err)
	}
	if err := compiledSchema.Validate(doc); err != nil {
		return types.Measurement{}, fmt.Errorf("%w: %v", ErrMeasurementInvalid, err)
	}

	// JSON unmarshal to the concrete type. Use a second Decoder to preserve
	// standard number handling (float64) for the struct fields.
	var m types.Measurement
	if err := json.Unmarshal(payload, &m); err != nil {
		return types.Measurement{}, fmt.Errorf("%w: %v", ErrMeasurementInvalid, err)
	}
	return m, nil
}

// RegistryKind is one entry from the quantity-kind registry.
type RegistryKind struct {
	Kind         string `json:"kind"`
	OwnerTenant  string `json:"owner_tenant"`
	UnitSystem   string `json:"unit_system"`
	Description  string `json:"description"`
	IntroducedIn int    `json:"introduced_in"`
	DeprecatedIn int    `json:"deprecated_in,omitempty"`
	MapsTo       any    `json:"maps_to,omitempty"`
}

// registryJSON is the shape of the generated registry/quantity-kinds.json.
type registryJSON struct {
	RegistryVersion int            `json:"registry_version"`
	ContentHash     string         `json:"content_hash"`
	GeneratedAt     string         `json:"generated_at"`
	GeneratedBy     string         `json:"generated_by"`
	Kinds           []RegistryKind `json:"kinds"`
}

// Registry holds a decoded quantity-kind registry. Use LoadRegistryJSON or
// LoadRegistryYAML to construct one.
type Registry struct {
	Version int
	kinds   map[string]RegistryKind // keyed by kind name
}

// LoadRegistryJSON loads a Registry from a generated JSON projection file
// (registry/quantity-kinds.json). Returns ErrRegistryParse or
// ErrRegistryInvalid on failure.
func LoadRegistryJSON(_ context.Context, r io.Reader) (*Registry, error) {
	var reg registryJSON
	if err := json.NewDecoder(r).Decode(&reg); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRegistryParse, err)
	}
	if reg.RegistryVersion < 1 {
		return nil, fmt.Errorf("%w: registry_version must be >= 1", ErrRegistryInvalid)
	}
	if len(reg.Kinds) == 0 {
		return nil, fmt.Errorf("%w: kinds list must not be empty", ErrRegistryInvalid)
	}
	return buildRegistry(reg.RegistryVersion, reg.Kinds)
}

// LoadRegistryJSONFile is a convenience wrapper around LoadRegistryJSON that
// opens the file at path and passes the reader to LoadRegistryJSON.
func LoadRegistryJSONFile(ctx context.Context, path string) (*Registry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: open %s: %v", ErrRegistryParse, path, err)
	}
	defer func() { _ = f.Close() }()
	return LoadRegistryJSON(ctx, f)
}

// rawYAMLRegistry mirrors quantity-kinds.yaml for decoding.
type rawYAMLRegistry struct {
	RegistryVersion int            `yaml:"registry_version"`
	Kinds           []RegistryKind `yaml:"kinds"`
}

// LoadRegistryYAML loads a Registry directly from the YAML governance source.
// Intended for use in tests and tooling; production consumers should prefer
// LoadRegistryJSON (the generated JSON projection).
func LoadRegistryYAML(_ context.Context, r io.Reader) (*Registry, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: read YAML: %v", ErrRegistryParse, err)
	}
	var raw rawYAMLRegistry
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: parse YAML: %v", ErrRegistryParse, err)
	}
	if raw.RegistryVersion < 1 {
		return nil, fmt.Errorf("%w: registry_version must be >= 1", ErrRegistryInvalid)
	}
	if len(raw.Kinds) == 0 {
		return nil, fmt.Errorf("%w: kinds list must not be empty", ErrRegistryInvalid)
	}
	return buildRegistry(raw.RegistryVersion, raw.Kinds)
}

func buildRegistry(version int, kinds []RegistryKind) (*Registry, error) {
	r := &Registry{
		Version: version,
		kinds:   make(map[string]RegistryKind, len(kinds)),
	}
	for _, k := range kinds {
		if k.Kind == "" {
			return nil, fmt.Errorf("%w: kind entry missing name", ErrRegistryInvalid)
		}
		if _, dup := r.kinds[k.Kind]; dup {
			return nil, fmt.Errorf("%w: duplicate kind %q", ErrRegistryInvalid, k.Kind)
		}
		r.kinds[k.Kind] = k
	}
	return r, nil
}

// ValidateKind checks that the Measurement's quantity_kind exists in the
// registry at the signal's claimed registry_version.
//
// Returns:
//   - ErrRegistryVersionMismatch when m.RegistryVersion > r.Version (the bridge
//     should defer-not-reject per §8 rule 2).
//   - ErrQuantityKindUnknown when the kind is not in the registry.
//   - nil on success.
//
// Extended kinds (ExtendedQuantityKind set, QuantityKind empty) bypass registry
// validation per §5.4: they are opaque to cross-tenant matching and are accepted
// as-is. If both fields are set, ErrMeasurementInvalid is returned.
func (r *Registry) ValidateKind(_ context.Context, m types.Measurement) error {
	if m.ExtendedQuantityKind != "" && m.QuantityKind != "" {
		return fmt.Errorf("%w: quantity_kind and extended_quantity_kind are mutually exclusive", ErrMeasurementInvalid)
	}
	// Extended kinds bypass registry check.
	if m.ExtendedQuantityKind != "" {
		return nil
	}
	if m.RegistryVersion > r.Version {
		return fmt.Errorf("%w: signal version %d, registry max %d",
			ErrRegistryVersionMismatch, m.RegistryVersion, r.Version)
	}
	k, ok := r.kinds[m.QuantityKind]
	if !ok {
		return fmt.Errorf("%w: %q not in registry v%d",
			ErrQuantityKindUnknown, m.QuantityKind, m.RegistryVersion)
	}
	// Check the kind was not introduced after the claimed emission version.
	if k.IntroducedIn > m.RegistryVersion {
		return fmt.Errorf("%w: %q introduced at v%d but signal claims v%d",
			ErrQuantityKindUnknown, m.QuantityKind, k.IntroducedIn, m.RegistryVersion)
	}
	return nil
}

// Kinds returns a copy of all RegistryKind entries. Ordered alphabetically by
// kind name (deterministic; same order as the generated JSON projection).
func (r *Registry) Kinds() []RegistryKind {
	out := make([]RegistryKind, 0, len(r.kinds))
	for _, k := range r.kinds {
		out = append(out, k)
	}
	return out
}
