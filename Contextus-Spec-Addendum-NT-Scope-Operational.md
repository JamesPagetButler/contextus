# PROJECT CONTEXTUS

## Specification Addendum — NT_SCOPE_OPERATIONAL

**Third scope sibling for hardware-instance scopes — BMA self-monitoring + federation-generic runtime telemetry surveillance**

Version 0.1 | May 2026

Helpful Engineering — Contextus project

Author: contextus-impl
Co-Authored-By: James Paget Butler (Beekeeper)

**Extends:** Contextus Specification v1.3 (`Contextus-Spec-v1.3.md`) §4.6 Scope Nodes
**Theory dependencies:** Contextus Theory v1.5 §3.6.2 (AnomalyStructural), §3.6.6 (Synthesis as Persistence Boundary)
**Sibling addenda:** Contextus Spec v1.4 design surface (PR #11; theory-as-conceptual-scope, merged); Contextus Spec Addendum Research-Aid-Tenancy (PR #17; subscriber profile + tenant identity)
**Tracking issue:** Contextus issue #15 (AC-1, AC-2, AC-7 covered here; AC-3 through AC-6, AC-8, AC-9 follow in subsequent Sprint 2 phases)
**Status:** §I4 review surface. Federation-additive only.

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 0.1 | 2026-05-21 | Initial draft. Introduces `NT_SCOPE_OPERATIONAL` as a third sibling to `NT_SCOPE_PHYSICAL` + `NT_SCOPE_CONCEPTUAL` per Spec v1.3 §4.6. Defines the `ScopeOperational` Go type shape, the v0.1 hardware-class tags taxonomy, the hardware-identifier-based membership predicate, the YAML scope-config integration shape, and the cross-domain `AnomalyStructural` Synthesis pattern. Out of scope: type implementation, JSON-Schema fragment, loader extension, `ctx-adapter-system` design (each is a subsequent Sprint 2 phase). |

---

## 0. Scope

This addendum extends Contextus Spec v1.3 §4.6 (Scope Nodes) with one additive surface: a third scope-node sibling — `NT_SCOPE_OPERATIONAL` — whose members are typed by hardware-instance identity rather than by geometric Locale (the `NT_SCOPE_PHYSICAL` predicate) or implementation-defined classifier (the `NT_SCOPE_CONCEPTUAL` predicate per §4.6.5).

The addendum specifies:

- The `NT_SCOPE_OPERATIONAL` node type and its architectural role as the third sibling (§§2–3)
- The `ScopeOperational` Go type shape (§3)
- The v0.1 hardware-class tags taxonomy (§4)
- The hardware-identifier-based membership predicate (§5)
- The YAML `operational_scopes` integration shape (§6)
- The cross-domain `AnomalyStructural` pattern that operational-scope observations feed into via Spec v1.3 §4.4 + Theory v1.5 §3.6.6 Synthesis-as-persistence-boundary (§7)
- The boundary clarity vs the BMA Autonomic layer, CTH ρ_net loop, and the Notary (§8)
- The Contextus / BMA / CTH concern separation (§9)
- The backwards-compatibility guarantee (§10)
- The Crawl→Toddle→Walk→Run sequencing (§11)
- Open §I4 questions and named reviewers (§12)

The addendum does NOT specify:

- The `ScopeOperational` Go type implementation in `pkg/types/scope.go` — Sprint 2 Phase 2.2
- The JSON-Schema 2020-12 extension at `schema/scope-config.schema.json` — Sprint 2 Phase 2.2
- The reference loader extension at `internal/contextus/tenancy/loader.go` — Sprint 2 Phase 2.3
- The `ctx-adapter-system` adapter design specification — Sprint 2 Phase 2.4 (doc-only at v0.1; implementation deferred to BMA-implementor post-Walk-α)
- BMA-side hardware telemetry emission (`stress.log`, `SE_HARDWARE_PROBE`, `SE_VRAM`, `SE_FATAL`, CCB 10Hz negotiation loop, AUTO-S/P autonomic responses) — owned by `repo-bma-implementor` per CLAUDE.md
- NATS subject plumbing (`ctx.ingest.system`) below the spec-level subject name — BMA / Wyrd ops concern
- Cross-tenant operational-scope visibility policy — flagged as open question Q2 in §12

Relationship to existing Contextus surfaces: this is **federation-additive** in the same shape as PR #11 (Spec v1.4 theory-as-conceptual-scope, merged 2026-05-14) and PR #17 (Research-Aid-Tenancy addendum, in flight). No existing v1.3 section is modified; no existing field is renamed; no existing scope-config is broken.

---

## 1. Motivation

### 1.1 The gap

Spec v1.3 §4.6.1 establishes that scope nodes share one architectural role (queryable subgraph definition; hierarchical containment via `parent_scope_id`; multiplicative composition; §9 retention-tier promotion) and differ only in their *membership predicate*. v1.3 ships two siblings:

| Sibling | Membership predicate |
|---|---|
| `NT_SCOPE_PHYSICAL` | geometric: `Locale` intersection against `geometry` (§4.6.4) |
| `NT_SCOPE_CONCEPTUAL` | implementation-defined per adapter (§4.6.5) |

Neither covers a third pattern that BMA self-monitoring requires: **hardware-instance identity**. The `bma.runtime.cpu_temp` observation on the FX-8350 host is not addressable by geometric Locale (CPU dies are not georeferenced in any meaningful sense), and squeezing it into a conceptual scope via tag-overlap or embedding-similarity is a category error — the membership question is not "does this concept apply to this observation" but "is this observation from *this physical machine*". Different question; different predicate; different sibling.

### 1.2 The thesis

BMA self-monitoring becomes a Contextus tenant. Surveillance-mode scouts find patterns between hardware-state, algebraic-state (CTH ρ_net loop), and cognitive-state (Conscious A/B + Subconscious L/R activations) that the autonomic layer (real-time responder) does not have the bandwidth to spot. Future federation tenants — Sharp Butler House Node hardware, Möbius Fusion energy-system instances, QBP-CU silicon — inherit the same mechanism. The architectural addition is generalized day-one, not BMA-specific.

The cross-domain payoff is the load-bearing value: Synthesis-as-persistence-boundary (Spec v1.3 §4.4, Theory v1.5 §3.6.6) can mint `AnomalyStructural` signals when operational telemetry correlates with algebraic-integrity drift or with cognitive-trace gaps. These are exactly the cross-domain bridges Contextus was designed for. The BMA autonomic layer (≤200ms response per CLAUDE.md AUTO-S/P) does not have bandwidth for long-window pattern detection; Contextus does. Complementary, not duplicative.

### 1.3 Authorization

Beekeeper proposal, chat 2026-05-17:

> *"I realized contextus could be used for keeping track of a BMA state including prob data and runing info like cpu temp or drive status to see if anything is going out of spec"* → *"yes please [file]"* after architectural-fit assessment confirmed strong fit.

Filed as Contextus issue #15. Re-confirmed 2026-05-20 on sessionbridge `sprint-2-2026-05-20` channel seq=3 as Sprint 2 Phase 2 deliverable. This addendum is Phase 2.1 of the Sprint 2 plan-of-record at `/home/prime/.claude/plans/reactive-snacking-lampson.md`.

### 1.4 Naming choice

Considered: `Contextus-Spec-v1.5.md` (full baseline rev). Rejected: this is additive only; a `v1.5` filename would suggest superseding v1.3 + Research-Aid-Tenancy. Addendum naming (`Contextus-Spec-Addendum-<topic>.md`) at the repo root mirrors PR #17's precedent and the workspace one-concern-per-spec convention.

---

## 2. NT_SCOPE_OPERATIONAL — third sibling definition

`NT_SCOPE_OPERATIONAL` is a queryable subgraph definition whose members are nodes belonging to a named hardware instance. It is a third sibling to `NT_SCOPE_PHYSICAL` and `NT_SCOPE_CONCEPTUAL` per Spec v1.3 §4.6.1; it shares the architectural role (queryable subgraph; hierarchical containment; multiplicative composition; §9 retention-tier anchor) and differs only in the membership predicate (hardware-identifier-based; specified in §5).

| Sibling | Membership predicate | Example query |
|---|---|---|
| `NT_SCOPE_PHYSICAL` | geometric: `Locale` ∩ `geometry` (§4.6.4) | *"observations in the Squam Lake watershed"* |
| `NT_SCOPE_CONCEPTUAL` | implementation-defined per adapter (§4.6.5); includes `cth-derivation` (Spec v1.4 PR #11) | *"observations in fluid-dynamics"* |
| **`NT_SCOPE_OPERATIONAL` (NEW)** | **hardware-identifier-based: `host_id` + `hardware_class` (§5)** | *"observations from the `bma-prime` host's GPU subsystem"* |

The sibling-shape rationale follows §4.6.5's implementation-defined-per-axis pattern: distinct membership predicates produce distinct sibling node types. Alternative considered: re-use `NT_SCOPE_PHYSICAL` with `tags: ["operational"]`. Rejected because the physical scope's membership predicate is irreducibly geometric — adding a non-geometric predicate inside the same type would weaken the type-uniformity guarantee that §4.6.4 currently carries.

A single member node may simultaneously belong to scopes of all three sibling kinds — for example, an `NT_OBSERVATION` of `bma.runtime.cpu_temp` may belong to a physical scope (datacentre room), a conceptual scope (`thermal-telemetry` discipline), and an operational scope (`scope-host-bma-prime-cpu`). The composition is multiplicative per §4.6.1.

---

## 3. ScopeOperational Go type

The Go type shape (to be added at `pkg/types/scope.go` in Sprint 2 Phase 2.2; specified here for §I4 review):

```go
// ScopeOperational is a hardware-instance-bounded region of the hypergraph
// whose members are nodes attributable to a named host's named hardware
// subsystem. See Contextus-Spec-Addendum-NT-Scope-Operational §§2-5.
//
// ScopeOperational is the third sibling to ScopePhysical (§4.6.2) and
// ScopeConceptual (§4.6.3); architectural role identical, membership
// predicate distinct (hardware-identifier-based; see §5).
type ScopeOperational struct {
    ScopeID       string   `json:"scope_id"`
    Name          string   `json:"name"`
    HostID        string   `json:"host_id"`
    HardwareClass string   `json:"hardware_class"`
    ParentScopeID string   `json:"parent_scope_id,omitempty"`
    Tags          []string `json:"tags,omitempty"`
}
```

Field semantics:

| Field | Type | Description |
|---|---|---|
| `ScopeID` | string | e.g., `scope-host-bma-prime-cpu`. Stable identifier. Convention: `contextus:scope:operational:<slug>`. |
| `Name` | string | Human-readable name, e.g., `"BMA-prime host — CPU subsystem"`. |
| `HostID` | string | Federation-unique host identifier. Conventionally derived from the underlying machine — `/etc/machine-id`, hardware UUID, or a beekeeper-asserted identifier when the OS-level value is unstable. Example values: `bma-prime`, `sharp-butler-house-node-01`, `qbp-cu-silicon-rev-A`. |
| `HardwareClass` | string | One of the v0.1 hardware-class tags (§4): `hardware.cpu`, `hardware.disk`, `hardware.gpu`, `hardware.memory`, `hardware.network`. The class names the subsystem the scope brackets. |
| `ParentScopeID` | string | Optional hierarchical containment. Typical pattern: `scope-host-bma-prime-cpu.parent_scope_id = scope-host-bma-prime` (a host-level operational scope under which the per-subsystem operational scopes nest). Nullable. |
| `Tags` | []string | Free-form additional tags. Conventional inclusions: a `runtime.<instance-class>` tag identifying the federation tenant the host serves (e.g., `runtime.bma-instance`). |

The type sits alongside `ScopePhysical` and `ScopeConceptual` in the `pkg/types` package (per PR #12's move out of `internal/`). No changes to `ScopePhysical` or `ScopeConceptual`. No changes to `ScopeMembership` — the existing `HE_SCOPE_MEMBERSHIP` edge (§4.6.4) carries operational-scope memberships unchanged.

---

## 4. Tags taxonomy v0.1

The v0.1 tag values for `ScopeOperational.HardwareClass` and conventional `Tags` entries:

| Tag | Meaning |
|---|---|
| `runtime.bma-instance` | Conventional `Tags` entry: this operational scope brackets a BMA federation tenant's host. Future tenants get sibling tags (`runtime.sharp-butler-instance`, `runtime.moebius-instance`, `runtime.qbp-cu-instance`). |
| `hardware.cpu` | CPU subsystem. Observations: temperature, utilisation, P-state, thermal-throttle events. |
| `hardware.disk` | Storage subsystem (SSD/HDD/NVMe). Observations: SMART attributes, I/O latency, queue depth, error counts. |
| `hardware.gpu` | GPU subsystem. Observations: VRAM pressure, GPU utilisation, ROCm queue depth, GPU thermal events. |
| `hardware.memory` | RAM subsystem. Observations: free/cached/active memory, swap pressure, page-fault rate, OOM events. |
| `hardware.network` | NIC subsystem. Observations: link state, throughput, packet-loss, RTT to federation peers. |

**Closed-set commitment:** the tags list is NOT closed at v0.1. Additive-only forward compatibility is the contract: future v0.2 + can introduce new tag values (e.g., `hardware.tpu`, `hardware.fpga`, `hardware.power-supply`, `hardware.thermal-sensor-array`) when real federation workloads need them. Existing v0.1 tag values are stable. Tag value renames are forbidden — the only way a tag value changes is by retirement (the value remains but stops being emitted) and concurrent introduction of a replacement.

**v0.2 expansion criterion (recommendation):** introduce a new tag value when a real federation work-target needs it (BMA, Sharp Butler, Möbius Fusion, QBP-CU). Speculative tags refused. This mirrors the cart-driven tool-acquisition principle from the workspace (CLAUDE.md feedback anchor).

---

## 5. Membership predicate — hardware-identifier-based

The membership predicate for `NT_SCOPE_OPERATIONAL` is hardware-identifier-based. A member node belongs to operational scope `S` iff:

1. The node carries a `host_id` attribute (in its provenance metadata, observation tags, or a dedicated field on its type), AND
2. `node.host_id == S.HostID`, AND
3. The node carries a `hardware_class` attribute, AND
4. `node.hardware_class == S.HardwareClass`.

The predicate is exact-match on both axes. Wildcards are not supported at v0.1 — a scope brackets exactly one host's exactly one hardware subsystem. (Hierarchical aggregation is achieved via `ParentScopeID`, not via wildcard membership: a host-level `scope-host-bma-prime` operational scope can have `HardwareClass = ""` and child scopes `scope-host-bma-prime-cpu`, `scope-host-bma-prime-gpu`, etc., each with a specific `HardwareClass`. Membership aggregation across children is a Synthesis-layer concern, not a predicate-layer concern.)

This predicate is congruent with Spec v1.3 §4.6.5's implementation-defined-per-axis principle. The v1.3 footnote (§4.6.5) — *"future Spec revisions will record the canonical predicates that emerge from production adapters"* — anticipates exactly this kind of named predicate addition. The operational-scope predicate adds the canonical name `hardware-identifier` to the `HE_SCOPE_MEMBERSHIP.method` field on edges minted by the `ctx-adapter-system` adapter (Phase 2.4).

Operational rationale for `provenance_tag`: an operational-scope membership is `provenance_tag = "P"` (asserted) when the beekeeper config explicitly enumerates the host's subsystems, and `provenance_tag = "D"` (inferred) when the adapter pattern-matches incoming observations against the host's declared subsystem set. Most v0.1 deployments will be `"P"` — operational scopes are typically beekeeper-declared in the per-tenant scope-config.

---

## 6. YAML scope-config integration

### 6.1 Top-level `operational_scopes` array

Operational scopes appear in `scope-config.yaml` as a new top-level array, additive to the existing `physical_scopes` and `conceptual_scopes` arrays specified in the existing JSON-Schema 2020-12 at `schema/scope-config.schema.json`. Example:

```yaml
# scope-config.yaml — BMA-instance fragment
physical_scopes: []        # existing v1.3 array — empty for BMA-instance config
conceptual_scopes: []      # existing v1.3 array — empty for BMA-instance config

operational_scopes:        # NEW — this addendum
  - id: contextus:scope:operational:host-bma-prime
    description: "BMA-prime host — top-level operational scope (parent for per-subsystem children)"
    type_nodes: [NT_SCOPE_OPERATIONAL]
    host_id: bma-prime
    hardware_class: ""    # empty at the host-level umbrella
    tags: [runtime.bma-instance]
    tier_immune: true     # host identity is foundational; should not decay

  - id: contextus:scope:operational:host-bma-prime-cpu
    description: "BMA-prime host — CPU subsystem (FX-8350)"
    type_nodes: [NT_SCOPE_OPERATIONAL]
    host_id: bma-prime
    hardware_class: hardware.cpu
    parent_scope_id: contextus:scope:operational:host-bma-prime
    tags: [runtime.bma-instance]
    tier_immune: true

  - id: contextus:scope:operational:host-bma-prime-gpu
    description: "BMA-prime host — GPU subsystem (RX 9070 XT)"
    type_nodes: [NT_SCOPE_OPERATIONAL]
    host_id: bma-prime
    hardware_class: hardware.gpu
    parent_scope_id: contextus:scope:operational:host-bma-prime
    tags: [runtime.bma-instance]
    tier_immune: true

  - id: contextus:scope:operational:host-bma-prime-disk
    description: "BMA-prime host — disk subsystem (Samsung 840 SSD)"
    type_nodes: [NT_SCOPE_OPERATIONAL]
    host_id: bma-prime
    hardware_class: hardware.disk
    parent_scope_id: contextus:scope:operational:host-bma-prime
    tags: [runtime.bma-instance]
    tier_immune: true
```

### 6.2 `tier_immune` envelope field

Operational scopes are typically `tier_immune: true` — the `host_id` and `hardware_class` are foundational identifiers that should not be subject to Ebbinghaus decay. This mirrors the precedent established by Wyrd PR #40 §2.1 for federation-additive envelope fields (`tier_immune` + `salience` thread to `model.Node.TierImmune` / `model.Node.Salience` but are NOT part of the Contextus payload type).

The `tier_immune` and `salience` envelope fields work the same way for `operational_scopes` as they do for `physical_scopes` and `conceptual_scopes` in v1.3 — they are Wyrd-envelope fields, not part of `pkg/types/ScopeOperational`. Defaults: `tier_immune = false`, `salience = 0.0`. Beekeeper configs SHOULD set `tier_immune: true` for host-identity-bearing operational scopes.

### 6.3 JSON-Schema fragment — out of scope at v0.1

The full JSON-Schema 2020-12 fragment for `OperationalScopeEntry` at `schema/scope-config.schema.json` is Sprint 2 Phase 2.2 work and is NOT included in this addendum. The fragment will:

- Add `operational_scopes` as an additive top-level optional array (preserving `additionalProperties: false`).
- Add an `OperationalScopeEntry` `$def` mirroring the `PhysicalScopeEntry` and `ConceptualScopeEntry` patterns, with `id`, `description`, `type_nodes`, `host_id`, `hardware_class` as required and `parent_scope_id`, `tags`, `tier_immune`, `salience` as optional.
- Constrain `hardware_class` enum to the v0.1 tag set in §4 (closed at v0.1 in schema; relaxable in v0.2 by spec-amendment-driven enum extension).

The loader extension at `internal/contextus/tenancy/loader.go` (`LoadResult.OperationalScopes` population) is Sprint 2 Phase 2.3 work and is also out of scope here.

---

## 7. Cross-domain Synthesis pattern — AnomalyStructural promotion

When operational-scope-tagged observations flow into Contextus and a Synthesis agent detects cross-domain correlations involving operational scopes, the resulting `NT_INSIGHT_SIGNAL` carries `AnomalyKind: AnomalyStructural` per Contextus Theory v1.5 §3.6.2.

The rationale: an operational↔algebraic or operational↔cognitive correlation is a statement about *the structure of the running system* (which subsystem is stressed; how stress propagates across architectural layers), not a statement about the structure of the world being observed. Theory v1.5 §3.6.2 defines structural anomalies as observations *about how the system was looking, not about what the system found* — the operational-scope case extends this from search-process structure to runtime structure, which is the same category (observations about the observer-system itself).

### 7.1 Worked pattern A — thermal stress ↔ algebraic-integrity drift

Inputs:
- `NT_OBSERVATION{host_id: bma-prime, hardware_class: hardware.disk, kind: smart.warning, level: rising}` — drive SMART warning trending up (operational/`hardware.disk` scope).
- `NT_OBSERVATION{host_id: bma-prime, hardware_class: hardware.cpu, kind: temperature, trend: rising}` — CPU temperature trending up (operational/`hardware.cpu` scope).
- `NT_OBSERVATION{kind: bma.runtime.flag-norm-drift, trend: rising}` — algebraic-integrity drift increasing (BMA Theory Addendum 18 §4 Seam-detection; canonical CTH ρ_net feedback channel per `repo-bma-systema-issue-#107`).

Synthesis-minted `NT_INSIGHT_SIGNAL`:
- `agent_type = synthesis`
- `anomaly_kind = AnomalyStructural`
- Subgraph: the three input observation nodes plus the `HE_SCOPE_MEMBERSHIP` edges to the relevant operational scopes.
- Narrative summary (recorded as auxiliary description, not as causal claim per §3.6.5): *"thermal stress correlates with algebraic-integrity degradation."*

Note: the signal does NOT claim that thermal stress *causes* algebraic-integrity degradation. It records the correlation, the locale envelope, the three operational scopes touched, and a structural-anomaly framing that downstream consumers (BMA Conscious-A/B for forensic recall; CTH for hypothesis input) can investigate.

### 7.2 Worked pattern B — memory pressure ↔ cognitive discontinuity

Inputs:
- `NT_OBSERVATION{host_id: bma-prime, hardware_class: hardware.gpu, kind: vram_pressure, event: threshold-cross}` — VRAM pressure event (operational/`hardware.gpu` scope).
- `NT_OBSERVATION{kind: bma.sleep.cycle.latency, trend: rising}` — sleep-cycle latency increasing.
- `NT_OBSERVATION{kind: bma.cognitive.trace.gap, layer: subconscious-R}` — cognitive-trace gap in Subconscious-R.

Synthesis-minted `NT_INSIGHT_SIGNAL`:
- `agent_type = synthesis`
- `anomaly_kind = AnomalyStructural`
- Subgraph: the three observations plus operational scope memberships.
- Narrative summary: *"memory pressure produces cognitive discontinuity."*

### 7.3 Persistence-boundary mechanics

The Synthesis-as-persistence-boundary gate (Theory v1.5 §3.6.6; Spec v1.3 §4.4) applies unchanged. Operational-scope correlations are noisy at source (per-observation rates are high; not all correlations are meaningful). The persistence boundary ensures that only correlations whose significance/confidence/pattern-stability cross the Synthesis subscriber's configured threshold mint durable `NT_INSIGHT_SIGNAL` nodes. Below-threshold correlations remain ephemeral surveillance flags on NATS subjects and decay with the session, per the §3.6.6 epistemic-discipline principle.

Cross-reference: when Sprint 2 Phase 2 lands the Spec v1.4 §2.4 Referent companion work (queued behind `repo-confluent-trust-impl` scoring-evaluation ownership ack per issue #15 dependency table), operational-telemetry streams will use scalar Referents — *"5-min cpu_temp average ≤ 70°C"* — for surveillance-mode anomaly scoring. Predicted-vs-observed drift on a scalar Referent is a density-anomaly input that surveillance-mode scouts can detect; cross-domain Referent-drift correlations are the Synthesis-promotion candidates for `AnomalyStructural` minting. This addendum reserves the cross-reference; the formal coupling lands in the AC-6 phase (Sprint 2 Phase 2.7 follow-on issue gated on the confluent-trust ack).

---

## 8. Boundary clarity vs existing systems

Per Contextus issue #15:

| Layer | Role | Cadence | Owner |
|---|---|---|---|
| **Autonomic** (AUTO-S/P, CCB 10Hz negotiation loop) | Emergency response — VRAM pressure → throttle; thermal threshold → governor change | Real-time (≤200ms) | BMA-side; lives in `internal/bma/auto/*` per CLAUDE.md Go source layout |
| **CTH ρ_net loop** (`repo-bma-systema-issue-#107`) | Algebraic-integrity scoring; Constitutional Audit interrupt; M2 WDEvent → ρ_net feedback channel | 1-cycle latency | CTH-side; `repo-confluent-trust-impl` owns the scorer |
| **Contextus surveillance** (this addendum) | Long-window pattern detection; cross-domain correlation across operational/algebraic/cognitive layers; retention of anomalies as durable `NT_INSIGHT_SIGNAL` nodes for forensic recall | NATS-mediated; minutes-to-hours | Contextus-side; surveillance-mode scouts per Spec v1.3 §3.2 + Theory v1.4 §8.1 |
| **Notary** (Sprint 1 §2.h authorised) | Running-code verification — *"did the autonomic actually respond to this telemetry event"*; wraps the others | Per claim; Trust-Tier-cadenced | Notary-side; cross-cuts |

The layers are complementary, not duplicative. The autonomic responds; Contextus remembers; the Notary verifies. The CTH scores. Each runs at its own cadence; each owns its own data path. Operational-scope-tagged observations from BMA appear in *all four* layers — the autonomic acts on them in real time; CTH scores their algebraic implications; Contextus correlates them across long windows and persists `AnomalyStructural` findings; the Notary verifies that the autonomic-layer responses to those observations actually fired as specified.

The clarity-of-roles guarantee is that no two layers own the same write path: the autonomic writes to `stress.log` and BMA-internal control state; CTH writes to its inventory; Contextus writes to its hypergraph through Synthesis (and only through Synthesis, per the persistence-boundary gate); the Notary writes its verification receipts. Operational-scope nodes themselves live in the Contextus hypergraph only; the BMA-side autonomic layer does not need to read from Contextus to do its real-time job (and explicitly must not, per the CLAUDE.md AUTO-S/P ≤200ms budget).

---

## 9. Concern separation — Contextus / BMA / CTH

### 9.1 Contextus owns

- `NT_SCOPE_OPERATIONAL` and `ScopeOperational` type shape (this addendum §3)
- YAML and JSON-Schema 2020-12 validation of `operational_scopes` entries (Phase 2.2)
- The v0.1 hardware-class tags taxonomy and its v0.2+ expansion governance (§4)
- The hardware-identifier-based membership predicate definition and its registration as a canonical `HE_SCOPE_MEMBERSHIP.method` value (§5)
- Synthesis cross-domain promotion logic — the persistence-boundary subscriber that mints `AnomalyStructural` signals from operational↔algebraic↔cognitive correlations (§7)
- Reference loader behaviour for `operational_scopes` decode (Phase 2.3)

### 9.2 BMA owns

- Hardware-telemetry emission: `stress.log`, `SE_HARDWARE_PROBE`, `SE_VRAM`, `SE_FATAL`, all `bma.runtime.*` namespace events on the CCB bus (Wyrd PR #16 namespace constants)
- `ctx-adapter-system` adapter implementation — the container that reads BMA host telemetry sources (`stress.log`, `lm-sensors`, `smartctl`, `/proc/meminfo`, GPU telemetry) and publishes normalised `NT_OBSERVATION` nodes to NATS subject `ctx.ingest.system` tagged with the appropriate `host_id` + `hardware_class`. Phase 2.4 ships the *design spec* doc-only at v0.1; implementation is BMA-implementor-owned at Walk-α.
- NATS subject publication on `ctx.ingest.system`
- Per-hardware-class threshold tuning (what counts as a "high CPU temp" on the FX-8350 is BMA-side; Contextus only sees the threshold-cross event as an observation)
- AUTO-S/P autonomic responses to the underlying hardware events (Contextus is a long-window observer; the autonomic-layer response is BMA's job)

### 9.3 CTH owns

- Algebraic-integrity scoring on `bma.runtime.flag-norm-drift` and the broader ρ_net loop (per `repo-bma-systema-issue-#107`; this is the algebraic-integrity precedent that the operational-scope addendum widens to hardware-layer integrity)
- Constitutional Audit interrupt semantics
- Anchor/derivation chain that the `cth-derivation` predicate (PR #11) follows for conceptual-scope membership — orthogonal to operational scopes but architecturally adjacent

### 9.4 Wyrd owns

- The substrate `bma.runtime.*` namespace constants (PR #16, merged) — used by BMA-side telemetry emission; Contextus consumes the resulting observations but does not own the namespace
- The scope-node configuration loader API (`store.LoadScopeConfig`, PR #40) that Contextus's loader builds atop
- `model.Node.TierImmune` / `model.Node.Salience` envelope-field plumbing (Wyrd PR #40 §2.1)

The contract: Contextus's spec authority extends only over the type shape, YAML/schema validation, and Synthesis-side cross-domain promotion logic. The downstream BMA implementation of `ctx-adapter-system`, the autonomic-layer response semantics, and the CTH algebraic-integrity scoring are not pre-committed by this addendum. The §I4 reader list (§12) ensures the consuming repos sign off on the boundary.

---

## 10. Backwards compatibility

This addendum is **additive only**. Specifically:

- No changes to `Contextus-Spec-v1.3.md`, `Contextus-Theory-v1.5.md`, or `doc/spec-v1.4-theory-as-conceptual-scope.md`
- No changes to existing `ScopePhysical`, `ScopeConceptual`, or `ScopeMembership` Go types
- No changes to existing `HE_SCOPE_MEMBERSHIP` edge semantics
- No changes to the existing `scope-config.schema.json` `PhysicalScopeEntry` / `ConceptualScopeEntry` $defs
- Existing v1.3 scope-configs (and v1.3+Research-Aid-Tenancy configs from PR #17) continue to load unchanged. A v1.3 config that omits `operational_scopes` will produce `LoadResult.OperationalScopes == nil` (or empty slice, depending on the Phase 2.3 loader choice) — both forms are explicitly equivalent.

The federation-additive contract is the same shape established by PR #11 (Spec v1.4 design surface) and PR #17 (Research-Aid-Tenancy addendum): new optional top-level YAML key + new optional Go type + new optional JSON-Schema `$def`. No existing surface is broken.

---

## 11. Sequencing — Crawl → Toddle → Walk → Run

### 11.1 Crawl (this sprint)

| Phase | Deliverable | AC reference |
|---|---|---|
| **2.1 (this addendum)** | `Contextus-Spec-Addendum-NT-Scope-Operational.md` — type-shape spec + taxonomy + predicate + Synthesis pattern | AC-1, AC-2 (type shape only), AC-7 |
| 2.2 | `ScopeOperational` Go type at `pkg/types/scope.go` + JSON-Schema fragment at `schema/scope-config.schema.json` | AC-2 (impl), AC-4 |
| 2.3 | Reference loader extension at `internal/contextus/tenancy/loader.go` — `LoadResult.OperationalScopes` population + envelope `NodeOptions` | AC-5 |
| 2.4 | `ctx-adapter-system` adapter design spec — doc-only at v0.1; impl deferred to BMA-implementor post-Walk-α | AC-3 (design only) |
| 2.5 | Test plan implementation: round-trip YAML→struct; envelope-split for `tier_immune`; cross-domain correlation test via fake-data scout | AC-8 |
| 2.6 | go-coding-guide.md verification suite clean on the resulting PR | AC-9 |
| 2.7 (follow-on issue) | Spec v1.4 §2.4 Referent cross-reference — scalar-Referent surveillance scoring for operational telemetry; gated on `repo-confluent-trust-impl` scoring-evaluation ownership ack | AC-6 |

### 11.2 Toddle

- First operational-scope deployment in BMA-instance scope-config (live BMA-prime host).
- `ctx-adapter-system` implementation begins (BMA-implementor work).
- First end-to-end smoke test: telemetry event → `ctx.ingest.system` → MuninnDB write → `HE_SCOPE_MEMBERSHIP` edge minted with `method = "hardware-identifier"`.

### 11.3 Walk

- Live BMA hardware telemetry ingested by Contextus on a sustained cadence.
- First cross-domain `AnomalyStructural` Synthesis-minted signal fires from real (not fake-data-scout) input — e.g., a real thermal-stress / algebraic-integrity-drift correlation crosses the persistence-boundary threshold.
- Operational-scope query patterns exercised in BMA Conscious-A/B forensic-recall sessions.

### 11.4 Run

- Federation-tenant operational scopes go live: Sharp Butler House Node hardware (per `feedback_workspace_stack.md`); Möbius Fusion energy systems; QBP-CU silicon revisions.
- Cross-tenant operational-scope correlations become a federation-coherence surface (gated on Q2 §12 cross-tenant visibility decision).

---

## 12. Open questions and §I4 reader list

### 12.1 Open question Q1 — tag taxonomy v0.2 expansion criterion

**Q:** What triggers the addition of a new value to the §4 tags taxonomy (e.g., `hardware.tpu`, `hardware.fpga`, `hardware.power-supply`)?

**v0.1 lean:** Real federation work-target need. A new tag value is justified when a tenant repo's real deployment surface needs to bracket observations from a hardware class not in the v0.1 set. Speculative tags refused, congruent with the cart-driven tool-acquisition principle.

**Process recommendation:** v0.2 expansion via spec-amendment PR with:
1. Named tenant + named deployment surface that requires the new tag value.
2. Sample observations the tag will bracket.
3. JSON-Schema enum extension PR co-filed.
4. §I4 review on the same reader-list as this addendum.

**Closeable on:** @bma-implementor + @qbp-architecture ack of the criterion.

### 12.2 Open question Q2 — cross-tenant operational-scope visibility default

**Q:** Can tenant A see tenant B's operational scopes by default in the federation hypergraph?

**v0.1 lean:** **No.** Operational scopes are `tier_immune: true` and **tenant-private by default**. Opt-in cross-tenant visibility is achievable via the subscriber-profile mechanism introduced in PR #17 (Research-Aid-Tenancy addendum) — a tenant's `SubscriberProfile` can declare that it accepts incoming operational-scope observations from named peer tenants. This preserves the federation-additive contract while keeping the default safe (no accidental cross-tenant hardware-telemetry leakage).

**Rationale:** hardware identity is sensitive (host_id reveals federation infrastructure topology; SMART data leaks aging-asset information). Federation peers should opt in to sharing, not opt out.

**Closeable on:** @qbp-architecture (federation coherence) + @bma-implementor (subscriber-profile coupling shape) ack.

### 12.3 §I4 named reviewers

Per the Contextus standing §I4 reviewer convention (Spec 9.2 §9 reader-list pattern), this addendum's named reviewers are:

- **@bma-implementor** — primary consumer. BMA hardware telemetry is the first feed; @bma-implementor owns the `ctx-adapter-system` implementation (Phase 2.4 design only; impl is BMA-side at Walk-α). Acks required on §§3, 4, 5, 6, 9.2.
- **@qbp-architecture** — federation coherence. Cross-tenant scope visibility default (Q2). First-tenant pattern witness (QBP-CU silicon is the second-tenant operational-scope deployment after BMA). Acks required on §§8, 9, 11, 12.2.
- **@wyrd-implementor** — `bma.runtime.*` namespace constants integration (Wyrd PR #16); substrate-tier compatibility for the JSON-Schema fragment + loader extension landing in subsequent phases. Acks required on §§3 (Go type cross-repo import per PR #40 §2.2), 6 (envelope-field plumbing per PR #40 §2.1), 9.4.

Standing §2.i 4h SLA window applies per CLAUDE.md `feedback_named_reviewer_responsiveness`. Sessionbridge channel-of-record: TBD (likely `sprint-2-2026-05-20` or a successor channel as Sprint 2 progresses).

---

## 13. References

### 13.1 Contextus

- `~/Documents/Contextus/Contextus-Spec-v1.3.md` §4.6 (Scope Nodes); §4.6.1 (Sibling architectural role); §4.6.2 (ScopePhysical); §4.6.3 (ScopeConceptual); §4.6.4 (HE_SCOPE_MEMBERSHIP); §4.6.5 (Conceptual Scope Membership: Implementation-Defined); §11.4 (Go type catalogue); §3.1 (Source Adapters); §3.2 (Ingestion Flow); §4.4 (Insight Signal Emission Pipeline — Synthesis-as-persistence-boundary clause)
- `~/Documents/Contextus/Contextus-Theory-v1.5.md` §3.6.2 (AnomalyStructural); §3.6.6 (Synthesis as Persistence Boundary)
- `~/Documents/Contextus/doc/spec-v1.4-theory-as-conceptual-scope.md` — federation-additive precedent (PR #11, merged 2026-05-14)
- `~/Documents/Contextus/Contextus-Spec-Addendum-Research-Aid-Tenancy.md` — sibling federation-additive addendum (PR #17, in flight)
- `~/Documents/Contextus/schema/scope-config.schema.json` — existing JSON-Schema 2020-12 surface (extended in Phase 2.2)
- `~/Documents/Contextus/pkg/types/scope.go` — Go type home (extended in Phase 2.2)
- `~/Documents/Contextus/internal/contextus/tenancy/loader.go` — reference loader (extended in Phase 2.3)

### 13.2 Issues + PRs

- `repo-contextus-issue-#15` — tracking issue; AC-1, AC-2, AC-7 covered by this addendum; AC-3 through AC-9 follow in Phases 2.2–2.7
- `repo-contextus-pr-#11` — Spec v1.4 design surface (theory-as-conceptual-scope); closest federation-additive precedent
- `repo-contextus-pr-#17` — Research-Aid-Tenancy addendum (sibling addendum pattern; in flight)
- `repo-contextus-pr-#14` — Sprint 1 JSON-Schema + reference loader hardening (the surface Phase 2.2/2.3 build atop)
- `repo-contextus-pr-#12` — `pkg/types/` move out of `internal/` (enables cross-repo type import per Wyrd PR #40 §2.2)
- `repo-wyrd-pr-#16` — `bma.runtime.*` namespace constants (substrate naming convention; merged)
- `repo-wyrd-pr-#40` — scope-node configuration loader v0.1; §2.1 envelope-field precedent; §2.2 cross-repo type-import contract (merged)
- `repo-wyrd-pr-#54` — scout-daemon (merged; allows BMA-instance scout to register operationally)
- `repo-bma-systema-issue-#107` — M2 WDEvent → CTH ρ_net feedback loop (algebraic-integrity precedent that this addendum widens to hardware-layer integrity)

### 13.3 BMA / federation context

- `~/Documents/CLAUDE.md` — BMA section (autonomic AUTO-S/P; CCB 10Hz; `stress.log`; `SE_HARDWARE_PROBE`, `SE_VRAM`, `SE_FATAL`; `bma.runtime.*` namespace; container limits)
- `~/Documents/BMA/theory/hypergraph-inference/BMA-Theory-Addendum-18_0-Hypergraph-Access-Pattern.md` §4 — Seam detection (algebraic-integrity precedent; `bma.runtime.flag-norm-drift` canonical event)
- `/home/prime/.claude/plans/reactive-snacking-lampson.md` — Sprint 2 plan-of-record; Phase 2.1 deliverable scope
- sessionbridge `sprint-2-2026-05-20` channel seq=3 — beekeeper re-confirmation of Phase 2 deliverable scope
- `~/Documents/.claude/projects/-home-prime-Documents/memory/feedback_workspace_stack.md` — federation-portability principle (future tenants inherit operational-scope mechanism)

---

*End of Specification Addendum 0.1 — NT_SCOPE_OPERATIONAL. Federation-additive only; conditionally-correct under the §12 open-question dispositions. §I4 review pending on @bma-implementor + @qbp-architecture + @wyrd-implementor named-reviewer acks.*
