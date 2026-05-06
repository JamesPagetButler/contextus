// Package contextus contains the core data types for the Contextus
// platform and the Contextus-CTH Bridge.
//
// Reconstruction status: NEAR-COMPLETE. Recovered from chat 2 transcript:
// https://claude.ai/chat/3175e792-70f2-48df-bd8e-58762dd066fc
//
// These types are reference definitions, not a complete compiling package.
// LocaleBounds, GrainSpec, HyperedgeRef, and q8.Addr are referenced but
// some field details may need filling in from the underlying MuninnDB
// schema or QBP Locale package.

package contextus

import (
	"time"

	"github.com/helpfulengineering/q8"
)

// =============================================================================
// Insight Signal Types
// =============================================================================

// AgentClass identifies the type of agent that produced a signal.
type AgentClass string

const (
	AgentScout       AgentClass = "scout"
	AgentCorrelation AgentClass = "correlation"
	AgentSynthesis   AgentClass = "synthesis"
)

// AnomalyKind classifies what makes a pattern interesting.
type AnomalyKind string

const (
	AnomalyDensity     AnomalyKind = "density"
	AnomalyCorrelation AnomalyKind = "correlation"
	AnomalyNarrative   AnomalyKind = "narrative"
)

// InsightSignal is the atomic unit of agent output.
// It represents a pattern detected in the hypergraph that an agent
// has flagged as anomalous, correlated, or narratively coherent.
type InsightSignal struct {
	SignalID           q8.Addr        `json:"signal_id"`
	AgentType          AgentClass     `json:"agent_type"`
	Subgraph           []HyperedgeRef `json:"subgraph"`
	TraversalPath      []q8.Addr      `json:"traversal_path"`
	AnomalyType        AnomalyKind    `json:"anomaly_type"`
	AnomalyScore       float64        `json:"anomaly_score"`
	LocaleEnvelope     LocaleBounds   `json:"locale_envelope"`
	LocaleGrain        GrainSpec      `json:"locale_grain"`
	Persistence        int            `json:"persistence"`
	Structural         bool           `json:"structural"`
	Confidence         float64        `json:"confidence"`
	ConfidenceVariance float64        `json:"confidence_variance"`
	FirstSeen          time.Time      `json:"first_seen"`
	LastSeen           time.Time      `json:"last_seen"`
	Version            int            `json:"version"`
	ClaimHistory       []ClaimVersion `json:"claim_history,omitempty"`
	Promoted           bool           `json:"promoted"`
	PromotionReceiptID *q8.Addr       `json:"promotion_receipt_id,omitempty"`
}

// ClaimVersion tracks the evolution of a signal's or anchor's description.
type ClaimVersion struct {
	Version   int       `json:"version"`
	Claim     string    `json:"claim"`
	Reason    string    `json:"reason"` // What prompted the revision
	Source    string    `json:"source"` // Which evidence triggered it
	Timestamp time.Time `json:"timestamp"`
}

// HyperedgeRef is a lightweight pointer to a hyperedge in MuninnDB.
type HyperedgeRef struct {
	EdgeID q8.Addr `json:"edge_id"`
	// Additional edge metadata fields TBD from MuninnDB schema
}

// LocaleBounds describes the spatial/temporal envelope of a pattern.
type LocaleBounds struct {
	SpatialMin q8.Addr   `json:"spatial_min"`
	SpatialMax q8.Addr   `json:"spatial_max"`
	TimeMin    time.Time `json:"time_min"`
	TimeMax    time.Time `json:"time_max"`
}

// GrainSpec describes the resolution at which the pattern is meaningful.
type GrainSpec struct {
	SpatialMeters float64       `json:"spatial_meters"`
	TemporalGrain time.Duration `json:"temporal_grain"`
}

// =============================================================================
// Bridge Types (Contextus -> CTH)
// =============================================================================

// TrustReceipt is the compact handoff object that seeds a CTH anchor
// from a promoted Contextus Insight Signal.
type TrustReceipt struct {
	ReceiptID          q8.Addr   `json:"receipt_id"`
	Claim              string    `json:"claim"`
	ActivationBaseline float64   `json:"activation_baseline"`
	Confidence         float64   `json:"confidence"`
	ProvenanceHash     []byte    `json:"provenance_hash"` // Content-addressed pointer into MuninnDB
	Rationale          string    `json:"rationale"`
	SourceNodes        []q8.Addr `json:"source_nodes"`
	SearchSignature    string    `json:"search_signature"` // Used by horizon scanner
	SurveyResults      []SurveyResult `json:"survey_results,omitempty"`
	MintedAt           time.Time `json:"minted_at"`
}

// SurveyResult is a single piece of evidence found during initial survey
// or horizon scanning, classified by stance toward the claim.
type SurveyResult struct {
	Source    string    `json:"source"`
	URL       string    `json:"url,omitempty"`
	Title     string    `json:"title"`
	Stance    Stance    `json:"stance"`
	Excerpt   string    `json:"excerpt,omitempty"`
	FoundAt   time.Time `json:"found_at"`
}

// Stance classifies the relationship between a piece of evidence and a claim.
type Stance string

const (
	StanceSupporting    Stance = "supporting"
	StanceContradicting Stance = "contradicting"
	StanceAdjacent      Stance = "adjacent"
)

// Heartbeat is the lightweight signal from Contextus to a coupled CTH anchor.
type Heartbeat struct {
	ReceiptID         q8.Addr   `json:"receipt_id"`
	ActivationCurrent float64   `json:"activation_current"`
	Timestamp         time.Time `json:"timestamp"`
}

// SeedAnchor extends a standard CTH anchor with heartbeat coupling and
// active evidence seeking.
type SeedAnchor struct {
	AnchorID           string            `json:"anchor_id"`
	Receipt            TrustReceipt      `json:"receipt"`
	Tier               int               `json:"tier"` // Initially 3
	Rho                float64           `json:"rho"`
	RhoFloor           float64           `json:"rho_floor"`
	HeartbeatActive    bool              `json:"heartbeat_active"`
	ThresholdFunc      ThresholdConfig   `json:"threshold_func"`
	InternalEvidence   []EvidenceRef     `json:"internal_evidence"`
	ClaimVersions      []ClaimVersion    `json:"claim_versions"`
	DeduplicationIndex map[string]string `json:"dedup_index"` // source_id -> channel
	ScanCadence        time.Duration     `json:"scan_cadence"`
}

// ThresholdConfig governs adaptive heartbeat sensitivity.
// As internal evidence accumulates, the threshold rises, so mature
// anchors require larger Contextus activation deltas to nudge rho.
type ThresholdConfig struct {
	BaseThreshold  float64 `json:"base_threshold"`  // Initial proposed: 0.15
	MaturityScalar float64 `json:"maturity_scalar"` // Initial proposed: 0.05 per evidence item
}

// EvidenceRef is a pointer to CTH-internal evidence supporting an anchor.
type EvidenceRef struct {
	Type  string  `json:"type"` // "derivation", "measurement", "confluence",
	                            // "survey_support", "survey_contradict"
	RefID string  `json:"ref_id"`
	Bits  float64 `json:"bits"` // Confirmed bits contributed (negative for contradicting)
}

// =============================================================================
// Initial Proposed Bridge Parameters (subject to empirical tuning)
//
// These were proposed in Synthetic Trust Receipts v0.1 but are NOT
// design commitments. The synthetic receipt exercise revealed that
// decay half-life, heartbeat bit fraction, and internal evidence
// resistance should derive from per-insight properties (evidence cycle
// time, evidence independence, evidence formality) rather than be
// global constants.
// =============================================================================

const (
	DefaultRhoFloor                    = 0.05  // Visible but non-actionable
	DefaultHeartbeatBitFraction        = 0.10  // Confirmed bits per strengthening heartbeat
	DefaultBaseThreshold               = 0.15  // 15% activation delta before heartbeat fires
	DefaultMaturityScalar              = 0.05  // Threshold rises 5% per internal evidence item
	DefaultDecayHalfLifeDays           = 90    // Time for rho to halve in silence
	DefaultInternalEvidenceResistance  = 0.30  // Each derivation reduces decay by 30%
)
