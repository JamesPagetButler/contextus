# Spec v1.4 Design Surface — Theory as Conceptual Scope

**Status:** Design only. §I4 review surface for the Spec v1.4 amendment.

**Date:** 2026-05-14

**Author:** contextus-impl

**Motivation:** Toddle-design meeting question (#toddle-design seq 8/9 context): *"Will Contextus dynamically enable monitoring for both theories and any defined context in the known universe?"* Spec v1.3 supports dynamic scope-based monitoring for arbitrary `NT_SCOPE_PHYSICAL` / `NT_SCOPE_CONCEPTUAL`, but theories — claims with their derivation chains living in the CTH inventory — are not a native Contextus node type. This document specifies how a theory becomes a first-class scope without introducing a new node type, by adding one canonical membership-predicate method to the §4.6.5 catalogue.

---

## 1. The Architectural Read

A **theory**, for Contextus's purposes, is a particular kind of conceptual scope: its members are claims and observations transitively in the derivation chain rooted at a CTH anchor.

- CTH owns the derivation chain (`cth.anchor.*` types in the CTH inventory).
- Contextus owns the scope mechanism (`NT_SCOPE_CONCEPTUAL` per Spec v1.3 §4.6.3).
- The bridge between them is a **membership predicate** named `cth-derivation`, registered against the §4.6.5 implementation-defined predicate catalogue.

No new node type. No new Wyrd extension. One canonical method value.

This preserves the Spec v1.3 §1 architectural commitment that *"Contextus never evaluates claims. CTH never monitors data. The Bridge translates between them."* — Contextus continues not to evaluate the theory; it scopes attention to the subgraph the theory's derivation chain spans.

---

## 2. The Membership Predicate

§4.6.5 currently lists five non-exhaustive example methods (`asserted`, `tag-overlap`, `embedding-similarity`, `classifier`, `citation-graph-hop`). Adding a sixth:

> **`cth-derivation`** — the member node is transitively reachable in the CTH inventory's derivation chain rooted at a CTH anchor whose ID is named in `NT_SCOPE_CONCEPTUAL.ontology_uri`. Concretely, the URI is `cth://anchor/<anchor-id>` and the inventory is queried at scope-load time and again whenever the inventory changes (CTH live-update API, CTH #51). Adapter implementations:
>
> - Subscribe to CTH inventory mutation events (NATS `contextus.heartbeat.{receipt_id}` already specified in Spec v1.3 §5.3, used in reverse here for read-only consumption).
> - On change, recompute the transitive closure of `cth.derivation.*` edges from the anchor and emit `HE_SCOPE_MEMBERSHIP` edges for newly-reached nodes (`provenance_tag = "D"`, `method = "cth-derivation"`).
> - On change reducing the closure, emit `HE_SCOPE_MEMBERSHIP` removal events.

The mechanism is unidirectional and read-only on the CTH side — preserves the §8.3 heartbeat invariant (no reverse channel CTH → Contextus that would compromise CTH judgment).

---

## 3. The User Story

A beekeeper or tenant-implementor declares:

```yaml
scopes:
  conceptual:
    - scope_id: scope-theory-qbp-exp-11
      name: "QBP-EXP-11 (GW-GRB temporal correlation)"
      ontology_uri: cth://anchor/qbp-exp-11
      tags: [qbp, gw-em, narrative-anomaly-candidate]
      # member predicate is implicit from ontology_uri scheme: cth-derivation
```

Activating this scope in a session does what Spec v1.3 §4.6.6 already specifies — promotes signals connected to scope members toward Core retention tier, focuses Edge Scout boundary cues on this scope's frontier, and lets Synthesis subscribe with elevated weighting on convergences that bridge into this scope from outside.

The capability the beekeeper experiences: *"monitor the QBP-EXP-11 theory"* → system does the right thing. No bespoke per-theory wiring; same mechanism as monitoring a watershed.

---

## 4. What Changes Between Spec v1.3 and v1.4

**Spec v1.4 §4.6 amendments (additive, no breakage):**

| § | Change |
|---|---|
| §4.6.3 | Add note that `ontology_uri` of form `cth://anchor/<id>` triggers the `cth-derivation` membership predicate by default. |
| §4.6.5 | Add `cth-derivation` to the non-exhaustive list of methods with the §2 specification above. |
| §4.6.6 | Add cross-reference: scope activation with a `cth-derivation` predicate also drives the CTH `compute.NetCompressionDetail` re-evaluation cadence (CTH-side concern; non-blocking for Contextus). |
| §11.4 | No Go type change. `ScopeConceptual` already carries `OntologyURI`; the scheme `cth://anchor/<id>` is recognised by the loader. |

**Plus the separately-tracked §2.4 Referent work** (queued behind cth-implementor's scoring-evaluation ownership ack). The two v1.4 items can land in one PR or two; my preference is two (this one is small and architectural-only; the Referent work is larger and depends on Wyrd predictions/ schema).

---

## 5. Dependencies

| Dependency | Where | Blocking |
|---|---|---|
| CTH live-update API (CTH issue #51) | confluent-trust | Yes for full live-monitoring; not for static scope-load at session start (which just reads the current inventory snapshot once) |
| Wyrd #33 scope-loader API | wyrd | Yes for the YAML-driven loader; not for direct in-Go `ScopeConceptual{}` construction |
| W-Toddle-1 (Wyrd PR #39 tier-immune + salience primitives) | wyrd | No — this design is orthogonal to retention-tier primitives |

**Practical sequencing:** ship this design surface now; implement the loader-recognises-`cth://anchor/`-scheme path after CTH #51 lands; until then, scopes can be declared in Go directly or via the loader fixture.

---

## 6. Test Plan

When implementation unblocks:

- Loader parses `ontology_uri: cth://anchor/<id>` and dispatches to the `cth-derivation` predicate.
- Predicate queries CTH inventory and emits `HE_SCOPE_MEMBERSHIP` edges with correct `provenance_tag` and `method`.
- Inventory mutation event triggers re-evaluation; new members added; removed members get their HE_SCOPE_MEMBERSHIP edges deleted (not orphaned with stale data).
- Activating the scope in a session correctly drives retention-tier promotion for members.
- Session deactivation reverses retention pressure (existing §4.6.6 mechanism — confirm it still works for `cth-derivation`-sourced members).

---

## 7. Open Questions for §I4 Review

1. **Should `cth-derivation` predicate's scope-loader cache the closure or always recompute?** Caching saves work but introduces staleness. Recommendation: recompute on activation (cheap if the inventory snapshot is bounded); subscribe to mutation events for live update. Compromise position; flag if reviewers disagree.

2. **What about `cth.confluence.*` and `cth.branch.*` node types — are they in the derivation closure?** Recommendation: yes for `confluence` (definitionally part of the theory's evidence structure); branches are scope-specific (the active branch is in; alternative branches are tagged but not in the closure unless the scope is broader). Flag for cth-implementor.

3. **Multi-anchor theories (a theory rooted at multiple CTH anchors)?** Recommendation: support via an **additive** extension — keep `OntologyURI string` (single anchor case) and add a new `OntologyURIs []string` field. Loader semantics: if `OntologyURIs` is populated, take the union of derivation closures across the anchors; else fall back to `OntologyURI`. This preserves the federation contract that Wyrd imports `contextus/internal/contextus/types/` directly (per Wyrd PR #40 — wyrd-implementor seq=15 on `#toddle-design`) — additive field changes are forward-wire-compatible; replacing `string` with `[]string` would be a breaking change. Defer if even the additive form is controversial.

---

## 8. References

- Spec v1.3 §4.6 (Scope Nodes); §8 (Contextus–CTH Bridge Integration)
- Theory v1.5 (PR #5 mergeable) §3.6.6 (Synthesis as Persistence Boundary)
- BMA Theory Addendum 18 §2.1 (Stance) — conceptual scope ≈ Stance per addendum-18-walk D9 vocab unification
- Wyrd PR #39 + issue #38 (W-Toddle-1 generic tier-immune + salience primitives) — orthogonal to this design
- toddle-design seq 8/9 — meeting context that surfaced the question
- CTH issue #51 — live inventory update API (load-bearing for the live-monitor path)
