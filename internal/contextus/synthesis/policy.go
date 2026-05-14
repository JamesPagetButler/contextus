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
}

// DefaultPolicy returns the v1.3 starting threshold set. These values are
// initial guesses; production deployment should tune against observed
// false-positive / false-negative rates.
func DefaultPolicy() Policy {
	return Policy{
		EdgeBoundarySignificanceThreshold:          0.85,
		CorpusDiversityVocabConcentrationThreshold: 0.80,
		BridgeInterventionConfidenceThreshold:      0.75,
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
