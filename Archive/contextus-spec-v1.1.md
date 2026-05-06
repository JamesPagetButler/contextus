# PROJECT CONTEXTUS

## Technical Specification

Version 1.1 | April 2026

Helpful Engineering

James Paget Butler

Classification: Open Source | Licence: TBD

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 1.0 | March 2026 | Initial specification: 8 sections, full hypergraph schema, agent architecture, 20-week roadmap |
| 1.1 | April 2026 | Added §2.3 (InsightSignal node/edge types); added §4.4 (InsightSignal emission pipeline); added §5.3 (Bridge integration NATS subjects); added §9 (Contextus–CTH Bridge Integration); added §10 (Proportional Data Retention); updated roadmap |

---

## 1. Integration Map: Contextus on the Existing Stack

Contextus is a domain-specific layer on infrastructure that already exists. This section maps every Contextus component to its existing counterpart in the BMA/Materia/War Table ecosystem.

### 1.1 Component Mapping

| Contextus Need | Existing Component | Integration Notes |
|---|---|---|
| Knowledge graph | MuninnDB | Hypergraph store. Contextus defines domain-specific node and edge types. No engine changes needed. |
| Message bus | NATS | Ingestion pipeline, agent communication, event-driven processing. Existing subjects extended with `ctx.*` prefix. Bridge subjects use `contextus.*` and `bridge.*` prefixes. |
| Container orchestration | Podman | Source adapters run as Podman containers on Pop!_OS. Same pattern as BMA Crawl-phase containers. |
| REST API | Go (port 8422) | Same pattern as Materia (8420). Contextus serves on 8422. Shared router library. |
| AI interface | FastMCP (port 8423) | MCP server exposing Contextus queries to AI agents. Same pattern as Materia (8421). |
| Multi-model access | BMA-BRIDGE | Four-layer transport (MCP, CLI, GitHub Issues, file drop). No changes needed. |
| Agent framework | War Table agents | Doctrine markdown + TOML sidecars. Agent spawning logic reused; doctrines are domain-specific. |
| Spatiotemporal math | QBP Locale | Quaternion-valued addressing. World-line rendering shares math with QBP visualisation. |
| Object storage | Local / S3-compatible | Heavy rasters, satellite imagery, large CSV archives. Indexed by MuninnDB metadata nodes. |
| Epistemic evaluation | CTH (via Bridge) | Contextus emits Insight Signals. The Bridge packages promoted signals as Trust Receipts for CTH evaluation. Contextus never evaluates claims directly. |

New code required: source adapters, domain schema definitions, visualisation renderer, hypothesis query patterns, agent doctrines, Insight Signal emission pipeline, and Bridge integration layer. Infrastructure: zero new components.

---

## 2. Hypergraph Schema

All data in Contextus is stored in MuninnDB as typed nodes connected by typed hyperedges. This section defines the complete schema.

### 2.1 Node Types

#### Biological/Ecological Nodes

| Type Tag | Example ID | Key Metadata Fields |
|---|---|---|
| `NT_SPECIES` | species-canis-lupus | taxonomy, common_name, conservation_status, population_estimate |
| `NT_INDIVIDUAL` | wolf-1118F | species_ref, sex, birth_year, pack_ref, collar_id (nullable) |
| `NT_PACK` | pack-lamar-canyon | species_ref, territory_locale, member_ids[], formation_date |
| `NT_VEGETATION` | veg-willow-lamar-transect-7 | species, density_class, height_m, locale, survey_date |

#### Geographic/Environmental Nodes

| Type Tag | Example ID | Key Metadata Fields |
|---|---|---|
| `NT_GEOTHERMAL` | geo-old-faithful | feature_type [geyser|hotspring|fumarole|mudpot], locale, eruption_interval (nullable) |
| `NT_STREAM` | stream-soda-butte | name, length_km, gauge_ids[], watershed_id |
| `NT_GAUGE` | gauge-usgs-06191500 | locale, parameter [discharge|temp|turbidity], agency, active_since |
| `NT_GRID_CELL` | grid-44.6N-110.4W-1km | centroid_locale, resolution_m, elevation_m, slope_deg, aspect_deg, land_cover_class |

#### Temporal and Contextual Nodes

| Type Tag | Example ID | Key Metadata Fields |
|---|---|---|
| `NT_SEASON` | season-2020-winter | year, season [winter|spring|summer|autumn], start_date, end_date |
| `NT_EVENT` | event-wolf-reintro-1995 | event_type, description, date_range, source_refs[], significance |
| `NT_OBSERVATION` | obs-20230615-cam-trap-047 | locale, timestamp, observer_type [human|camera|satellite|sensor], raw_data_ref, confidence |
| `NT_DATA_SOURCE` | src-nps-wolf-project | agency, url, format, update_frequency, licence, date_range |
| `NT_HUMAN_ACTIVITY` | human-visitor-count-2020-03 | activity_type [visitation|road_use|backcountry_permit], count, locale_region, period |

#### Agent-Generated Nodes

| Type Tag | Example ID | Key Metadata Fields |
|---|---|---|
| `NT_INSIGHT_SIGNAL` | sig-20260414-scout-001 | See §2.3 for full definition |

### 2.2 Hyperedge Types

Hyperedges connect two or more nodes and carry their own metadata. Provenance metadata (source, confidence, method, timestamp) is mandatory on every edge.

#### Provenance Envelope (required on all edges)

| Field | Type | Description |
|---|---|---|
| `source_ref` | string | NT_DATA_SOURCE ID or agent ID that created this edge |
| `confidence` | float64 | 0.0–1.0 |
| `method` | string | How the relationship was established |
| `timestamp` | time.Time | When this edge was created or last updated |
| `provenance_tag` | string | `T` (theoretical), `P` (measured/observed), `D` (derived/inferred), `I` (imported) |

#### Edge Type Catalogue

| Type Tag | Arity | Description | Provenance Tag |
|---|---|---|---|
| `HE_OBSERVATION` | 2–N | Links an observation to the entities observed and the context (grid cell, season, conditions) | P |
| `HE_MOVEMENT` | 2 | Links an individual to a sequence of Locales (trajectory) | P |
| `HE_INTERACTION` | 2–N | Links entities whose world-lines converged (potential interaction) | D |
| `HE_TROPHIC` | 2 | Links predator to prey species (trophic relationship) | P or I |
| `HE_HABITAT` | 2–N | Links a species or individual to geographic/environmental nodes (habitat use) | D |
| `HE_CAUSAL_CLAIM` | 2–N | Links a cause node to an effect node with supporting evidence nodes | D |
| `HE_TEMPORAL_SEQUENCE` | 2+ | Orders events or observations in time within a Locale | P or D |
| `HE_AGENT_INFERENCE` | 2–N | Links an NT_INSIGHT_SIGNAL to the subgraph it describes. Must carry `provenance_tag = "D"`. | D |
| `HE_CONTEXT` | 2–N | Links an observation to environmental context (grid cell, weather, season) | D |
| `HE_SIGNAL_CONVERGENCE` | 2+ | Links multiple NT_INSIGHT_SIGNAL nodes when independent agents converge on overlapping subgraphs | D |

**Critical constraint:** Agent-generated edges (`HE_AGENT_INFERENCE`, `HE_SIGNAL_CONVERGENCE`) must always carry `provenance_tag = "D"` to distinguish them from observation-derived edges. This prevents inference loops where agents treat other agents' inferences as ground truth.

### 2.3 Insight Signal Schema

The Insight Signal (Theory §3.6) is both a node and a constellation of edges in the hypergraph.

#### NT_INSIGHT_SIGNAL Node Fields

| Field | Type | Description |
|---|---|---|
| `signal_id` | q8.Addr | Quaternion address in MuninnDB |
| `agent_type` | AgentClass | scout, correlation, synthesis |
| `anomaly_type` | AnomalyKind | density, correlation, narrative |
| `anomaly_score` | float64 | Information-theoretic surprise (0.0–1.0) |
| `confidence` | float64 | Bayesian confidence in signal validity |
| `confidence_variance` | float64 | Width of confidence distribution |
| `persistence` | int | Consolidation cycles survived |
| `structural` | bool | Landscape (true) or exploration (false) pattern |
| `locale_envelope` | LocaleBounds | Spatiotemporal bounds where pattern holds |
| `locale_grain` | GrainSpec | Finest resolution at which pattern is detectable |
| `first_seen` | time.Time | First emission |
| `last_seen` | time.Time | Most recent strengthening |
| `version` | int | Claim version (incremented on mutation) |
| `claim_history` | []ClaimVersion | Version trail for mutations |
| `promoted` | bool | Whether this signal has been exported to the Bridge |
| `promotion_receipt_id` | q8.Addr | Trust Receipt ID if promoted (nullable) |

#### Edges Created Per Signal

Each NT_INSIGHT_SIGNAL creates:
- One `HE_AGENT_INFERENCE` edge connecting the signal node to every node in its flagged subgraph
- Traversal path stored as ordered metadata on the inference edge
- If convergence detected: one `HE_SIGNAL_CONVERGENCE` edge linking the converging signals

---

## 3. Ingestion Pipeline

### 3.1 Source Adapters

Each data source is ingested by a dedicated adapter running as a Podman container. Adapters normalise data into the Contextus schema and publish to NATS.

| Adapter | Source | NATS Subject | Output Node Types |
|---|---|---|---|
| `ctx-adapter-usgs` | USGS Water Services API | `ctx.ingest.usgs` | NT_GAUGE, NT_OBSERVATION |
| `ctx-adapter-collar` | Movebank / CSV collar data | `ctx.ingest.collar` | NT_INDIVIDUAL, NT_OBSERVATION, HE_MOVEMENT |
| `ctx-adapter-nps` | NPS visitation and permit data | `ctx.ingest.nps` | NT_HUMAN_ACTIVITY |
| `ctx-adapter-satellite` | Sentinel/Landsat NDVI, land cover | `ctx.ingest.satellite` | NT_GRID_CELL updates |
| `ctx-adapter-literature` | Structured paper ingestion | `ctx.ingest.literature` | NT_DATA_SOURCE, HE_CAUSAL_CLAIM |
| `ctx-adapter-edna` | eDNA metabarcoding pipeline | `ctx.ingest.edna` | NT_OBSERVATION (multi-species), NT_SPECIES |
| `ctx-adapter-arxiv` | arXiv/journal preprint feeds | `ctx.ingest.arxiv` | NT_DATA_SOURCE (for horizon scanning) |

### 3.2 Ingestion Flow

```
External Source → Adapter Container → NATS (ctx.ingest.*) → MuninnDB Writer → Context Builder Agent → Enriched Nodes
```

The Context Builder agent enriches raw observations with environmental context: links to grid cells, elevation, land cover, nearest water, snow depth, season. It transforms raw GPS fixes into richly contextualised nodes.

---

## 4. Agent Architecture

### 4.1 Agent Definition

Each agent is defined by:
- A **doctrine document** (markdown): analytical perspective, priorities, search strategy
- A **parameter sidecar** (TOML): operational configuration (traversal depth, anomaly thresholds, trigger conditions)

Agents are spawned dynamically via the War Table pattern. No agents are hard-coded.

### 4.2 Agent Taxonomy

| Agent Type | Trigger | Behaviour |
|---|---|---|
| Scout | Continuous traversal | Traverses the hypergraph looking for anomalies. Flags but does not interpret. Writes `HE_AGENT_INFERENCE` edges with `inference_type = "anomaly"`. Emits `NT_INSIGHT_SIGNAL` with `anomaly_type = density`. |
| Correlation | On demand (researcher query) or triggered by scout alert | Tests a specific hypothetical relationship by querying the hypergraph for supporting and contradicting evidence. Reports balance of evidence, identifies confounders. Emits `NT_INSIGHT_SIGNAL` with `anomaly_type = correlation`. |
| Synthesis | Triggered by scout alerts or researcher request | Combines outputs from multiple scouts/correlators to propose narrative explanations. Outputs always flagged as hypotheses with falsifiable predictions. Emits `NT_INSIGHT_SIGNAL` with `anomaly_type = narrative`. |
| Provenance | On every write to MuninnDB | Validates provenance metadata completeness, checks for circular inference chains, flags confidence degradation. Does not produce Insight Signals; maintains data integrity. |
| Context Builder | On new NT_OBSERVATION | Links raw observations to environmental context: grid cell, elevation, land cover, nearest water, snow depth, season. Does not produce Insight Signals. |

### 4.3 Metadata Construction Pipeline

New observations arrive via NATS. The Context Builder agent enriches them. Scout agents monitor for anomalies. When scouts flag anomalies (by emitting Insight Signals), synthesis agents may be triggered to propose explanations. The provenance agent validates every write. This pipeline is asynchronous; agents communicate via NATS subjects and share state through MuninnDB.

### 4.4 Insight Signal Emission Pipeline

When an agent detects a pattern worth flagging:

1. **Subgraph identification:** The agent identifies the set of nodes and hyperedges constituting the pattern.
2. **Anomaly scoring:** The agent computes the anomaly score against its current baseline expectations.
3. **Confidence estimation:** The agent estimates Bayesian confidence from provenance chain quality.
4. **Locale computation:** The agent computes the spatiotemporal envelope and grain from the subgraph's Locale data.
5. **Signal creation:** The agent creates an `NT_INSIGHT_SIGNAL` node with all fields populated.
6. **Edge creation:** The agent creates `HE_AGENT_INFERENCE` edges connecting the signal to its subgraph.
7. **NATS publication:** The agent publishes the signal to `ctx.signal.emitted`.
8. **Convergence check:** The system checks for overlapping subgraphs with existing signals. If found, creates `HE_SIGNAL_CONVERGENCE` edges.

---

## 5. API and MCP Interface

### 5.1 Go REST API (Port 8422)

The REST API follows the Materia pattern. All responses include provenance metadata. Endpoints support Locale-based spatial queries natively.

| Endpoint | Method | Description |
|---|---|---|
| `/nodes` | GET | List/search nodes by type, Locale, or text query |
| `/nodes/{id}` | GET | Get a specific node with full metadata |
| `/edges` | GET | List/search edges by type, connected nodes, or provenance |
| `/edges/{id}` | GET | Get a specific edge with provenance envelope |
| `/traverse` | POST | Execute a hypergraph traversal from a starting node |
| `/locale/query` | POST | Spatial query: find all nodes/edges within a Locale region |
| `/locale/worldline` | GET | Get an entity's world-line (trajectory) over a time range |
| `/signals` | GET | List active Insight Signals, filterable by type, score, persistence |
| `/signals/{id}` | GET | Get a specific signal with full subgraph and history |
| `/signals/{id}/promote` | POST | Manually promote a signal to the Bridge triage gate |

### 5.2 MCP Interface (Port 8423)

FastMCP tools for AI agent access:

| Tool Name | Description |
|---|---|
| `ctx_query_locale` | Find nodes and edges within a spatiotemporal region |
| `ctx_traverse` | Walk the hypergraph from a starting node following specified edge types |
| `ctx_hypothesis_test` | Submit a hypothetical relationship for correlation agent testing |
| `ctx_anomaly_report` | Get current scout agent anomaly flags |
| `ctx_worldline` | Get an entity's trajectory with behavioural state segmentation |
| `ctx_signal_list` | List active Insight Signals above a threshold |
| `ctx_signal_inspect` | Get full detail on a specific Insight Signal including subgraph |
| `ctx_signal_promote` | Promote a signal to the Bridge triage gate |

### 5.3 NATS Subject Hierarchy

#### Internal Subjects

| Subject | Direction | Description |
|---|---|---|
| `ctx.ingest.*` | Adapters → MuninnDB Writer | Raw data ingestion |
| `ctx.agent.scout.*` | Scout → System | Scout traversal events and alerts |
| `ctx.agent.correlation.*` | Correlation → System | Hypothesis test results |
| `ctx.agent.synthesis.*` | Synthesis → System | Narrative hypothesis proposals |
| `ctx.agent.provenance.*` | Provenance → System | Validation results and warnings |
| `ctx.signal.emitted` | Agents → System | New Insight Signal created |
| `ctx.signal.strengthened` | System → Subscribers | Existing signal strengthened by new evidence |
| `ctx.signal.decayed` | System → Subscribers | Signal decay event (consolidation cycle) |
| `ctx.signal.mutated` | System → Subscribers | Signal claim version changed |
| `ctx.signal.absorbed` | System → Subscribers | Two signals merged |

#### Bridge Subjects (Contextus ↔ CTH)

| Subject | Direction | Description |
|---|---|---|
| `contextus.insight.candidate` | Contextus → Triage Gate | Promoted Insight Signal for evaluation |
| `contextus.insight.receipt` | Triage Gate → CTH | Trust Receipt minted from signal |
| `contextus.heartbeat.{receipt_id}` | Contextus → CTH | Activation change for coupled signal |
| `bridge.survey.request` | Triage Gate → Survey Service | Search signature for initial survey |
| `bridge.survey.result` | Survey Service → Triage Gate | Classified search results |
| `bridge.horizon.{receipt_id}` | Horizon Scanner → CTH | New evidence from horizon scan |
| `bridge.horizon.schedule` | CTH → Horizon Scanner | Scan cadence updates |

---

## 6. Visualisation Specification

### 6.1 Rendering Modes

| Mode | Use Case | Primary Channel |
|---|---|---|
| Spatial overview | "Where is everything?" | WebGL map with grid cell colouring, observation points, Insight Signal indicators |
| World-line playback | "What did this entity do?" | Animated trajectory with behavioural state colouring |
| Cascade navigator | "What caused what?" | Interactive causal chain following HE_CAUSAL_CLAIM edges |
| Signal browser | "What has the system noticed?" | Insight Signals ranked by anomaly score, filterable by type and Locale |

### 6.2 Perceptual Design Rules

- Maximum three active data layers at any time (hard cap).
- Confidence encoded as opacity (low = ghost, high = solid).
- Anomaly-flagged regions indicated by subtle colour shift, not jarring alerts.
- All colour encodings pass Okabe-Ito colourblind validation.
- Shape and pattern redundantly encode colour channel information.
- Animation speed: smooth transitions with ease-in-out curves.
- Insight Signals render as subtle perceptual cues whose prominence scales with anomaly score × confidence.

---

## 7. Data Sensitivity and Security

### 7.1 Access Tiers

| Tier | Data Resolution | Access Requirement |
|---|---|---|
| PUBLIC | Aggregated only. Population estimates, regional trends, coarse maps. | None |
| RESEARCHER | Individual-level with temporal lag. Non-sensitive precise coordinates. | Institutional affiliation + data use agreement |
| STEWARD | Real-time, full-resolution. Sensitive locations. | Domain authority + approved protocol + audit logging |

### 7.2 Enforcement

Sensitivity is enforced at the MuninnDB query layer, not the API layer. Every node carries a `sensitivity_tier` field. Queries are filtered by the authenticated user's tier before results are returned. This prevents accidental exposure through API bugs or new endpoints.

### 7.3 Vulnerability Disclosure

When an agent or researcher identifies a domain vulnerability (poaching corridor, pollution pathway), the finding is automatically classified as STEWARD sensitivity. Designated authorities are notified via secure channel. A 90-day embargo applies before reclassification. The provenance agent logs the full disclosure chain.

---

## 8. Contextus–CTH Bridge Integration

This section specifies how Contextus connects to the Confluent Trust Hypergraph for epistemic evaluation. Full details are in the Contextus–CTH Bridge Specification v0.1.

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

The Bridge maintains ongoing search signatures for each active seed anchor. A Horizon Scanner service periodically searches external sources (arxiv, journals, datasets) for new evidence. Results are classified and routed to the relevant CTH anchor.

---

## 9. Proportional Data Retention

For large-scale corpus monitoring (e.g., ~200M scientific papers), storage cost per item is proportional to epistemic proximity to active research flows.

### 9.1 Retention Tiers

| Tier | Proximity | Stored Data | Approximate Cost |
|---|---|---|---|
| **Core** | Direct evidence for active seed anchor; deep-dive completed | Full text, parsed claims, stance classification, citation graph, all metadata | ~500 KB–2 MB |
| **Near** | Matches search signature of active seed anchor | Abstract, claim fingerprint, citation links, compressed embedding, full metadata | ~10–50 KB |
| **Peripheral** | Weakly connected to active flow | DOI, title, author list, publication date, compressed embedding, claim fingerprint | ~1–5 KB |
| **Distant** | No current connection to any active flow | DOI, claim fingerprint, single compressed embedding vector | ~200–500 bytes |
| **Skeleton** | Never co-activated with any active flow | DOI and publication year only | ~50–100 bytes |

### 9.2 Tier Transitions

Transitions are driven by co-activation with active research flows:
- Upward promotion: co-activation strength crosses threshold → promote to next tier
- Downward demotion: Ebbinghaus decay in absence of co-activation → demote
- Core papers never drop below Peripheral (deep-dive investment preserved)
- Skeleton papers are never deleted (DOI remains findable)

### 9.3 Storage Sentinel

A Sentinel service monitors storage health:
- **70% capacity:** Advisory to beekeeper with tier distribution report
- **80% capacity:** Warning. Auto-increase decay rates. Reduce deep-dive sensitivity.
- **85% capacity:** Critical. Suspend new deep dives. Monitor-only mode.
- **90% capacity:** Emergency. Suspend horizon scanning. Read-only maintenance mode.

Continuous throttle: `throttle_factor = 1.0 - max(0, (usage_pct - 0.60) / 0.30)`

---

## 10. Development Roadmap

### 10.1 Phase 1: Foundation (Weeks 1–4)

- Schema registration: Define all node and edge types in MuninnDB, including NT_INSIGHT_SIGNAL. Write Go type definitions matching the schema.
- First adapter: `ctx-adapter-usgs` for USGS Water Services. Proves the ingestion pipeline end-to-end.
- Locale integration: Import QBP Locale library. Verify quaternion spatiotemporal addressing works for GPS coordinate conversion.
- REST API scaffold: Port 8422 with `/nodes`, `/edges`, and `/signals` endpoints. Authentication and sensitivity filtering.

### 10.2 Phase 2: Data Population (Weeks 5–8)

- GPS collar adapter: `ctx-adapter-collar` for Movebank/CSV wolf and elk collar data.
- NPS visitor adapter: `ctx-adapter-nps` for visitation and permit data.
- Context Builder agent: First agent deployment. Links raw observations to grid cells and environmental context.
- MCP server: FastMCP on port 8423 with `ctx_query_locale` and `ctx_traverse` tools.

### 10.3 Phase 3: Agent Intelligence (Weeks 9–12)

- Scout agents: Population anomaly detection, spatial displacement detection. Insight Signal emission pipeline implemented.
- Correlation agent: First hypothesis testing capability (wolf cascade).
- Provenance agent: Inference chain validation and confidence tracking.
- Signal lifecycle: Emission, strengthening, decay, and convergence detection operational.
- COVID module: Specific adapter and analysis pipeline for the COVID bear hypothesis.

### 10.4 Phase 4: Visualisation and Bridge (Weeks 13–16)

- Spatial overview: WebGL map with grid cell colouring, observation points, and Insight Signal indicators.
- World-line playback: Animated trajectories with behavioural state colouring.
- Cascade navigator: Interactive causal chain exploration.
- Signal browser: Insight Signal listing with filtering and inspection.
- Bridge integration: Triage gate, Trust Receipt minting, heartbeat coupling, NATS subjects operational.
- Perceptual validation: Colourblind testing, animation speed calibration, cognitive load assessment.

### 10.5 Phase 5: Integration and Scaling (Weeks 17–20)

- Synthesis agents: Higher-order pattern detection and narrative hypothesis generation.
- Satellite adapter: NDVI and land cover change detection from Sentinel/Landsat.
- Literature adapter: Structured ingestion of published papers.
- Horizon Scanner: Automated search against seed anchor search signatures.
- eDNA adapter: Metabarcoding pipeline integration.
- Proportional retention: Five-tier storage system activated for corpus monitoring.
- Storage sentinel: Capacity monitoring and predictive throttling.
- Cross-domain validation: Run system against at least two non-ecology domains to verify domain-agnosticism.

---

## 11. Go Data Structures

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
    EdgeID   string `json:"edge_id"`
    EdgeType string `json:"edge_type"`
}

// LocaleBounds defines the spatiotemporal region where a pattern holds.
type LocaleBounds struct {
    SpatialMin  q8.Addr   `json:"spatial_min"`
    SpatialMax  q8.Addr   `json:"spatial_max"`
    TemporalMin time.Time `json:"temporal_min"`
    TemporalMax time.Time `json:"temporal_max"`
}

// GrainSpec defines the finest resolution at which a pattern is detectable.
type GrainSpec struct {
    SpatialGrain  float64       `json:"spatial_grain_m"`   // Metres
    TemporalGrain time.Duration `json:"temporal_grain"`    // Duration
}
```

### 11.2 Bridge Integration Types

See Contextus–CTH Bridge Specification v0.1 for full Trust Receipt, Heartbeat, SearchSignature, and SeedAnchor type definitions.

---

## 12. Related Documents

| Document | Version | Scope |
|---|---|---|
| Contextus Theory | v1.3 | Philosophical foundations, design principles, Insight Signal theory |
| Contextus–CTH Bridge Specification | v0.1 | Trust Receipt, triage gate, heartbeat coupling, active evidence seeking |
| Contextus–CTH Bridge Synthetic Receipts | v0.1 | Parameter calibration using QBP and non-QBP test cases |
| CTH Engine | v2.0 | Confluent Trust Hypergraph evaluation framework |
| MuninnDB Specification | — | Hypergraph storage, Hebbian co-activation, Ebbinghaus decay |
| War Table Specification | v8 | Agent doctrine and spawning patterns |
| QBP Locale | v0.3 | Quaternion spatiotemporal addressing |

---

*End of Specification v1.1. Implementation language: Go. Primary implementation targets: Pop!_OS 22.04, AMD FX-8350, PowerColor Red Devil RX 9070 XT.*
