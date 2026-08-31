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

*End — worked bundle for craft↔Contextus completeness stress. Design-only; reconcile against qbp-architecture's craft-side render/navigate walk-through. Findings A/B are candidate scalar-contract amendments to fold into the seam-contract v0 → v0.1.*
