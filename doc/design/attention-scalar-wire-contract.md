# Attention Scalar — emit-side wire contract (craft↔Contextus seam v0.1)

**Status:** Design/contract-of-record. The authoritative **emit-side** shape Contextus hands the craft-server (inter#96). Answers hutchins's seq=971 §7 render-side asks. **Contract-firming only — NOT the build:** authoring these into `pkg/types/` + wiring the live emit pipeline (NATS publish, Wyrd valid-time stamping) is the post-Sprint-3-close build window per qbp-architecture seq=966. This doc settles the wire so the render side can conform now (closes the "proven≠wired" gap hutchins flagged — the contract exists even though the code doesn't yet).

**Date:** 2026-09-01 · **Author:** contextus-impl (emit side) · **Co-owner:** hutchins (render side) · **Seam:** v0.1 PINNED (see `squam-bennett-brook-craft-bundle.md §4`).

---

## 1. The load-bearing finding (hutchins seq=971) — ACKED, correct
The seam is design-pinned at v0.1 but the emit Go types are still v0: `ScalarReferent` lacks `{rate, projected_peak_time, projection_confidence}`; `InsightSignal` has no scalar-level `lead_time`/`event_time`/`scope_ids`; **there is no `AttentionScalar` struct at all.** So the WS/wire JSON is undefined in code on both sides. This doc defines it. proven≠wired → the contract is now wired *as contract*; the Go/pipeline build follows in the window.

## 2. The wire types (proposed — emit-side source of truth)

```go
// AttentionScalar is the per-signal payload the craft renders. Emitted by Contextus,
// READ-ONLY to the craft. One frame = the activated scope's signal set (see §4).
type AttentionScalar struct {
    // --- magnitude (what deserves attention) ---
    Salience           float64  `json:"salience"`             // [0,1] default rank = AnomalyScore (operational: MaxScore of Referents)
    Divergence         *float64 `json:"divergence"`           // ScalarReferent.Score; EXPLICIT null for conceptual/membership (non-monitor) signals
    Persistence        int      `json:"persistence"`          // landscape-stability count
    Structural         bool     `json:"structural"`           // landscape (true) vs exploration (false)
    Confidence         float64  `json:"confidence"`           // [0,1]
    ConfidenceVariance float64  `json:"confidence_variance"`  // present-value uncertainty
    Promote            bool     `json:"promote"`              // crossed Synthesis persistence-boundary / AnomalyStructural (Seam-surprise)
    LeadTime           *string  `json:"lead_time"`            // DERIVED summary (ISO-8601 duration); null if no trajectory
    // --- placement (where / at what zoom) ---
    Locale             Envelope `json:"locale"`               // spatial+temporal envelope = LocaleEnvelope (scope span, NOT the instant)
    Grain              string   `json:"grain"`                // LocaleGrain — render zoom/resolution
    // --- anchor (identity + traversal) ---
    ScopeIDs           []ScopeRef `json:"scope_ids"`          // multi-manifold membership; NEVER a bare id
    Locus              []string   `json:"locus"`              // TraversalPath []Addr — the causal-walk ENTRY-SET (craft walks edges from here)
    SignalID           string     `json:"signal_id"`          // round-trip handle
    DetailRef          string     `json:"detail_ref"`         // path for Evidence[]+ClaimHistory[] (see §7.Q5)
    // --- bitemporal (Finding D) ---
    EventTime          *string    `json:"event_time"`         // RAW valid-time instant (RFC3339), hoisted like Locus; null if NO valid-time datum — NEVER a capture-time fallback
}

type ScopeRef struct {
    ScopeID string `json:"scope_id"`
    Kind    string `json:"kind"`             // DOMAIN-KIND (render selector), NOT raw scope-kind: physical|theory|cognition|code|operational. Derived emit-side — see §2.1.
    Parent  string `json:"parent,omitempty"` // hierarchical grain: the enclosing scope id (point ⊂ sub-watershed). Emit-side populates from ScopePhysical.ParentScopeID; empty when top-level. (Absorbed from craft #96 render — hutchins seq=1013.)
}

// ScalarReferent v0.1 additions (trajectory — stays inside the referent = monitoring, not evaluation):
//   Rate                 *float64 `json:"rate,omitempty"`                 // divergence rate-of-change
//   ProjectedPeakTime    *string  `json:"projected_peak_time,omitempty"`  // RFC3339 valid-time projection
//   ProjectionConfidence *float64 `json:"projection_confidence,omitempty"`// distinct from present-value Confidence (forecast glyph vs measurement glyph)
```

Casing is **snake_case** throughout (matches every existing Contextus JSON tag — `scope_id`, `provenance_tag`, `signal_id`).

### 2.1 `kind` = domain_kind — emit-side derivation (architect-ruled, seq 1005)

`ScopeRef.Kind` carries the **domain-kind** (the render's canopy selector), **not** the raw Go scope-kind. Reason: **theory and code are BOTH `ScopeConceptual`**, so raw scope-kind can't distinguish a proof-lattice (CTH theory) from a type-lattice (Edda code) — the render must pick the canopy off domain-kind. Per `inter/craft-domain-kinds-and-edda-seam.md`, domain-kind is keyed on *state-source*, and the derivation is **emit-side** (Contextus knows the scope type + membership predicate):

| Contextus scope | membership predicate / marker | `domain_kind` | state source |
|---|---|---|---|
| `ScopePhysical` | (geometry) | `physical` | Contextus pipeline |
| `ScopeOperational` | `hardware-identifier` | `operational` | Contextus pipeline (host telemetry) |
| `ScopeConceptual` | `cth-derivation` (`ontology_uri: cth://anchor/…`) | `theory` | CTH (bridge) |
| `ScopeConceptual` | edda-module predicate | `code` | Edda compiler (bridge) |
| `ScopeConceptual` | BMA-instance predicate | `cognition` | BMA (bridge) |

Fallback: a `ScopeConceptual` with no domain-specific predicate is emitted with an explicit `kind` from its config (no silent default — an un-kinded manifold-exit would be a render ambiguity). Naming: `scalar.DomainKind{Physical,Theory,Cognition,Code,Operational}` constants (mirrors hutchins's render-side constants). **For Squam:** the `brownfield-eligibility` scope (`cth-derivation`) emits `kind: "theory"` → the `barn-soil-ddt` node's `scope_ids` = `[{…watershed, physical}, {…brownfield, theory}]`, so its `reframe` exits are `{physical, theory}` — two domain-kinds from day one.

### 2.2 Null vs empty — the two distinct "nothing"s (wire convention)

The contract uses **two different encodings for "nothing," and they are not interchangeable** — the render depends on the distinction (verified against the craft render, seq 1409/1413):

- **Absent *scalar* → explicit `null`** (a pointer field, never `omitempty`): `divergence`, `event_time`, `lead_time`. `null` is a *meaningful* value — "no monitor referent" / "no valid-time datum, cannot place on the causal-lag axis" / "no trajectory" — semantically distinct from a zero (`0.0` divergence ≠ no divergence; epoch ≠ no event-time). So the key is always present with an explicit `null`.
- **Empty *collection* → `[]`** (never `null`): the `detail_ref` collections (`referents`, `evidence`, `claim_history`) and the cascade `edges` list. An empty collection is the natural "none" and marshals as `[]` so the render reads `.length` without a null-guard. Emit MUST NOT send `null` for these.

Rule of thumb: **a missing value is `null`; a missing set is `[]`.** (The craft render already implements this — `[]`-not-`null` on the detail collections + cascade edges; this documents it emit-side so the live `detail_ref`/frame endpoints match the stub exactly.)

## 3. Answers to the §7 render-side asks

**Q1 — Wire envelope.**
- **Shape:** yes, this is the emit shape (reconcile your `craft-96-render-design.md §2` draft against §2 above — I own emit, you own render; let's diff).
- **Frame model:** **snapshot-then-delta.** On scope *activation* (a discrete §4.6.6 event) Contextus emits a **full snapshot** = the activated scope's current `AttentionScalar` set; thereafter **per-signal deltas** (`upsert`/`remove` by `signal_id`) as signals change. Rationale: the activated set is bounded (focus-shaped), deltas keep the WS cheap, and it mirrors the sessionbridge subscribe-with-recent-history pattern.
- **Casing:** snake_case (§2).
- **`divergence` absent-vs-null:** **explicit `null`**, never field-absent. Absent is ambiguous ("N/A" vs "not-yet-computed"); explicit `null` = "this signal is a conceptual/membership signal with no monitor referent" (e.g. `sig:regulatory-bridge`). The render must distinguish no-referent from zero-divergence, so the key is always present.

**Q2 — Transport.** Contextus publishes `AttentionScalar` frames on **NATS `ctx.craft.attention.<scope_id>`** (snapshot on subscribe + deltas), consistent with the existing `ctx.*` subject hierarchy. **The craft-SERVER subscribes there** and relays to the browser over its WebSocket — the WS is craft-server↔browser (your layer), NOT browser↔Contextus. Initial snapshot is delivered NATS-side on subscribe (no history replay). You do **not** pull raw Wyrd for the frame — Contextus assembles it (retention-tier + salience logic is Contextus's).

**Q3 — `event_time`.** CONFIRMED: raw valid-time instant, hoisted scalar-level (Finding D, ratified). And CONFIRMED the guard, emphatically: signals lacking a valid-time datum emit **`event_time: null`** — **never a capture-time fallback.** A capture-time fallback silently reintroduces the 3-clock bug (renders a 40-yr lag as ~0). `null` means "no valid-time — cannot place on the causal-lag axis," which is a *correct* render state, not a defaulted one.

**Q4 — Cascade edges.** CONFIRMED: you query Wyrd hyperedges **directly from the `locus` node-set** — that's your traversal side (I emit nodes, you walk edges; the ratified cascade-split). **Emit-side wiring requirement I own (surfacing it now so it's not a build-time surprise):** for the endpoint-`event_time` lag-diff to work *at traversal time*, valid-time must be persisted as a **first-class property on the Wyrd member/observation node**, not only hoisted into the `AttentionScalar`. So the node-writer (ctx-adapter / Synthesis) must stamp `event_time` on the node. I'll carry that as the emit-side build task; you can assume each `locus` endpoint node exposes its valid-time when you walk it.

**Q5 — `detail_ref` round-trip.** **Contextus serves it**; the **craft-server proxies**. The foveal detail bundle is:

```
detail_ref → {
  referents:     []ScalarReferent,  // INCLUDING the v0.1 trajectory {rate, projected_peak_time, projection_confidence}
  evidence:      []EvidencePointer, // wire shape {ref, kind, note?} — see below
  claim_history: []ClaimVersion,    // {version, claim, reason, source, timestamp} — matches pkg/types.ClaimVersion 1:1
}
```

**`EvidencePointer` wire shape (pinned — closes the render-side gap, craft #96):** the wire element is `{ref string, kind string, note string (omitempty)}`, NOT Contextus's internal tiered `EvidencePointer` (which carries `Locator`/`LocatorKind` + tier-conditional fields per Spec §5.4). **Emit-side maps down at the `detail_ref` boundary:** `ref` = the resolved `Locator`, `kind` = `LocatorKind` normalised to `observation|citation|derivation`, `note` = the tier-conditional note when the retention tier carries one (else omitted). Tier resolution stays Contextus-side (the craft renders what it's handed); the craft never sees raw tier internals.

Rationale: Evidence-pointer *tier resolution* + ClaimHistory assembly are Contextus semantics (retention tiers, provenance envelope), not raw Wyrd reads — resolving craft-server-side would duplicate Contextus logic and drift. So the round-trip terminates at Contextus. (`detail_ref` is carried on the scalar so you never construct the path yourself.)

**Trajectory delivery (closes hutchins seq=986 open point) — CONFIRMED via `detail_ref`, not hoisted.** The full `trajectory` stays inside `ScalarReferent` and arrives on the **foveal pull** alongside `evidence`/`claim_history` — because trajectory is **foveal-detail** (needed only to render the forecast glyph at foveal engagement), exactly the register-dual discipline: peripheral = scalar (`lead_time` summary drives peripheral glyph character); foveal = `detail_ref` pull (full `{rate, projected_peak_time, projection_confidence}` drives the fuzzy-edged forecast glyph, `projection_confidence` setting the fuzz distinctly from present-value `confidence`). Hoisting full per-referent trajectory into every frame would stream foveal-detail to every peripheral viewer — the exact cost the dual exists to avoid. So: `lead_time` on the scalar (peripheral); full trajectory on `detail_ref.referents[]` (foveal). Monitoring-pure either way (trajectory never leaves the referent).

## 4. Seam invariants held
Emit-side stays OBSERVE-only: no craft/render state flows back; `divergence`/trajectory are monitoring quantities (never verdicts); the eligibility/action layers stay craft/judge-side. `scope_ids` carries *membership*, not manifold-*rulings*.

---

*End — emit-side wire contract for the attention scalar v0.1. Design/contract-only; the Go-type authoring + live-pipeline wiring (NATS publish, Wyrd valid-time stamping per Q4) is the post-Sprint-3-close build, riding the normal §I4/PR path (push beekeeper-gated). Reconcile against hutchins's `craft-96-render-design.md §2`.*
