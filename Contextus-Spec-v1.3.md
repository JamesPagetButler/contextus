# PROJECT CONTEXTUS

## Technical Specification

Version 1.3 | May 2026

Helpful Engineering

James Paget Butler

Classification: Open Source | Licence: TBD

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 1.0 | March 2026 | Initial specification: 8 sections, full hypergraph schema, agent architecture, 20-week roadmap |
| 1.1 | April 2026 | Added §2.3 (InsightSignal node/edge types); added §4.4 (InsightSignal emission pipeline); added §5.3 (Bridge integration NATS subjects); added §9 (Contextus–CTH Bridge Integration); added §10 (Proportional Data Retention); updated roadmap |
| 1.2 | April 2026 | Added three new agent types (Edge Scout, Corpus Edge Scout, Bridge Agent) per Theory v1.4 §§8.2–8.6; added §4.5 (Session-Scoped Agents); added new NATS subjects, MCP tools, visualisation overlay, and Go types |
| 1.3 | May 2026 | Closes the four open architectural questions from v1.2: added §4.6 (Scope Nodes — `NT_SCOPE_PHYSICAL`, `NT_SCOPE_CONCEPTUAL`, `HE_SCOPE_MEMBERSHIP`); added §5.4 (Evidence Pointer Discipline with tier-conditional fields and cap-per-tier eviction); added §4.4 Synthesis-as-persistence-boundary clause; added `EvidencePointer` Go type and `AnomalyStructural` constant in §11.1; reference rows in §2.1 and §2.2 pointing to §4.6. Resolves Wyrd issue [#6](https://github.com/JamesPagetButler/wyrd/issues/6) `SignalSource` enum (corrected to `scout | correlation | synthesis`). |

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

New code required: source adapters, domain schema definitions, visualisation renderer, hypothesis query patterns, agent doctrines, Insight Signal emission pipeline, Bridge integration layer, edge scout session manager, corpus diversity monitor, bridge agent session state. Infrastructure: zero new components.

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

#### Focus-Area Nodes

| Type Tag | Example ID | Key Metadata Fields |
|---|---|---|
| `NT_SCOPE_PHYSICAL` | scope-squam-lake-watershed | See §4.6 for full definition (geometry, elevation_range, temporal_range, grain, parent_scope_id) |
| `NT_SCOPE_CONCEPTUAL` | scope-physics-fluid-dynamics | See §4.6 for full definition (ontology_uri, parent_scope_id, related_scope_ids, tags) |

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
| `HE_SCOPE_MEMBERSHIP` | 2 | Links a scope node (`NT_SCOPE_PHYSICAL` or `NT_SCOPE_CONCEPTUAL`) to a member node within the scope. See §4.6 for membership-predicate semantics. | P or D |

**Critical constraint:** Agent-generated edges (`HE_AGENT_INFERENCE`, `HE_SIGNAL_CONVERGENCE`) must always carry `provenance_tag = "D"` to distinguish them from observation-derived edges. This prevents inference loops where agents treat other agents' inferences as ground truth. `HE_SCOPE_MEMBERSHIP` carries `P` when membership is asserted by beekeeper config and `D` when inferred by an adapter (e.g., spatial intersection of a `Locale` against `NT_SCOPE_PHYSICAL.geometry`).

### 2.3 Insight Signal Schema

The Insight Signal (Theory §3.6) is both a node and a constellation of edges in the hypergraph.

#### NT_INSIGHT_SIGNAL Node Fields

| Field | Type | Description |
|---|---|---|
| `signal_id` | q8.Addr | Quaternion address in MuninnDB |
| `agent_type` | AgentClass | scout, correlation, synthesis |
| `anomaly_type` | AnomalyKind | density, correlation, narrative, structural (see §11.1; `structural` is provisional in v1.3 — Theory v1.5 will formalize) |
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
| `evidence` | []EvidencePointer | Tier-conditional pointers to evidence outside the signal's flagged subgraph. See §5.4 (Evidence Pointer Discipline) for the population rules and cap-per-tier eviction. |

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

| Agent Type | Scope | Trigger | Behaviour |
|---|---|---|---|
| Scout | Global | Continuous traversal | Traverses the hypergraph looking for anomalies. Flags but does not interpret. Writes `HE_AGENT_INFERENCE` edges with `inference_type = "anomaly"`. Emits `NT_INSIGHT_SIGNAL` with `anomaly_type = density`. |
| Correlation | Global | On demand (researcher query) or triggered by scout alert | Tests a specific hypothetical relationship by querying the hypergraph for supporting and contradicting evidence. Reports balance of evidence, identifies confounders. Emits `NT_INSIGHT_SIGNAL` with `anomaly_type = correlation`. |
| Synthesis | Global | Triggered by scout alerts or researcher request | Combines outputs from multiple scouts/correlators to propose narrative explanations. Outputs always flagged as hypotheses with falsifiable predictions. Emits `NT_INSIGHT_SIGNAL` with `anomaly_type = narrative`. |
| Provenance | Global | On every write to MuninnDB | Validates provenance metadata completeness, checks for circular inference chains, flags confidence degradation. Does not produce Insight Signals; maintains data integrity. |
| Context Builder | Global | On new NT_OBSERVATION | Links raw observations to environmental context: grid cell, elevation, land cover, nearest water, snow depth, season. Does not produce Insight Signals. |
| Edge Scout | Session | Active researcher session; continuous within session | Monitors boundary between researcher's explored subgraph and unexplored remainder. Compares explored nodes' full connectivity against queried nodes. Ranks unexplored connections by statistical significance and recency. Renders as peripheral boundary cues in visualisation. Does not emit Insight Signals; publishes to `ctx.edge.boundary`. |
| Corpus Edge Scout | Per-search | Before search results are returned to researcher or querying agent | Evaluates domain distribution of search results. If convergence detected, injects broadening queries for adjacent domains. Publishes diversity assessment to `ctx.corpus.diversity`. |
| Bridge Agent | Session | Active researcher session; concurrent with all searches | Maintains concurrent view of scout flags and active search trajectory. Detects when surveillance finding and active search converge on the same physical process from different domain perspectives (local minimum signature). Publishes cross-domain connection candidates to `ctx.bridge.intervention`. Does not emit Insight Signals. |

### 4.3 Metadata Construction Pipeline

New observations arrive via NATS. The Context Builder agent enriches them. Scout agents monitor for anomalies. When scouts flag anomalies (by emitting Insight Signals), synthesis agents may be triggered to propose explanations. The provenance agent validates every write. Within a researcher session, edge scouts and bridge agents run concurrently. This pipeline is asynchronous; agents communicate via NATS subjects and share state through MuninnDB.

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

#### Synthesis as Persistence Boundary

The session-scoped agents introduced in §4.5 (Edge Scout, Corpus Edge Scout, Bridge Agent) emit ephemeral NATS events on `ctx.edge.boundary.{session_id}`, `ctx.corpus.diversity`, and `ctx.bridge.intervention.{session_id}`. These events are session-scoped and never persisted as `NT_INSIGHT_SIGNAL` nodes — §4.5 is unambiguous that *"session state is not persisted."*

When a session-scoped finding warrants graph-level persistence — for example, a Bridge Agent convergence that exceeds a configured confidence threshold and represents a genuine cross-domain pattern rather than a transient boundary effect — promotion happens through a **Synthesis agent** subscribed to the relevant ephemeral subject, not by direct write from the session-scoped agent.

> **Synthesis agents may subscribe to `ctx.bridge.intervention.*`, `ctx.edge.boundary.*`, and `ctx.corpus.diversity` and mint Insight Signals when ephemeral session findings warrant persistence. The Synthesis agent is the sole persistence boundary for session-scoped agent outputs.**

This preserves the §4.2 separation of concerns (only Scout / Correlation / Synthesis emit Insight Signals; Edge Scout / Corpus Edge Scout / Bridge Agent do not) while creating exactly one well-defined gate at which ephemeral findings can become hypergraph-resident records. The minted signal's `agent_type` is `synthesis`; its `anomaly_type` is typically `narrative` for Bridge Agent convergences and `structural` for Edge Scout boundary blindspots or Corpus Edge Scout diversity findings (see §11.1 — `structural` is provisional in v1.3 pending Theory v1.5 formalization).

### 4.5 Session-Scoped Agents

Edge scouts and bridge agents operate per researcher session, not globally. This section specifies their lifecycle and state management.

#### Edge Scout Session Lifecycle

1. **Session start:** One edge scout spawned per researcher session. Parameters: `significance_threshold` (minimum statistical significance for a boundary flag), `domain_weight` (weighting by recency of change vs. static connectivity).
2. **Exploration tracking:** The edge scout maintains a session-local set of explored node IDs, updated as the researcher queries or traverses nodes.
3. **Boundary computation:** Continuously computed as the set of nodes in the explored set with edges to unexplored domains. Weighted by the significance of the unexplored connections.
4. **Flag publication:** When a boundary connection's significance exceeds threshold, the edge scout publishes an `EdgeScoutFlag` to `ctx.edge.boundary.{session_id}`. Flags are ranked, not streamed — the scout maintains a ranked list and publishes updates only when the top-N ranking changes.
5. **Rendering trigger:** The visualisation layer subscribes to the session's boundary subject and updates node highlight state.
6. **Session end:** Edge scout de-spawned. Session state is not persisted (boundary flags are session-relative, not hypergraph state).

**Design requirement (Theory §8.4):** Every researcher session must have at least one active edge scout.

#### Corpus Edge Scout Per-Search Lifecycle

1. **Search intercept:** Before results are returned to the researcher or to any querying agent, results pass through the corpus edge scout.
2. **Diversity evaluation:** The scout computes domain distribution across results: number of distinct domains, vocabulary concentration (ratio of domain-specific to domain-neutral terms), institutional source spread.
3. **Convergence detection:** If domain concentration exceeds threshold (configurable, default: >70% of results from a single domain), convergence is flagged.
4. **Query injection:** On convergence, the scout generates one to three broadening queries targeting adjacent domains not represented in the current results. These are issued automatically.
5. **Result augmentation:** Broadening query results are appended to the original results with a `source: "corpus_broadening"` tag, visible to the researcher.
6. **Diversity report:** Published to `ctx.corpus.diversity` for session monitoring.

#### Bridge Agent Session Lifecycle

1. **Session start:** One bridge agent spawned per researcher session.
2. **State maintenance:** The bridge agent maintains two concurrent inventories:
   - *Surveillance inventory:* All active scout flags with their subgraph summaries (physical/process domain tags).
   - *Search trajectory:* Summary of the current search's node visits and result domains.
3. **Overlap detection:** Continuously compares the physical/process domain tags of surveillance flags against the search trajectory. Looks for flags and searches that share physical process nodes (water leaving ground, energy transfer, population dynamics) but are separated in knowledge domain (hydrology vs. plant physiology).
4. **Intervention trigger:** When overlap is detected with sufficient confidence, publishes a `BridgeIntervention` to `ctx.bridge.intervention.{session_id}` with: the matching scout flag, the active search context, the proposed cross-domain query.
5. **Rendering:** Bridge interventions render as a distinct visual cue (configurable; default: a connecting arc between the flagged subgraph and the active search nodes, visible only when the researcher hovers over either).
6. **Session end:** Bridge agent de-spawned.

### 4.6 Scope Nodes

The Locale framework (Theory §2.2) addresses individual observations spatiotemporally. But Contextus also needs first-class addressing of *focus areas* — the bounded regions in either physical space or topical space within which a researcher, agent, or downstream consumer (BMA, Sharp Butler) wants to operate. A focus area is not a single observation; it is a *queryable subgraph definition* that constrains traversal, anchors agent attention, and drives §9 retention-tier transitions for member nodes.

This section introduces **scope nodes** as that addressing layer.

#### 4.6.1 Two Sibling Node Types, One Architectural Role

A focus area can be physical (a watershed, a lake, a national park) or conceptual (a topic, an academic discipline, a research area). These are siblings, not different things:

- **Both** are queryable subgraph definitions.
- **Both** support hierarchical containment via a `parent_scope_id` reference (a sub-watershed contains a stream reach; a sub-discipline like *fluid dynamics* sits inside *physics*).
- **Both** compose multiplicatively — a query can be in the *Squam Lake watershed* scope AND the *fluid dynamics* scope simultaneously, and the result is the intersection of the two scope memberships.
- **Both** anchor §9 retention-tier transitions — activating a scope promotes member signals toward Core; deactivating demotes them toward Skeleton.

What differs is the *membership predicate* — how the system decides whether a given node is a member. Physical scopes use geometric intersection of `Locale` against `geometry`. Conceptual scopes use a weaker, adapter-defined predicate (see §4.6.4). The architectural role is identical, so the node-type pair shares a shape.

The query `"what is known about fluid dynamics in the Squam Lake watershed in spring 2024?"` becomes a Locale-bounded traversal from the intersection of two scope nodes — not a custom search engine.

#### 4.6.2 NT_SCOPE_PHYSICAL

A region of space-time with a name, geometry, and queryable membership.

| Field | Type | Description |
|---|---|---|
| `scope_id` | string | e.g., `scope-squam-lake-watershed`. Stable identifier. |
| `name` | string | Human-readable name. |
| `geometry` | GeoJSON | Polygon or multipolygon bounding the scope spatially. |
| `elevation_range_m` | [float64, float64] | Optional vertical bounds. `null` if not relevant (e.g. surface-only scopes). |
| `temporal_range` | LocaleBounds | When this scope is "active". May be open-ended (e.g. `{start: "1995-01-15", end: null}` for an ongoing reintroduction zone). |
| `grain` | GrainSpec | Finest spatial/temporal resolution at which membership predicates are evaluated within this scope. |
| `parent_scope_id` | string | Hierarchical containment (nullable). |
| `tags` | []string | e.g., `["watershed", "freshwater", "north-america"]`. |

#### 4.6.3 NT_SCOPE_CONCEPTUAL

A topic or domain with a name and queryable membership.

| Field | Type | Description |
|---|---|---|
| `scope_id` | string | e.g., `scope-physics-fluid-dynamics`. Stable identifier. |
| `name` | string | Human-readable name. |
| `ontology_uri` | string | Optional pointer to an external ontology entry (e.g. a Wikidata QID, a domain-specific URI). Used by adapters that infer membership from external classification. |
| `parent_scope_id` | string | Hierarchical containment (nullable). |
| `related_scope_ids` | []string | Cross-references for traversal — e.g. `scope-physics-fluid-dynamics` relates to `scope-mathematics-pde`. Not hierarchical; used for adjacent-domain corpus edge scouting. |
| `tags` | []string | e.g., `["physics", "fluid-dynamics", "transport-phenomena"]`. |

#### 4.6.4 HE_SCOPE_MEMBERSHIP

A binary edge connecting a scope node to a member node within the scope. Catalogued in §2.2.

| Field | Type | Description |
|---|---|---|
| `scope_id` | string | The `NT_SCOPE_PHYSICAL` or `NT_SCOPE_CONCEPTUAL` ID. |
| `member_id` | string | Any node in the graph (typically `NT_OBSERVATION`, `NT_INSIGHT_SIGNAL`, `NT_INDIVIDUAL`, `NT_GAUGE`, `NT_DATA_SOURCE`). |
| `since` | time.Time | When membership was first established. |
| `confidence` | float64 | 0.0–1.0; for inferred memberships, reflects predicate confidence. |
| `provenance_tag` | string | `P` for asserted memberships (beekeeper config); `D` for inferred memberships (adapter computation). |

A single member may belong to multiple scopes (e.g., an observation in *Squam Lake watershed* may also be in *fluid dynamics* and in *spring 2024*). Each membership is a separate edge.

**Membership predicate** — varies by scope type and adapter:

- **Physical scopes** use geometric intersection: a node has membership iff its `Locale` (or `LocaleBounds` if it has spatial extent) intersects the scope's `geometry` and its timestamp falls within `temporal_range`. This is well-defined and adapter-uniform.
- **Conceptual scopes** use an **implementation-defined predicate** (see §4.6.5). v1.3 ships the schema; the predicate ships per-adapter.

#### 4.6.5 Conceptual Scope Membership: Implementation-Defined

The membership predicate for conceptual scopes is deferred to the adapter that establishes the membership. v1.3 documents the contract; v1.x will document the canonical predicates as they emerge.

**The contract:** an adapter creating an `HE_SCOPE_MEMBERSHIP` edge to a conceptual scope must populate `confidence` and `method` (per the standard Provenance Envelope in §2.2) reflecting how the membership was determined. Acceptable methods include but are not limited to:

- **Asserted (`provenance_tag = "P"`):** beekeeper configuration explicitly listed the member in the scope.
- **Tag overlap:** the member node carries tags that overlap with the scope's `tags` field above a configured threshold.
- **Embedding similarity:** the member's content embedding has cosine similarity with the scope's `ontology_uri`-dereferenced text above a configured threshold.
- **Classifier output:** an external classifier (LLM, ontology mapper) confirms the member's relevance to the scope.
- **Citation graph proximity:** for `NT_DATA_SOURCE` (paper) members, the citation graph distance to scope-anchor papers is below a threshold.

The expected pattern is that early adapters will use simple methods (tag overlap, asserted), and Walk-phase adapters will graduate to embedding-similarity or classifier methods as the corpus matures. Spec v1.3 does not commit to one mechanism because the right mechanism depends on the corpus in question.

> *Footnote: future Spec revisions will record the canonical predicates that emerge from production adapters, with named `MembershipMethod` values added to the Provenance Envelope `method` field.*

#### 4.6.6 Scope Activation and Retention

Scope nodes are the operational lever for Spec §9 retention-tier transitions. Activating a scope (e.g. `"focus on Squam Lake watershed"`) signals the retention layer that signals connected to that scope via `HE_SCOPE_MEMBERSHIP` should be promoted toward Core. Deactivating reverses the signal.

In practice, scope activation happens at the session layer:

- A researcher session declares its active scopes (zero or more physical, zero or more conceptual).
- The retention layer treats `HE_SCOPE_MEMBERSHIP` edges to active scopes as a strong promotion signal alongside the existing co-activation strength from §9.2.
- When all sessions referencing a scope close, the scope deactivates and member promotion pressure fades; §9.2 demotion proceeds normally.

This is what makes Contextus "universe-shaped in principle, focus-shaped in practice": the schema supports a universe of scope memberships, but the active scope set at any moment determines what fits on disk.

#### 4.6.7 Worked Example: `{Squam Lake watershed × fluid dynamics}` in Spring 2024

A researcher studying turbulent mixing in stratified lake systems opens a session and activates two scopes:

- `scope-squam-lake-watershed` (physical) — geometry is the watershed polygon; `parent_scope_id` is `scope-new-hampshire-lakes`.
- `scope-physics-fluid-dynamics` (conceptual) — `ontology_uri` points to a discipline taxonomy entry; `tags = ["physics", "fluid-dynamics", "transport-phenomena"]`.

Plus an implicit Locale-bounds filter on Spring 2024 (a temporal Locale envelope, not a scope per se).

Query: *"What density anomalies have agents flagged in this combined scope?"*

The traversal:
1. Find all `HE_SCOPE_MEMBERSHIP` edges to either scope, where `since` ≤ Spring 2024 endpoint.
2. Take the intersection of member sets.
3. Filter members by Locale bounds.
4. From this filtered set, find connected `NT_INSIGHT_SIGNAL` nodes with `anomaly_type = density`.
5. Rank by anomaly score × confidence × persistence per §4.4.

Concrete signals that might land in the result: gauge-temperature anomalies during stratification onset; eDNA shifts concurrent with thermocline movement; a Bridge Agent–promoted signal connecting hydrology and plant-physiology nodes (a Hogan-style cross-domain finding).

The same architecture handles a QBP-physics session activating `scope-magnetar-spectroscopy` (conceptual) with no physical scope active, querying for narrative anomalies across published spectra. One mechanism, two domains.

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
| `/session/{id}/boundary` | GET | Get current edge scout boundary flags for a session |
| `/session/{id}/bridge` | GET | Get current bridge agent interventions for a session |

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
| `ctx_edge_boundary` | Get ranked unexplored boundary connections for the current session. Returns top-N EdgeScoutFlags sorted by significance. |
| `ctx_corpus_diversity` | Get the diversity assessment for the most recent search. Returns domain distribution and any injected broadening queries. |
| `ctx_bridge_interventions` | Get active bridge agent interventions for the current session. Returns BridgeIntervention list with cross-domain connection candidates. |

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
| `ctx.edge.boundary.{session_id}` | Edge Scout → Visualisation | Updated boundary flag ranking for session |
| `ctx.corpus.diversity` | Corpus Scout → System | Search diversity assessment and broadening queries |
| `ctx.bridge.intervention.{session_id}` | Bridge Agent → Visualisation | Cross-domain connection candidate for session |

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

### 5.4 Evidence Pointer Discipline

The §9 Proportional Data Retention tiers (Core / Near / Peripheral / Distant / Skeleton) already prescribe what *content* is held at each tier. The lower tiers — Distant and Skeleton — are essentially pointer-only: a Skeleton-tier paper is a DOI and a year (50–100 bytes); a Distant-tier paper adds a claim fingerprint and a single embedding vector (200–500 bytes). The pattern is already there; this section names it and generalises it across all evidence kinds, not just paper corpora.

#### 5.4.1 The Where-Not-What Principle

> **Contextus is the *index* of evidence, not the evidence itself.**

The Wyrd graph holds what Contextus *knows*: relationships, anomalies, hypotheses, scope memberships. The actual data — satellite imagery, gauge timeseries, journal articles, eDNA reads, magnetar X-ray spectra — lives wherever it was sourced. Contextus knows where it lives and can fetch it on demand. Most queries do not need the raw data; they need the metadata-and-relationships layer.

This is the same posture web search engines adopt: Google does not store a copy of every page on the web; it stores an index that points back at every page. The index fits on Google's servers; the web does not.

For Contextus, the consequence is concrete. The deployment hardware (Crawl-phase: a 250 GB SATA SSD already at ~70% capacity) cannot hold the universe of relevant evidence. It does not need to. It needs to know how to *reach* that evidence and to verify, on dereference, that what it reaches is what it indexed.

#### 5.4.2 The EvidencePointer Type

The atomic unit of "where, not what" is the `EvidencePointer`. Full Go definition is in §11.1; the schema-level fields are:

| Field | Type | Description | Populated at tier |
|---|---|---|---|
| `Locator` | string | URL, file path, archive coordinate, NATS subject, telescope observation ID, DOI, etc. | All tiers (mandatory) |
| `LocatorKind` | string | `"https"` / `"file"` / `"archive"` / `"nats"` / `"telescope_obs_id"` / `"doi"` / `"summary"` / etc. | All tiers (mandatory) |
| `Hash` | []byte | SHA-256 of the dereferenced content. Verifies that the locator still resolves to the same content as when first read. | Peripheral and above |
| `SizeBytes` | int64 | Size of the dereferenced content. | Peripheral and above |
| `LoadedAt` | time.Time | First-dereference timestamp. | Peripheral and above |
| `AccessHint` | string | `"cold"` / `"warm"` / `"hot"` — adapter hint for the retention layer about expected dereference cost. | Peripheral and Near (dropped at Distant per §5.4.5) |
| `Note` | string | Optional human-readable provenance ("Hogan et al. 2026 §3.2 figure 4 caption"). | Core only |

A signal carries a slice of pointers, `evidence []EvidencePointer`. The slice is governed by §5.4.4 cap-per-tier rules.

#### 5.4.3 Tier-Conditional Field Population

Spec §9.1 sizes Skeleton at 200–500 bytes and Distant at 200–500 bytes. A naive `EvidencePointer` carrying every field above is ~135 bytes of structural overhead before the `Locator` (which is itself frequently 50–200 bytes for URLs and DOIs) and consumes ~80% of the Skeleton budget by itself. To preserve the §9 tier budget contract, fields populate conditionally:

- **Skeleton**: `Locator` + `LocatorKind` only. ~110 bytes worst case.
- **Distant**: `Locator` + `LocatorKind` only. ~110 bytes (drops `Note`, `AccessHint` from the original architecture-doc proposal — see §5.4.5).
- **Peripheral**: + `Hash`, `SizeBytes`, `LoadedAt`, `AccessHint`. ~135 bytes worst case.
- **Near**: + same. ~135 bytes; cap allows up to 20 pointers (~2.7 KB total).
- **Core**: + `Note`. ~135–500 bytes per pointer including notes; cap allows up to 50.

The Go struct (§11.1) defines all fields uniformly. Tier policy is enforced by the retention layer at write time — when a signal demotes from Peripheral to Distant, the retention layer zeros the per-tier-disallowed fields. Promotion in the other direction populates fields lazily on first dereference.

This keeps the Go API surface uniform (one type, one shape) while preserving the §9.1 byte-budget contract. The alternative — separate types per tier — was considered and rejected as premature separation-of-concerns.

#### 5.4.4 Cap-per-Tier with Summary-Pointer Eviction

Theory §3.6.3 names *Strengthening* and *Mutation* as ongoing lifecycle events. A long-lived signal can accumulate evidence pointers indefinitely as Synthesis or Bridge Agent finds additional corroborating sources. Without a bound, a Core-tier signal could grow an unbounded `evidence` list — equivalent to holding the corpus we explicitly said we would not hold.

The bound is a tier-conditional cap:

| Tier | Cap | Total pointer bytes (worst case) | Pointer headroom against §9.1 budget |
|---|---|---|---|
| Skeleton | 1 | ~110 B | within 50–100 B budget; tight but fits |
| Distant | 5 | ~550 B | exceeds the original 200–500 B; tightening pointer (§5.4.5) reclaims headroom |
| Peripheral | 5 | ~675 B | well within the 1–5 KB budget |
| Near | 20 | ~2.7 KB | well within the 10–50 KB budget |
| Core | 50 | ~6.75 KB | trivial against the 500 KB–2 MB content budget |

When the cap is exceeded — for example, a Near-tier signal that has accumulated 21 evidence pointers during Strengthening — eviction does not delete pointers. Eviction *summarises* them.

**Summary-pointer eviction:**

1. Sort the current evidence list by `confidence` (ascending).
2. Take the lowest-confidence (`current_count - cap + 1`) pointers.
3. Delegate them to an aggregate index: write the list to an external locator (S3 object, file, NATS subject) and obtain its locator.
4. Replace those pointers with a single new pointer of `LocatorKind = "summary"` whose `Locator` is the aggregate index, and whose `Note` (Core only) records `"N evicted pointers from Strengthening cycles X–Y"`.

The summary pointer occupies one cap slot. The signal retains the *fact* of corroboration — `"this signal was supported by N evidence items at peak"` — without holding the pointers themselves. On promotion, the summary pointer is dereferenceable: the aggregate index can be fetched and the original pointers reconstructed.

This preserves Strengthening semantics (signals that gathered more evidence during their lifetime are not penalised by losing that history on demotion) while bounding the on-disk footprint per signal.

#### 5.4.5 Distant Tier: Tightened Pointer

The original architecture-doc proposal had `EvidencePointer` carrying `Hash`, `LoadedAt`, `AccessHint`, and `Note` at all tiers. Byte arithmetic against §9.1 showed that a Distant-tier signal with 5 pointers under that proposal would exceed the §9.1 budget of 200–500 bytes (`5 × 135 B = 675 B` plus the rest of the Distant-tier node body). Two corrections were possible:

- **(a) Tighten the pointer at Distant** — drop `Note` and `AccessHint` at Distant, keeping only `Locator`, `LocatorKind`, and the always-populated structural fields.
- **(b) Recalibrate §9.1 upward** — admit that Distant-tier entries with rich pointer fingerprints are worth 5–10 KB.

v1.3 adopts **(a)** — preserves the existing §9.1 budget contract, keeps Skeleton/Distant in the "pointer-only" regime where they were originally placed, and avoids cascading recalibration of the other tier budgets. The choice is reversible in v1.4 if implementation experience shows Distant entries genuinely need richer pointer metadata.

#### 5.4.6 The Signal-Level Evidence Slot Supplements, Does Not Duplicate

A reader may notice that the signal's flagged subgraph (Theory §3.6.2 "The Subgraph") already contains `NT_OBSERVATION` nodes, which themselves link via `HE_OBSERVATION` to `NT_DATA_SOURCE` nodes that carry their own URLs. Why does the signal also carry an `evidence` slot?

Because the signal's `evidence` slot holds pointers to evidence *outside* its immediate subgraph — corroborating papers a Synthesis agent cited during Strengthening; an external dataset a Correlation agent referenced during hypothesis testing; a Bridge Agent's literature find that supports the cross-domain convergence. These are evidentiary references that the signal, not the underlying observations, asserts.

The two layers do not duplicate. The subgraph's evidence chain (observation → data source → URL) records *what was observed*. The signal's `evidence` slot records *what additionally supports the agent's interpretation of the pattern*. Both are useful; neither is redundant with the other.

When a signal is promoted to a CTH Trust Receipt (§8.2), both layers feed the search-signature construction: the subgraph evidence chain anchors the claim's empirical basis, and the signal-level evidence pointers anchor the agent's reasoning trace.

---

## 6. Visualisation Specification

### 6.1 Rendering Modes

| Mode | Use Case | Primary Channel |
|---|---|---|
| Spatial overview | "Where is everything?" | WebGL map with grid cell colouring, observation points, Insight Signal indicators |
| World-line playback | "What did this entity do?" | Animated trajectory with behavioural state colouring |
| Cascade navigator | "What caused what?" | Interactive causal chain following HE_CAUSAL_CLAIM edges |
| Signal browser | "What has the system noticed?" | Insight Signals ranked by anomaly score, filterable by type and Locale |
| Boundary overlay | "What am I not asking?" | Edge scout boundary flags rendered as peripheral hue shifts on nodes at the explored/unexplored frontier |
| Bridge overlay | "Am I in a local minimum?" | Bridge agent interventions rendered as arcs connecting scout flags to active search nodes; visible on hover |

### 6.2 Perceptual Design Rules

- Maximum three active data layers at any time (hard cap).
- Confidence encoded as opacity (low = ghost, high = solid).
- Anomaly-flagged regions indicated by subtle colour shift, not jarring alerts.
- All colour encodings pass Okabe-Ito colourblind validation.
- Shape and pattern redundantly encode colour channel information.
- Animation speed: smooth transitions with ease-in-out curves.
- Insight Signals render as subtle perceptual cues whose prominence scales with anomaly score × confidence.
- **Edge scout boundary flags** render as a distinct peripheral hue shift on boundary nodes. Shift intensity scales with the significance of unexplored connections. Never renders as a modal or alert.
- **Bridge agent interventions** render as thin connecting arcs visible only on hover. They do not interrupt the researcher's current focus; they are discoverable on demand.

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
- **Edge Scout:** Session manager, explored subgraph tracking, boundary computation, NATS publication. `ctx_edge_boundary` MCP tool. Boundary overlay in visualisation.
- **Corpus Edge Scout:** Search intercept, domain diversity evaluation, broadening query injection. Integrated into all search paths (REST, MCP, BMA-BRIDGE queries).

### 10.4 Phase 4: Visualisation and Bridge (Weeks 13–16)

- Spatial overview: WebGL map with grid cell colouring, observation points, and Insight Signal indicators.
- World-line playback: Animated trajectories with behavioural state colouring.
- Cascade navigator: Interactive causal chain exploration.
- Signal browser: Insight Signal listing with filtering and inspection.
- Bridge integration: Triage gate, Trust Receipt minting, heartbeat coupling, NATS subjects operational.
- **Bridge Agent:** Session-scoped surveillance/search overlap detection. Physical process tag inference. Intervention publication and hover-reveal arc visualisation.
- Perceptual validation: Colourblind testing, animation speed calibration, cognitive load assessment. Specific validation of edge scout boundary intensity and bridge agent arc discoverability.

### 10.5 Phase 5: Integration and Scaling (Weeks 17–20)

- Synthesis agents: Higher-order pattern detection and narrative hypothesis generation.
- Satellite adapter: NDVI and land cover change detection from Sentinel/Landsat.
- Literature adapter: Structured ingestion of published papers.
- Horizon Scanner: Automated search against seed anchor search signatures.
- eDNA adapter: Metabarcoding pipeline integration.
- Proportional retention: Five-tier storage system activated for corpus monitoring.
- Storage sentinel: Capacity monitoring and predictive throttling.
- Cross-domain validation: Run system against at least two non-ecology domains to verify domain-agnosticism. **Specifically include a domain where biology/ecology intersects with physics/engineering to test bridge agent performance on the canonical failure mode.**

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
    AgentScout         AgentClass = "scout"
    AgentCorrelation   AgentClass = "correlation"
    AgentSynthesis     AgentClass = "synthesis"
)

// AnomalyKind classifies what makes a pattern interesting.
//
// AnomalyStructural is provisional in Spec v1.3 and will be formalized in
// Theory v1.5. It is the kind used by Synthesis-promoted signals minted from
// session-scoped agent output (Edge Scout boundary blindspots, Corpus Edge
// Scout diversity findings) where neither density / correlation / narrative
// fits cleanly. See §4.4 (Synthesis as Persistence Boundary) for context.
//
// Note: AnomalyKind is distinct from the InsightSignal.Structural bool field,
// which records whether the pattern is landscape (persistent under
// perturbation) or exploration (contingent). The two carry different
// information; both may be set independently on the same signal.
type AnomalyKind string

const (
    AnomalyDensity     AnomalyKind = "density"
    AnomalyCorrelation AnomalyKind = "correlation"
    AnomalyNarrative   AnomalyKind = "narrative"
    AnomalyStructural  AnomalyKind = "structural" // provisional in v1.3; see Theory v1.5 (forthcoming)
)

// InsightSignal is the atomic unit of agent output.
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
    Evidence           []EvidencePointer `json:"evidence,omitempty"` // tier-conditional; see §5.4
}

// EvidencePointer is the load-bearing primitive for the where-not-what
// principle (§5.4): a signal carries pointers to evidence outside its
// flagged subgraph, not the evidence itself. Field population is
// tier-conditional and enforced by the retention layer at write time.
//
// Tier rules (§5.4.3):
//   Skeleton/Distant: Locator + LocatorKind only (~110 bytes).
//   Peripheral:       + Hash, SizeBytes, LoadedAt, AccessHint.
//   Near:             same as Peripheral.
//   Core:             + Note.
//
// At Distant the AccessHint and Note fields are dropped to preserve the
// §9.1 byte budget (§5.4.5).
type EvidencePointer struct {
    Locator     string    `json:"locator"`                 // URL, file path, archive coordinate, NATS subject, DOI, telescope_obs_id
    LocatorKind string    `json:"locator_kind"`            // "https" | "file" | "archive" | "nats" | "doi" | "telescope_obs_id" | "summary" | ...
    Hash        []byte    `json:"hash,omitempty"`          // SHA-256 of dereferenced content (Peripheral+)
    SizeBytes   int64     `json:"size_bytes,omitempty"`    // size of dereferenced content (Peripheral+)
    LoadedAt    time.Time `json:"loaded_at,omitempty"`     // first-dereference timestamp (Peripheral+)
    AccessHint  string    `json:"access_hint,omitempty"`   // "cold" | "warm" | "hot" — adapter hint (Peripheral, Near; dropped at Distant)
    Note        string    `json:"note,omitempty"`          // human-readable provenance (Core only)
}

// EvidenceTierCap returns the per-tier cap on the number of EvidencePointers
// that may live on a single signal. Eviction beyond the cap is by
// summary-pointer eviction per §5.4.4 — the spec, not this helper, defines
// the eviction algorithm; this is the numeric contract only.
func EvidenceTierCap(tier string) int {
    switch tier {
    case "core":
        return 50
    case "near":
        return 20
    case "peripheral":
        return 5
    case "distant":
        return 5
    case "skeleton":
        return 1
    default:
        return 0
    }
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
    SpatialGrain  float64       `json:"spatial_grain_m"`
    TemporalGrain time.Duration `json:"temporal_grain"`
}
```

#### Cross-reference: SedenionResult

Insight Signals from QBP-domain sources (e.g., a Scout agent monitoring sedenion-typed BMA hypergraph state) may carry a `SedenionResult` payload as part of their subgraph metadata. The canonical definition lives in `qbp-compute-unit/doc/wyrd-integration.md`:

```go
// (canonical: qbp-compute-unit/doc/wyrd-integration.md)
type SedenionResult struct {
    Value     [16]float64
    ZDClass   uint8       // 0 = NoZD, 1 = CrossCopySymbolic (ZDCHK.SYM), 2 = GeneralFullMultiply
    ZDIndices [4]uint8
}
```

Contextus does not redefine this type; it imports the canonical version. Signals carrying a `SedenionResult` payload typically operate at `Wyrd.model.TierSedenion` and are emitted only by adapters with QBP-CU bridging context. See §8 (Contextus–CTH Bridge Integration) for how sedenion-typed signals interact with the bridge.

### 11.2 Bridge Integration Types

See Contextus–CTH Bridge Specification v0.1 for full Trust Receipt, Heartbeat, SearchSignature, and SeedAnchor type definitions.

### 11.3 Session-Scoped Agent Types

```go
package contextus

import "time"

// EdgeScoutFlag is a ranked boundary connection flagged by the Edge Scout.
// Published to ctx.edge.boundary.{session_id} when top-N ranking changes.
type EdgeScoutFlag struct {
    SessionID          string    `json:"session_id"`
    BoundaryNodeID     string    `json:"boundary_node_id"`     // Node in explored set
    UnexploredDomain   string    `json:"unexplored_domain"`    // Domain not yet visited
    ConnectionCount    int       `json:"connection_count"`     // Edges to unexplored domain
    Significance       float64   `json:"significance"`         // Statistical significance (0–1)
    RecentChangeCount  int       `json:"recent_change_count"`  // Connections changed recently
    Rank               int       `json:"rank"`                 // Position in current top-N
    ComputedAt         time.Time `json:"computed_at"`
}

// CorpusDiversityReport summarises a search's domain distribution.
// Published to ctx.corpus.diversity after each search.
type CorpusDiversityReport struct {
    SearchID            string           `json:"search_id"`
    OriginalQuery       string           `json:"original_query"`
    DomainDistribution  map[string]float64 `json:"domain_distribution"` // domain → fraction
    VocabConcentration  float64          `json:"vocab_concentration"`  // 0=diverse, 1=monoculture
    ConvergenceDetected bool             `json:"convergence_detected"`
    BroadeningQueries   []string         `json:"broadening_queries,omitempty"`
    ComputedAt          time.Time        `json:"computed_at"`
}

// BridgeIntervention is a cross-domain connection candidate from the Bridge Agent.
// Published to ctx.bridge.intervention.{session_id} when local minimum detected.
type BridgeIntervention struct {
    SessionID           string    `json:"session_id"`
    ScoutFlagID         string    `json:"scout_flag_id"`        // The surveillance flag involved
    ScoutPhysicalDomain string    `json:"scout_physical_domain"` // Physical process tag of flag
    SearchPhysicalDomain string   `json:"search_physical_domain"` // Physical process tag of search
    KnowledgeDomainGap  string    `json:"knowledge_domain_gap"`  // The domain boundary separating them
    ProposedQuery       string    `json:"proposed_query"`        // Cross-domain query to inject
    Confidence          float64   `json:"confidence"`            // Confidence this is a local minimum
    DetectedAt          time.Time `json:"detected_at"`
}
```

### 11.4 Scope Node Types

Scope nodes (§4.6) are first-class node types in v1.3. The Go types live alongside the InsightSignal types in package `contextus`.

```go
package contextus

import "time"

// ScopePhysical is a region of space-time with a name and queryable
// membership. See §4.6.2.
type ScopePhysical struct {
    ScopeID          string       `json:"scope_id"`
    Name             string       `json:"name"`
    Geometry         []byte       `json:"geometry"`              // GeoJSON bytes; opaque to the type
    ElevationRangeM  *[2]float64  `json:"elevation_range_m,omitempty"`
    TemporalRange    LocaleBounds `json:"temporal_range"`
    Grain            GrainSpec    `json:"grain"`
    ParentScopeID    string       `json:"parent_scope_id,omitempty"`
    Tags             []string     `json:"tags,omitempty"`
}

// ScopeConceptual is a topic or domain with a name and queryable
// membership. See §4.6.3.
type ScopeConceptual struct {
    ScopeID          string   `json:"scope_id"`
    Name             string   `json:"name"`
    OntologyURI      string   `json:"ontology_uri,omitempty"`
    ParentScopeID    string   `json:"parent_scope_id,omitempty"`
    RelatedScopeIDs  []string `json:"related_scope_ids,omitempty"`
    Tags             []string `json:"tags,omitempty"`
}

// ScopeMembership is the on-edge metadata for HE_SCOPE_MEMBERSHIP. See §4.6.4.
// This is the metadata payload; the edge itself is a Wyrd model.Hyperedge of
// arity 2 connecting the scope node to the member node.
type ScopeMembership struct {
    ScopeID       string    `json:"scope_id"`
    MemberID      string    `json:"member_id"`
    Since         time.Time `json:"since"`
    Confidence    float64   `json:"confidence"`     // 0.0–1.0
    ProvenanceTag string    `json:"provenance_tag"` // "P" (asserted) | "D" (inferred)
    Method        string    `json:"method"`         // adapter-defined; see §4.6.5
}
```

The membership predicate (§4.6.5) is implementation-defined per adapter. The `Method` field carries the adapter-specific predicate name (e.g. `"asserted"`, `"tag-overlap"`, `"embedding-cosine"`, `"citation-graph-hop"`). Spec v1.3 does not enumerate canonical method values; future revisions will record those that emerge from production adapters.

---

## 12. Related Documents

| Document | Version | Scope |
|---|---|---|
| Contextus Theory | v1.4 | Philosophical foundations, design principles, Insight Signal theory, worked failure (Colorado River), surveillance mode |
| Contextus Theory (forthcoming) | v1.5 | Formalisation of `AnomalyStructural` for Synthesis-promoted session findings (Edge Scout boundary blindspots, Corpus Edge Scout diversity findings); v1.3 ships the kind provisionally |
| Contextus ↔ Wyrd Integration Architecture | 2026-05-05 | Resolves Wyrd issue #6; closes the four open architectural questions from v1.2; inputs to this v1.3 |
| Contextus–CTH Bridge Specification | v0.1 | Trust Receipt, triage gate, heartbeat coupling, active evidence seeking |
| Contextus–CTH Bridge Synthetic Receipts | v0.1 | Parameter calibration using QBP and non-QBP test cases |
| CTH Engine | v2.0 | Confluent Trust Hypergraph evaluation framework |
| MuninnDB Specification | — | Hypergraph storage, Hebbian co-activation, Ebbinghaus decay |
| War Table Specification | v8 | Agent doctrine and spawning patterns |
| QBP Locale | v0.3 | Quaternion spatiotemporal addressing |
| QBP-CU Wyrd Integration | v0.2 | Canonical `SedenionResult` definition cited by §11.1 |
| qbp-compute-unit ADR-003 | — | §I4 design-doc-as-S-01-review-surface invariant under which v1.3 is reviewed |

---

*End of Specification v1.3. Implementation language: Go. Primary implementation targets: Pop!_OS 22.04, AMD FX-8350, PowerColor Red Devil RX 9070 XT. Reviewed under qbp-compute-unit ADR-003 §I4 (design-doc-as-S-01-review-surface).*
