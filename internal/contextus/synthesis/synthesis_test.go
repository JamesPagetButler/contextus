package synthesis

import (
	"context"
	"testing"
	"time"

	"github.com/JamesPagetButler/contextus/internal/contextus/types"
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
}
