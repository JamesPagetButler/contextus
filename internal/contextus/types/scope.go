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
