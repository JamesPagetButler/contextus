# Contextus ↔ Wyrd Integration Architecture

**Date:** 2026-05-05
**Status:** Architecture decisions; ready for spec absorption (Spec v1.3) and Wyrd issue #6 correction
**Author:** Claude Opus 4.7 (architecture instance, `~/Documents/QBP-Compute-Unit/`)
**Audience:** Contextus implementer instance, Wyrd implementer instance, BMA Implementor (informed)
**Triggered by:** Contextus implementer review of Wyrd [issue #6](https://github.com/JamesPagetButler/wyrd/issues/6) and the two open questions in `wyrd/doc/integration/contextus.md`
**Companion docs:**
- [`Contextus-Spec-v1.2.md`](Contextus-Spec-v1.2.md) §2.1, §2.3, §4.2, §4.5, §9, §11.1, §11.3 — authoritative source
- [`Contextus-Theory-v1.4.md`](Contextus-Theory-v1.4.md) §3.6 — Insight Signal definition
- `wyrd/doc/integration/contextus.md` — Wyrd-side integration sketch (needs correction per §3 below)
- Wyrd issue [#6](https://github.com/JamesPagetButler/wyrd/issues/6) — Walk-phase glue (needs correction per §3 below)

---

## 1. Executive Summary

The Contextus implementer flagged four issues. All four are resolvable. Three are spec-clarification (Spec v1.2 already has the right answer; the Wyrd issue and integration doc lag); one is genuinely new architectural surface (scope nodes for "focus areas").

**Net result:** small set of concrete changes that tighten the integration without expanding scope. The bigger philosophical question — Contextus as universe-scale index on a 250 GB SATA SSD — has a clean architectural answer that the existing Spec v1.2 §9 already partially supports.

---

## 2. Resolutions to the Four Open Items

### 2.1 Source-enum inconsistency — Spec v1.2 wins

**Verdict:** Wyrd issue #6 is wrong. Spec v1.2 §11.1 is right.

Spec v1.2 §11.1 already defines:

```go
type AgentClass string

const (
    AgentScout         AgentClass = "scout"
    AgentCorrelation   AgentClass = "correlation"
    AgentSynthesis     AgentClass = "synthesis"
)
```

And §2.3 explicitly types `agent_type AgentClass` on `NT_INSIGHT_SIGNAL`. The three values are **scout, correlation, synthesis** — the *global* agents.

§4.2 is unambiguous about who emits Insight Signals and who does not:

| Agent | Emits InsightSignal? | What instead |
|---|---|---|
| Scout (global) | ✅ `anomaly_type = density` | — |
| Correlation (global) | ✅ `anomaly_type = correlation` | — |
| Synthesis (global) | ✅ `anomaly_type = narrative` | — |
| Provenance (global) | ❌ | Validates writes |
| Context Builder (global) | ❌ | Enriches NT_OBSERVATION |
| **Edge Scout (session)** | ❌ | NATS to `ctx.edge.boundary.{session_id}` |
| **Corpus Edge Scout (per-search)** | ❌ | NATS to `ctx.corpus.diversity` |
| **Bridge Agent (session)** | ❌ | NATS to `ctx.bridge.intervention` |

§4.5 closes it: *"Session state is not persisted (boundary flags are session-relative, not hypergraph state)."*

**Fix to Wyrd issue #6 / `contextus.md`:** the `SignalSource` enum should mirror `AgentClass`:

```go
type SignalSource string
const (
    SourceScout       SignalSource = "scout"
    SourceCorrelation SignalSource = "correlation"
    SourceSynthesis   SignalSource = "synthesis"
)
```

The Edge / Corpus / Bridge agents do **not** appear here. They are not Wyrd-resident state.

### 2.2 Open Question (a): single shared graph vs. per-agent graphs

**Decision: single shared `wyrd.model.Graph` with `Node.Type` distinguishing source.**

Spec §4.3: *"agents communicate via NATS subjects and share state through MuninnDB."* MuninnDB is the single graph; `Node.Type` does the source discrimination.

Concrete `Node.Type` taxonomy (extends Spec §2.1):

| `Node.Type` | Backed by | Created by |
|---|---|---|
| `contextus.signal.scout` | `NT_INSIGHT_SIGNAL` with `agent_type = "scout"` | Scout agent |
| `contextus.signal.correlation` | `NT_INSIGHT_SIGNAL` with `agent_type = "correlation"` | Correlation agent |
| `contextus.signal.synthesis` | `NT_INSIGHT_SIGNAL` with `agent_type = "synthesis"` | Synthesis agent |
| `contextus.scope.physical` | NEW — see §2.5 | Beekeeper / config |
| `contextus.scope.conceptual` | NEW — see §2.5 | Beekeeper / config / Synthesis agent |

The `Node.Type` namespace `contextus.*` keeps Contextus state distinguishable from QBP-CU, BMA, and other co-tenants in the same Wyrd graph. This matches the pattern Wyrd already uses for `cth.*` and the BMA-side `cth_id` convention from BMA's 2026-05-03 reply.

### 2.3 Open Question (b): hypothesis-tier candidates ("fish-on-the-line") in Wyrd?

**Decision: no. Session-scoped agent output stays ephemeral. Synthesis is the persistence boundary.**

Two reasons:

1. **The spec already says so.** §4.5: *"Session state is not persisted."* Promoting EdgeScoutFlag / BridgeIntervention / CorpusDiversityReport to Wyrd would contradict the spec.
2. **Storage budget argues for it.** Spec §9.3 Storage Sentinel triggers advisories at 70% capacity (the Samsung 840 is already at this threshold). Inflating Wyrd with every fish-on-the-line candidate would accelerate the trip toward 80% / 85% / 90% thresholds.

**The promotion mechanism the spec needs to formalize** (currently implicit in §4.4):

> When the Bridge Agent's converged intervention exceeds confidence threshold, a **Synthesis agent** subscribed to `ctx.bridge.intervention` validates the convergence and mints an `NT_INSIGHT_SIGNAL` of `agent_type = "synthesis"` with `anomaly_type = "narrative"` and a description that names the bridge convergence. That signal is what enters Wyrd.

Same pattern for Edge Scout and Corpus Edge Scout outputs that warrant persistence: a Synthesis agent listens, validates, and mints. The Bridge Agent / Edge Scout / Corpus Edge Scout themselves never write Wyrd state.

**Spec change (small):** add to §4.4 step list:

> *"Synthesis agents may subscribe to `ctx.bridge.intervention`, `ctx.edge.boundary.*`, and `ctx.corpus.diversity` and mint Insight Signals when ephemeral session findings warrant persistence. The Synthesis agent is the sole persistence boundary for session-scoped agent outputs."*

### 2.4 Type mapping (Spec §11 → Wyrd)

| Spec §11 type | Wyrd disposition | Lifetime |
|---|---|---|
| `InsightSignal` | `model.Node` of `Type = "contextus.signal.<agent>"` | Persistent (subject to §9 retention tiering) |
| `EdgeScoutFlag` | **ephemeral** — NATS only | Session-scoped; ranked + re-published when top-N changes |
| `BridgeIntervention` | **ephemeral** — NATS only | Session-scoped; ranked + re-published |
| `CorpusDiversityReport` | **ephemeral** — NATS only | Per-search; ranked + re-published |
| Cross-domain match candidate (N≥3 Insight Signals) | `model.Hyperedge` (arity ≥ 3) | Persistent; Theorem 2 irreducibility applies |
| **Physical scope** (e.g., Squam Lake watershed) | `model.Node` of `Type = "contextus.scope.physical"` | Persistent (rarely written) |
| **Conceptual scope** (e.g., physics) | `model.Node` of `Type = "contextus.scope.conceptual"` | Persistent (rarely written) |

**Bridge Agent's converged intervention** that warrants persistence becomes a Synthesis-authored Insight Signal per §2.3 above, which then can be the source of an N-arity Hyperedge if 3+ domains agree.

### 2.5 New addition: Scope Nodes

The user's framing concern — *"Contextus could hold a concept area like physics and a physical area like earth"* — is unaddressed by the current spec. Locale (Theory §2.2) handles spatiotemporal addressing of *individual observations*, but there's no first-class node for "the focus area my session is constrained to."

This becomes a new node-type pair in §2.1:

#### `NT_SCOPE_PHYSICAL`

A region of space-time with a name and queryable membership.

| Field | Type | Description |
|---|---|---|
| `scope_id` | string | e.g., `scope-squam-lake-watershed` |
| `geometry` | GeoJSON | Polygon / multipolygon |
| `elevation_range_m` | [float64, float64] | Optional; `null` if not relevant |
| `temporal_range` | LocaleBounds | When this scope is "active"; may be open-ended |
| `grain` | GrainSpec | Finest resolution within this scope |
| `name` | string | Human-readable name |
| `parent_scope_id` | string | Hierarchical containment (nullable) |
| `tags` | []string | e.g., `["watershed", "freshwater", "north-america"]` |

#### `NT_SCOPE_CONCEPTUAL`

A topic / domain with name and queryable membership.

| Field | Type | Description |
|---|---|---|
| `scope_id` | string | e.g., `scope-physics-fluid-dynamics` |
| `ontology_uri` | string | Optional pointer to an external ontology entry |
| `name` | string | Human-readable name |
| `parent_scope_id` | string | Hierarchical containment (nullable) |
| `related_scope_ids` | []string | Cross-references for traversal |
| `tags` | []string | e.g., `["physics", "fluid-dynamics", "transport-phenomena"]` |

#### `HE_SCOPE_MEMBERSHIP`

| Field | Type | Description |
|---|---|---|
| `scope_id` | string | NT_SCOPE_PHYSICAL or NT_SCOPE_CONCEPTUAL |
| `member_id` | string | Any node in the graph |
| `since` | time.Time | When membership was established |
| `confidence` | float64 | Inferred (e.g., spatial intersection) or asserted (beekeeper config) |
| `provenance_tag` | string | `T`, `P`, `D`, or `I` |

#### Why both kinds use the same shape

A **concept scope** and a **physical scope** are both *queryable subgraph definitions*. They differ in the kind of membership predicate they apply (geometric intersection vs. topical relevance) but the architectural role is identical: define a focus area; let queries traverse from it; let signals link to it. Treating them as sibling node types — with hierarchical parent edges and an unrestricted cross-product (any signal can have N physical scopes AND M conceptual scopes) — keeps the schema small and the query semantics uniform.

A query like *"what's known about fluid dynamics in the Squam Lake watershed in spring 2024?"* becomes a Locale-bounded traversal from the intersection of two scope nodes, not a custom search engine.

### 2.6 New addition: Evidence Pointer Discipline

Spec §9 already has the storage-tier ladder (Core / Near / Peripheral / Distant / Skeleton, sized 50 bytes to 2 MB). The Distant and Skeleton tiers are essentially **pointer-only** entries — the spec already discovered the pattern, just hasn't named it.

This pattern should be **explicitly applied to all evidence**, not just paper corpora. Add a §5 (or §2.4) "Evidence Pointer Discipline":

```go
// EvidencePointer is the load-bearing primitive that lets Contextus
// span evidence the local disk cannot store. Every fact carries WHERE
// it came from, not WHAT it said.
type EvidencePointer struct {
    Locator     string    `json:"locator"`     // URL, file path, archive coordinate, NATS subject
    LocatorKind string    `json:"locator_kind"` // "https" / "file" / "archive" / "nats" / "telescope_obs_id" / "doi"
    Hash        []byte    `json:"hash,omitempty"` // SHA-256 of dereferenced content (verification)
    SizeBytes   int64     `json:"size_bytes,omitempty"`
    LoadedAt    time.Time `json:"loaded_at,omitempty"`
    AccessHint  string    `json:"access_hint,omitempty"` // "cold" / "warm" / "hot"
    Note        string    `json:"note,omitempty"`         // human-readable provenance
}
```

Add this to `NT_INSIGHT_SIGNAL` as a `[]EvidencePointer` field.

The discipline:
- **Default:** an Insight Signal holds metadata + EvidencePointers, not raw evidence.
- **Promotion:** raw evidence is dereferenced lazily on demand (deep-dive moves a paper from Peripheral → Core per §9.2).
- **Demotion:** Ebbinghaus decay can move a Core paper back to Near and eventually drop the full text, keeping the EvidencePointer for re-fetch on demand.

This generalizes §9's paper-corpus story to apply to satellite imagery, USGS gauge timeseries, magnetar X-ray observations, snow depth measurements, eDNA metabarcoding output, and the rest of the universe-scale catalogue.

---

## 3. Concrete Changes Required

### 3.1 Spec v1.2 → v1.3

Five small additions/clarifications. Total: ~50 lines of new text.

| § | Change | Driven by |
|---|---|---|
| §2.1 | Add `NT_SCOPE_PHYSICAL` and `NT_SCOPE_CONCEPTUAL` to node-type catalogue | §2.5 above |
| §2.2 | Add `HE_SCOPE_MEMBERSHIP` to edge-type catalogue | §2.5 above |
| §2.3 | Add `evidence []EvidencePointer` to `NT_INSIGHT_SIGNAL` field table | §2.6 above |
| §4.4 | Add Synthesis-as-persistence-boundary clause at end of step list | §2.3 above |
| §11.1 | Add `EvidencePointer` Go type | §2.6 above |

No changes required to §9 — the existing retention-tier ladder absorbs the EvidencePointer discipline naturally.

### 3.2 Wyrd issue #6 corrections

The integration sketch at `wyrd/doc/integration/contextus.md` and the issue body both need:

| Item | Current | Should be |
|---|---|---|
| `SignalSource` enum | `edge | corpus | bridge` | `scout | correlation | synthesis` |
| "Mapping Contextus types to Wyrd types" table — row for "Edge Scout / Corpus Scout / Bridge Agent" | "not stored — runtime services that produce edges" | Replace with: "**Session-scoped agents** — outputs are NATS-only, never persisted to Wyrd. Bridge intervention promotion happens via Synthesis." |
| Type sketch | `FromSignal(InsightSignal) model.Node` | Same shape, but with a note that `evidence []EvidencePointer` flows through to `Node.Properties` |
| Open Questions section | Two unanswered questions | Both answered: single shared graph; ephemeral stays ephemeral. Reference this doc. |
| Out-of-scope item "fish-on-the-line filtering" | "lives in Contextus repo" | Same, but note explicitly: filtering happens at the Synthesis agent that promotes session-scoped findings, before anything reaches Wyrd. |

Add to the type-mapping table:

| Contextus concept | Wyrd type | Notes |
|---|---|---|
| Physical scope | `model.Node` of `Type = "contextus.scope.physical"` | Geometric scoping for queries |
| Conceptual scope | `model.Node` of `Type = "contextus.scope.conceptual"` | Topical scoping for queries |
| Scope membership | `model.Hyperedge` of `Type = "contextus.scope.member"` | Connects scope to signals/observations within it |

### 3.3 BMA-side implications

None immediate. BMA imports Wyrd; whatever Wyrd exposes becomes available. The scope-node addition gives BMA's `internal/bma/cth/projection.go` (planned) a way to constrain CTH inventories to a focus area without hand-rolling a query layer.

---

## 4. The Bigger Picture: Contextus as Index, Not Corpus

The user's framing concern: *"Contextus could be taking into account all of the universe down to the current state of the earth weather patterns and locations of animals humans houses and so on. But in reality the SATA drive is already 70% and old in terms of writes."*

This tension is real and the architecture has a clean answer.

### 4.1 The principle

> **Contextus is the *index* of evidence, not the evidence itself.**

The Wyrd graph holds what we *know*: relationships, anomalies, hypotheses, scopes. The actual data — satellite imagery, gauge timeseries, journal articles, eDNA reads — lives wherever it was sourced. Contextus knows where it lives and can fetch it on demand. Most queries don't need the raw data; they need the metadata-and-relationships layer.

This is the same principle web search engines use: Google does not store a copy of every page on the web; it stores an index that points back at every page. The index fits on Google's servers; the web does not.

### 4.2 How the existing spec already supports this

§9 Retention Tiers is the load-bearing mechanism. Translated for any evidence kind (not just papers):

| Tier | What lives in Wyrd | What's deferred to the locator |
|---|---|---|
| **Core** | Full enriched node + recent dereferenced data | Original archival source |
| **Near** | Node + indexed fingerprint + EvidencePointer | Full content |
| **Peripheral** | Node + EvidencePointer + small embedding | Everything else |
| **Distant** | Node + EvidencePointer | Everything else |
| **Skeleton** | Node + EvidencePointer (ID + locator only) | Everything else |

The architecture is already sized for "Distant + Skeleton dominate; Core is rare and valuable."

### 4.3 How scope nodes complete the picture

Adding `NT_SCOPE_PHYSICAL` and `NT_SCOPE_CONCEPTUAL` (§2.5) gives the system **focus**:

- *"Activate the Squam Lake watershed scope"* → all signals not linked to that scope get demoted; signals within the scope get promoted to Near or Core.
- *"Activate the magnetar spectroscopy scope"* → promote QBP-domain signals about specific Galactic magnetars; demote everything else.
- *"Deactivate the watershed scope after the seasonal report ships"* → tier transitions reverse; storage capacity recovers.

Scope activation becomes the operational lever that controls what fits on the Samsung 840. **Contextus is universe-shaped in principle, focus-shaped in practice.**

### 4.4 Why concept and physical scopes use the same architecture

A Squam Lake watershed scope and a fluid dynamics scope are both *focus areas* that:

- Constrain queries
- Trigger §9 tier transitions
- Anchor session-scoped agent attention
- Expose a hierarchical containment structure
- Compose with each other (you can be in *both* the watershed and fluid-dynamics scopes simultaneously)

If they were different node types architecturally, you'd need two query engines, two retention adapters, two membership models. Sibling node types with shared shape collapse this into one mechanism.

This is also where the `EvidencePointer` discipline (§2.6) compounds: a signal in the {watershed × fluid-dynamics} intersection might point to a NetCDF dataset (NOAA), a USGS gauge timeseries (HTTPS), a satellite tile (S3), and a journal article (DOI) — four locators, all for the same signal, all dereferenced lazily. The signal node stays small; the universe stays out of disk.

### 4.5 What this rules out

- **Storing every observation indefinitely.** Skeleton tier is the floor; below that, observations are dropped (never deleted from the *original* source).
- **Treating Wyrd as a CDN for raw data.** Wyrd is the index; CDN-style caching of raw data is a separate concern outside Contextus's scope.
- **Building physical and conceptual scopes as parallel-but-different systems.** They share the Node + Hyperedge + Locale architecture and use it differently, not separately.

### 4.6 What this enables

- **A fresh BMA instance booted in 2027 can subscribe to the Squam Lake watershed scope** and immediately query relevant signals without re-ingesting the whole observation history.
- **A QBP physics workflow can activate the magnetar-spectroscopy concept scope** and pull in adjacent measurement-agent flags from telescopes worldwide via EvidencePointers.
- **Sharp Butler at Holderness focuses on the watershed scope by default** — its relevant signals are always Core-tier; everything else hovers at Peripheral or below.
- **The Storage Sentinel's 70%/80%/85%/90% thresholds** become natural choke points where the system tightens scope-membership criteria automatically rather than failing.

---

## 5. Open Items Tracked Elsewhere

- Spec v1.3 drafting cadence — Contextus implementer's decision when to absorb these changes (§3.1)
- BMA-side `internal/bma/cth/projection.go` may eventually want a "scope filter" parameter — track when M1 work begins
- The Synthesis-as-persistence-boundary mechanism (§2.3) requires the Synthesis agent to subscribe to three new NATS subjects; small implementation work for the Contextus implementer

---

## 6. Concrete Next Moves

For the Contextus implementer:

1. **Spec v1.3 absorption** of the five changes in §3.1. Estimated: ~1 hour writing + cross-reference checking.
2. **Synthesis agent subscriber wiring** (§2.3) — adds three NATS subscriptions, one for each session-scoped agent's output. Estimated: ~1 day.
3. **Scope node implementation** — schema definitions in MuninnDB, basic CRUD, scope-bounded query helper. Estimated: ~3 days.
4. **EvidencePointer discipline rollout** — apply to all existing adapter outputs (USGS gauge, satellite, etc. per §3.1 of spec). Estimated: incremental; ~1 day per adapter.

For the Wyrd implementer:

1. **Issue #6 corrections** per §3.2. Estimated: ~30 minutes.
2. **Integration doc update** at `wyrd/doc/integration/contextus.md`. Estimated: ~30 minutes.
3. **Acceptance criteria for `contextus/` subpackage** — extend to include scope-node round-trip and EvidencePointer round-trip tests. Estimated: ~1 hour.

For the architecture instance (me) — done. This document is the deliverable.

---

## 7. References

- [`Contextus-Spec-v1.2.md`](Contextus-Spec-v1.2.md) §2.1, §2.3, §4.2, §4.4, §4.5, §9, §11.1, §11.3 — authoritative source
- [`Contextus-Theory-v1.4.md`](Contextus-Theory-v1.4.md) §2.2 (Locale), §3.6 (Insight Signal), §8 (Surveillance Mode)
- Wyrd integration doc: `wyrd/doc/integration/contextus.md`
- Wyrd issue [#6](https://github.com/JamesPagetButler/wyrd/issues/6) — Walk-phase glue (this doc resolves)
- BMA handoff [`2026-05-03-bma-reply-cth-qbp.md`](../BMA/doc/handoff/2026-05-03-bma-reply-cth-qbp.md) Q3 — `cth_id` namespace pattern that `contextus.*` mirrors
- BMA handoff [`2026-05-05-architecture-update-to-bma.md`](../BMA/doc/handoff/2026-05-05-architecture-update-to-bma.md) — sister architecture update from same session

**Attribution carrying forward.** Contextus theory builds on prior work in cross-domain insight discovery, knowledge graphs, and cognitive ergonomics; specific citations carried forward from Theory v1.4 (Yellowstone trophic cascade case studies, whale shark confluence work, black hole hidden correlation literature). Wyrd quaternion-native hypergraph: substrate authored by James Paget Butler with Lean 4 formal soundness corpus. Spec v1.2 § references are by authority of the implementer who maintains it.

---

*Status: ARCHITECTURE DOC | Next: Contextus implementer absorbs into Spec v1.3; Wyrd implementer corrects issue #6 + integration doc | Cadence: respond to either implementer's questions per §6 priority order*
