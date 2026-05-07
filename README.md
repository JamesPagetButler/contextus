# Contextus

**Helpful Engineering — Quaternion-Based Physics Programme**
Principal Investigator: James Paget Butler

Contextus is the cross-domain pattern-matching layer for the Helpful Engineering / QBP programme. Where BMA is a single instance's cognitive substrate, CTH measures epistemic health, and Wyrd is the typed-hypergraph database, **Contextus is the index of evidence across domains** — a focus-area-aware lens on the universe of observations that lets researchers, agents, and Sharp Butler instances find cross-domain patterns without holding the corpus on disk.

## Status

**Crawl phase.** Theory v1.4 + Spec v1.2 are the canonical references on disk. **Spec v1.3 is being drafted** on `feat/spec-v1.3` — see issue [#1](https://github.com/JamesPagetButler/contextus/issues/1) for scope and the §I4 review-surface framing per [qbp-compute-unit ADR-003](https://github.com/JamesPagetButler/qbp-compute-unit/blob/feat/issue-7-lean2rom/architecture/adr-003-m1-wdevent-observer-invariants.md).

All four open architectural questions from Spec v1.2 are resolved:

- **Placement** — discrete new `§4.6 Scope Nodes` + `§5.x Evidence Pointer Discipline` sections (`#live-test` seq=42; reverts the architect's earlier distributed-into-existing-sections recommendation, which carried Walk-phase Wyrd-shape import bias).
- **EvidencePointer sizing** — tier-conditional field population: Skeleton/Distant carry `Locator` + `LocatorKind` only; richer fields populate at Peripheral and above. Distant tier tightened (drops `Note` + `AccessHint`) to preserve the existing §9.1 budget contract.
- **Evidence list growth bound** — cap-per-tier with summary-pointer eviction (lowest-confidence-N pointers collapse into a single summary pointer); preserves Strengthening semantics from Theory §3.6.3.
- **`SignalSource` enum** — corrected to `scout | correlation | synthesis` per Spec v1.2 §11.1 `AgentClass`. Edge Scout / Corpus Edge Scout / Bridge Agent do **not** emit Insight Signals — Synthesis is the persistence boundary for ephemeral session-scoped agent output.

A fourth `structural` anomaly_type is added provisionally in v1.3 for Synthesis-promoted signals from session-scoped agent output where `narrative` doesn't fit; **Theory v1.5** will formalize.

## Read order

1. [`Contextus-Theory-v1.4.md`](Contextus-Theory-v1.4.md) — what Contextus is, why agents are data, the Locale concept, the surveillance / search-mode distinction, three case studies (Yellowstone, whale sharks, black holes), the Colorado River failure analysis.
2. [`Contextus-Spec-v1.2.md`](Contextus-Spec-v1.2.md) — current canonical spec: hypergraph schema, agent taxonomy (Scout / Correlation / Synthesis as global authors of Insight Signals; Edge Scout / Corpus Edge Scout / Bridge Agent as session-scoped non-authors), NATS subjects, MCP interface, retention tiers, CTH bridge integration. **Superseded by v1.3 once it lands.**
3. [`contextus-wyrd-integration-architecture-2026-05-05.md`](contextus-wyrd-integration-architecture-2026-05-05.md) — the architecture-instance integration doc that resolves Wyrd's issue #6 (SignalSource correction, scope-node taxonomy, EvidencePointer discipline). Inputs to Spec v1.3.
4. [`doc/contextus-impl-onboarding-prompt.md`](doc/contextus-impl-onboarding-prompt.md) — bootstrap prompt for fresh contextus-impl sessions; captures decisions taken to date.
5. [`MANIFEST.md`](MANIFEST.md) — document inventory.

## Integration with sibling programmes

| Project | Repo | Relationship |
|---|---|---|
| **BMA** | [`bma-systema`](https://github.com/JamesPagetButler/bma-systema) | Walk-phase consumer. BMA's `internal/bma/cth/projection.go` walks the BMA hypergraph and emits CTH inventories scoped by Contextus scope nodes. |
| **CTH** | [`confluent-trust`](https://github.com/JamesPagetButler/confluent-trust) | Bridge target. Confirmed Insight Signals at sufficient confidence promote to CTH anchors via `compute.Bridge.Promote`. |
| **Wyrd** | [`wyrd`](https://github.com/JamesPagetButler/wyrd) | Storage substrate. Insight Signals are `model.Node` of `Type = "contextus.signal.<agent>"`; cross-domain matches are `model.Hyperedge`; scope nodes (`contextus.scope.physical` / `contextus.scope.conceptual`) are first-class. |
| **QBP-CU** | [`qbp-compute-unit`](https://github.com/JamesPagetButler/qbp-compute-unit) | Future-Walk consumer. QBP-CU's WDEvent stream may surface as scout density anomalies if BMA's hypergraph state is being scouted, but Contextus stays on the "detects, never evaluates" side of the surveillance-lane boundary. |

## Architectural principle

> **Contextus is the *index* of evidence, not the evidence itself.**

The universe is too big to hold on a 250 GB SATA SSD. Spec §9 (Proportional Data Retention) tiers signals from Core (full content, 500 KB–2 MB) down through Skeleton (DOI + year, 50–100 bytes). Combined with `EvidencePointer` discipline (`Locator` + `LocatorKind` + tier-conditional metadata) and scope-node activation (`Squam Lake watershed`, `physics`, `magnetar spectroscopy`), the system becomes universe-shaped in principle, focus-shaped in practice.

## Contributing

Contextus is currently a single-implementor + single-architect collaboration; the bridge (`#live-test` and `#qbp-cu-walk` on `~/.claude/mcp-servers/sessionbridge/`) is where cross-instance coordination happens. Spec changes flow through architecture instance review per the design-doc-as-S-01-review-surface invariant captured in qbp-compute-unit ADR-003 §I4.

## License

TBD — pending alignment with sibling programme licenses (BMA: Apache 2.0; CTH: Apache 2.0; Wyrd: TBD; QBP-CU: TBD).

## Attribution

Theory builds on prior work in cross-domain insight discovery, knowledge graphs, cognitive ergonomics, and quaternion algebra. Specific case-study citations carried forward from Theory v1.4 (Yellowstone trophic cascade work, whale shark confluence research, black hole hidden-correlation literature).
