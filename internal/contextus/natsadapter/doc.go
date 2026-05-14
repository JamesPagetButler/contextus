// Package natsadapter wires a synthesis.Agent to a NATS broker.
//
// The adapter is intentionally a thin glue layer: it owns the NATS connection,
// decodes JSON payloads into the typed Spec v1.3 §11.3 structs (EdgeScoutFlag,
// CorpusDiversityReport, BridgeIntervention), and dispatches each into the
// matching Agent.Handle* method.
//
// The adapter does NOT own the persistence-boundary policy (that lives in
// synthesis.Policy) and does NOT own the Wyrd write side (that lives behind
// synthesis.Persister). Swapping the broker, the policy, or the persister is
// each a concern of one layer.
//
// Subscriber is the consumer-side interface (per go-coding-guide.md "Accept
// interfaces, return concrete types"). The real NATS client implements it; a
// fake satisfies it for tests so we don't need an embedded NATS server in the
// test binary.
//
// Spec cross-references:
//   - §5.3 NATS Subject Hierarchy (the three subjects this adapter listens on)
//   - §4.4 Synthesis as Persistence Boundary (the dispatch target)
//   - §11.3 EdgeScoutFlag, CorpusDiversityReport, BridgeIntervention payloads
package natsadapter
