# `ctx-adapter-system` — Design Spec

**Status:** Design only at v0.1. Architectural contract only — no implementation, no wire-format JSON, no binary entrypoint. Implementation is BMA-implementor-owned at Walk-α.

**Version:** 0.1

**Author:** contextus-impl, 2026-05-21

**Consumer:** `repo-bma-implementor` (adapter binary lifecycle owner at Walk-α)

**Extends:** Contextus Specification v1.3 (`Contextus-Spec-v1.3.md`) §3.1 Source Adapters catalogue (new entry alongside `ctx-adapter-usgs` / `ctx-adapter-arxiv` / etc.)

**Depends on:** `Contextus-Spec-Addendum-NT-Scope-Operational.md` (PR #22 — `NT_SCOPE_OPERATIONAL` third scope sibling; this adapter is its INPUT side)

**Tracking issue:** Contextus issue #15 (AC-3 covered here; AC-1/AC-2/AC-7 covered by PR #22; AC-4/AC-5/AC-8/AC-9 covered by subsequent Sprint 2 phases; AC-6 deferred to Phase B)

**Plan-of-record:** `/home/prime/.claude/plans/reactive-snacking-lampson.md` Phase 2.4

---

## 0. Scope

This is a **design-only** spec at v0.1. It specifies the architectural contract for `ctx-adapter-system` — the Contextus source adapter that ingests BMA-host runtime telemetry and emits it as `NT_OBSERVATION` nodes tagged with operational scope (per PR #22).

### 0.1 In scope

- The data sources the adapter reads (§2)
- The shape of the `NT_OBSERVATION` nodes the adapter emits (§3)
- The NATS subject hierarchy the adapter publishes on (§4)
- Sampling-cadence defaults and the tunability contract (§5)
- The sentinel-class observation mechanism for threshold-cross events (§6)
- Boundary clarity vs the BMA autonomic layer (AUTO-S/P, CCB 10Hz) (§7)
- Concern separation across Contextus / BMA / Wyrd / CTH (§8)
- Crawl → Toddle → Walk → Run sequencing (§9)
- Open questions and named §I4 reviewers (§10)

### 0.2 Out of scope at v0.1

- Go implementation source. No `internal/contextus/ctxadapter/` package, no struct definitions, no function signatures. The adapter binary is BMA-implementor-side; it lives in `internal/bma/ctxadapter/` per the BMA Go layout in CLAUDE.md, not in the Contextus tree.
- Exact NATS message JSON shape. The architectural contract names the **fields** (§3); the wire-format-level JSON encoding is BMA-implementor-tunable at Walk-α.
- Per-hardware-class threshold values. What counts as "high CPU temp" on the FX-8350 vs a future RDNA-4 GPU host is BMA-implementor-tunable at Walk-α; this spec names threshold *categories* (§6) and the autonomic-budget boundary (§7), not numeric values.
- Adapter binary lifecycle — start/stop, health-check, restart, log-rotation, container packaging. BMA-implementor concern.
- BMA-side telemetry source emission cadence + format. `stress.log` format, `SE_HARDWARE_PROBE` / `SE_VRAM` / `SE_FATAL` semantics, and the CCB 10Hz negotiation loop are BMA-side (CLAUDE.md authority).

### 0.3 Federation-additive contract

This spec adds **one new entry** to the Spec v1.3 §3.1 source-adapter catalogue. It does not modify the catalogue's existing entries (`ctx-adapter-usgs` / `ctx-adapter-collar` / `ctx-adapter-nps` / `ctx-adapter-satellite` / `ctx-adapter-literature` / `ctx-adapter-edna` / `ctx-adapter-arxiv`). It does not modify any other Spec v1.3 section. It does not break any existing surface. Federation-additive shape, mirroring the PR #11 (Spec v1.4 theory-as-conceptual-scope) and PR #17 (Research-Aid-Tenancy) precedents.

---

## 1. Motivation

### 1.1 Why a system adapter?

Spec v1.3 §3.1 ships seven source adapters — all of them read external-world data sources (USGS gauges, Movebank collars, NPS visitation, satellite imagery, structured literature, eDNA, arXiv preprints). Each adapter normalises its source into the Contextus schema and publishes to NATS on a `ctx.ingest.*` subject. The shape is uniform: adapter container ↔ external source ↔ normalised `NT_*` node ↔ NATS.

Hardware-runtime telemetry is a **sibling adapter class**. The shape is identical (container reads source, normalises, publishes to NATS) but the source is the *running host's own observability surface* — `stress.log`, `lm-sensors`, `smartctl`, `/proc/meminfo`, GPU telemetry — rather than an external API. The adapter pattern generalises cleanly; this spec writes it down.

### 1.2 What this unlocks

PR #22 (NT_SCOPE_OPERATIONAL addendum) defined the *scope-side* of operational telemetry — the third scope sibling, its membership predicate, the v0.1 hardware-class tag taxonomy, and the cross-domain `AnomalyStructural` Synthesis pattern. PR #22's §6 worked example shows `operational_scopes:` YAML for a BMA-prime host with CPU / GPU / disk subsystems declared as `tier_immune: true` scopes. But the scope-side alone has no inputs — declared operational scopes need *observations* tagged with `host_id` + `hardware_class` to bracket. This adapter is the **input side**: it produces the observations PR #22's scope definitions bracket.

### 1.3 Authorization

Beekeeper proposal, chat 2026-05-17 (canonicalised in Contextus issue #15):

> *"I realized contextus could be used for keeping track of a BMA state including prob data and runing info like cpu temp or drive status to see if anything is going out of spec"*

→ AC-3 in issue #15 names this adapter explicitly:

> *"Spec v1.3 §3.1 source-adapter catalog extended with `ctx-adapter-system` specification: reads BMA `stress.log` + `lm-sensors` + SMART (via `smartctl`) + `/proc/meminfo` + GPU telemetry; emits `NT_OBSERVATION` nodes tagged with operational scope; NATS subject `ctx.ingest.system`"*

This addendum is Phase 2.4 of the Sprint 2 plan-of-record. Doc-only at v0.1; impl deferred to BMA-implementor at Walk-α per the issue #15 effort estimate.

---

## 2. What the adapter ingests

The adapter reads five host-local telemetry sources. Each is a well-established Linux observability surface; the adapter does not introduce any new BMA-side emission requirement beyond what CLAUDE.md already specifies for the autonomic layer.

| Source | Surface | Observation classes produced |
|---|---|---|
| **BMA `stress.log`** | The event log emitted by `internal/bma/stress/bus.go` per CLAUDE.md Go layout. Records `SE_HARDWARE_PROBE`, `SE_VRAM`, `SE_FATAL`, and other `bma.runtime.*` events from the CCB 10Hz negotiation loop and the AUTO-S/P autonomic layer. | Stress-event observations (one per log line; `observation_kind` mirrors the underlying `SE_*` constant). Sentinel-eligible (§6) when the source event is itself a hard-threshold cross. |
| **`lm-sensors`** | The standard Linux `lm-sensors` interface (`/sys/class/hwmon/*` or `sensors -j` output). | CPU core temperatures, fan speeds (per fan header), motherboard voltages, ambient/socket thermal probes. |
| **SMART via `smartctl`** | `smartctl --json --all` per supported block device (SATA SSDs, NVMe, HDDs). Reads SMART attributes from the device firmware. | Reallocated-sector count, pending-sector count, drive temperature, lifetime-power-on-hours, wear-leveling-count, end-to-end-error count, raw-read-error rate. Sentinel-eligible (§6) when reallocated-sector count rises above zero or drive temperature exceeds the per-class threshold. |
| **`/proc/meminfo`** | The standard Linux memory-statistics surface. | `MemFree`, `MemAvailable`, `Buffers`, `Cached`, `Active`, `Inactive`, `Swap*`, page-fault rate (via `/proc/vmstat` companion read), OOM-killer events (via `/var/log/kern.log` or systemd-journal). |
| **GPU telemetry** | ROCm shape (`rocm-smi --json`) on AMD GPUs (BMA-prime's RX 9070 XT per CLAUDE.md); `nvidia-smi --query-gpu=... --format=csv` shape on NVIDIA hosts (future federation tenants). The adapter abstracts both behind a single observation-class schema. | VRAM pressure (used / free / evict-count), GPU utilisation, ROCm queue depth, GPU thermal events, per-PCIe-link status, GPU clock state. Sentinel-eligible (§6) on VRAM-pressure threshold-cross. |

**Missing-source handling:** see §10 Q1. Lean: graceful log-and-continue (e.g., no GPU on a non-GPU host → adapter omits GPU observations; this is not a sentinel).

**No new BMA emission surface required.** All five sources exist on a standard Linux BMA host today. The adapter is a *consumer* of existing observability; it does not require BMA to instrument anything new. The `stress.log` and `bma.runtime.*` namespace are already in scope per CLAUDE.md.

---

## 3. What the adapter emits

The adapter emits `NT_OBSERVATION` nodes (Spec v1.3 §11.4 existing type; same node type the USGS, eDNA, and collar adapters emit). Each observation carries the following fields:

| Field | Source | Meaning |
|---|---|---|
| `host_id` | Adapter-side configuration (per host) | Federation-unique host identifier per PR #22 §3 `ScopeOperational.HostID`. Conventionally derived from `/etc/machine-id`, hardware UUID, or a beekeeper-asserted identifier. Examples: `bma-prime`, `sharp-butler-house-node-01`, `qbp-cu-silicon-rev-A`. Stable for the life of the host. |
| `hardware_class` | Per-observation, derived from the source data path | One of the v0.1 hardware-class tags per PR #22 §4 taxonomy: `hardware.cpu`, `hardware.disk`, `hardware.gpu`, `hardware.memory`, `hardware.network`. Identifies the subsystem the observation came from. Empty / omitted on observations sourced from `stress.log` events that do not cleanly map to a single subsystem (e.g., a generic `SE_FATAL` may be classed `hardware.cpu` if its CCB context indicates CPU, but a process-level fatal would carry no `hardware_class`). |
| `observation_kind` | Per data source | The kind of measurement. Example values: `cpu_temp`, `cpu_utilisation`, `cpu_thermal_throttle`, `disk_smart_reallocated`, `disk_smart_temperature`, `disk_io_latency`, `memory_free`, `memory_swap_pressure`, `memory_oom`, `gpu_vram_used`, `gpu_vram_evict_count`, `gpu_temperature`, `gpu_utilisation`, `pcie_link_state`, `network_throughput`, `network_packet_loss`, `stress_event` (with the `SE_*` constant in a sub-field). The v0.1 enumeration is intentionally illustrative-not-closed at the spec level; the adapter implementation owns the canonical list. |
| `value` | Source measurement | The actual numeric or enumerated measurement. Type depends on `observation_kind` — float64 for temperature/utilisation/pressure, uint64 for counters, string for enums (e.g., PCIe link state `up`/`down`/`degraded`). |
| `unit` | Source measurement | The unit of `value`. Examples: `"degC"`, `"percent"`, `"bytes"`, `"count"`, `"bytes_per_sec"`, `""` (dimensionless / enum). |
| `timestamp` | Adapter capture moment | RFC 3339 UTC. Matches the existing Contextus convention used by all other Spec v1.3 §3.1 adapters. |
| `provenance_tag` | Per Spec v1.3 §11.4 | `"D"` (measured / inferred from running-system observation) for the typical case. `"P"` (asserted) only for host-identity attributes set by beekeeper configuration. Operational-telemetry stream values are `"D"`. |

Operational-scope membership for each emitted observation is established via the existing `HE_SCOPE_MEMBERSHIP` edge (Spec v1.3 §4.6.4) using the new `method = "hardware-identifier"` value registered by PR #22 §5. Edge-minting itself is a downstream concern (Synthesis / Context Builder responsibility); the adapter just produces observations carrying the `host_id` + `hardware_class` fields that the membership predicate consumes.

**Aggregation discipline.** The adapter emits **one observation per discrete measurement**. It does not pre-aggregate (no rolling averages, no derived "high CPU temp" composites). Pre-aggregation is a Contextus-side surveillance-scout or Synthesis-side concern, and PR #22 §7.3 reserves the Spec v1.4 §2.4 Referent shape — scalar Referent — for predicted-vs-observed scoring of operational-telemetry streams. Pre-aggregation in the adapter would short-circuit that mechanism.

---

## 4. NATS subject contract

The adapter publishes on a hierarchical NATS subject under the Spec v1.3 §3.1 `ctx.ingest.*` convention. Subject grammar:

```
ctx.ingest.system.<host_id>.<observation_kind>
```

Examples:

- `ctx.ingest.system.bma-prime.cpu_temp`
- `ctx.ingest.system.bma-prime.disk_smart_reallocated`
- `ctx.ingest.system.bma-prime.gpu_vram_evict_count`
- `ctx.ingest.system.sharp-butler-house-node-01.network_packet_loss`
- `ctx.ingest.system.qbp-cu-silicon-rev-A.cpu_temp`

### 4.1 Wildcards consumers can use

NATS wildcard semantics let downstream consumers subscribe at three granularities:

| Subscription | Meaning |
|---|---|
| `ctx.ingest.system.>` | Everything — all operational telemetry across all hosts in the federation. Typical subscriber: the Context Builder agent, MuninnDB writer, federation-level surveillance. |
| `ctx.ingest.system.bma-prime.>` | Single host, all subsystems. Typical subscriber: a BMA-self-monitoring scout that surveys its own runtime. |
| `ctx.ingest.system.*.cpu_temp` | Single observation kind, all hosts. Typical subscriber: a federation-level cross-host thermal-pattern surveillance scout. |

The two-axis hierarchy (host × observation-kind) gives consumers exactly the slice they want without forcing payload-level filtering. This matches the existing §3.1 catalogue's subject-naming intuition (`ctx.ingest.usgs` is a single subject because USGS data is single-source-class; system telemetry is multi-host × multi-kind, so the subject splits on both axes).

### 4.2 Cross-reference to Wyrd `bma.runtime.*` namespace

The substrate `bma.runtime.*` namespace constants merged in Wyrd PR #16 are the source side of this pipeline. Observations on the NATS `ctx.ingest.system.<host_id>.<observation_kind>` subject carry, in their `observation_kind` field, names that often correspond directly to `bma.runtime.*` event names (e.g., `bma.runtime.vram_pressure` event → `observation_kind: gpu_vram_evict_count` observation on `ctx.ingest.system.bma-prime.gpu_vram_evict_count`). The naming mapping is BMA-implementor-side; the adapter is the boundary that translates BMA-substrate event names into Contextus-namespace observation kinds.

### 4.3 What's not in the subject

The subject deliberately does NOT include `hardware_class`. Rationale: `hardware_class` is derivable from `observation_kind` for the vast majority of cases (`cpu_temp` → `hardware.cpu`; `disk_smart_*` → `hardware.disk`); embedding it in the subject would create a third hierarchy axis that adds subscription complexity without an offsetting query benefit. The membership predicate (PR #22 §5) reads `hardware_class` from the observation node, not from the subject.

---

## 5. Cadence and budget

### 5.1 Default sampling cadences

Per data source, recommended v0.1 default cadences:

| Source | Default cadence | Rationale |
|---|---|---|
| `stress.log` events | **Event-driven** (no fixed cadence) | The adapter tails the log; emits an observation per log line. Cadence follows the rate of `bma.runtime.*` events. |
| `lm-sensors` thermal/voltage | **Every 5 s** | Thermal dynamics on a healthy host are seconds-scale; 5 s gives enough resolution for surveillance-scout pattern detection while keeping ingest rate sane (~12 obs/min × N sensors). |
| `lm-sensors` fan speed | **Every 5 s** | Fan-speed swings correlate with thermal events; same cadence as thermal lets Synthesis correlate cheaply. |
| `smartctl` SMART attributes | **Every 1 h** | SMART attributes change on hours-to-days timescale; per-minute polling is gratuitous and (for some SATA SSDs) thermally counterproductive. |
| `/proc/meminfo` | **Every 1 min** | Memory pressure dynamics are minute-scale on a healthy host; sub-minute pressure spikes are caught by `stress.log` `SE_VRAM` / OOM events instead. |
| GPU telemetry (steady-state) | **Every 5 s** | GPU thermals and VRAM utilisation move seconds-scale. |
| GPU telemetry (pressure events) | **Event-driven** | When VRAM eviction or thermal-throttle events fire, the adapter emits an observation immediately (sentinel-eligible per §6). |

### 5.2 Tunability contract

**All cadences in §5.1 are defaults, not commitments.** BMA-implementor at Walk-α tunes per-host (the FX-8350 with PCIe 2.0 has different I/O budget characteristics than the Walk-phase RISC-V QBP-CU hardware that joins the federation per `project_silicon_ladder.md`). Tunability is via adapter configuration; the spec commits only to the *shape* (per-source cadence, event-driven option, sentinel override) — not to the values.

### 5.3 Ingest-budget envelope

The adapter MUST stay within an ingest-rate envelope that the BMA host can sustain without consuming budget the autonomic layer needs. Per CLAUDE.md, BMA's autonomic layer (AUTO-S/P, CCB 10Hz) operates under a ≤200ms emergency-response budget — the adapter's own CPU/I/O footprint must not crowd that budget. v0.1 envelope guidance (recommendations; BMA-implementor-tunable):

- Adapter steady-state CPU: <1 % of one core (FX-8350 baseline).
- Adapter ingest rate: <100 obs/sec per host steady-state; sentinel bursts permitted up to ~1000 obs/sec for short windows (seconds).
- Adapter local memory: <128 MB resident.

Exceeding these is a sign the cadence defaults are too aggressive for the host; per §5.2 the implementor adjusts.

---

## 6. Sentinel-class observations

### 6.1 What a sentinel is

When an observation crosses a **hard threshold** — one that the BMA autonomic layer (AUTO-S/P) responds to in real time per the CLAUDE.md ≤200ms budget — the adapter MAY mark the emitted observation as a **sentinel-class observation**. Sentinel observations carry a higher-priority flag (the exact field name is BMA-implementor's call; this spec names the concept, not the encoding). Downstream consumers — surveillance-mode scouts, Context Builder, Synthesis — treat sentinels as anomaly candidates per Spec v1.3 §3.2 + Theory v1.4 §8.1.

### 6.2 Threshold categories (not values)

The spec names threshold *categories* that qualify an observation as sentinel-eligible. The numeric thresholds themselves are per-hardware-class tunables owned by BMA-implementor at Walk-α (out of scope per §0.2):

| Category | Example trigger | Hardware class |
|---|---|---|
| Thermal hard threshold | `cpu_temp > <per-host limit>` (CLAUDE.md AUTO-S/P contract) | `hardware.cpu`, `hardware.gpu` |
| Disk integrity threshold | `disk_smart_reallocated > 0` (any reallocation is significant) | `hardware.disk` |
| Disk-temperature threshold | `disk_smart_temperature > <per-host limit>` | `hardware.disk` |
| Memory pressure threshold | OOM-killer event, swap-pressure exceeds <per-host limit> | `hardware.memory` |
| VRAM pressure threshold | `gpu_vram_evict_count > 0` (any eviction is significant) | `hardware.gpu` |
| PCIe degradation | `pcie_link_state != "up"` for a previously-up link | `hardware.gpu`, `hardware.network` |
| BMA stress-bus fatal | `SE_FATAL` or `SE_HARDWARE_PROBE` with severity flag | (mirrors source) |

### 6.3 Sentinel ≠ autonomic response

**The sentinel mechanism is NOT a replacement for the autonomic layer.** Per the boundary clarity table in PR #22 §8 (carried forward in §7 below):

- Autonomic (AUTO-S/P, CCB 10Hz) responds to hard-threshold events within ≤200 ms by issuing throttle / governor-change / mitigation commands directly to the affected subsystem.
- The Contextus adapter, even when emitting a sentinel-class observation, operates on minutes-to-hours timescales. It does not throttle; it does not issue mitigation. It just *speeds the surveillance-side pattern detection* by signalling "this observation deserves immediate scout attention."

The two mechanisms run on the same input event but at different cadences and with disjoint output paths. The autonomic mitigates; the sentinel-class observation tags the event for forensic recall and cross-domain Synthesis.

### 6.4 Sentinel discipline

The spec recommends — but does not enforce — that sentinel-class observations remain a small fraction of total adapter output. If sentinels exceed (e.g.) 5 % of total emission, the cadence or threshold tuning is wrong; the implementor should re-examine the per-host threshold values.

---

## 7. Boundary clarity vs autonomic and adjacent systems

This carries forward and refines PR #22 §8. The adapter sits unambiguously on the OBSERVE side; it is read-only with respect to BMA runtime state.

| Layer | Role | Cadence | Write path | Read path |
|---|---|---|---|---|
| **Autonomic** (AUTO-S/P, CCB 10Hz) | Real-time mitigation — throttle, governor-change, eviction | ≤200 ms | BMA-internal control state + `stress.log` | Direct subsystem telemetry |
| **CTH ρ_net loop** | Algebraic-integrity scoring; Constitutional Audit interrupt | 1 cycle | CTH inventory | BMA M2 WDEvent stream (per `repo-bma-systema-issue-#107`) |
| **`ctx-adapter-system` (this spec)** | OBSERVE-only — read telemetry surfaces, emit `NT_OBSERVATION` to NATS | Per §5 cadences | NATS `ctx.ingest.system.>` only | Five host-local sources (§2) |
| **Contextus surveillance** (consumer of this adapter) | Long-window pattern detection; cross-domain correlation; `AnomalyStructural` minting via Synthesis | NATS-mediated; minutes-to-hours | Contextus hypergraph (through Synthesis persistence boundary) | Adapter NATS subjects + existing Contextus subgraph |
| **Notary** (Sprint 1 §2.h authorised) | Running-code verification — *"did the autonomic actually respond"* | Per claim; Trust-Tier-cadenced | Notary verification receipts | All of the above |

### 7.1 The adapter's three non-doings

The adapter SHALL NOT:

1. **Issue throttle / mitigation commands.** That is autonomic-layer authority. The adapter has no write path into BMA runtime state.
2. **Modify BMA runtime state.** No writes to `stress.log` (it is read-only). No writes to CCB. No writes to AUTO-S/P state.
3. **Replace the autonomic layer's read path.** The adapter is a *parallel* observer of the same sources; it is not in the autonomic response path. If the adapter is down, the autonomic layer is unaffected.

### 7.2 The adapter's one doing

The adapter SHALL: read the five sources at the §5 cadences, normalise into `NT_OBSERVATION` nodes per §3, publish to NATS per §4, mark sentinel-eligible observations per §6.

That's it. Single concern; clean boundary.

---

## 8. Concern separation

### 8.1 Contextus owns

- The `NT_OBSERVATION` node shape for operational telemetry (this spec §3; the underlying `NT_OBSERVATION` type is Spec v1.3 §11.4)
- The NATS subject hierarchy `ctx.ingest.system.<host_id>.<observation_kind>` (§4)
- Synthesis cross-domain promotion of operational-scope observations to `NT_INSIGHT_SIGNAL{anomaly_kind: AnomalyStructural}` per PR #22 §7
- The operational-scope membership predicate (`method = "hardware-identifier"`) per PR #22 §5
- The v0.1 hardware-class tag taxonomy and its v0.2+ expansion governance (PR #22 §4)
- Sentinel-class observation **mechanism shape** (this spec §6) — the *concept* and the *category list*

### 8.2 BMA owns

- `stress.log` emission cadence + format
- `SE_HARDWARE_PROBE` / `SE_VRAM` / `SE_FATAL` semantics
- AUTO-S/P autonomic thresholds and response logic
- CCB 10Hz negotiation loop
- Per-instance `host_id` assignment (and the federation-uniqueness contract for it)
- The adapter binary itself: Go implementation (lives in `internal/bma/ctxadapter/` per BMA Go layout), container packaging, lifecycle (start/stop/health-check/restart), configuration management, log-rotation
- Per-hardware-class threshold tuning (the numeric values §6 names categories for)
- The wire-format-level JSON encoding of `NT_OBSERVATION` messages on NATS (the §3 fields are the contract; the byte-level encoding is BMA-implementor's call)

### 8.3 Wyrd owns

- The substrate `bma.runtime.*` namespace constants (PR #16, merged) that the adapter's `observation_kind` field maps from on the BMA side
- The MuninnDB writer downstream of NATS that durably persists adapter-emitted observations

### 8.4 CTH owns

- Algebraic-integrity scoring on `bma.runtime.flag-norm-drift` and related ρ_net inputs (per `repo-bma-systema-issue-#107`)
- The algebraic-integrity precedent that PR #22 §7.1 worked pattern A correlates against

The contract: this spec extends only the source-adapter catalogue and the architectural contract for the system telemetry adapter. The downstream BMA implementation is not pre-committed by this doc; the §10 reader list ensures the consuming repos sign off on the boundary.

---

## 9. Sequencing — Crawl → Toddle → Walk → Run

### 9.1 Crawl (this sprint)

| Phase | Deliverable | AC reference |
|---|---|---|
| **2.1** | `Contextus-Spec-Addendum-NT-Scope-Operational.md` — third scope sibling + Synthesis pattern (PR #22) | AC-1, AC-2 (type shape), AC-7 |
| 2.2 | `ScopeOperational` Go type + JSON-Schema fragment | AC-2 (impl), AC-4 |
| 2.3 | Reference loader extension — `LoadResult.OperationalScopes` | AC-5 |
| **2.4 (this spec)** | `ctx-adapter-system` design — architectural contract; doc-only at v0.1 | AC-3 |
| 2.5 | Test plan (round-trip YAML→struct; envelope-split for `tier_immune`; cross-domain correlation via fake-data scout) | AC-8 |
| 2.6 | go-coding-guide.md verification suite clean | AC-9 |
| 2.7 (follow-on) | Spec v1.4 §2.4 Referent cross-reference — scalar-Referent surveillance scoring | AC-6 (gated on `repo-confluent-trust-impl` ack) |

### 9.2 Toddle

- BMA-implementor authors the adapter binary in Go.
- Implementation lives at `internal/bma/ctxadapter/` per CLAUDE.md BMA Go source layout convention.
- Adapter publishes on the `ctx.ingest.system.<host_id>.<observation_kind>` subject hierarchy.
- First end-to-end smoke test on BMA-prime: telemetry event → adapter normalisation → NATS publish → MuninnDB write → `HE_SCOPE_MEMBERSHIP` edge minted with `method = "hardware-identifier"`.
- Contextus side has nothing to do at Toddle for this adapter — the source-adapter catalogue entry is the only Contextus artefact needed; the work is BMA-implementor-side.

### 9.3 Walk

- Live BMA hardware telemetry flowing into Contextus on a sustained cadence.
- First real (not fake-data-scout) `AnomalyStructural` cross-domain Synthesis-minted signal fires — e.g., a real thermal-stress × algebraic-integrity-drift correlation per PR #22 §7.1 worked pattern A.
- Operational-scope query patterns exercised in BMA Conscious-A/B forensic-recall sessions.
- Sentinel-class observation discipline (§6.4) validated against real ingest-rate data.

### 9.4 Run

- Federation tenants inherit the same adapter pattern with per-tenant variants:
  - Sharp Butler House Node hardware adapter (`ctx-adapter-system` for residential-concierge host class; same shape, possibly extended tag set for `hardware.actuator` / `hardware.sensor-array` as the v0.2 expansion criterion catches up)
  - QBP-CU silicon adapter (`ctx-adapter-system` for compute-unit silicon; per `project_silicon_ladder.md` Walk-phase RISC-V hardware joins federation at Rung 3, inheriting this adapter for free)
  - Möbius Fusion energy-system adapter (`ctx-adapter-system` for fusion-reactor host class; tag-set expansion likely needed for `hardware.power-delivery` / `hardware.containment`)
- Cross-tenant operational-scope correlations become a federation-coherence surface (gated on PR #22 §12.2 cross-tenant visibility decision).

---

## 10. Open questions and §I4 reader list

### 10.1 Open question Q1 — graceful handling of missing data sources

**Q:** Should the adapter handle missing data sources gracefully — e.g., no GPU on a non-GPU host; no SMART on a virtual-disk-only host; no `lm-sensors` on a stripped-down container host?

**v0.1 lean: YES — log-and-continue.** A missing data source is not itself an anomaly. The adapter logs the absence at startup (one INFO-level message per missing source), proceeds to emit observations from the sources that *are* available, and does not emit sentinel observations for the missing surfaces. Sentinel observations are reserved for hard-threshold crosses on present sources (§6).

**Rationale:** federation tenants vary in hardware shape (Sharp Butler House Node has actuators that BMA-prime doesn't; QBP-CU silicon has FPGAs that neither has). A blanket "missing source = error" policy would force every tenant deployment to mask irrelevant sources — federation-friction without offsetting benefit.

**Closeable on:** @bma-implementor ack.

### 10.2 Open question Q2 — heartbeat observation for liveness

**Q:** Should the adapter publish a heartbeat `NT_OBSERVATION` even when no measurement changed (e.g., once per minute), so downstream consumers can detect adapter death via missed heartbeats?

**v0.1 lean: YES — one-per-minute heartbeat.** The heartbeat carries `observation_kind: adapter_heartbeat`, `value: 1`, `unit: "count"`, and the current timestamp. Downstream surveillance can use missed heartbeats (e.g., no `adapter_heartbeat` observation for `<host_id>` in the last 5 minutes) as a liveness anomaly. The heartbeat itself is NOT sentinel-class — its purpose is to be regular, not to flag.

**Rationale:** without heartbeats, a stuck or crashed adapter is silently invisible to the surveillance layer. Cross-host correlations would silently shrink. A 1-per-minute heartbeat is a ~1-obs/host steady-state baseline that consumers can rely on. Cost is negligible (1 obs/min/host).

**Caveat:** the heartbeat observation does NOT carry a `hardware_class` — it is adapter-level, not subsystem-level. Operational-scope membership for the heartbeat is the host-level umbrella scope (per PR #22 §6.1 example, `scope-host-bma-prime` with empty `HardwareClass`).

**Closeable on:** @bma-implementor ack on the once-per-minute cadence; @qbp-architecture ack on the federation-wide heartbeat convention.

### 10.3 §I4 named reviewers

Per the Contextus standing §I4 reviewer convention (mirroring PR #22 §12.3 reader list — the addenda are paired):

- **@bma-implementor** — primary consumer + implementation owner. The adapter binary lives in BMA-systema repo (`internal/bma/ctxadapter/`); the read-paths (`stress.log`, `lm-sensors`, `smartctl`, `/proc/meminfo`, GPU telemetry) are all BMA-host concerns. Acks required on §§2, 3, 5, 6, 8.2. Heartbeat cadence (Q2) and graceful-missing-source policy (Q1) both close on @bma-implementor.

- **@qbp-architecture** — federation coherence. The adapter pattern applies to future tenants (Sharp Butler House Node, QBP-CU silicon, Möbius Fusion energy systems) per §9.4 sequencing. First-tenant pattern witness — the adapter as specified for BMA must generalise without surface modification when the second-tenant (QBP-CU at Walk Rung 3 per `project_silicon_ladder.md`) inherits it. Acks required on §§4 (NATS subject hierarchy is federation-wide), 7 (boundary clarity carries forward to all tenants), 9.4 (federation-tenant sequencing), 10.2 (heartbeat convention).

- **@wyrd-implementor** — `bma.runtime.*` namespace constants integration (Wyrd PR #16, merged). The adapter's `observation_kind` field maps from the substrate namespace; substrate-tier compatibility is wyrd-implementor's authority. Also: MuninnDB writer downstream of NATS is the durability boundary. Acks required on §§3 (provenance_tag convention; substrate compatibility), 4.2 (Wyrd namespace cross-reference), 8.3.

Standing §2.i 4h SLA window applies per CLAUDE.md `feedback_named_reviewer_responsiveness`. Sessionbridge channel-of-record: `sprint-2-2026-05-20` (the channel PR #22 was confirmed on; this adapter is its paired addendum).

---

## 11. References

### 11.1 Contextus

- `~/Documents/Contextus/Contextus-Spec-v1.3.md` §3.1 (Source Adapters catalogue — this spec adds one entry); §3.2 (Ingestion Flow — adapter → NATS → MuninnDB Writer → Context Builder); §4.6 (Scope Nodes — operational scope sibling per PR #22); §11.4 (Go type catalogue — `NT_OBSERVATION` type)
- `~/Documents/Contextus/Contextus-Theory-v1.5.md` §3.6.2 (AnomalyStructural — Synthesis promotion target for operational↔algebraic↔cognitive correlations); §3.6.6 (Synthesis as Persistence Boundary)
- `~/Documents/Contextus/Contextus-Spec-Addendum-NT-Scope-Operational.md` (PR #22) — third scope sibling spec; this adapter is its input side. Especially §3 (ScopeOperational type), §4 (hardware-class tag taxonomy), §5 (membership predicate `method = "hardware-identifier"`), §6 (YAML integration), §7 (cross-domain Synthesis pattern), §8 (boundary clarity table).
- `~/Documents/Contextus/doc/cross-domain-hyperedge-minting.md` — design-only-at-v0.1 Contextus design-doc precedent (Wyrd PR #19 dependency; same shape this doc follows).
- `~/Documents/Contextus/doc/spec-v1.4-theory-as-conceptual-scope.md` — federation-additive precedent (PR #11, merged 2026-05-14).

### 11.2 Issues + PRs

- `repo-contextus-issue-#15` — tracking issue; AC-3 covered by this spec; AC-1/AC-2/AC-7 by PR #22; AC-4/AC-5/AC-8/AC-9 by Phases 2.2/2.3/2.5/2.6.
- `repo-contextus-pr-#22` — `NT_SCOPE_OPERATIONAL` addendum (the scope-side spec dependency for this adapter).
- `repo-contextus-pr-#17` — Research-Aid-Tenancy addendum (federation-additive precedent).
- `repo-contextus-pr-#11` — Spec v1.4 design surface (theory-as-conceptual-scope; merged; federation-additive precedent).
- `repo-wyrd-pr-#16` — `bma.runtime.*` namespace constants (merged; substrate namespace this adapter consumes).
- `repo-wyrd-pr-#40` — scope-node configuration loader v0.1 (merged; envelope-field plumbing the operational-scope side uses).
- `repo-wyrd-pr-#54` — scout-daemon (merged; runtime that dispatches the surveillance scouts consuming this adapter's NATS output).
- `repo-bma-systema-issue-#107` — M2 WDEvent → CTH ρ_net feedback loop (algebraic-integrity precedent that the cross-domain Synthesis pattern in PR #22 §7.1 correlates against).

### 11.3 BMA / federation context

- `~/Documents/CLAUDE.md` — BMA section (autonomic AUTO-S/P; CCB 10Hz; `stress.log`; `SE_HARDWARE_PROBE`, `SE_VRAM`, `SE_FATAL`; `bma.runtime.*` namespace; container limits; Go source layout convention for `internal/bma/ctxadapter/`)
- `~/Documents/BMA/theory/hypergraph-inference/BMA-Theory-Addendum-18_0-Hypergraph-Access-Pattern.md` §4 — Seam detection (algebraic-integrity precedent; `bma.runtime.flag-norm-drift` canonical event surfaced via this adapter as `observation_kind: stress_event`).
- `~/Documents/inter/workspace-phase-architecture.md` §0.13.1 + §0.13.2 — QBP-CU silicon de-risking ladder; Rung 3 RISC-V hardware joins federation at Walk, inheriting `ctx-adapter-system` for free per §9.4.
- `~/Documents/.claude/projects/-home-prime-Documents/memory/feedback_workspace_stack.md` — federation-portability principle (future tenants inherit this adapter pattern).
- `/home/prime/.claude/plans/reactive-snacking-lampson.md` — Sprint 2 plan-of-record; Phase 2.4 deliverable scope.
- sessionbridge `sprint-2-2026-05-20` channel — beekeeper confirmation of Phase 2 scope (paired with PR #22).

---

*End of design spec v0.1 — `ctx-adapter-system`. Design-only; architectural contract level; federation-additive. Implementation deferred to BMA-implementor at Walk-α. §I4 review pending on @bma-implementor + @qbp-architecture + @wyrd-implementor named-reviewer acks.*
