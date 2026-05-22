package types

import "time"

// ScopePhysical is a region of space-time with a name and queryable
// membership. Spec v1.3 §4.6.2 + §11.4.
type ScopePhysical struct {
	ScopeID         string       `json:"scope_id"`
	Name            string       `json:"name"`
	Geometry        []byte       `json:"geometry"` // GeoJSON bytes; opaque to the type
	ElevationRangeM *[2]float64  `json:"elevation_range_m,omitempty"`
	TemporalRange   LocaleBounds `json:"temporal_range"`
	Grain           GrainSpec    `json:"grain"`
	ParentScopeID   string       `json:"parent_scope_id,omitempty"`
	Tags            []string     `json:"tags,omitempty"`
}

// ScopeConceptual is a topic or domain with a name and queryable
// membership. Spec v1.3 §4.6.3 + §11.4.
type ScopeConceptual struct {
	ScopeID         string   `json:"scope_id"`
	Name            string   `json:"name"`
	OntologyURI     string   `json:"ontology_uri,omitempty"`
	ParentScopeID   string   `json:"parent_scope_id,omitempty"`
	RelatedScopeIDs []string `json:"related_scope_ids,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

// ScopeOperational is a hardware-instance-bounded region of the hypergraph
// whose members are nodes attributable to a named host's named hardware
// subsystem. See Contextus-Spec-Addendum-NT-Scope-Operational §§2-5.
//
// ScopeOperational is the third sibling to ScopePhysical (Spec v1.3 §4.6.2)
// and ScopeConceptual (§4.6.3); architectural role identical, membership
// predicate distinct (hardware-identifier-based — exact match on HostID +
// HardwareClass; hierarchical aggregation via ParentScopeID; no wildcards at
// v0.1 per spec §5).
type ScopeOperational struct {
	ScopeID       string   `json:"scope_id"`
	Name          string   `json:"name"`
	HostID        string   `json:"host_id"`
	HardwareClass string   `json:"hardware_class"`
	ParentScopeID string   `json:"parent_scope_id,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// HardwareClass values per Contextus-Spec-Addendum-NT-Scope-Operational §4
// v0.1 taxonomy. Additive-only forward-compatibility commitment; v0.2
// expansion candidates documented in spec §4 (e.g. hardware.tpu, hardware.fpga,
// hardware.power-supply, hardware.thermal-sensor-array, hardware.actuator for
// Sharp Butler, hardware.power-delivery + hardware.containment for Möbius).
// Speculative tags refused at v0.1 — real federation work-target need is the
// trigger (cart-driven tool-acquisition principle).
//
// HardwareClassRuntimeBMAInstance is a conventional Tags entry (NOT a value
// for ScopeOperational.HardwareClass); it brackets a BMA federation tenant's
// host. Sibling runtime.* tags will accompany new tenant classes
// (runtime.sharp-butler-instance, runtime.moebius-instance, runtime.qbp-cu-instance).
const (
	HardwareClassRuntimeBMAInstance = "runtime.bma-instance"
	HardwareClassHardwareCPU        = "hardware.cpu"
	HardwareClassHardwareDisk       = "hardware.disk"
	HardwareClassHardwareGPU        = "hardware.gpu"
	HardwareClassHardwareMemory     = "hardware.memory"
	HardwareClassHardwareNetwork    = "hardware.network"
)

// ScopeMembership is the on-edge metadata for HE_SCOPE_MEMBERSHIP. The edge
// itself is a Wyrd model.Hyperedge of arity 2 connecting a scope node to a
// member node. Spec v1.3 §4.6.4 + §11.4.
//
// Method is adapter-defined per Spec §4.6.5 (the conceptual scope membership
// predicate is implementation-defined; v1.3 ships the schema, v1.x will
// record canonical method values as they emerge).
type ScopeMembership struct {
	ScopeID       string    `json:"scope_id"`
	MemberID      string    `json:"member_id"`
	Since         time.Time `json:"since"`
	Confidence    float64   `json:"confidence"`     // 0.0–1.0
	ProvenanceTag string    `json:"provenance_tag"` // "P" (asserted) | "D" (inferred)
	Method        string    `json:"method"`         // adapter-defined; see §4.6.5
}

// ProvenanceTag values for ScopeMembership. Mirrors the Provenance Envelope
// values in Spec §2.2.
const (
	ProvenanceTagAsserted    = "P"
	ProvenanceTagInferred    = "D"
	ProvenanceTagTheoretical = "T"
	ProvenanceTagImported    = "I"
)

// TenantProfile names a tenant's federation identity and declares its
// subscriber profile against the BMA Research-Aid Protocol (Spec 9.4) and
// the BMA cross-tenant autonomic signal bus (A22 §3 rule 2).
//
// Contextus owns the type shape and YAML/JSON-Schema validation; BMA owns
// the routing semantics. See Contextus-Spec-Addendum-Research-Aid-Tenancy §5.
type TenantProfile struct {
	TenantID          string            `json:"tenant_id"`
	SubscriberProfile SubscriberProfile `json:"subscriber_profile"`
	TenantSubgraphRef TenantSubgraphRef `json:"tenant_subgraph_ref"`
}

// SubscriberProfile names the BMA Research-Aid output classes this tenant
// accepts. Field semantics are defined by BMA Spec 9.4 §3 and §4; this struct
// is the syntactic carrier only.
type SubscriberProfile struct {
	AcceptedScaffoldTypes    []string `json:"accepted_scaffold_types"`
	AcceptedCorpusClasses    []string `json:"accepted_corpus_classes"`
	IntendedConsumersDefault []string `json:"intended_consumers_default"`
}

// TenantSubgraphRef references the CTH subgraph into which BMA writes
// scaffold output for this tenant. Convention-derivable from TenantID;
// explicit URI permitted for tenants whose subgraph anchor diverges from
// the convention.
//
// Validation discipline:
//   - YAML-load path: JSON Schema 2020-12 enforces the "^cth://" pattern at
//     decode time (schema/scope-config.schema.json + the embedded schemaJSON
//     constant in internal/contextus/tenancy/loader.go; TestSchemaFileMatchesEmbedded
//     enforces those two stay in sync).
//   - Programmatic construction (e.g., a future BMA-implementor adapter at
//     Walk-α): the Go type carries no compile-time guard on URI shape, so
//     callers constructing TenantSubgraphRef values directly should validate
//     the "cth://" scheme before publishing the value into the hypergraph.
//     IsValid() method discipline tracked as Walk-α follow-up per
//     @qbp-architecture PR #20 §I4 observation (2026-05-21).
type TenantSubgraphRef struct {
	URI string `json:"uri"` // canonical form: cth://tenant/<tenant_id>/subgraph
}

// AcceptedCorpusClasses values — mirror of BMA Spec 9.4 §3.3.1 corpus_class
// table. Contextus validates the value is in this set; BMA interprets routing.
const (
	CorpusClassPhysicsPreprint   = "PHYSICS_PREPRINT"
	CorpusClassJournalArticle    = "JOURNAL_ARTICLE"
	CorpusClassDatasetDescriptor = "DATASET_DESCRIPTOR"
	CorpusClassCodeRepo          = "CODE_REPO"
	CorpusClassContractPrecedent = "CONTRACT_PRECEDENT"
	CorpusClassRegulatoryText    = "REGULATORY_TEXT"
	CorpusClassBeekeeperNote     = "BEEKEEPER_NOTE"
	CorpusClassOther             = "OTHER"
)

// AcceptedScaffoldTypes values — provisional v0.1 list pending @bma-implementor
// canonical taxonomy. Contextus accepts these values and rejects others with
// ErrScopeConfigInvalid (see internal/contextus/tenancy/errors.go).
// The list will be extended (additive-only commitment) when bma-implementor
// publishes the canonical taxonomy per BMA Spec 9.4 §3.1.
const (
	ScaffoldTypePrecedentGraph             = "PRECEDENT_GRAPH"
	ScaffoldTypeEvidenceLattice            = "EVIDENCE_LATTICE"
	ScaffoldTypeAlgebraicStructureScaffold = "ALGEBRAIC_STRUCTURE_SCAFFOLD"
	ScaffoldTypeSourceLocationHypothesis   = "SOURCE_LOCATION_HYPOTHESIS"
)
