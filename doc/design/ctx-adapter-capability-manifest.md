# `ctx-adapter-system` capability-manifest — Design Spec

**Status:** Design only at v0.1. Architectural contract only — no implementation. Walk-α deliverable; impl is BMA-implementor-owned (adapter binary owner), exactly as `ctx-adapter-system.md` scopes the adapter itself.

**Version:** 0.1 · **Author:** contextus-impl · **Date:** 2026-10-17

**Resolves:** contextus#28 (housekeeping / Walk-α) — @qbp-architecture nudge, `sprint-2-2026-05-20` seq=26, PR #23 Q1 review.

**Extends:** `doc/design/ctx-adapter-system.md` §10 Q1 (graceful missing-source handling). This is the Walk-α upgrade of that spec's Crawl-phase *log-and-continue* lean.

---

## 1. The gap (from #28)

`ctx-adapter-system` (PR #23) uses **log-and-continue** for a missing data source — correct at Crawl, but at Walk, when cross-host correlation runs, it conflates two meaningfully different states:

1. **Absent by design** — the host doesn't have that hardware/sensor (a non-GPU host has no GPU). Expected; not an error.
2. **Failed silently** — the adapter *tried* to read a source it should have and the read failed. Unexpected; signals degraded telemetry.

Log-and-continue makes both look identical from the outside: **no `gpu_*` observations arrive**. A consumer correlating across hosts cannot tell "this host has no GPU" from "this host's GPU telemetry is broken" — so it can neither trust the silence nor flag the failure. Capability must be **declared**, not **inferred from signal presence/absence**.

## 2. The manifest node

A per-host **capability manifest**, emitted by the adapter, declaring every source it knows about and that source's current state. Consumers read the manifest for host capability; they never infer it from signal absence.

**Node:** a Wyrd `model.Node`, `Type = "contextus.adapter.capability_manifest"` (`NT_ADAPTER_CAPABILITY_MANIFEST`) — Contextus-defined per the `wyrd:Node` consumer convention (consumer owns its `NodeType`; Wyrd reserves only the `wyrd.` prefix). One current manifest per host; superseded on each new emission (prior manifests decay per retention tier).

**Payload fields:**

| Field | Meaning |
|---|---|
| `host_id` | Federation-unique host id; matches `ScopeOperational.HostID`. |
| `host_class` | The host's declared class (e.g. `runtime.bma-instance`, future `sharp-butler-house-node`, `qbp-cu-silicon`) — the key to *which sources are expected*. |
| `emitted_at` | RFC 3339 UTC. |
| `reason` | `startup` \| `capability_change`. |
| `sources[]` | One entry per known source (§3). |

**`sources[]` entry:**

| Field | Meaning |
|---|---|
| `source` | The `ctx-adapter-system` source: `stress_log` \| `lm_sensors` \| `smartctl` \| `proc_meminfo` \| `gpu`. |
| `hardware_class` | The v0.1 tag it feeds (`hardware.cpu/disk/gpu/memory/network`), or empty for host-level (`stress_log`). |
| `expected` | `bool` — is this source expected present on `host_class`? (The disambiguator.) |
| `state` | `present` \| `absent_by_design` \| `failed` (§3). |
| `detail` | Optional free-form (e.g. `"rocm-smi: no device"`); omitted when `present`. |

## 3. State semantics — the disambiguation

| `expected` | source emitting? | `state` | Meaning |
|---|---|---|---|
| true | yes | `present` | Source read OK; observations flowing. |
| false | no | `absent_by_design` | Host class has no such source (non-GPU host). **Not an error.** |
| true | no | `failed` | Source *should* be there and the read failed. **Degraded telemetry — actionable.** |

`absent_by_design` vs `failed` is exactly `expected`: the field #28 asks for. `absent_by_design` is never sentinel-eligible; `failed` **is** a surveillance signal (a `host_class`-expected source going dark is itself an anomaly a consumer should raise).

## 4. Emission

- **At adapter startup:** the full manifest (every known source with its initial state). This replaces the Crawl "one INFO log line per missing source" (`ctx-adapter-system` §10 Q1) with a machine-readable declaration.
- **On capability change:** a fresh manifest (`reason: capability_change`) when any source transitions state — e.g. GPU `present → failed` (driver crash) or `absent_by_design → present` (hot-plug). Not per-observation; only on a state edge.
- Scoped to the host's **host-level umbrella `ScopeOperational`** (`HardwareClass=""`) via `HE_SCOPE_MEMBERSHIP`, so a consumer resolving a host's operational scope also finds its manifest.

Distinct from the **heartbeat** (`ctx-adapter-system` §10 Q2): heartbeat = *liveness* ("adapter alive"); manifest = *capability* ("which sources, in which state"). Both are adapter-level, host-umbrella-scoped; they compose (a missed heartbeat = adapter down → the whole manifest is stale; a `failed` entry = adapter up, one source down).

## 5. Cross-host correlation (the consumer contract)

Surveillance scouts / Synthesis / Context Builder **MUST read the manifest, not infer from signal presence/absence**:

- "No `gpu_temp` from `host-X` for 5 min" → look up `host-X` manifest → `gpu.state == absent_by_design` ⇒ *correct silence, no anomaly*; `gpu.state == failed` ⇒ *degraded-telemetry anomaly, raise it*; `gpu.state == present` but no observations ⇒ *a genuine gap worth a sentinel* (adapter thinks it's reading but nothing lands).
- A cross-host thermal-pattern scout (`ctx.ingest.system.*.cpu_temp` subscriber) weights only hosts whose manifest says `lm_sensors: present` — it does not treat an `absent_by_design` host as "0 observations = cool."

This closes the federation-level false read where silence was ambiguous. The `ctx-adapter-system.md` §5 cross-host correlation text gains a one-line pointer: *capability is read from the manifest, never inferred from signal presence.*

## 6. Concern separation (carries `ctx-adapter-system` §8)

- **Contextus owns:** the `NT_ADAPTER_CAPABILITY_MANIFEST` node shape + state semantics (this spec); the consumer contract (§5).
- **BMA owns:** emitting it from the adapter binary at Walk-α (startup + capability-change hooks); the per-`host_class` `expected`-source table (which sources each class should have — a tuning table, like the §6 threshold values).
- **Wyrd owns:** persistence of the manifest node.

## 7. Acceptance criteria (#28)

- [x] **Manifest schema designed** — §2 (node + payload + `sources[]`) and §3 (states).
- [x] **Adapter emits a manifest node at startup** — §4 (`NT_ADAPTER_CAPABILITY_MANIFEST`, startup + capability-change); *implementation* Walk-α / BMA-implementor, per the issue's own "no Crawl action required" scope.
- [x] **Cross-host correlation references the manifest, not signal-inference** — §5, plus the pointer to fold into `ctx-adapter-system.md` §5.
- [ ] **§I4 reader-list:** @wyrd-implementor (node type / substrate) · @qbp-architecture (federation coherence). ← review gate, on the PR.

## 8. References

- `doc/design/ctx-adapter-system.md` — §3 (NT_OBSERVATION shape), §5 (cadence), §6 (sentinels), §8 (concern separation), §10 Q1/Q2 (missing-source + heartbeat). This spec is its Walk-α capability-declaration layer.
- `Contextus-Spec-Addendum-NT-Scope-Operational.md` — `ScopeOperational` (HostID / HardwareClass), host-level umbrella scope (`HardwareClass=""`).
- `pkg/types/scope.go` — `ScopeOperational`, `HardwareClass*` tags.
- contextus#28 — this issue.

*Design-only v0.1. Impl deferred to BMA-implementor at Walk-α. §I4 pending on @wyrd-implementor + @qbp-architecture.*
