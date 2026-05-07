package types

import "testing"

func TestEvidenceTierCap(t *testing.T) {
	// Verbatim against Spec v1.3 §5.4.4 cap table.
	cases := []struct {
		tier string
		want int
	}{
		{TierCore, 50},
		{TierNear, 20},
		{TierPeripheral, 5},
		{TierDistant, 5},
		{TierSkeleton, 1},
		{"unknown-tier", 0},
		{"", 0},
	}
	for _, c := range cases {
		t.Run(c.tier, func(t *testing.T) {
			got := EvidenceTierCap(c.tier)
			if got != c.want {
				t.Errorf("EvidenceTierCap(%q): want %d, got %d", c.tier, c.want, got)
			}
		})
	}
}

func TestAnomalyKindConstants(t *testing.T) {
	// Locks the wire format. Changing these values would silently break
	// existing serialised signals on disk and over NATS.
	if AnomalyDensity != "density" {
		t.Errorf("AnomalyDensity wire value drifted: %s", AnomalyDensity)
	}
	if AnomalyCorrelation != "correlation" {
		t.Errorf("AnomalyCorrelation wire value drifted: %s", AnomalyCorrelation)
	}
	if AnomalyNarrative != "narrative" {
		t.Errorf("AnomalyNarrative wire value drifted: %s", AnomalyNarrative)
	}
	if AnomalyStructural != "structural" {
		t.Errorf("AnomalyStructural wire value drifted: %s", AnomalyStructural)
	}
}

func TestAgentClassConstants(t *testing.T) {
	// Spec v1.3 §11.1 + Wyrd issue #6 SignalSource correction:
	// only scout / correlation / synthesis are valid; never edge / corpus / bridge.
	if AgentScout != "scout" {
		t.Errorf("AgentScout wire value drifted: %s", AgentScout)
	}
	if AgentCorrelation != "correlation" {
		t.Errorf("AgentCorrelation wire value drifted: %s", AgentCorrelation)
	}
	if AgentSynthesis != "synthesis" {
		t.Errorf("AgentSynthesis wire value drifted: %s", AgentSynthesis)
	}
}
