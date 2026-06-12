package types

import "time"

// EdgeScoutFlag is a ranked boundary connection flagged by the Edge Scout.
// Published to ctx.edge.boundary.{session_id} when top-N ranking changes.
// Spec v1.3 §11.3.
//
// EdgeScoutFlag is ephemeral — never persisted to Wyrd. Synthesis is the
// persistence boundary (§4.4); when Synthesis decides a boundary blindspot
// warrants persistence, it mints an InsightSignal of AnomalyStructural.
type EdgeScoutFlag struct {
	SessionID         string    `json:"session_id"`
	BoundaryNodeID    string    `json:"boundary_node_id"`    // Node in explored set
	UnexploredDomain  string    `json:"unexplored_domain"`   // Domain not yet visited
	ConnectionCount   int       `json:"connection_count"`    // Edges to unexplored domain
	Significance      float64   `json:"significance"`        // Statistical significance (0–1)
	RecentChangeCount int       `json:"recent_change_count"` // Connections changed recently
	Rank              int       `json:"rank"`                // Position in current top-N
	ComputedAt        time.Time `json:"computed_at"`
}

// CorpusDiversityReport summarises a search's domain distribution.
// Published to ctx.corpus.diversity after each search. Spec v1.3 §11.3.
//
// Ephemeral; same persistence-boundary discipline as EdgeScoutFlag.
type CorpusDiversityReport struct {
	SearchID            string             `json:"search_id"`
	OriginalQuery       string             `json:"original_query"`
	DomainDistribution  map[string]float64 `json:"domain_distribution"` // domain → fraction
	VocabConcentration  float64            `json:"vocab_concentration"` // 0=diverse, 1=monoculture
	ConvergenceDetected bool               `json:"convergence_detected"`
	BroadeningQueries   []string           `json:"broadening_queries,omitempty"`
	ComputedAt          time.Time          `json:"computed_at"`
}

// BridgeIntervention is a cross-domain connection candidate from the Bridge
// Agent. Published to ctx.bridge.intervention.{session_id} when a local
// minimum is detected. Spec v1.3 §11.3.
//
// Ephemeral; same persistence-boundary discipline as EdgeScoutFlag. Synthesis
// promotes durable convergences as InsightSignal of AnomalyNarrative
// (cross-domain causal story) or AnomalyStructural (boundary-bridging without
// a clear narrative).
type BridgeIntervention struct {
	SessionID            string    `json:"session_id"`
	ScoutFlagID          string    `json:"scout_flag_id"`          // The surveillance flag involved
	ScoutPhysicalDomain  string    `json:"scout_physical_domain"`  // Physical process tag of flag
	SearchPhysicalDomain string    `json:"search_physical_domain"` // Physical process tag of search
	KnowledgeDomainGap   string    `json:"knowledge_domain_gap"`   // The domain boundary separating them
	ProposedQuery        string    `json:"proposed_query"`         // Cross-domain query to inject
	Confidence           float64   `json:"confidence"`             // Confidence this is a local minimum
	DetectedAt           time.Time `json:"detected_at"`
}

// OperationalCorrelation is a cross-domain correlation event produced by
// surveillance-mode scouts monitoring operational-scope telemetry streams.
// Published to ctx.operational.correlation when a scout detects that
// scalar-referent divergence on one or more hardware subsystems correlates
// with a pattern of interest (algebraic-integrity drift, cognitive-trace gap,
// or co-occurrence across multiple hardware classes).
//
// Ephemeral; same persistence-boundary discipline as EdgeScoutFlag. Synthesis
// promotes durable cross-domain findings as InsightSignal of AnomalyStructural.
// See Contextus-Spec-Addendum-NT-Scope-Operational §7 + Spec v1.4 §2.4.
//
// CorrelationID is the surveillance-session-scoped dedup key. A scout that
// re-evaluates the same cross-domain signature must use the same CorrelationID
// so Synthesis can apply the idempotency contract.
type OperationalCorrelation struct {
	// CorrelationID is the dedup key for this correlation event.
	CorrelationID string `json:"correlation_id"`

	// Referents is the set of scalar-referent observations that constitute
	// the correlation. At least one referent is required; a cross-domain
	// correlation requires referents from two or more distinct ScopeIDs.
	Referents []ScalarReferent `json:"referents"`

	// MaxScore is the maximum ScalarReferent.Score across all Referents.
	// Pre-computed by the scout for Policy evaluation without iterating.
	MaxScore float64 `json:"max_score"`

	// DetectedAt is the scout's detection timestamp.
	DetectedAt time.Time `json:"detected_at"`
}
