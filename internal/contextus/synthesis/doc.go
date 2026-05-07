// Package synthesis implements the Synthesis agent's persistence-boundary
// role per Contextus-Spec-v1.3 §4.4.
//
// The session-scoped agents (Edge Scout, Corpus Edge Scout, Bridge Agent) emit
// ephemeral NATS events on:
//
//	ctx.edge.boundary.{session_id}     — EdgeScoutFlag
//	ctx.corpus.diversity               — CorpusDiversityReport
//	ctx.bridge.intervention.{session_id} — BridgeIntervention
//
// These events are session-scoped and never persisted as InsightSignal nodes.
// When a session-scoped finding warrants graph-level persistence (e.g., a
// Bridge Agent convergence above a confidence threshold representing a
// genuine cross-domain pattern rather than a transient boundary effect), the
// Synthesis agent — and only the Synthesis agent — mints an InsightSignal.
//
// This is the persistence-boundary discipline:
//
//   - Edge Scout / Corpus Edge Scout / Bridge Agent never write InsightSignals.
//   - Synthesis is the sole gate at which ephemeral findings become hypergraph-
//     resident records.
//   - Minted signals carry agent_type=synthesis and anomaly_type=narrative
//     (Bridge convergence) or anomaly_type=structural (boundary blindspot,
//     diversity convergence).
//
// The threshold policy and minting logic are local to this package; the Wyrd
// write side is delegated to a Persister interface (see persister.go) so the
// substrate can be swapped without changing the policy.
//
// Spec cross-references:
//   - §4.4 Insight Signal Emission Pipeline (Synthesis as persistence boundary)
//   - §4.5 Session-Scoped Agents (the three ephemeral subjects)
//   - §11.1 InsightSignal, AnomalyKind, AgentClass
//   - §11.3 EdgeScoutFlag, CorpusDiversityReport, BridgeIntervention
//
// Theory cross-references:
//   - Theory v1.4 §3.6.2 (anomaly classes; v1.5 will formalise structural)
//   - Theory v1.4 §8 (surveillance mode; the source of session-scoped findings)
package synthesis
