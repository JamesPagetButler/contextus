# Contextus–CTH Bridge Specification

**Version:** 0.1 (Draft)
**Date:** 2026-04-11
**Author:** James Paget Butler (beekeeper), Claude (red team)
**Status:** Initial working draft

---

## 1. Purpose

This document specifies the formal bridge between Contextus (domain-agnostic insight discovery platform) and the Confluent Trust Hypergraph (CTH) epistemic health framework. The bridge enables insights surfaced by Contextus to be seeded into new or existing CTH instances for structured investigation, while retaining provenance to the originating activation pattern without importing full Contextus context.

## 2. Problem Statement

Contextus and CTH share a MuninnDB substrate but currently operate as decoupled systems. Contextus generates insight candidates via co-activation patterns; CTH evaluates epistemic claims via trust scoring and derivation chains. No mechanism exists to:

- Translate a Contextus insight into a CTH-compatible claim structure.
- Triage which insights warrant structured investigation.
- Maintain a lightweight coupling between a seeded claim and its originating evidence.
- Prevent heartbeat traffic from overwhelming CTH at scale.

## 3. Core Concepts

### 3.1 Trust Receipt

When a Contextus insight crosses the triage threshold (§4), Contextus mints a **Trust Receipt** — a compact, self-contained object that seeds a CTH anchor without importing the full activation graph.

A Trust Receipt contains:

| Field | Type | Description |
|---|---|---|
| `receipt_id` | quaternion address | Unique identifier in MuninnDB address space |
| `claim` | string | Natural-language statement of the insight |
| `activation_baseline` | float64 | Co-activation strength at time of minting |
| `confidence` | float64 | Derived from activation strength and pattern stability |
| `provenance_hash` | []byte | Content-addressed pointer into MuninnDB linking to the full activation subgraph |
| `rationale` | string | Short summary of why the co-activation fired |
| `source_nodes` | []quaternion | The MuninnDB nodes whose co-activation produced the insight |
| `minted_at` | timestamp | Time of receipt creation |

The `provenance_hash` is the key architectural decision. It makes the full Contextus context *reachable* but not *present* — CTH can operate on the claim without needing to represent or store the upstream evidence.

### 3.2 Seed Anchor

A Trust Receipt maps to a **Tier 3 seed anchor** in CTH. The seed anchor inherits the receipt's claim, is assigned an initial ρ derived from the confidence score, and enters the normal CTH evaluation pipeline.

The seed anchor carries one additional property not present on standard anchors: a **heartbeat coupling** to its originating Contextus pattern (§5).

### 3.3 Anchor Lifecycle

```
Contextus insight
        │
        ▼
  [Triage Gate] ──── below threshold ──→ (no action)
        │
   above threshold
        │
        ▼
  Trust Receipt minted
        │
        ▼
  Seed Anchor (Tier 3)
        │
        ├── heartbeat coupling active
        │   (Contextus → CTH, unidirectional)
        │
        ▼
  Investigation
        │
        ├── evidence accumulates → ρ increases
        ├── no evidence, no heartbeat → gentle decay toward floor
        ├── heartbeat strengthens → fractional ρ bump
        └── heartbeat weakens → accelerated decay toward floor
        │
        ▼
  Promotion to Tier 2 (on own evidence)
        │
        └── heartbeat coupling goes dormant
```

## 4. Triage Gate

Not every Contextus co-activation warrants CTH tracking. The triage gate filters insights before Trust Receipt minting.

### 4.1 Triage Criteria

A Contextus insight crosses the triage threshold when it satisfies **all** of:

1. **Activation strength:** Above a configurable minimum (prevents noise).
2. **Pattern stability:** The co-activation has persisted across multiple MuninnDB consolidation cycles (prevents transient spikes).
3. **Novelty:** The insight does not duplicate an existing CTH anchor (checked via claim similarity against active CTH instances).
4. **Expressibility:** The co-activation pattern is reducible to a natural-language claim. Diffuse, high-dimensional activations that resist summarisation are deferred.

### 4.2 Triage Output

Insights that pass the gate produce a Trust Receipt. Insights that fail are not discarded — they remain in Contextus and may re-qualify as activation patterns evolve.

## 5. Heartbeat Coupling

### 5.1 Mechanism

A seed anchor maintains a lightweight, unidirectional coupling to its originating Contextus activation pattern. Contextus monitors the activation strength of the source nodes and emits a heartbeat signal to the coupled CTH anchor when a significance threshold is crossed.

**Directionality is strictly enforced:** Contextus pushes to CTH, never the reverse. CTH epistemic judgments must not contaminate Contextus activation patterns, as this would create a feedback loop where CTH could artificially sustain patterns that should naturally fade.

### 5.2 Relative Threshold

The seed anchor records `activation_baseline` at minting. A heartbeat is emitted only when:

```
|activation_current - activation_baseline| > threshold(t)
```

Where `threshold(t)` is an adaptive function of the anchor's maturity:

- **Early life** (low CTH-internal evidence): Low threshold — the anchor is hungry for signal.
- **Maturing** (accumulating derivations and evidence): Rising threshold — the heartbeat matters less relative to direct evidence.
- **Tier 2 promoted:** Threshold effectively infinite — heartbeat coupling dormant.

This adaptive threshold naturally limits heartbeat traffic as the system scales. A mature CTH instance with hundreds of promoted anchors generates near-zero heartbeat load.

### 5.3 Heartbeat Dynamics

Three regimes govern how the heartbeat affects seed anchor ρ:

| Regime | Condition | Effect on ρ |
|---|---|---|
| **Silence** | Delta below threshold | Gentle asymptotic decay on internal clock; ρ approaches floor but never reaches zero |
| **Strengthening** | Activation rises above threshold | Fractional ρ bump (< 1 confirmed bit); structurally equivalent to weak observational confirmation |
| **Weakening** | Activation falls below threshold (downward) | Accelerated decay toward floor |

### 5.4 Floor Constraint

A seed anchor's ρ never decays to zero. The original observation was real — the co-activation fired. Decaying to zero would be rewriting history. The floor should be a small positive value representing "this was observed but not confirmed."

### 5.5 Double-Counting Prevention

Once a seed anchor is promoted to Tier 2 on its own CTH-internal evidence (derivations, independent measurements, confluence), the heartbeat coupling goes dormant. This prevents double-counting where Contextus activation strength and CTH-internal evidence both contribute to ρ for a claim that has already been independently validated.

The coupling is dormant, not severed. If the anchor is later demoted (evidence retracted, derivation chain broken), the heartbeat can reactivate to provide background signal during re-evaluation.

## 6. Active Evidence Seeking

The heartbeat coupling (§5) is a passive mechanism — it listens for changes in Contextus activation patterns. But real scientific investigation is not passive. When a researcher has an inkling, they search for existing evidence before waiting for new data to arrive. The bridge must support this.

Active evidence seeking operates in two phases: **initial survey** (pre-seeding, runs once) and **horizon scanning** (post-seeding, runs continuously).

### 6.1 Initial Survey

When a Trust Receipt passes the triage gate but *before* the seed anchor is minted, the bridge executes an automated initial survey. This is the equivalent of an afternoon on Google Scholar — a bounded search to establish what evidence already exists in the world.

The initial survey:

1. Extracts a **search signature** from the Trust Receipt — a set of keywords, entity names, and conceptual terms derived from the claim, rationale, and provenance fields.
2. Searches available sources: arxiv, journal databases, preprint servers, public datasets, domain-specific repositories (e.g., GWOSC for gravitational wave data, GenBank for genomic data).
3. Classifies results into four categories:
   - **Supporting:** evidence consistent with the claim
   - **Contradicting:** evidence inconsistent with the claim
   - **Adjacent:** related work that neither confirms nor denies but contextualises
   - **Null:** no relevant results found
4. Enriches the Trust Receipt with survey results before minting.

The initial survey directly affects the seed anchor's starting ρ:

- Supporting evidence raises initial ρ above the Trust Receipt's raw confidence score.
- Contradicting evidence lowers it — and if overwhelming, can *block minting entirely*. An insight that's already refuted in the literature shouldn't enter CTH as a Tier 3 anchor.
- Adjacent evidence doesn't change ρ but populates the seed anchor's context, reducing the beekeeper's effort to understand the claim's landscape.
- Null results are informative — they mean the claim is in genuinely unexplored territory, which affects the expected evidence cycle time (§5 of the validation document).

### 6.2 Search Signature

The search signature is a structured object generated from the Trust Receipt, designed to produce useful search queries across heterogeneous sources. It is *not* a simple keyword list — it encodes the claim at multiple levels of abstraction to catch both direct evidence and conceptual neighbours.

A search signature contains:

| Field | Description | Example (SR-07: IGF1/ESR1) |
|---|---|---|
| `primary_terms` | Core entities and relationships in the claim | `["IGF1", "ESR1", "LCORL", "sighthound", "longevity"]` |
| `broadening_terms` | Conceptual generalisations | `["growth hormone signaling", "breed-specific aging", "body size lifespan tradeoff"]` |
| `narrowing_terms` | Specific testable predictions | `["Irish Wolfhound IGF1 expression", "ESR1 LCORL independence"]` |
| `exclusion_terms` | Terms to filter noise | `["IGF1 diabetes", "IGF1 cancer therapy"]` |
| `source_hints` | Preferred repositories | `["pubmed", "genbank", "biorxiv"]` |
| `provenance_authors` | Authors from the Trust Receipt's provenance chain | Used for citation graph traversal |

Search signature generation is a candidate for LLM-assisted automation (via BMA-BRIDGE). The claim and rationale fields provide enough context for an LLM to produce a reasonable signature, but the beekeeper should be able to review and refine it.

### 6.3 Horizon Scanning

After the seed anchor is minted, horizon scanning monitors for *new* evidence that appears in the world. This is the automated version of the standing instruction to check arxiv for Goodwin's VLA follow-up publications — except generalised and applied to every active seed anchor.

Horizon scanning differs from the heartbeat in three ways:

| Property | Heartbeat | Horizon Scan |
|---|---|---|
| Source | Contextus activation patterns | External publications and datasets |
| Trigger | Activation delta crosses threshold | New search result matches signature |
| What it detects | Internal pattern changes | External evidence arriving in the world |
| Direction | Contextus → CTH | External world → CTH (via search) |

Horizon scanning runs on a configurable cadence tied to the seed anchor's estimated evidence cycle time. An anchor expecting weekly evidence gets scanned weekly. An anchor expecting annual evidence gets scanned monthly (not annually — you want to catch it when it arrives, not a year late).

### 6.4 Scan Result Processing

When a horizon scan returns results, they require classification before they can affect ρ. Not every search hit is evidence.

**Classification pipeline:**

```
New search result
       │
       ▼
  [Relevance filter] ─── irrelevant ──→ discard
       │
    relevant
       │
       ▼
  [Novelty check] ─── already known ──→ discard
       │
     novel
       │
       ▼
  [Stance classification]
       │
       ├── Supporting ──→ EvidenceRef (type: "survey_support")
       ├── Contradicting ──→ EvidenceRef (type: "survey_contradict")
       └── Adjacent ──→ Context enrichment (no ρ change)
```

Relevance filtering and stance classification are judgment calls. In the early system, these should be LLM-assisted (via BMA-BRIDGE) with beekeeper review. As the system matures and accumulates training data from beekeeper corrections, the classification can become more autonomous.

### 6.5 Claim Mutation

Adjacent search results sometimes reveal that the claim itself needs refinement — not confirmed, not denied, but *reshaped*. A search for IGF1/ESR1 independence might surface a paper showing partial coupling under specific conditions. This doesn't kill the claim, but the original formulation is no longer accurate.

Claim mutation is the hardest problem in active evidence seeking. Three options:

1. **Amend in place:** Update the seed anchor's claim text and adjust ρ. Simple but loses history — you can't see what the original claim was.
2. **Fork:** Create a new seed anchor with the refined claim, linked to the original. Preserves history but risks inventory proliferation.
3. **Version:** The seed anchor carries a claim version history. Current claim is always the latest version. ρ trajectory is continuous across versions.

**Recommendation:** Versioning (option 3). It preserves history, avoids proliferation, and mirrors how real scientific hypotheses evolve — they're refined, not replaced wholesale. The version history also provides a natural narrative of how the insight matured, which is valuable for the beekeeper's understanding.

### 6.6 Interaction with Heartbeat

Active evidence seeking and heartbeat coupling are complementary but must not double-count:

- If a horizon scan detects a new publication, and that same publication also causes Contextus source nodes to activate (triggering a heartbeat), the ρ contribution should only be counted once.
- Resolution: horizon scan results are tagged with source identifiers. Before applying a heartbeat ρ bump, CTH checks whether the heartbeat's activation was caused by data already processed through the scan pipeline. If so, the heartbeat is suppressed.
- This requires a lightweight deduplication index on the seed anchor, mapping source identifiers (DOIs, dataset IDs) to the channel through which they entered.

## 7. Transport

### 7.1 NATS Subjects

All bridge communication uses the existing NATS bus.

| Subject | Direction | Payload |
|---|---|---|
| `contextus.insight.candidate` | Contextus → Triage Gate | Raw insight with activation data |
| `contextus.insight.receipt` | Triage Gate → CTH | Trust Receipt (enriched by initial survey) |
| `contextus.heartbeat.{receipt_id}` | Contextus → CTH | `{receipt_id, activation_current, timestamp}` |
| `bridge.survey.request` | Triage Gate → Survey Service | Search signature + receipt ID |
| `bridge.survey.result` | Survey Service → Triage Gate | Classified search results |
| `bridge.horizon.{receipt_id}` | Horizon Scanner → CTH | New classified evidence |
| `bridge.horizon.schedule` | CTH → Horizon Scanner | Scan cadence updates |

### 7.2 Heartbeat Frequency

Heartbeats are event-driven (threshold crossing), not polled. No threshold crossing, no message. This is critical for scalability — the system is silent by default.

### 7.3 Horizon Scan Frequency

Horizon scans are cadence-driven, tied to the seed anchor's estimated evidence cycle time. The schedule is managed by CTH and published to the Horizon Scanner via `bridge.horizon.schedule`. Scan frequency decreases as the anchor matures (parallel to heartbeat threshold increasing).

## 8. Data Structures (Go)

```go
package bridge

import (
    "time"
    "github.com/helpfulengineering/q8"
)

// TrustReceipt is the compact handoff object from Contextus to CTH.
type TrustReceipt struct {
    ReceiptID              q8.Addr         `json:"receipt_id"`
    Claim                  string          `json:"claim"`
    ActivationBaseline     float64         `json:"activation_baseline"`
    Confidence             float64         `json:"confidence"`
    ProvenanceHash         []byte          `json:"provenance_hash"`
    Rationale              string          `json:"rationale"`
    SourceNodes            []q8.Addr       `json:"source_nodes"`
    MintedAt               time.Time       `json:"minted_at"`
    SearchSignature        SearchSignature `json:"search_signature"`
    InitialSurveyResults   []SurveyResult  `json:"initial_survey_results,omitempty"`
    EstimatedEvidenceCycle time.Duration   `json:"estimated_evidence_cycle"`
}

// SearchSignature drives both initial survey and ongoing horizon scanning.
type SearchSignature struct {
    PrimaryTerms     []string `json:"primary_terms"`
    BroadeningTerms  []string `json:"broadening_terms"`
    NarrowingTerms   []string `json:"narrowing_terms"`
    ExclusionTerms   []string `json:"exclusion_terms"`
    SourceHints      []string `json:"source_hints"`
    ProvenanceAuthors []string `json:"provenance_authors"`
}

// SurveyResult is a classified search result from initial survey or horizon scan.
type SurveyResult struct {
    SourceID    string    `json:"source_id"`    // DOI, dataset ID, URL
    Title       string    `json:"title"`
    Stance      Stance    `json:"stance"`       // supporting, contradicting, adjacent
    Summary     string    `json:"summary"`
    Confidence  float64   `json:"confidence"`   // Classification confidence
    FoundAt     time.Time `json:"found_at"`
    Channel     string    `json:"channel"`      // "initial_survey" or "horizon_scan"
}

// Stance represents the classified relationship of evidence to a claim.
type Stance string

const (
    StanceSupporting    Stance = "supporting"
    StanceContradicting Stance = "contradicting"
    StanceAdjacent      Stance = "adjacent"
)

// Heartbeat is the lightweight signal from Contextus to a coupled CTH anchor.
type Heartbeat struct {
    ReceiptID         q8.Addr   `json:"receipt_id"`
    ActivationCurrent float64   `json:"activation_current"`
    Timestamp         time.Time `json:"timestamp"`
}

// SeedAnchor extends a standard CTH anchor with heartbeat coupling and
// active evidence seeking.
type SeedAnchor struct {
    AnchorID           string          `json:"anchor_id"`
    Receipt            TrustReceipt    `json:"receipt"`
    Tier               int             `json:"tier"` // Initially 3
    Rho                float64         `json:"rho"`
    RhoFloor           float64         `json:"rho_floor"`
    HeartbeatActive    bool            `json:"heartbeat_active"`
    ThresholdFunc      ThresholdConfig `json:"threshold_func"`
    InternalEvidence   []EvidenceRef   `json:"internal_evidence"`
    ClaimVersions      []ClaimVersion  `json:"claim_versions"`
    DeduplicationIndex map[string]string `json:"dedup_index"` // source_id → channel
    ScanCadence        time.Duration   `json:"scan_cadence"`
}

// ClaimVersion tracks the evolution of a seed anchor's claim over time.
type ClaimVersion struct {
    Version   int       `json:"version"`
    Claim     string    `json:"claim"`
    Reason    string    `json:"reason"`     // What prompted the revision
    Source    string    `json:"source"`     // Which evidence triggered it
    Timestamp time.Time `json:"timestamp"`
}

// ThresholdConfig governs adaptive heartbeat sensitivity.
type ThresholdConfig struct {
    BaseThreshold  float64 `json:"base_threshold"`
    MaturityScalar float64 `json:"maturity_scalar"` // Multiplied by evidence count
}

// EvidenceRef is a pointer to CTH-internal evidence for the anchor.
type EvidenceRef struct {
    Type   string  `json:"type"` // "derivation", "measurement", "confluence",
                                 // "survey_support", "survey_contradict"
    RefID  string  `json:"ref_id"`
    Bits   float64 `json:"bits"` // Confirmed bits contributed (negative for contradicting)
}
```

## 9. Open Questions

1. **Threshold calibration:** What are sensible defaults for `BaseThreshold` and `MaturityScalar`? These need empirical tuning — possibly using QBP's existing CTH inventory as a retroactive test case.

2. **Multi-CTH routing:** If multiple CTH instances exist (one per research programme), does a single Trust Receipt seed one or many? Current assumption: one, targeting the most relevant CTH instance by domain.

3. **Rationale generation:** Who writes the natural-language rationale in the Trust Receipt? If automated, this requires an LLM pass over the activation subgraph. If manual (beekeeper), it becomes a bottleneck.

4. **Fractional bit quantification:** How much ρ does a strengthening heartbeat contribute? A full confirmed bit is too much — this is weak evidence. Proposal: define a `heartbeat_bit_fraction` parameter, likely in the range 0.05–0.20.

5. **Floor value:** What is the appropriate ρ floor? Too high and dead insights clutter the inventory. Too low and they're functionally invisible. Proposal: ρ_floor = 0.01 (present but not actionable without new evidence).

6. **Reactivation hysteresis:** When a demoted Tier 2 anchor reactivates its heartbeat coupling, should the threshold reset to early-life sensitivity or retain its mature value?

7. **Composability with Species Hypergraph:** The CTH-Species Hypergraph isomorphism is established. Does the bridge need special handling for biological insights routed through Materia-Bio, or does the standard Trust Receipt suffice?

8. **Search signature generation:** Should signature generation be fully automated (LLM via BMA-BRIDGE), beekeeper-curated, or hybrid? Automation scales but risks poor search queries for novel domains. Hybrid (auto-generate, beekeeper reviews) is safest but adds latency.

9. **Stance classification reliability:** How accurate does automated stance classification need to be before the system can adjust ρ without beekeeper review? A false "contradicting" classification could prematurely suppress a valid insight. Consider a confidence threshold below which results are flagged for human review rather than applied automatically.

10. **Initial survey blocking vs. async:** Should the initial survey block Trust Receipt minting (safer — ensures enriched receipt) or run asynchronously (faster — seeds the anchor immediately, survey results arrive later)? Blocking is correct for high-confidence triage gates; async may be needed if survey latency is high.

11. **Claim mutation governance:** Who authorises a claim version bump? If automated, an LLM rewriting claims without beekeeper approval could introduce drift. If manual, it's a bottleneck. Threshold proposal: adjacent evidence triggers a *suggested* mutation that requires beekeeper approval before taking effect.

12. **Survey depth bounds:** How many sources should the initial survey examine? Unbounded search is expensive and slow. Proposed default: top 20 results across all source hints, with an option for the beekeeper to request deeper survey on high-priority insights.

## 10. Relationship to Existing Systems

| System | Role in Bridge |
|---|---|
| **Contextus** | Source — generates insights, emits heartbeats |
| **CTH** | Sink — receives seeds, evaluates claims |
| **MuninnDB** | Shared substrate — stores activation graphs, provenance targets |
| **NATS** | Transport — carries receipts, heartbeats, survey results, and horizon signals |
| **BMA-BRIDGE** | Executor — automates search signature generation, stance classification, triage decisions |
| **Materia-Bio** | Domain specialisation — biological insights may carry additional metadata |
| **Horizon Scanner** | New service — executes periodic searches against seed anchor search signatures |

## 11. Next Steps

- [ ] Review open questions (§9) and make initial parameter decisions
- [ ] Implement `TrustReceipt`, `SearchSignature`, `SurveyResult`, and `Heartbeat` types in the `bridge` package
- [ ] Implement initial survey service as NATS consumer on `bridge.survey.request`
- [ ] Define stance classification pipeline (LLM-assisted with confidence thresholds)
- [ ] Define triage gate as NATS consumer on `contextus.insight.candidate`
- [ ] Extend CTH Engine v2.0 with `SeedAnchor` support, heartbeat listener, and horizon scan consumer
- [ ] Implement Horizon Scanner service with configurable cadence
- [ ] Implement deduplication index for heartbeat/horizon scan overlap
- [ ] Retroactive validation: attempt to mint Trust Receipts for known QBP insights and verify CTH behaviour matches expected dynamics
- [ ] Cross-domain validation: run synthetic receipts SR-06 through SR-08 through the survey pipeline to verify search signature quality outside physics
- [ ] Stress test heartbeat + horizon scan traffic at projected scale

---

*This document is a working draft under the Systema framework. Theory Cart status: initial. Engineering Cart status: pre-implementation.*
