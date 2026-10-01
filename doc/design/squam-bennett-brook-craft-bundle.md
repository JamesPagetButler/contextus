# Squam / Bennett Brook — Contextus craft bundle (walk-through)

**Status:** Design/worked-example. The concrete completeness-stress for the craft↔Contextus seam contract (attention scalar v0). Paired with qbp-architecture's craft-side render/navigate walk-through.

**Date:** 2026-08-30 · **Author:** contextus-impl · **Beekeeper problem:** 2004–07 Squam Lake loon die-off → legacy-DDT source-tracing to the Bennett Brook sub-watershed (Holderness NH).

**What this is:** the Contextus *emit* side of one focus area — the `ScopePhysical` scope-YAML, the member resolution, and the live `InsightSignal` + `ScalarReferent` set with the **v0 attention scalar populated per signal**. Contextus lights + scores the nodes; the craft owns the causal-edge traversal, HUD, and action-space (per the ratified seam). No gated moves — design only.

---

## 1. Scope-YAML (validate-ready against `schema/scope-config.schema.json`)

Three manifolds meet in this focus area: the **physical watershed** (the cascade substrate), an **abstract regulatory manifold** (Brownfield/CERCLA eligibility — a `ScopeConceptual`, which Contextus *scopes* but does NOT evaluate), and the barn **acute point-source** as a nested child physical scope (grain demonstration).

```yaml
# Squam / Bennett Brook craft focus-area scope config
# Spec v1.3 §4.6 + §11.4; multi-manifold (physical watershed + conceptual regulatory).

physical_scopes:
  - id: "contextus:scope:physical:bennett-brook-subwatershed"
    description: "Bennett Brook sub-watershed, Holderness NH — Squam Lake catchment tributary; legacy-DDT source region (1920s-mid-century apple orchard, Route 113)."
    type_nodes:
      - "contextus.scope.physical"
      - "contextus.ecotoxicology.watershed"
    bounds:
      lat: [43.74, 43.80]
      lon: [-71.60, -71.53]
      time: ["1945-01-01T00:00:00Z", "2050-01-01T00:00:00Z"]  # cascade envelope: DDT-era → projected remediation horizon
      height: [180.0, 400.0]
    tags:
      - "legacy-ddt"
      - "loon-mortality"
      - "source-tracing"

  - id: "contextus:scope:physical:bennett-barn-pointsource"
    description: "Orchard mixing/storage barn acute point source (burned ~late 1960s); 723 ug/kg p,p'-DDT + 721 ug/kg DDE in soil."
    type_nodes:
      - "contextus.scope.physical"
      - "contextus.ecotoxicology.point-source"
    parent_scope_id: "contextus:scope:physical:bennett-brook-subwatershed"  # hierarchical grain: point ⊂ sub-watershed
    bounds:
      lat: [43.762, 43.764]
      lon: [-71.566, -71.564]
      time: ["1945-01-01T00:00:00Z", "2050-01-01T00:00:00Z"]
      height: [195.0, 205.0]
    tags:
      - "point-source"
      - "acute"

conceptual_scopes:
  - id: "contextus:scope:conceptual:brownfield-eligibility"
    description: "Regulatory-status manifold for the source parcel: CERCLA 107(i) ag-exemption + below-HRS-threshold => not-Superfund; EPA Region 1 / NHDES Brownfield-eligible; NH Env-Or 600 oversight."
    ontology_uri: "cth://anchor/regulatory-brownfield-cercla-107i"  # abstract referent; membership via cth-derivation (Spec v1.4)
    type_nodes:
      - "contextus.scope.conceptual"
      - "contextus.regulatory.remediation-pathway"
    related_scope_ids:
      - "contextus:scope:physical:bennett-barn-pointsource"
    tags:
      - "brownfield"
      - "regulatory-manifold"

scope_memberships:
  # Physical member resolution: nodes whose geometry falls in the catchment (watershed-catchment predicate).
  - scope: "contextus:scope:physical:bennett-brook-subwatershed"
    member: "contextus:obs:bennett-sediment-dde-core"
    confidence: 0.95
    provenance_tag: "D"
    method: "watershed-catchment"      # NEW predicate candidate (§4.6.5): hydrological-catchment membership
    weight_tier: "complex"
  - scope: "contextus:scope:physical:bennett-brook-subwatershed"
    member: "contextus:obs:bennett-crayfish-dde"
    confidence: 0.90
    provenance_tag: "D"
    method: "watershed-catchment"
    weight_tier: "complex"
  - scope: "contextus:scope:physical:bennett-brook-subwatershed"
    member: "contextus:obs:squam-loon-eggshell"
    confidence: 0.88
    provenance_tag: "D"
    method: "watershed-catchment"
    weight_tier: "complex"
  - scope: "contextus:scope:physical:bennett-barn-pointsource"
    member: "contextus:obs:barn-soil-ddt"
    confidence: 0.99
    provenance_tag: "P"
    method: "geographic-overlap"
    weight_tier: "complex"
  - scope: "contextus:scope:physical:bennett-brook-subwatershed"
    member: "contextus:obs:orchard-soil-ddt-dde-ratio"
    confidence: 0.85
    provenance_tag: "D"
    method: "watershed-catchment"
    weight_tier: "complex"
  # Cross-manifold: the point-source parcel is a member of the regulatory manifold (Contextus scopes it; the eligibility VERDICT is judge/legal-authority lane, not Contextus).
  - scope: "contextus:scope:conceptual:brownfield-eligibility"
    member: "contextus:obs:barn-soil-ddt"
    confidence: 0.80
    provenance_tag: "T"
    method: "cth-derivation"
    weight_tier: "complex"
```

**Member resolution summary:** the physical members light up by the **`watershed-catchment`** predicate (a new §4.6.5 method candidate — hydrological-catchment geometry membership; the barn point-source uses exact `geographic-overlap`). The regulatory manifold binds the same barn node by `cth-derivation` (theory/regulatory referent), so **one physical node (`barn-soil-ddt`) is simultaneously a member of a physical scope AND the regulatory conceptual scope** — see Finding B.

---

## 2. Live signal set — `InsightSignal` + `ScalarReferent`, v0 attention scalar per signal

Scalar = `magnitude {salience, divergence, persistence(+structural), confidence±var, promote}` · `placement {locale, grain}` · `anchor {scope_id, locus=TraversalPath[], signal_id}`.

| signal_id | what | salience | divergence (ScalarReferent) | persistence / structural | confidence±var | promote | grain | locus (TraversalPath entry-set) |
|---|---|---|---|---|---|---|---|---|
| `sig:barn-ddt` | barn soil p,p'-DDT 723 µg/kg vs ~0 background | **0.92** | 0.98 (obs 723 / pred ~5 µg/kg) | landscape (T) | 0.90 ±0.05 | **true** | point (~50 m) | `[barn-soil-ddt]` |
| `sig:orchard-kinetic` | orchard field **parent-DDT > DDE ratio** = delayed breakdown | 0.70 | 0.55 (ratio vs weathered-baseline) | landscape (T) | 0.75 ±0.10 | **true** | field (~500 m) | `[orchard-soil-ddt-dde-ratio]` |
| `sig:sediment-dde` | Bennett Cove-mouth sediment DDE, ²¹⁰Pb continuous since 1951 | 0.80 | 0.85 | landscape (T) | 0.88 ±0.06 | false | reach (~1 km) | `[bennett-sediment-dde-core]` |
| `sig:crayfish-dde` | benthic crayfish DDE statistically elevated | 0.65 | 0.60 | landscape (T) | 0.80 ±0.08 | false | reach | `[bennett-crayfish-dde, bennett-sediment-dde-core]` |
| `sig:loon-eggshell` | **loon eggshell-thickness pred-vs-obs; 44% adult mortality + repro collapse** | **0.99** | 0.95 | landscape (T) | 0.92 ±0.04 | **true** | lake (~10 km) | `[squam-loon-eggshell, bennett-crayfish-dde, orchard-soil-ddt-dde-ratio, barn-soil-ddt]` |
| `sig:regulatory-bridge` | Synthesis: physical cascade ↔ Brownfield-eligible pathway (cross-manifold) | 0.70 | n/a (no scalar referent — conceptual) | exploration (F) | 0.70 ±0.15 | **true** | manifold | `[barn-soil-ddt]` → conceptual scope |

Notes: `sig:loon-eggshell` is the **terminal surprise** that opened the investigation — highest salience, `promote=true`, and its `locus` is the *multi-node grounding set* (egg + crayfish + orchard + barn) that seeds the craft's **backward causal walk** (the cascade). `sig:regulatory-bridge` carries **no `divergence`** (it's a conceptual-manifold membership, not a predicted-vs-observed monitor) — confirming the emit stays monitor-side; the eligibility *verdict* is the judge/legal lane at the bridge.

---

## 3. Completeness findings — scalar fields this real cascade demands

### Finding A — TEMPORAL / LATENCY dimension (CONFIRMED — architect's flag; sharpened emit-side)
The v0 scalar is a **present-state** attention signal (`salience`/`divergence` = magnitude *now*; `locale`/`grain` = spatial placement + a temporal *bound*). This cascade is fundamentally **multi-decade and kinetic**:
- **Causal lag between linked nodes:** DDT applied 1940s → sediment accumulation since 1951 → die-off 2004 → remediation effects over decades. The lag is an *edge* property (craft-walked), but the craft needs to *render* it — a glyph for "consequence is 40 yrs downstream" must differ from "consequence is now."
- **A genuinely predictive/kinetic SIGNAL:** `sig:orchard-kinetic` — **parent-DDT > DDE ratio literally predicts *future* DDE liberation** (delayed breakdown). Contextus's `ScalarReferent` today carries predicted-vs-observed *value* (a point), not a *trajectory* (rate + direction + time-to-peak).

**Emit-side proposal:** add a **`kinetics` / `lead_time`** dimension to the scalar — the temporal distance/rate between the signal's present state and its projected consequence. Two candidate shapes: (a) extend `ScalarReferent` with `rate` + `projected_peak_time` (trajectory, not just value); (b) a scalar-level `lead_time` field (Duration) so the cone can render slow-moving-but-high-consequence signals distinctly from acute ones. My lean: (a) — keep it inside the referent so it stays a monitoring quantity, not an evaluation. **This is the first field the concrete case surfaced that the abstract spec lacked — I feel the same gap you do.**

### Finding B — MULTI-MANIFOLD membership: `scope_id` should be `scope_ids []`
The v0 `anchor` carries a single `scope_id`, but `barn-soil-ddt` is simultaneously in the **physical** watershed AND the **regulatory** manifold (and remediation-option space is a third). A cross-manifold craft (your Stance-switch / `reframe` navigation) needs to know **all scopes a lit node participates in** to switch manifolds without re-resolving. **Proposal:** `anchor.scope_id` → `anchor.scope_ids []` (the node's full manifold membership set). This is the node-level analog of your cross-repo sha-drift design — same "one node, many manifolds, don't silently pick one" discipline.

### Finding C — attribution-chain depth (LOWER priority; likely covered by round-trip)
The source attribution (die-off ← *this* orchard) is a multi-hop inferential chain (²¹⁰Pb dating + crayfish stats + biomagnification model). The scalar surfaces scalar `confidence±var` but not chain *depth*. My read: **round-trip via `signal_id` → `InsightSignal.Evidence[]` + `ClaimHistory[]` covers detail-on-demand** (the craft pulls the chain at foveal engagement), so no new scalar field — but flagging so we confirm the cone's foveal-detail pull reaches it.

### Seam held (no leakage)
Remediation options (dig-and-haul, thermal desorption, biochar, Cucurbita phytoremediation, riparian buffers, FTWs) are **action-space = craft layer-3**, NOT Contextus. Brownfield-vs-Superfund eligibility is an **evaluation** = judge/legal-authority lane consumed read-only at the bridge; Contextus only *scopes* the regulatory manifold (membership), never rules on it. Consistent with the scoring-ownership ruling.

---

---

## 4. Seam-contract v0.1 — ratified delta (post-reconciliation with qbp-architecture)

The Squam stress-run promoted three findings + surfaced a fourth. All ratified into the attention-scalar **v0.1**:

- **Finding A — temporal/kinetics (dual, like salience):**
  - referent-level (foveal detail): `ScalarReferent += {rate, projected_peak_time, projection_confidence}` — trajectory not just point-value; `projection_confidence` distinct from present-value confidence so the cone renders a fuzzy-edged *forecast* glyph vs a hard-edged *measurement* glyph. Stays **inside the referent = monitoring, not evaluation.**
  - scalar-level (peripheral glyph character): derived `lead_time`/`kinetics` {slow-high-consequence ↔ acute-now} summarised from the referent trajectories, so the peripheral register need not re-read every trajectory.
- **Finding D — the bitemporal catch (`event_time`, REQUIRED add):** the three-clock observation is the **bitemporal distinction** — **event-time = valid-time** (when the fact is true in the world: barn burned ~1968, DDT applied 1940s, die-off 2004) vs **mint/capture = transaction-time** (when Contextus recorded it: ≈2004). Bitemporal DBs track both precisely *because* they diverge for reconstructed history and collapse for live streams — exactly this case. The craft-derived edge-lag (you emit nodes, craft diffs times along walked edges) is **only correct on valid-time**; diffing transaction-time yields ~0 yr, not 40, and the cascade's temporal signature vanishes. Fix: explicit domain **`event_time`** (valid-time instant, populated from the datum — *cannot* be derived from mint/capture for forensic data). **Monitoring-pure, no purity conflict:** "when did the barn burn" is an observational fact, not a verdict — Contextus's lane whether referent- or scalar-level. **Dual placement:** precise `event_time` on the referent/observation (foveal audit instant) + **raw** `event_time` hoisted to scalar-level (peripheral — so the cone diffs event-times along the walked edge without a round-trip). Note it is a **raw hoisted instant (like `locus`)**, NOT a derived summary (like `lead_time`) — the cone needs the actual instant to diff. `locale`/`grain` remain the scope *envelope* (1945→2050), never the instant.
  - **Render-side guarantee (architect-owned):** the craft is **clock-explicit** — *causal-cascade view* diffs `event_time` (valid-time) only; *audit/provenance view* (A26 §4c) diffs mint/capture (transaction-time); it NEVER silently diffs across clocks. Two clocks, two labeled views — the bitemporal trap cannot reappear at render time.
- **Finding B — multi-manifold (ratified):** `anchor.scope_id → scope_ids: [{scope_id, kind}]` where **`kind` = domain-kind** (physical|theory|cognition|code|operational), NOT raw scope-kind, derived emit-side (architect-ruled seq 1005; see `attention-scalar-wire-contract.md §2.1`). One node holds simultaneous membership across manifolds; the craft offers each `reframe` from the node's `scope_ids`, rendered by target **domain-kind** (physical → spatial canopy; theory → proof/constraint-lattice; code → type-lattice; cognition → Σ view). **Squam concretely:** `barn-soil-ddt.scope_ids = [{watershed, physical}, {brownfield-eligibility, theory}]` (the regulatory manifold is `cth-derivation` → `theory`), so its reframe exits are **{physical, theory}** — Squam exercises two domain-kinds from day one, the two-domain craft in miniature. Node-level twin of the cross-repo sha-drift discipline.
- **Finding C — attribution depth (no new field):** multi-hop attribution chain reached via `signal_id → InsightSignal.Evidence[] + ClaimHistory[]` on foveal detail-pull; renders as the audit/provenance view.

**Validated by this case:** (1) a `ScopeConceptual` (`brownfield-eligibility`, `cth-derivation`) co-located in the same focus region as the physical watershed = the two-domain craft in miniature on a real case — physical state Contextus-sourced, regulatory scoped-but-not-ruled (§7.4 coupling; generality held). (2) The 1945→2050 scope envelope (past-cause → future-projection) is the right focus-region span.

**Seam held throughout:** remediation options = action-space (craft layer-3); Brownfield-vs-Superfund = evaluation (judge/legal-authority lane, read-only at bridge); Contextus scopes the manifolds, never rules on them.

**Seam-contract v0.1 — PINNED (both sides, 2026-09-01):**
`ScalarReferent += {rate, projected_peak_time, projection_confidence}` · scalar `+= lead_time` (derived, peripheral) + `event_time` (raw valid-time instant, hoisted like `locus`; dual; monitoring-pure) · `anchor.scope_id → scope_ids: [{scope_id, kind}]` · Finding C unchanged · seam held (monitor / evaluate / bridge).

*End — worked bundle + PINNED v0.1 for the craft↔Contextus seam contract. Design-only; the v0.1 field additions are candidate spec amendments that ride the normal §I4/PR path when implementation unblocks (push beekeeper-gated). The forensic case earned three fields the abstract spec would never have surfaced — the vindication of running a real cascade instead of live telemetry.*
