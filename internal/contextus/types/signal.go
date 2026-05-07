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
}
