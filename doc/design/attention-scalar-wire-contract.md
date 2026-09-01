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
    Kind    string `json:"kind"` // "physical" | "conceptual" | "operational" — craft renders each manifold-exit by kind
}

// ScalarReferent v0.1 additions (trajectory — stays inside the referent = monitoring, not evaluation):
//   Rate                 *float64 `json:"rate,omitempty"`                 // divergence rate-of-change
//   ProjectedPeakTime    *string  `json:"projected_peak_time,omitempty"`  // RFC3339 valid-time projection
//   ProjectionConfidence *float64 `json:"projection_confidence,omitempty"`// distinct from present-value Confidence (forecast glyph vs measurement glyph)
```

Casing is **snake_case** throughout (matches every existing Contextus JSON tag — `scope_id`, `provenance_tag`, `signal_id`).

## 3. Answers to the §7 render-side asks

**Q1 — Wire envelope.**
- **Shape:** yes, this is the emit shape (reconcile your `craft-96-render-design.md §2` draft against §2 above — I own emit, you own render; let's diff).
- **Frame model:** **snapshot-then-delta.** On scope *activation* (a discrete §4.6.6 event) Contextus emits a **full snapshot** = the activated scope's current `AttentionScalar` set; thereafter **per-signal deltas** (`upsert`/`remove` by `signal_id`) as signals change. Rationale: the activated set is bounded (focus-shaped), deltas keep the WS cheap, and it mirrors the sessionbridge subscribe-with-recent-history pattern.
- **Casing:** snake_case (§2).
- **`divergence` absent-vs-null:** **explicit `null`**, never field-absent. Absent is ambiguous ("N/A" vs "not-yet-computed"); explicit `null` = "this signal is a conceptual/membership signal with no monitor referent" (e.g. `sig:regulatory-bridge`). The render must distinguish no-referent from zero-divergence, so the key is always present.

**Q2 — Transport.** Contextus publishes `AttentionScalar` frames on **NATS `ctx.craft.attention.<scope_id>`** (snapshot on subscribe + deltas), consistent with the existing `ctx.*` subject hierarchy. **The craft-SERVER subscribes there** and relays to the browser over its WebSocket — the WS is craft-server↔browser (your layer), NOT browser↔Contextus. Initial snapshot is delivered NATS-side on subscribe (no history replay). You do **not** pull raw Wyrd for the frame — Contextus assembles it (retention-tier + salience logic is Contextus's).

**Q3 — `event_time`.** CONFIRMED: raw valid-time instant, hoisted scalar-level (Finding D, ratified). And CONFIRMED the guard, emphatically: signals lacking a valid-time datum emit **`event_time: null`** — **never a capture-time fallback.** A capture-time fallback silently reintroduces the 3-clock bug (renders a 40-yr lag as ~0). `null` means "no valid-time — cannot place on the causal-lag axis," which is a *correct* render state, not a defaulted one.

**Q4 — Cascade edges.** CONFIRMED: you query Wyrd hyperedges **directly from the `locus` node-set** — that's your traversal side (I emit nodes, you walk edges; the ratified cascade-split). **Emit-side wiring requirement I own (surfacing it now so it's not a build-time surprise):** for the endpoint-`event_time` lag-diff to work *at traversal time*, valid-time must be persisted as a **first-class property on the Wyrd member/observation node**, not only hoisted into the `AttentionScalar`. So the node-writer (ctx-adapter / Synthesis) must stamp `event_time` on the node. I'll carry that as the emit-side build task; you can assume each `locus` endpoint node exposes its valid-time when you walk it.

**Q5 — `detail_ref` round-trip.** **Contextus serves it** (`detail_ref` → `InsightSignal.Evidence[] + ClaimHistory[]`); the **craft-server proxies**. Rationale: Evidence-pointer *tier resolution* + ClaimHistory assembly are Contextus semantics (retention tiers, provenance envelope), not raw Wyrd reads — resolving craft-server-side against Wyrd would duplicate Contextus logic and drift. So the round-trip terminates at Contextus. (`detail_ref` is carried on the scalar so you never construct the path yourself.)

## 4. Seam invariants held
Emit-side stays OBSERVE-only: no craft/render state flows back; `divergence`/trajectory are monitoring quantities (never verdicts); the eligibility/action layers stay craft/judge-side. `scope_ids` carries *membership*, not manifold-*rulings*.

---

*End — emit-side wire contract for the attention scalar v0.1. Design/contract-only; the Go-type authoring + live-pipeline wiring (NATS publish, Wyrd valid-time stamping per Q4) is the post-Sprint-3-close build, riding the normal §I4/PR path (push beekeeper-gated). Reconcile against hutchins's `craft-96-render-design.md §2`.*
