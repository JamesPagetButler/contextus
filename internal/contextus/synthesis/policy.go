package synthesis

import (
	"github.com/JamesPagetButler/contextus/pkg/types"
)

// Policy decides whether an ephemeral session-scoped finding warrants
// promotion to a persistent InsightSignal, and what AnomalyKind to assign.
//
// Policy is the persistence-boundary contract per Spec v1.3 §4.4.
// Policy is intentionally simple in v1.3: a single confidence/significance
// threshold per source kind. The Theory v1.5 work (issue #3) may surface
// richer policy as worked examples accumulate.
type Policy struct {
	// EdgeBoundarySignificanceThreshold is the minimum
	// EdgeScoutFlag.Significance below which findings stay ephemeral.
	// Default 0.85 — calibrated against the surveillance-mode false-positive
	// budget; adjust based on production telemetry.
	EdgeBoundarySignificanceThreshold float64

	// CorpusDiversityVocabConcentrationThreshold is the minimum
	// CorpusDiversityReport.VocabConcentration above which a convergence
	// is durable enough to warrant persistence. Default 0.80 (= 80%
	// vocabulary monoculture).
	CorpusDiversityVocabConcentrationThreshold float64

	// BridgeInterventionConfidenceThreshold is the minimum
	// BridgeIntervention.Confidence below which a candidate cross-domain
	// connection stays ephemeral. Default 0.75.
	BridgeInterventionConfidenceThreshold float64

	// OperationalCorrelationScoreThreshold is the minimum
	// OperationalCorrelation.MaxScore above which a cross-domain
	// operational-telemetry pattern warrants persistence as an
	// AnomalyStructural signal. Default 0.70.
	//
	// 0.70 was chosen as a conservatively lower threshold than the edge-
	// boundary (0.85) because operational-telemetry divergence is a noisier
	// signal than exploration-boundary significance — the surveillance scout
	// should let more candidates through; Synthesis can use confidence to
	// attenuate the mint's downstream weight. Tune based on production
	// false-positive rates per Contextus-Spec-Addendum-NT-Scope-Operational §7.3.
	OperationalCorrelationScoreThreshold float64
}

// DefaultPolicy returns the v1.3/v1.4 starting threshold set. These values are
// initial guesses; production deployment should tune against observed
// false-positive / false-negative rates.
func DefaultPolicy() Policy {
	return Policy{
		EdgeBoundarySignificanceThreshold:          0.85,
		CorpusDiversityVocabConcentrationThreshold: 0.80,
		BridgeInterventionConfidenceThreshold:      0.75,
		OperationalCorrelationScoreThreshold:       0.70,
	}
}

// EvaluateEdgeBoundary returns the AnomalyKind for the minted signal, or the
// empty string if the flag does not warrant promotion. AnomalyStructural
// reflects that boundary blindspots are observations about the structure of
// the researcher's exploration boundary, not about the data (Theory v1.5,
// forthcoming, will formalise).
func (p Policy) EvaluateEdgeBoundary(flag types.EdgeScoutFlag) types.AnomalyKind {
	if flag.Significance >= p.EdgeBoundarySignificanceThreshold {
		return types.AnomalyStructural
	}
	return ""
}

// EvaluateCorpusDiversity returns the AnomalyKind for the minted signal, or
// the empty string if the report does not warrant promotion. Convergence-
// without-detection (ConvergenceDetected=false) never promotes regardless of
// concentration.
func (p Policy) EvaluateCorpusDiversity(report types.CorpusDiversityReport) types.AnomalyKind {
	if !report.ConvergenceDetected {
		return ""
	}
	if report.VocabConcentration >= p.CorpusDiversityVocabConcentrationThreshold {
		return types.AnomalyStructural
	}
	return ""
}

// EvaluateBridgeIntervention returns the AnomalyKind for the minted signal,
// or the empty string if the intervention does not warrant promotion.
// AnomalyNarrative for cross-domain causal stories is the typical mint;
// AnomalyStructural is reserved for cases where the bridging is purely
// structural (a topological connection across a domain boundary without an
// implied causal narrative). v1.3 ships only the narrative path; structural
// promotion via Bridge Agent is a Theory v1.5 follow-up.
func (p Policy) EvaluateBridgeIntervention(iv types.BridgeIntervention) types.AnomalyKind {
	if iv.Confidence >= p.BridgeInterventionConfidenceThreshold {
		return types.AnomalyNarrative
	}
	return ""
}

// EvaluateOperationalCorrelation returns AnomalyStructural if the correlation's
// MaxScore meets or exceeds the OperationalCorrelationScoreThreshold, or the
// empty string if the correlation stays ephemeral.
//
// An OperationalCorrelation must carry at least one referent and the MaxScore
// must be in [0.0, 1.0]. A correlation with no referents is always skipped
// (defensive; the scout should not emit empty-referent correlations).
//
// Per §7.4.2 of the NT_SCOPE_OPERATIONAL addendum and the OperationalCorrelation
// type doc, a cross-domain claim requires referents from at least 2 distinct
// ScopeIDs. A correlation whose referents all share the same ScopeID describes
// a single-domain anomaly and does not warrant an AnomalyStructural mint — it
// stays ephemeral so a higher-specificity signal kind (future work) can handle it.
//
// AnomalyStructural reflects that a cross-domain hardware↔cognitive or
// hardware↔algebraic correlation is a statement about the structure of the
// running system (Theory v1.5 §3.6.2 + NT_SCOPE_OPERATIONAL spec §7).
func (p Policy) EvaluateOperationalCorrelation(c types.OperationalCorrelation) types.AnomalyKind {
	if len(c.Referents) == 0 {
		return ""
	}
	if distinctScopeIDs(c.Referents) < 2 {
		return ""
	}
	if c.MaxScore >= p.OperationalCorrelationScoreThreshold {
		return types.AnomalyStructural
	}
	return ""
}

// distinctScopeIDs counts the number of unique ScopeID values across the
// given referent slice. Used to enforce the cross-domain requirement: an
// OperationalCorrelation must span ≥2 distinct ScopeIDs to qualify as a
// cross-domain claim (§7.4.2).
func distinctScopeIDs(refs []types.ScalarReferent) int {
	seen := make(map[string]struct{}, len(refs))
	for _, r := range refs {
		seen[r.ScopeID] = struct{}{}
	}
	return len(seen)
}
