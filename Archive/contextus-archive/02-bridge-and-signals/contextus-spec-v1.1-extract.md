# PROJECT CONTEXTUS — Technical Specification v1.1

**Reconstruction status:** NEAR-COMPLETE for v1.0 → v1.1 changes (which are well-documented in chat 2). Inherited content from v1.0 is referenced rather than duplicated.
**Source chat:** https://claude.ai/chat/3175e792-70f2-48df-bd8e-58762dd066fc

---

# PROJECT CONTEXTUS

## Technical Specification

**Version 1.1 | April 2026**
**Helpful Engineering**
**James Paget Butler**

*Classification: Open Source | Licence: TBD*

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 1.0 | March 2026 | Initial specification: 8 sections, full hypergraph schema, agent architecture, 20-week roadmap |
| 1.1 | April 2026 | Added §2.3 (InsightSignal node/edge types); added §4.4 (InsightSignal emission pipeline); added §5.3 (Bridge integration NATS subjects); added §8 (Contextus–CTH Bridge Integration); added §9 (Proportional Data Retention); added §11 (Go data structures); updated roadmap |

---

## 1. Integration Map

> See `01-contextus-theory-spec/contextus-spec-v1.0-extract.md` §1 for the full mapping.

The mapping table is updated in v1.1 to add one row:

| Contextus Need | Existing Component | Integration Notes |
|---|---|---|
| Epistemic evaluation | CTH (via Bridge) | Contextus emits Insight Signals. The Bridge packages promoted signals as Trust Receipts for CTH evaluation. Contextus never evaluates claims directly. |

---

## 2. Hypergraph Schema

> Inherited from v1.0 (14 node types, 9 hyperedge types). v1.1 adds:

### 2.3 InsightSignal Node and Edge Types (NEW)

**InsightSignal** — a new node type representing a structured agent observation.

Fields (full Go struct in §11):
- `signal_id` — quaternion address
- `agent_type` — scout / correlation / synthesis
- `subgraph` — list of HyperedgeRef
- `traversal_path` — sequence of quaternion addresses
- `anomaly_type` — density / correlation / narrative
- `anomaly_score` — float64
- `locale_envelope` — bounding LocaleBounds
- `locale_grain` — temporal/spatial resolution
- `persistence` — count of consolidation cycles survived
- `structural` — boolean flag for structural patterns
- `confidence`, `confidence_variance` — float64
- `first_seen`, `last_seen` — timestamps
- `version`, `claim_history` — versioned claim evolution
- `promoted` — flag set when Trust Receipt minted
- `promotion_receipt_id` — link to Trust Receipt if promoted

**New edge types:**
- `SIG_DESCRIBES` — links InsightSignal to the hyperedges it describes
- `SIG_PRODUCED_BY` — links InsightSignal to the Agent that produced it
- `SIG_ABSORBED_INTO` — version edge when one signal is absorbed by another
- `SIG_CONVERGES_WITH` — soft edge between signals describing related patterns

---

## 3. Ingestion Pipeline

> Inherited from v1.0.

---

## 4. Agent Architecture

> Inherited from v1.0 (5-agent taxonomy, doctrine + TOML pattern). v1.1 adds:

### 4.4 Insight Signal Emission Pipeline (NEW)

Eight-step pipeline:

1. **Detection.** Agent (scout, correlation, or synthesis) walks the hypergraph and identifies a candidate pattern.
2. **Scoring.** Anomaly score computed from co-activation history; confidence computed from evidence strength.
3. **Locale binding.** The pattern's spatial/temporal envelope is computed.
4. **Signal construction.** A new InsightSignal node is created with subgraph, traversal path, scores, locale.
5. **Provenance binding.** SIG_PRODUCED_BY edge created. Source data lineage attached.
6. **Persistence check.** If pattern previously existed, link to prior version; increment persistence counter.
7. **Convergence check.** If subgraph overlaps existing signals, evaluate for absorption or convergence-edge creation.
8. **Publication.** Signal published to `ctx.signal.created` NATS subject. Visualisation subscribers notified.

---

## 5. REST and MCP Tool Catalogue

> Inherited from v1.0 base. v1.1 extends:

### 5.1 REST Endpoints (extensions)

- `GET /signals` — list active signals with filters (agent, anomaly type, locale)
- `GET /signals/{id}` — fetch signal detail with full subgraph
- `POST /signals/{id}/promote` — manually promote a signal to Trust Receipt
- `GET /signals/{id}/history` — claim version history

### 5.2 MCP Tools (extensions)

- `ctx_signal_list` — agent-callable signal enumeration
- `ctx_signal_inspect` — full signal detail for a given ID
- `ctx_signal_promote` — request promotion to CTH (subject to triage gate)

### 5.3 NATS Subjects (extensions)

**Internal Contextus signals:**
- `ctx.signal.created` — new signal emitted
- `ctx.signal.mutated` — signal version updated
- `ctx.signal.absorbed` — signal merged into another
- `ctx.signal.decayed` — signal activation crossed decay threshold

**Bridge integration:**
- `contextus.insight.candidate` — promoted signal proposed to Bridge
- `contextus.insight.receipt` — Trust Receipt minted (Bridge → CTH)
- `contextus.heartbeat.{receipt_id}` — Contextus → CTH activation update
- `bridge.survey.request`, `bridge.survey.result` — active evidence seeking

---

## 6. Visualisation Spec

> Inherited from v1.0. v1.1 note: Insight Signals render as subtle perceptual cues whose prominence scales with anomaly score × confidence.

---

## 7. Data Sensitivity and Security

> Inherited from v1.0 (PUBLIC / RESEARCHER / STEWARD tiers, MuninnDB-layer enforcement, 90-day vulnerability embargo).

---

## 8. Contextus–CTH Bridge Integration (NEW)

This section specifies how Contextus connects to the Confluent Trust Hypergraph for epistemic evaluation. Full details in the Contextus–CTH Bridge Specification v0.1.

### 8.1 Separation of Concerns

| System | Role | What It Produces |
|---|---|---|
| **Contextus** | Pattern detection | Insight Signals — structured descriptions of anomalous patterns |
| **Bridge** | Claim translation | Trust Receipts — compact claims with provenance and search signatures |
| **CTH** | Epistemic evaluation | Trust scores (ρ), derivation chains, confluence assessments |

Contextus never evaluates claims. CTH never monitors data. The Bridge translates between them.

### 8.2 Promotion Pipeline

When an Insight Signal's confidence and persistence exceed the triage threshold:

1. Signal published to `contextus.insight.candidate`
2. Bridge executes initial survey (automated literature/data search)
3. Survey results classified as supporting/contradicting/adjacent
4. Trust Receipt minted with: claim, confidence, provenance hash, search signature, survey results
5. Receipt published to `contextus.insight.receipt`
6. CTH creates seed anchor (Tier 3) with heartbeat coupling

### 8.3 Heartbeat Coupling

After promotion, Contextus monitors the originating activation pattern. When activation delta crosses a significance threshold, Contextus emits a heartbeat to the coupled CTH anchor. The heartbeat is unidirectional (Contextus → CTH, never reverse) to prevent CTH judgments from contaminating Contextus patterns.

### 8.4 Horizon Scanning

The Bridge runs a horizon scanner against each seed anchor's search signature on a configurable cadence. New evidence appearing in scanned sources triggers automated stance classification and posting to CTH for derivation chain integration.

---

## 9. Proportional Data Retention (NEW)

Corpus management for ingested papers, observations, and external sources operates on a five-tier model. Retention scales with epistemic proximity to active research flows.

| Tier | Retention | Trigger for tier transition |
|---|---|---|
| Core | Full text, all metadata, all derived artifacts, embeddings | Active citation in current research flow |
| Active | Full text, key metadata, embeddings | Recent activation, in-flight investigation |
| Watching | Abstract, key metadata, citation graph | Periodic activation |
| Background | Abstract only, citation graph | Rare activation |
| Skeleton | DOI + 1-line summary | Floor — never below this |

**Ebbinghaus decay** governs tier transitions. The decay clock is per-source: a paper cited yesterday holds Core tier; one untouched for six months drops toward Background.

**Storage sentinel** monitors aggregate corpus size against capacity targets. When capacity pressure rises, the sentinel applies a continuous throttle function: lower tiers are demoted gradually rather than hitting hard walls. The sentinel is part of Contextus infrastructure, not CTH.

**The Skeleton tier is the floor.** A paper that was once relevant to research never disappears from the corpus. Fading to zero would be erasing history.

---

## 10. Twenty-Week Phased Roadmap

> Inherited from v1.0 with v1.1 additions:

- Weeks 1-4: Phase 1 — Foundation (unchanged)
- Weeks 5-8: Phase 2 — Multiple adapters + context builder agent (unchanged)
- Weeks 9-12: Phase 3 — Scout and correlation agents + **Insight Signal emission pipeline (NEW)**
- Weeks 13-16: Phase 4 — Synthesis agents + visualisation prototype + **Bridge integration (NEW)**
- Weeks 17-20: Phase 5 — Integration, scaling, additional adapters (literature, satellite, eDNA), cross-domain validation + **proportional retention activation + horizon scanner (NEW)**

**Crawl exit gate:** 50 Trust Receipts through full lifecycle + completion of synthetic receipt calibration (SR-01 through SR-08).

---

## 11. Go Data Structures (NEW)

### 11.1 Insight Signal Types

```go
package contextus

import (
    "time"
    "github.com/helpfulengineering/q8"
)

// AgentClass identifies the type of agent that produced a signal.
type AgentClass string

const (
    AgentScout       AgentClass = "scout"
    AgentCorrelation AgentClass = "correlation"
    AgentSynthesis   AgentClass = "synthesis"
)

// AnomalyKind classifies what makes a pattern interesting.
type AnomalyKind string

const (
    AnomalyDensity     AnomalyKind = "density"
    AnomalyCorrelation AnomalyKind = "correlation"
    AnomalyNarrative   AnomalyKind = "narrative"
)

// InsightSignal is the atomic unit of agent output.
// It represents a pattern detected in the hypergraph that an agent
// has flagged as anomalous, correlated, or narratively coherent.
type InsightSignal struct {
    SignalID           q8.Addr        `json:"signal_id"`
    AgentType          AgentClass     `json:"agent_type"`
    Subgraph           []HyperedgeRef `json:"subgraph"`
    TraversalPath      []q8.Addr      `json:"traversal_path"`
    AnomalyType        AnomalyKind    `json:"anomaly_type"`
    AnomalyScore       float64        `json:"anomaly_score"`
    LocaleEnvelope     LocaleBounds   `json:"locale_envelope"`
    LocaleGrain        GrainSpec      `json:"locale_grain"`
    Persistence        int            `json:"persistence"`
    Structural         bool           `json:"structural"`
    Confidence         float64        `json:"confidence"`
    ConfidenceVariance float64        `json:"confidence_variance"`
    FirstSeen          time.Time      `json:"first_seen"`
    LastSeen           time.Time      `json:"last_seen"`
    Version            int            `json:"version"`
    ClaimHistory       []ClaimVersion `json:"claim_history,omitempty"`
    Promoted           bool           `json:"promoted"`
    PromotionReceiptID *q8.Addr       `json:"promotion_receipt_id,omitempty"`
}

// ClaimVersion tracks the evolution of a signal's description.
type ClaimVersion struct {
    Version   int       `json:"version"`
    Claim     string    `json:"claim"`
    Reason    string    `json:"reason"`
    Source    string    `json:"source"`
    Timestamp time.Time `json:"timestamp"`
}

// HyperedgeRef is a lightweight pointer to a hyperedge in MuninnDB.
type HyperedgeRef struct {
    EdgeID q8.Addr `json:"edge_id"`
    // ... edge metadata
}

// LocaleBounds describes the spatial/temporal envelope of a pattern.
type LocaleBounds struct {
    SpatialMin q8.Addr   `json:"spatial_min"`
    SpatialMax q8.Addr   `json:"spatial_max"`
    TimeMin    time.Time `json:"time_min"`
    TimeMax    time.Time `json:"time_max"`
}

// GrainSpec describes the resolution at which the pattern is meaningful.
type GrainSpec struct {
    SpatialMeters float64       `json:"spatial_meters"`
    TemporalGrain time.Duration `json:"temporal_grain"`
}
```

### 11.2 Bridge Types (cross-reference)

> Full bridge data structures (TrustReceipt, SeedAnchor, Heartbeat, ThresholdConfig, EvidenceRef, ClaimVersion, Stance) are in `insight-signal-go-types.go` in this directory.

---

*End of Specification v1.1.*
*Implementation language: Go. Primary implementation targets: Pop!_OS 22.04, AMD FX-8350, PowerColor Red Devil RX 9070 XT.*
