package types

import "time"

// AgentClass identifies the type of agent that produced a signal.
// Spec v1.3 §11.1.
//
// Only the three global agents emit signals; session-scoped agents
// (Edge Scout, Corpus Edge Scout, Bridge Agent) never appear here.
// Spec §4.2 + §4.4 cite this constraint.
type AgentClass string

const (
	AgentScout       AgentClass = "scout"
	AgentCorrelation AgentClass = "correlation"
	AgentSynthesis   AgentClass = "synthesis"
)

// AnomalyKind classifies what makes a pattern interesting.
//
// AnomalyStructural is provisional in Spec v1.3 and will be formalised in
// Theory v1.5 (issue #3). It is the kind used by Synthesis-promoted signals
// minted from session-scoped agent output (Edge Scout boundary blindspots,
// Corpus Edge Scout diversity findings) where neither density / correlation
// / narrative fits cleanly. See §4.4 (Synthesis as Persistence Boundary).
//
// Note: AnomalyKind is distinct from the InsightSignal.Structural bool field,
// which records whether the pattern is landscape (persistent under
// perturbation) or exploration (contingent). The two carry different
// information; both may be set independently on the same signal.
type AnomalyKind string

const (
	AnomalyDensity     AnomalyKind = "density"
	AnomalyCorrelation AnomalyKind = "correlation"
	AnomalyNarrative   AnomalyKind = "narrative"
	AnomalyStructural  AnomalyKind = "structural" // provisional; see Theory v1.5 (forthcoming, issue #3)
)

// HyperedgeRef is a lightweight pointer to a hyperedge in MuninnDB / Wyrd.
// Spec v1.3 §11.1.
type HyperedgeRef struct {
	EdgeID   string `json:"edge_id"`
	EdgeType string `json:"edge_type"`
}

// ClaimVersion tracks the evolution of a signal's description across
// Strengthening and Mutation lifecycle events (Theory §3.6.3).
// Spec v1.3 §11.1.
type ClaimVersion struct {
	Version   int       `json:"version"`
	Claim     string    `json:"claim"`
	Reason    string    `json:"reason"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

// ScalarReferent records a predicted-vs-observed scalar value pair and a
// divergence score for a named surveillance target on an operational scope.
// Spec v1.4 §2.4.
//
// ScalarReferent is the Contextus-side carrier for surveillance-mode anomaly
// scoring on operational-telemetry streams (e.g. "5-min cpu_temp ≤ 70°C").
// Predicted is the expected value per the active baseline; Observed is the
// telemetry reading at detection time; Score is the normalised divergence
// (0.0 = no divergence, 1.0 = maximum divergence) computed by the
// surveillance-mode scout as |observed - predicted| / max(|predicted|, 1).
//
// ScalarReferent does NOT carry the raw time-series — only the summary tuple
// (predicted, observed, score) at the moment the Synthesis persistence-boundary
// threshold was crossed. Full telemetry lives in Wyrd via the ctx-adapter-system
// adapter (NT_SCOPE_OPERATIONAL spec addendum §9.2).
//
// Ownership note: ScalarReferent score computation is Contextus-side
// (surveillance agent owns the divergence formula). CTH's ScorePrediction
// primitive (CTH issue #53) is a separate scoring surface for algebraic-
// integrity claims; the two may be coupled at Walk-phase by a Synthesis
// cross-referencing step, but that coupling is not implemented here.
// See Contextus-Spec-Addendum-NT-Scope-Operational §7.3 cross-reference.
type ScalarReferent struct {
	// Label is the surveillance target name, e.g. "cpu_temp_5min_avg_celsius".
	Label string `json:"label"`

	// Predicted is the expected value for this metric per the active baseline.
	Predicted float64 `json:"predicted"`

	// Observed is the telemetry value at detection time.
	Observed float64 `json:"observed"`

	// Score is the normalised divergence in [0.0, 1.0]:
	//   min(|observed-predicted| / max(|predicted|, 1.0), 1.0)
	Score float64 `json:"score"`

	// ScopeID references the NT_SCOPE_OPERATIONAL scope node whose hardware
	// subsystem this referent tracks. Matches ScopeOperational.ScopeID.
	ScopeID string `json:"scope_id"`
}

// InsightSignal is the atomic unit of agent output (Theory §3.6).
// Spec v1.3 §11.1.
//
// Tier-conditional field population is enforced by the retention layer at
// write time, not by this struct. See Spec §5.4.3 for the per-tier rules.
type InsightSignal struct {
	SignalID           Addr              `json:"signal_id"`
	AgentType          AgentClass        `json:"agent_type"`
	Subgraph           []HyperedgeRef    `json:"subgraph"`
	TraversalPath      []Addr            `json:"traversal_path"`
	AnomalyType        AnomalyKind       `json:"anomaly_type"`
	AnomalyScore       float64           `json:"anomaly_score"`
	LocaleEnvelope     LocaleBounds      `json:"locale_envelope"`
	LocaleGrain        GrainSpec         `json:"locale_grain"`
	Persistence        int               `json:"persistence"`
	Structural         bool              `json:"structural"` // landscape vs exploration; distinct from AnomalyStructural
	Confidence         float64           `json:"confidence"`
	ConfidenceVariance float64           `json:"confidence_variance"`
	FirstSeen          time.Time         `json:"first_seen"`
	LastSeen           time.Time         `json:"last_seen"`
	Version            int               `json:"version"`
	ClaimHistory       []ClaimVersion    `json:"claim_history,omitempty"`
	Promoted           bool              `json:"promoted"`
	PromotionReceiptID *Addr             `json:"promotion_receipt_id,omitempty"`
	Evidence           []EvidencePointer `json:"evidence,omitempty"` // tier-conditional; see §5.4

	// Referents carries the predicted-vs-observed scalar referent pairs for
	// surveillance-mode AnomalyStructural signals emitted from operational-scope
	// telemetry streams. Nil for non-operational signals. Populated by the
	// surveillance-mode scout that detects the cross-domain correlation per
	// Contextus-Spec-Addendum-NT-Scope-Operational §7 + Spec v1.4 §2.4.
	Referents []ScalarReferent `json:"referents,omitempty"`
}
