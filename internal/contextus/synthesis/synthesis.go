package synthesis

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

// Subjects under which session-scoped agents publish ephemeral findings.
// Verbatim from Spec v1.3 §5.3 / §11.3.
const (
	SubjectEdgeBoundary       = "ctx.edge.boundary.>"       // wildcard for {session_id}
	SubjectCorpusDiversity    = "ctx.corpus.diversity"      // global
	SubjectBridgeIntervention = "ctx.bridge.intervention.>" // wildcard for {session_id}
)

// Agent is the Synthesis subscriber that watches the three ephemeral
// subjects and mints InsightSignals for findings that warrant persistence.
//
// Agent does not own its NATS connection; the caller wires Subscribe()
// against whichever subscriber abstraction is in play (real NATS in
// production; a fake in tests).
type Agent struct {
	Policy    Policy
	Persister Persister
	Now       func() time.Time // injectable clock for deterministic tests
}

// New returns an Agent with the default policy and a wall-clock Now.
func New(persister Persister) *Agent {
	return &Agent{
		Policy:    DefaultPolicy(),
		Persister: persister,
		Now:       time.Now,
	}
}

// HandleEdgeBoundary processes an EdgeScoutFlag from
// ctx.edge.boundary.{session_id}. Returns the minted signal's SignalID and
// true if the flag warranted promotion; the zero Addr and false otherwise.
func (a *Agent) HandleEdgeBoundary(ctx context.Context, flag types.EdgeScoutFlag) (types.Addr, bool, error) {
	kind := a.Policy.EvaluateEdgeBoundary(flag)
	if kind == "" {
		return types.Addr{}, false, nil
	}
	sig := a.signalFromEdgeBoundary(flag, kind)
	id, err := a.Persister.MintSignal(ctx, sig)
	if err != nil {
		return types.Addr{}, false, fmt.Errorf("synthesis: edge-boundary mint: %w", err)
	}
	return id, true, nil
}

// HandleCorpusDiversity processes a CorpusDiversityReport from
// ctx.corpus.diversity. Returns the minted signal's SignalID and true if
// the report warranted promotion.
func (a *Agent) HandleCorpusDiversity(ctx context.Context, report types.CorpusDiversityReport) (types.Addr, bool, error) {
	kind := a.Policy.EvaluateCorpusDiversity(report)
	if kind == "" {
		return types.Addr{}, false, nil
	}
	sig := a.signalFromCorpusDiversity(report, kind)
	id, err := a.Persister.MintSignal(ctx, sig)
	if err != nil {
		return types.Addr{}, false, fmt.Errorf("synthesis: corpus-diversity mint: %w", err)
	}
	return id, true, nil
}

// SubjectOperationalCorrelation is the NATS subject on which surveillance-mode
// scouts publish OperationalCorrelation events for Synthesis evaluation.
// Format: ctx.operational.correlation (global; no session suffix — operational
// surveillance is continuous, not session-scoped).
const SubjectOperationalCorrelation = "ctx.operational.correlation"

// HandleOperationalCorrelation processes an OperationalCorrelation from
// ctx.operational.correlation. Returns the minted signal's SignalID and true
// if the correlation warranted promotion as an AnomalyStructural signal.
//
// This is the Spec v1.4 §2.4 Referent path: surveillance-mode scouts compute
// scalar-referent divergence on operational-telemetry streams and publish
// OperationalCorrelation events; Synthesis decides whether the cross-domain
// pattern crosses the persistence-boundary threshold.
//
// See Contextus-Spec-Addendum-NT-Scope-Operational §7 + Spec v1.4 §2.4 for
// the full cross-domain Synthesis pattern.
func (a *Agent) HandleOperationalCorrelation(ctx context.Context, c types.OperationalCorrelation) (types.Addr, bool, error) {
	kind := a.Policy.EvaluateOperationalCorrelation(c)
	if kind == "" {
		return types.Addr{}, false, nil
	}
	sig := a.signalFromOperationalCorrelation(c, kind)
	id, err := a.Persister.MintSignal(ctx, sig)
	if err != nil {
		return types.Addr{}, false, fmt.Errorf("synthesis: operational-correlation mint: %w", err)
	}
	return id, true, nil
}

// HandleBridgeIntervention processes a BridgeIntervention from
// ctx.bridge.intervention.{session_id}. Returns the minted signal's SignalID
// and true if the intervention warranted promotion.
func (a *Agent) HandleBridgeIntervention(ctx context.Context, iv types.BridgeIntervention) (types.Addr, bool, error) {
	kind := a.Policy.EvaluateBridgeIntervention(iv)
	if kind == "" {
		return types.Addr{}, false, nil
	}
	sig := a.signalFromBridgeIntervention(iv, kind)
	id, err := a.Persister.MintSignal(ctx, sig)
	if err != nil {
		return types.Addr{}, false, fmt.Errorf("synthesis: bridge-intervention mint: %w", err)
	}
	return id, true, nil
}

// signalFromEdgeBoundary converts an EdgeScoutFlag into an InsightSignal.
// Anomaly score is the flag's significance (per Theory §3.6.2 information-
// theoretic interpretation). Confidence equals significance for v1.3; future
// revisions may decouple them.
func (a *Agent) signalFromEdgeBoundary(flag types.EdgeScoutFlag, kind types.AnomalyKind) types.InsightSignal {
	now := a.Now()
	return types.InsightSignal{
		SignalID:           deterministicAddr("edge-boundary", flag.SessionID, flag.BoundaryNodeID, flag.UnexploredDomain),
		AgentType:          types.AgentSynthesis,
		AnomalyType:        kind,
		AnomalyScore:       flag.Significance,
		Confidence:         flag.Significance,
		ConfidenceVariance: 0,
		Persistence:        0,    // first emission
		Structural:         true, // boundary-blindspot is structural-by-construction
		FirstSeen:          now,
		LastSeen:           now,
		Version:            1,
		Promoted:           false,
		Subgraph:           []types.HyperedgeRef{}, // populated when Wyrd contextus/ subpackage lands
		TraversalPath:      []types.Addr{},         // session-id encoded in SignalID derivation; explicit path is post-Phase-3
	}
}

func (a *Agent) signalFromCorpusDiversity(report types.CorpusDiversityReport, kind types.AnomalyKind) types.InsightSignal {
	now := a.Now()
	return types.InsightSignal{
		SignalID:           deterministicAddr("corpus-diversity", report.SearchID, report.OriginalQuery),
		AgentType:          types.AgentSynthesis,
		AnomalyType:        kind,
		AnomalyScore:       report.VocabConcentration,
		Confidence:         report.VocabConcentration,
		ConfidenceVariance: 0,
		Persistence:        0,
		Structural:         true, // diversity convergence describes search-process structure
		FirstSeen:          now,
		LastSeen:           now,
		Version:            1,
		Promoted:           false,
		Subgraph:           []types.HyperedgeRef{},
		TraversalPath:      []types.Addr{},
	}
}

func (a *Agent) signalFromBridgeIntervention(iv types.BridgeIntervention, kind types.AnomalyKind) types.InsightSignal {
	now := a.Now()
	return types.InsightSignal{
		SignalID:           deterministicAddr("bridge-intervention", iv.SessionID, iv.ScoutFlagID, iv.KnowledgeDomainGap),
		AgentType:          types.AgentSynthesis,
		AnomalyType:        kind,
		AnomalyScore:       iv.Confidence,
		Confidence:         iv.Confidence,
		ConfidenceVariance: 0,
		Persistence:        0,
		Structural:         false, // narrative bridge is exploration-tier, not landscape
		FirstSeen:          now,
		LastSeen:           now,
		Version:            1,
		Promoted:           false,
		Subgraph:           []types.HyperedgeRef{},
		TraversalPath:      []types.Addr{},
	}
}

// signalFromOperationalCorrelation converts an OperationalCorrelation into an
// InsightSignal. AnomalyScore is MaxScore (the peak referent divergence across
// the correlated hardware subsystems). Confidence equals AnomalyScore for v1.4;
// future revisions may factor in observation count and window length.
//
// The Referents field on the returned signal carries the full
// predicted-vs-observed pairs per Spec v1.4 §2.4 — the downstream consumer
// (BMA Conscious-A/B forensic recall; CTH for hypothesis input) can inspect
// both the summary score and the individual referent details.
//
// Structural=true reflects that operational↔cognitive or operational↔algebraic
// correlations are about the structure of the running system (Theory v1.5 §3.6.2).
func (a *Agent) signalFromOperationalCorrelation(c types.OperationalCorrelation, kind types.AnomalyKind) types.InsightSignal {
	now := a.Now()
	return types.InsightSignal{
		SignalID:           deterministicAddr("operational-correlation", c.CorrelationID),
		AgentType:          types.AgentSynthesis,
		AnomalyType:        kind,
		AnomalyScore:       c.MaxScore,
		Confidence:         c.MaxScore,
		ConfidenceVariance: 0,
		Persistence:        0,
		Structural:         true, // runtime-structural: observation about the running system
		FirstSeen:          now,
		LastSeen:           now,
		Version:            1,
		Promoted:           false,
		Subgraph:           []types.HyperedgeRef{},
		TraversalPath:      []types.Addr{},
		Referents:          c.Referents,
	}
}

// deterministicAddr derives a 16-byte content address from the source-kind
// and the source identifiers, so re-evaluation of the same ephemeral finding
// produces the same SignalID. This is the v1.3 implementation of the
// MintSignal idempotency contract (see persister.go).
//
// Replace with a real q8.Addr derivation once the QBP Locale library is
// imported.
func deterministicAddr(parts ...string) types.Addr {
	h := sha256.New()
	for i, p := range parts {
		// Length-prefix each part to avoid concatenation collisions.
		var lenBuf [8]byte
		binary.BigEndian.PutUint64(lenBuf[:], uint64(len(p)))
		_, _ = h.Write(lenBuf[:])
		_, _ = h.Write([]byte(p))
		// Index-prefix to avoid swap collisions across parts.
		var idxBuf [4]byte
		binary.BigEndian.PutUint32(idxBuf[:], uint32(i))
		_, _ = h.Write(idxBuf[:])
	}
	sum := h.Sum(nil)
	var addr types.Addr
	copy(addr[:], sum[:16]) // first 128 bits of SHA-256
	return addr
}
