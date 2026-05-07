# Cross-Domain Hyperedge Minting — Design Spec

**Status:** Design only. Implementation blocked on Wyrd `contextus/` subpackage (Wyrd issue [#6](https://github.com/JamesPagetButler/wyrd/issues/6) corrections + Wyrd PR [#19](https://github.com/JamesPagetButler/wyrd/pull/19) bridge-batch).

**Author:** contextus-impl, 2026-05-07

**Closes the gap:** Spec v1.3 §4.4 makes Synthesis the persistence boundary for ephemeral session-scoped findings. When a Bridge Agent convergence promotes, Spec v1.3 mints a single `NT_INSIGHT_SIGNAL` of `agent_type = synthesis`. But the architectural value of a Bridge Agent convergence is precisely that *N domains agreed on a cross-domain pattern* — the irreducible joint information per `Wyrd.HolographicHypergraphHigherArity.theorem2_irreducibility_n_arity` (Phase 4 v1.5). A single signal does not capture that. This document specifies how to *also* mint an arity-N `model.Hyperedge` connecting the involved Insight Signals so that Theorem 2's irreducibility surface holds end-to-end.

---

## 1. Why arity ≥ 3 matters

Wyrd's [`model.Hyperedge`](https://github.com/JamesPagetButler/wyrd/blob/main/model/hyperedge.go) godoc is explicit:

> Soundness — irreducibility: per `Wyrd.HolographicHypergraph.theorem2_irreducibility` (Phase 4 v1.4) and the higher-arity generalisation `Wyrd.HolographicHypergraphHigherArity.theorem2_irreducibility_n_arity` (Phase 4 v1.5), a hyperedge of arity k ≥ 3 carries information that CANNOT be recovered from any decomposition into smaller-arity edges. Consumers MUST NOT silently split a k-edge into k(k−1)/2 pair edges at storage or transit time — the joint constraint is information that pair edges cannot encode.

A Bridge Agent intervention by construction involves three or more domain perspectives converging on the same physical or process node. Examples:

- *Hydrology + plant physiology + climate* converging on the snowpack-streamflow correlation gap (the Colorado River failure, retrospectively).
- *Marine biology + oceanography + thermodynamics* converging on the chlorophyll-a × dissolved-oxygen confluence for whale shark neonate habitat.
- *Trophic ecology + stream geomorphology + behavioural ecology* converging on the wolf reintroduction cascade.

In each case, the *joint* observation across N domains is what makes the convergence valuable. Decomposing it into pair edges (hydrology↔plant, plant↔climate, climate↔hydrology) loses the irreducible-joint information that the arity-N edge encodes. The `BridgeIntervention.ScoutFlagID + Active Search Trajectory` together identify N specific Insight Signals; these must be connected by a single arity-N hyperedge.

---

## 2. The minting trigger

Synthesis already evaluates `BridgeIntervention.Confidence ≥ 0.75` (Spec v1.3 §4.4 / `synthesis.Policy`). When that threshold is met, the current implementation (PR [#4](https://github.com/JamesPagetButler/contextus/pull/4)) mints **one** `NT_INSIGHT_SIGNAL` of `anomaly_type = narrative`. This design extends the same trigger to **also** mint:

1. The original `NT_INSIGHT_SIGNAL` (unchanged).
2. **Plus:** one `model.Hyperedge` of arity ≥ 3, connecting the N Insight Signals that the Bridge Agent's surveillance + search trajectory inventories identified as the convergence participants.

The new Hyperedge has:

| Field | Value |
|---|---|
| `Type` | `"contextus.bridge.convergence"` |
| `Nodes` | `[]NodeID` — the N participating Insight Signal IDs (sorted by SignalID for deterministic equality) |
| `Weight` | `model.NewComplexWeight(score, 0)` where `score = BridgeIntervention.Confidence` |
| `IsSymmetric` | `true` — the convergence is direction-agnostic; participating signals are siblings, not source/sink |
| `Created` | `time.Now()` (or the deterministic clock under test) |

The Hyperedge is **not** a separate persistence layer; it is part of the same atomic mint. Either both the signal and the hyperedge land, or neither does.

---

## 3. Atomicity contract

The Hyperedge mint must be atomic with the Insight Signal mint. Three failure modes the design must prevent:

- **Signal-without-edge:** signal lands, hyperedge fails. Consumers see a `synthesis`-authored signal with no convergence record. Bad: violates the architectural contract that arity-N convergence is part of the signal's evidence.
- **Edge-without-signal:** hyperedge lands referencing a signal NodeID that doesn't exist. Bad: reference integrity broken.
- **Partial-N-edge:** hyperedge minted with k < N nodes because some referenced signals weren't yet persisted. Bad: silently violates Theorem 2 irreducibility — the resulting edge is decomposable.

The mechanism is `compute.Bridge.PromoteBatch` (Wyrd PR #19, design-gated under §I4 review). The batch contains:
1. The new `NT_INSIGHT_SIGNAL` node.
2. The arity-N `model.Hyperedge` node (with all N participating signal IDs).

The batch primitive's all-or-nothing semantics (Wyrd PR #19; §I4 ratification by bma + bma-implementor) is exactly what this design needs. Non-atomic batch creates the *apparent-completion-without-completion* window that I3 prohibits — same reasoning as bma's seq=22 framing.

**Open question for Wyrd:** the existing N participating signals must already be persisted in the graph before this batch can reference them by NodeID. If a Bridge Agent convergence references a signal that's still being minted in a parallel-running Synthesis goroutine, the batch will fail with `ErrBridgeUnknownEdge` (or its node-side equivalent). Two acceptable handlings:

- **(a) Two-phase mint.** Round 1 of the batch mints all signals; round 2 mints the hyperedge. Wyrd would need `Bridge.PromoteSequence([]Batch)` rather than `PromoteBatch(Batch)` — heavier.
- **(b) Single-phase mint with all participants in one batch.** The batch contains all N signal nodes plus the arity-N hyperedge that references them. Wyrd's `Graph.AddHyperedge` either is willing to accept node references that are also in the same batch, or `PromoteBatch` is documented as resolving forward-references within the batch.

(b) is simpler. Recommend that wyrd-implementor's `bridge-batch.md` v0.1 design doc clarify that `PromoteBatch` resolves intra-batch references before committing.

---

## 4. Identifying the N participating signals

The Bridge Agent's `BridgeIntervention` does not currently carry a list of participating signal IDs — only `ScoutFlagID` (one flag) and `KnowledgeDomainGap` (one boundary). A v1.3 spec gap.

Two ways to close it:

- **Synthesis re-derives.** When Synthesis evaluates a `BridgeIntervention`, it traverses its own active-signal store, finds Insight Signals whose subgraphs intersect the `ScoutPhysicalDomain` and `SearchPhysicalDomain` tags, and uses the union as the N participants. Cost: one synthesis-side traversal per evaluation. Risk: depends on Synthesis having a current view of the active-signal store.
- **Bridge Agent enriches.** Extend `BridgeIntervention` with a `ParticipatingSignalIDs []string` field populated when the intervention is computed. Costs: BridgeIntervention struct change (Spec v1.3 §11.3 schema delta) + Bridge Agent implementation work. Benefits: the participants are determined by the agent that actually identified the convergence; Synthesis is not re-doing the work.

**Recommendation: Bridge Agent enriches.** The Bridge Agent already maintains the surveillance and search inventories that name the convergence participants (Spec v1.3 §4.5 Bridge Agent Session Lifecycle step 2). Adding the field is a small schema delta and a cleaner separation of concerns.

This is a v1.4 spec amendment. Filing as forthcoming work.

---

## 5. Implementation interface

The implementation lives behind the existing `synthesis.Persister` interface, augmented with a single optional method:

```go
// MintBridgeConvergence writes both the synthesis-authored Insight Signal AND
// the arity-N hyperedge connecting the participating signals. Either both
// land or neither does (compute.Bridge.PromoteBatch atomicity).
//
// Implementations that do not support batch promotion (e.g. NoopPersister
// in v1.3) may return ErrConvergenceUnsupported; callers must then fall back
// to MintSignal-only behaviour.
type ConvergenceMinter interface {
    MintBridgeConvergence(
        ctx context.Context,
        sig types.InsightSignal,
        participants []types.Addr,
    ) (signalID types.Addr, edgeID string, err error)
}
```

This is a **separate interface from `Persister`**, not an extension of it, so:

- Existing `Persister` implementations (NoopPersister) don't need to change.
- The Wyrd-backed Persister (when it lands) implements both `Persister` and `ConvergenceMinter`.
- Synthesis tests for arity-N minting can use a separate fake without disturbing the existing `synthesis_test.go` fixtures.

The synthesis dispatch in `Agent.HandleBridgeIntervention` becomes:

```go
func (a *Agent) HandleBridgeIntervention(ctx context.Context, iv types.BridgeIntervention) (types.Addr, bool, error) {
    kind := a.Policy.EvaluateBridgeIntervention(iv)
    if kind == "" {
        return types.Addr{}, false, nil
    }
    sig := a.signalFromBridgeIntervention(iv, kind)

    // If the Persister can mint convergence (i.e. supports arity-N edge),
    // and the intervention has ≥3 participants, prefer that path.
    if cm, ok := a.Persister.(ConvergenceMinter); ok && len(iv.ParticipatingSignalIDs) >= 3 {
        participants := signalIDsToAddrs(iv.ParticipatingSignalIDs)
        id, _, err := cm.MintBridgeConvergence(ctx, sig, participants)
        if err == nil {
            return id, true, nil
        }
        if !errors.Is(err, ErrConvergenceUnsupported) {
            return types.Addr{}, false, fmt.Errorf("synthesis: bridge convergence mint: %w", err)
        }
        // ErrConvergenceUnsupported: fall through to plain MintSignal.
    }

    id, err := a.Persister.MintSignal(ctx, sig)
    if err != nil {
        return types.Addr{}, false, fmt.Errorf("synthesis: bridge-intervention mint: %w", err)
    }
    return id, true, nil
}
```

Convergence path is opt-in — no existing test breaks; no plain `Persister` implementation is forced to grow.

---

## 6. Test plan

When implementation unblocks (Wyrd subpackage available), tests cover:

- `MintBridgeConvergence` writes both the signal and the arity-N hyperedge atomically (assert via Wyrd graph inspection).
- Atomicity failure: a fake Persister that fails on the hyperedge step leaves no signal in the graph (rollback works).
- Arity-2 path: when `len(iv.ParticipatingSignalIDs) < 3`, the convergence mint is skipped and the plain `MintSignal` path is taken (avoids degenerate hyperedges that violate Theorem 2 irreducibility's k≥3 constraint).
- Forward reference: a batch containing the signal node + the arity-N hyperedge that references it succeeds (validates the §3 (b) recommendation against Wyrd's actual `PromoteBatch` semantics).
- Idempotency: re-evaluating the same `BridgeIntervention` twice produces the same `signalID` and `edgeID` (deterministic derivation).

---

## 7. Sequencing dependencies

| Dependency | Status | Blocking |
|---|---|---|
| Wyrd issue #6 corrections (`SignalSource` enum + integration doc updates) | filed [comment](https://github.com/JamesPagetButler/wyrd/issues/6#issuecomment-4392405393); awaiting wyrd-implementor cadence | yes |
| Wyrd `contextus/` subpackage (`FromSignal`, `PromoteToCTH`) | not yet started | yes |
| Wyrd PR #19 `bridge-batch.md` v0.1 → implementation PR | v0.1 design surface up; awaiting bma + bma-implementor §I4 sign-off; implementation PR after | yes |
| Wyrd `PromoteBatch` resolves intra-batch references (recommendation §3) | ask in Wyrd PR #19 review | yes |
| Spec v1.4: `BridgeIntervention.ParticipatingSignalIDs` field (recommendation §4) | not yet filed; this doc is the source | yes |

When all five clear, implementation can land in a single PR. Estimated ~1 day of work (the design is simple; the interface boundary is clean).

---

## 8. References

- Spec v1.3 [§4.4 Synthesis as Persistence Boundary](../Contextus-Spec-v1.3.md#44-insight-signal-emission-pipeline) — the trigger
- Spec v1.3 §4.5 Bridge Agent Session Lifecycle — the data source
- Spec v1.3 §11.3 `BridgeIntervention` — the input struct (needs v1.4 field add)
- Theory v1.5 [§7.6](../Contextus-Theory-v1.5.md) — worked example of structural anomaly; same Bridge-style convergence framing
- [Wyrd `model.Hyperedge`](https://github.com/JamesPagetButler/wyrd/blob/main/model/hyperedge.go) — Theorem 2 irreducibility godoc
- [Wyrd issue #6](https://github.com/JamesPagetButler/wyrd/issues/6) — `contextus/` subpackage (corrections filed)
- [Wyrd PR #19](https://github.com/JamesPagetButler/wyrd/pull/19) — bridge-batch design surface
- [qbp-compute-unit ADR-003 §I4](https://github.com/JamesPagetButler/qbp-compute-unit/blob/feat/issue-7-lean2rom/architecture/adr-003-m1-wdevent-observer-invariants.md) — design-doc-as-S-01-review-surface (this doc qualifies)
