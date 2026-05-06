# Contextus Specification v1.0 — Reconstruction

**Reconstruction status:** STRUCTURAL — section list and §1 mapping recovered with high confidence. Section bodies are partial.
**Source chat:** https://claude.ai/chat/a0dc742b-cb62-4322-8cc1-3434114ddc43

> NOTE: The fuller Spec v1.1 (in `02-bridge-and-signals/`) was reconstructed in chat 2 based on this v1.0 with new sections added. For implementation reference, prefer the v1.1 extract.

---

# PROJECT CONTEXTUS

## Technical Specification

**Version 1.0 | March 2026**
**Helpful Engineering**
**James Paget Butler**

*Classification: Open Source | Licence: TBD*

---

## 1. Integration Map: Contextus on the Existing Stack

Contextus is a domain-specific layer on infrastructure that already exists. This section maps every Contextus component to its existing counterpart in the BMA / Materia / War Table ecosystem.

### 1.1 Component Mapping

| Contextus Need | Existing Component | Integration Notes |
|---|---|---|
| Knowledge graph | MuninnDB | Hypergraph store. Contextus defines domain-specific node and edge types. No engine changes needed. |
| Message bus | NATS | Ingestion pipeline, agent communication, event-driven processing. Existing subjects extended with `ctx.*` prefix. |
| Container orchestration | Podman | Source adapters run as Podman containers on Pop!_OS. Same pattern as BMA Crawl-phase containers. |
| REST API | Go (port 8422) | Same pattern as Materia (8420). Contextus serves on 8422. Shared router library. |
| AI interface | FastMCP (port 8423) | MCP server exposing Contextus queries to AI agents. Same pattern as Materia (8421). |
| Multi-model access | BMA-BRIDGE | Four-layer transport (MCP, CLI, GitHub Issues, file drop). No changes needed. |
| Agent framework | War Table agents | Doctrine markdown + TOML sidecars. Agent spawning logic reused; doctrines are domain-specific. |
| Spatiotemporal math | QBP Locale | Quaternion-valued addressing. World-line rendering shares math with QBP visualisation. |
| Object storage | Local / S3-compatible | Heavy rasters, satellite imagery, large CSV archives. Indexed by MuninnDB metadata nodes. |

**New code required:** source adapters, domain schema definitions, visualisation renderer, hypothesis query patterns, agent doctrines.
**Infrastructure required:** zero new components.

---

## 2. Hypergraph Schema

[Reconstruction limited. Key points recovered:]

- 14 node types covering: Observation, Locale, Organism, Event, Hypothesis, Agent, Provenance, [+7 unrecovered]
- 9 hyperedge types covering: HE_OBSERVATION_LINK, HE_AGENT_INFERENCE (mandatory "D" tag), HE_HYPOTHESIS_PATH, HE_PROVENANCE_BIND, [+5 unrecovered]
- Mandatory T/P/D/I provenance envelopes (T=Telemetry, P=Published, D=Derived, I=Inferred — these correspond to existing BMA/Materia provenance tags)

> See Spec v1.1 §2.3 for the InsightSignal schema added in chat 2.

---

## 3. Ingestion Pipeline

[Reconstruction limited. Key points recovered:]

Five initial source adapters planned:
1. USGS stream gauges
2. GPS collar feeds
3. Satellite imagery indices (Sentinel/Landsat NDVI)
4. NPS visitor counts
5. [Fifth adapter — possibly weather stations or eDNA, unrecovered]

Pattern: each adapter normalises its source into observation messages on NATS. Enrichment service subscribes, performs geocoding and temporal alignment, tags provenance, writes to MuninnDB. Functionally identical to BMA's cognitive pipeline.

---

## 4. Agent Architecture

[Reconstruction limited. Key points recovered:]

**Five-agent taxonomy:**
1. **Source Adapters** — domain-aware ingestion (one per data source)
2. **Context Builder** — enrichment, geocoding, temporal alignment, hyperedge construction. Identified as the linchpin agent.
3. **Scout agents** — anomaly detection in subgraphs
4. **Correlation agents** — cross-domain pattern matching
5. **Synthesis agents** — narrative hypothesis generation from correlated patterns

**Pattern:** Each agent defined as doctrine markdown + TOML sidecar (reused from War Table). The agent IS a configuration; the runtime is shared.

---

## 5. REST and MCP Tool Catalogue

[Reconstruction limited. Key points recovered:]

- REST API on Go, port 8422
- FastMCP on port 8423
- Tool catalogue exposes Contextus queries to AI agents
- Pattern matches Materia (8420 / 8421)

> See Spec v1.1 §5.1-5.3 for the expanded endpoints (signals, bridge integration) added in chat 2.

---

## 6. Visualisation Spec

[Reconstruction limited. Key points recovered:]

- Three-layer maximum (hard cap)
- Confidence-as-opacity
- Okabe-Ito colourblind safety with redundant shape and pattern channels
- Ease-in-out animation curves, never abrupt jumps
- Perception-native rendering (motion, colour, narrative work with human perceptual strengths)

---

## 7. Data Sensitivity

[Recovered with high confidence:]

| Tier | Data Resolution | Access Requirement |
|---|---|---|
| PUBLIC | Aggregated only. Population estimates, regional trends, coarse maps. | None |
| RESEARCHER | Individual-level with temporal lag. Non-sensitive precise coordinates. | Institutional affiliation + data use agreement |
| STEWARD | Real-time, full-resolution. Sensitive locations. | Domain authority + approved protocol + audit logging |

**Enforcement at MuninnDB query layer**, not API layer. Every node carries `sensitivity_tier` field. Queries filtered by authenticated user's tier before results returned. Prevents accidental exposure through API bugs or new endpoints.

**Vulnerability disclosure:** Domain vulnerabilities (poaching corridors, pollution pathways) auto-classified STEWARD. Designated authorities notified via secure channel. 90-day embargo before reclassification.

---

## 8. Twenty-Week Phased Roadmap

[Reconstruction limited. Phase structure recovered:]

- Weeks 1-4: Phase 1 — Foundation (MuninnDB schema, first source adapter, basic REST/MCP)
- Weeks 5-8: Phase 2 — Multiple adapters + context builder agent
- Weeks 9-12: Phase 3 — Scout and correlation agents
- Weeks 13-16: Phase 4 — Synthesis agents + visualisation prototype
- Weeks 17-20: Phase 5 — Integration, scaling, additional adapters (literature, satellite, eDNA), cross-domain validation

> See Spec v1.1 §10 for the updated roadmap with Bridge integration and proportional retention added.

---

*End of Spec v1.0. See Theory v1.2 for philosophical foundations.*
*Updated to v1.1 in chat 2 — see `02-bridge-and-signals/contextus-spec-v1.1-extract.md`.*
