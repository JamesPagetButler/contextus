package synthesis

import (
	"context"
	"testing"
	"time"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

func fixedNow() time.Time { return time.Date(2026, 5, 6, 22, 0, 0, 0, time.UTC) }

func newTestAgent() (*Agent, *NoopPersister) {
	p := &NoopPersister{}
	a := New(p)
	a.Now = fixedNow
	return a, p
}

func TestEdgeBoundary_BelowThreshold_Skipped(t *testing.T) {
	a, p := newTestAgent()
	flag := types.EdgeScoutFlag{
		SessionID:      "s1",
		BoundaryNodeID: "n-stream-soda-butte",
		Significance:   0.50, // below default 0.85
	}
	id, ok, err := a.HandleEdgeBoundary(context.Background(), flag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Errorf("expected promoted=false; got true with id %x", id)
	}
	if p.Calls != 0 {
		t.Errorf("expected Persister.MintSignal not called; got %d calls", p.Calls)
	}
}

func TestEdgeBoundary_AboveThreshold_MintsStructural(t *testing.T) {
	a, p := newTestAgent()
	flag := types.EdgeScoutFlag{
		SessionID:        "s1",
		BoundaryNodeID:   "n-stream-soda-butte",
		UnexploredDomain: "vegetation-phenology",
		Significance:     0.92,
	}
	id, ok, err := a.HandleEdgeBoundary(context.Background(), flag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected promoted=true; got false")
	}
	if p.Calls != 1 {
		t.Fatalf("expected exactly one MintSignal call; got %d", p.Calls)
	}
	got := p.LastSignal
	if got.AgentType != types.AgentSynthesis {
		t.Errorf("agent_type: want synthesis, got %s", got.AgentType)
	}
	if got.AnomalyType != types.AnomalyStructural {
		t.Errorf("anomaly_type: want structural, got %s", got.AnomalyType)
	}
	if got.AnomalyScore != 0.92 {
		t.Errorf("anomaly_score: want 0.92, got %v", got.AnomalyScore)
	}
	if got.Confidence != 0.92 {
		t.Errorf("confidence: want 0.92, got %v", got.Confidence)
	}
	if !got.Structural {
		t.Errorf("structural bool: want true (boundary blindspot), got false")
	}
	if got.FirstSeen != fixedNow() {
		t.Errorf("first_seen: want fixedNow, got %v", got.FirstSeen)
	}
	if got.SignalID != id {
		t.Errorf("returned id %x != minted signal id %x", id, got.SignalID)
	}
}

func TestEdgeBoundary_DeterministicID(t *testing.T) {
	// Idempotency contract: two evaluations of the same flag produce the
	// same SignalID. The Persister should treat that as the dedup key.
	a, _ := newTestAgent()
	flag := types.EdgeScoutFlag{
		SessionID:        "s1",
		BoundaryNodeID:   "n-stream-soda-butte",
		UnexploredDomain: "vegetation-phenology",
		Significance:     0.92,
	}
	id1, _, _ := a.HandleEdgeBoundary(context.Background(), flag)
	id2, _, _ := a.HandleEdgeBoundary(context.Background(), flag)
	if id1 != id2 {
		t.Errorf("expected deterministic SignalID across evaluations; got %x vs %x", id1, id2)
	}
}

func TestCorpusDiversity_NotConverged_Skipped(t *testing.T) {
	a, p := newTestAgent()
	report := types.CorpusDiversityReport{
		SearchID:            "search-1",
		ConvergenceDetected: false,
		VocabConcentration:  0.95, // high concentration but convergence flag off
	}
	_, ok, err := a.HandleCorpusDiversity(context.Background(), report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected promoted=false when ConvergenceDetected=false; got true")
	}
	if p.Calls != 0 {
		t.Errorf("expected MintSignal not called; got %d", p.Calls)
	}
}

func TestCorpusDiversity_Converged_BelowThreshold_Skipped(t *testing.T) {
	a, p := newTestAgent()
	report := types.CorpusDiversityReport{
		SearchID:            "search-1",
		ConvergenceDetected: true,
		VocabConcentration:  0.50, // below default 0.80
	}
	_, ok, _ := a.HandleCorpusDiversity(context.Background(), report)
	if ok {
		t.Error("expected promoted=false below threshold; got true")
	}
	if p.Calls != 0 {
		t.Errorf("expected MintSignal not called; got %d", p.Calls)
	}
}

func TestCorpusDiversity_Converged_AboveThreshold_MintsStructural(t *testing.T) {
	a, p := newTestAgent()
	report := types.CorpusDiversityReport{
		SearchID:            "search-1",
		OriginalQuery:       "where is the missing water in the Colorado River?",
		ConvergenceDetected: true,
		VocabConcentration:  0.85,
	}
	_, ok, err := a.HandleCorpusDiversity(context.Background(), report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected promoted=true; got false")
	}
	got := p.LastSignal
	if got.AgentType != types.AgentSynthesis {
		t.Errorf("agent_type: want synthesis, got %s", got.AgentType)
	}
	if got.AnomalyType != types.AnomalyStructural {
		t.Errorf("anomaly_type: want structural, got %s", got.AnomalyType)
	}
	if got.AnomalyScore != 0.85 {
		t.Errorf("anomaly_score: want 0.85, got %v", got.AnomalyScore)
	}
}

func TestBridgeIntervention_BelowThreshold_Skipped(t *testing.T) {
	a, p := newTestAgent()
	iv := types.BridgeIntervention{
		SessionID:  "s1",
		Confidence: 0.50, // below default 0.75
	}
	_, ok, _ := a.HandleBridgeIntervention(context.Background(), iv)
	if ok {
		t.Error("expected promoted=false below threshold; got true")
	}
	if p.Calls != 0 {
		t.Errorf("expected MintSignal not called; got %d", p.Calls)
	}
}

func TestBridgeIntervention_AboveThreshold_MintsNarrative(t *testing.T) {
	a, p := newTestAgent()
	iv := types.BridgeIntervention{
		SessionID:            "s1",
		ScoutFlagID:          "flag-snowpack-streamflow-2000",
		ScoutPhysicalDomain:  "hydrology",
		SearchPhysicalDomain: "plant-physiology",
		KnowledgeDomainGap:   "vegetation-as-passive-vs-active",
		ProposedQuery:        "what extracts groundwater in stratified watersheds besides wells?",
		Confidence:           0.88,
	}
	_, ok, err := a.HandleBridgeIntervention(context.Background(), iv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected promoted=true; got false")
	}
	got := p.LastSignal
	if got.AgentType != types.AgentSynthesis {
		t.Errorf("agent_type: want synthesis, got %s", got.AgentType)
	}
	if got.AnomalyType != types.AnomalyNarrative {
		t.Errorf("anomaly_type: want narrative (cross-domain causal story), got %s", got.AnomalyType)
	}
	if got.AnomalyScore != 0.88 {
		t.Errorf("anomaly_score: want 0.88, got %v", got.AnomalyScore)
	}
	if got.Structural {
		t.Errorf("structural bool: want false (narrative bridge is exploration-tier), got true")
	}
}

func TestPolicy_DefaultThresholds(t *testing.T) {
	p := DefaultPolicy()
	if p.EdgeBoundarySignificanceThreshold != 0.85 {
		t.Errorf("EdgeBoundarySignificanceThreshold: want 0.85, got %v", p.EdgeBoundarySignificanceThreshold)
	}
	if p.CorpusDiversityVocabConcentrationThreshold != 0.80 {
		t.Errorf("CorpusDiversityVocabConcentrationThreshold: want 0.80, got %v", p.CorpusDiversityVocabConcentrationThreshold)
	}
	if p.BridgeInterventionConfidenceThreshold != 0.75 {
		t.Errorf("BridgeInterventionConfidenceThreshold: want 0.75, got %v", p.BridgeInterventionConfidenceThreshold)
	}
	if p.OperationalCorrelationScoreThreshold != 0.70 {
		t.Errorf("OperationalCorrelationScoreThreshold: want 0.70, got %v", p.OperationalCorrelationScoreThreshold)
	}
}

// TestOperationalCorrelation_BelowThreshold_Skipped verifies that an
// OperationalCorrelation with MaxScore below the policy threshold stays
// ephemeral (no mint, promoted=false). Uses two distinct ScopeIDs so the
// ≥2-ScopeID cross-domain guard passes and the score gate is the deciding
// condition.
//
// AC-6: cross-domain correlation test — below-threshold case.
func TestOperationalCorrelation_BelowThreshold_Skipped(t *testing.T) {
	a, p := newTestAgent()
	corr := types.OperationalCorrelation{
		CorrelationID: "corr-bma-prime-below",
		Referents: []types.ScalarReferent{
			{
				Label:     "cpu_temp_5min_avg_celsius",
				Predicted: 65.0,
				Observed:  67.0,
				Score:     0.031, // |67-65|/65 ≈ 0.031 — well below 0.70
				ScopeID:   "contextus:scope:operational:host-bma-prime-cpu",
			},
			{
				Label:     "disk_smart_reallocated_sectors",
				Predicted: 0.0,
				Observed:  0.0,
				Score:     0.0, // no reallocated sectors — also below 0.70
				ScopeID:   "contextus:scope:operational:host-bma-prime-disk",
			},
		},
		MaxScore: 0.031,
	}
	_, ok, err := a.HandleOperationalCorrelation(context.Background(), corr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected promoted=false below threshold; got true")
	}
	if p.Calls != 0 {
		t.Errorf("expected MintSignal not called; got %d", p.Calls)
	}
}

// TestOperationalCorrelation_NoReferents_Skipped verifies that an
// OperationalCorrelation with an empty Referents slice is always skipped,
// regardless of MaxScore.
//
// AC-6: cross-domain correlation test — empty-referent guard.
func TestOperationalCorrelation_NoReferents_Skipped(t *testing.T) {
	a, p := newTestAgent()
	corr := types.OperationalCorrelation{
		CorrelationID: "corr-empty-referents",
		Referents:     nil,
		MaxScore:      0.99,
	}
	_, ok, err := a.HandleOperationalCorrelation(context.Background(), corr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected promoted=false for empty-referent correlation; got true")
	}
	if p.Calls != 0 {
		t.Errorf("expected MintSignal not called; got %d", p.Calls)
	}
}

// TestOperationalCorrelation_MintsStructuralWithReferent verifies the full
// AC-6 path: when a cross-domain OperationalCorrelation (cpu_temp trending
// with algebraic-integrity drift as a referent proxy) has MaxScore above the
// threshold, Synthesis mints an AnomalyStructural InsightSignal and the
// signal carries Referents with non-zero scores.
//
// This is the Spec v1.4 §2.4 Referent promotion path: predicted-vs-observed
// scalar divergence on operational-telemetry streams drives AnomalyStructural
// minting. See Contextus-Spec-Addendum-NT-Scope-Operational §7 worked
// pattern A (thermal stress ↔ algebraic-integrity drift).
//
// AC-6: cross-domain correlation test — above-threshold case.
func TestOperationalCorrelation_MintsStructuralWithReferent(t *testing.T) {
	a, p := newTestAgent()

	// Worked Pattern A (spec addendum §7.1): thermal stress ↔ algebraic-integrity drift.
	// CPU temperature well above predicted; disk SMART warning trending up.
	corr := types.OperationalCorrelation{
		CorrelationID: "corr-bma-prime-thermal-algebraic-2026-06-12",
		Referents: []types.ScalarReferent{
			{
				Label:     "cpu_temp_5min_avg_celsius",
				Predicted: 65.0,
				Observed:  82.0, // thermal stress: |82-65|/65 ≈ 0.26 → score 0.26
				Score:     0.26,
				ScopeID:   "contextus:scope:operational:host-bma-prime-cpu",
			},
			{
				Label:     "disk_smart_reallocated_sectors",
				Predicted: 0.0,
				Observed:  3.0, // non-zero reallocated sectors → capped score 1.0
				Score:     1.0,
				ScopeID:   "contextus:scope:operational:host-bma-prime-disk",
			},
		},
		MaxScore: 1.0, // max across referents
	}

	id, ok, err := a.HandleOperationalCorrelation(context.Background(), corr)
	if err != nil {
		t.Fatalf("HandleOperationalCorrelation: unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected promoted=true for above-threshold operational correlation; got false")
	}
	if p.Calls != 1 {
		t.Fatalf("expected exactly one MintSignal call; got %d", p.Calls)
	}

	got := p.LastSignal

	// Signal metadata checks.
	if got.AgentType != types.AgentSynthesis {
		t.Errorf("agent_type: want %q, got %q", types.AgentSynthesis, got.AgentType)
	}
	if got.AnomalyType != types.AnomalyStructural {
		t.Errorf("anomaly_type: want %q (runtime-structural), got %q", types.AnomalyStructural, got.AnomalyType)
	}
	if got.AnomalyScore != 1.0 {
		t.Errorf("anomaly_score: want 1.0 (MaxScore), got %v", got.AnomalyScore)
	}
	if !got.Structural {
		t.Error("structural bool: want true (runtime-structural: about the running system), got false")
	}
	if got.FirstSeen != fixedNow() {
		t.Errorf("first_seen: want fixedNow, got %v", got.FirstSeen)
	}
	if got.SignalID != id {
		t.Errorf("returned id %x != minted signal id %x", id, got.SignalID)
	}

	// Referent checks — this is the AC-6 core: the signal must carry the
	// predicted-vs-observed pairs so downstream consumers can inspect both
	// the aggregate score and the per-subsystem divergence.
	if got.Referents == nil {
		t.Fatal("Referents: want non-nil slice on AnomalyStructural operational signal; got nil")
	}
	if len(got.Referents) != 2 {
		t.Fatalf("Referents: want 2 entries (cpu + disk); got %d", len(got.Referents))
	}

	// Index by ScopeID for order-independent assertions.
	refByScope := make(map[string]types.ScalarReferent, len(got.Referents))
	for _, r := range got.Referents {
		refByScope[r.ScopeID] = r
	}

	cpu, ok := refByScope["contextus:scope:operational:host-bma-prime-cpu"]
	if !ok {
		t.Error("Referents: missing cpu scope referent")
	} else {
		if cpu.Score <= 0 {
			t.Errorf("cpu referent score: want > 0, got %v", cpu.Score)
		}
		if cpu.Observed <= cpu.Predicted {
			t.Errorf("cpu referent: Observed (%v) should exceed Predicted (%v) for thermal-stress pattern", cpu.Observed, cpu.Predicted)
		}
	}

	disk, ok := refByScope["contextus:scope:operational:host-bma-prime-disk"]
	if !ok {
		t.Error("Referents: missing disk scope referent")
	} else {
		if disk.Score != 1.0 {
			t.Errorf("disk referent score: want 1.0 (capped, non-zero reallocated sectors), got %v", disk.Score)
		}
	}
}

// TestOperationalCorrelation_DeterministicID verifies the §7.4.2 idempotency
// contract: two evaluations of the same OperationalCorrelation (same
// CorrelationID) produce the same SignalID. Synthesis must treat the SignalID
// as the dedup key.
//
// The correlation must be above-threshold AND span ≥2 distinct ScopeIDs so
// that deterministicAddr is actually invoked and both ids are real minted
// addresses — not the zero Addr from a short-circuit skip.
//
// AC-6: deterministic SignalID for operational correlations (C1 fix: was
// MaxScore 0.277 below the 0.70 threshold — test passed vacuously on zero
// Addr equality without ever exercising deterministicAddr).
func TestOperationalCorrelation_DeterministicID(t *testing.T) {
	a, _ := newTestAgent()
	// Pattern-A correlation (thermal stress ↔ algebraic-integrity drift):
	// MaxScore 1.0 > 0.70 threshold; two distinct ScopeIDs → cross-domain guard passes.
	corr := types.OperationalCorrelation{
		CorrelationID: "corr-bma-prime-determinism-test",
		Referents: []types.ScalarReferent{
			{Label: "cpu_temp", Predicted: 65.0, Observed: 83.0, Score: 0.26,
				ScopeID: "contextus:scope:operational:host-bma-prime-cpu"},
			{Label: "disk_smart_reallocated_sectors", Predicted: 0.0, Observed: 3.0, Score: 1.0,
				ScopeID: "contextus:scope:operational:host-bma-prime-disk"},
		},
		MaxScore: 1.0,
	}
	id1, ok1, err1 := a.HandleOperationalCorrelation(context.Background(), corr)
	if err1 != nil {
		t.Fatalf("first evaluation: unexpected error: %v", err1)
	}
	if !ok1 {
		t.Fatal("first evaluation: expected promoted=true (above-threshold, cross-domain); got false — deterministicAddr was never exercised")
	}
	id2, ok2, err2 := a.HandleOperationalCorrelation(context.Background(), corr)
	if err2 != nil {
		t.Fatalf("second evaluation: unexpected error: %v", err2)
	}
	if !ok2 {
		t.Fatal("second evaluation: expected promoted=true; got false")
	}
	if id1 != id2 {
		t.Errorf("expected deterministic SignalID across evaluations; got %x vs %x", id1, id2)
	}
}

// TestOperationalCorrelation_SingleScopeID_Skipped verifies that a correlation
// whose referents all share the same ScopeID is rejected as a single-domain
// observation — it does not qualify as a cross-domain AnomalyStructural claim.
//
// §7.4.2 + OperationalCorrelation type doc: cross-domain correlation requires
// referents from ≥2 distinct ScopeIDs. Enforcement is in
// EvaluateOperationalCorrelation (G1 fix).
//
// AC-6: ≥2-distinct-ScopeID guard.
func TestOperationalCorrelation_SingleScopeID_Skipped(t *testing.T) {
	a, p := newTestAgent()
	// Both referents report from the same ScopeID (single hardware domain).
	// MaxScore is well above the 0.70 score threshold to confirm that the
	// ScopeID guard fires before the score gate.
	corr := types.OperationalCorrelation{
		CorrelationID: "corr-single-scope-rejection-test",
		Referents: []types.ScalarReferent{
			{Label: "cpu_temp_5min_avg", Predicted: 65.0, Observed: 83.0, Score: 0.26,
				ScopeID: "contextus:scope:operational:host-bma-prime-cpu"},
			{Label: "cpu_temp_10min_avg", Predicted: 66.0, Observed: 84.0, Score: 0.27,
				ScopeID: "contextus:scope:operational:host-bma-prime-cpu"}, // same ScopeID
		},
		MaxScore: 0.90,
	}
	_, ok, err := a.HandleOperationalCorrelation(context.Background(), corr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected promoted=false for single-ScopeID correlation (not cross-domain); got true")
	}
	if p.Calls != 0 {
		t.Errorf("expected MintSignal not called for single-ScopeID correlation; got %d calls", p.Calls)
	}
}
