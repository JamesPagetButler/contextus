# Contextus Theory Document — Reconstruction

**Reconstruction status:** PARTIAL — title page, §1 opening, §5 case studies, §7 design principles substantially recovered. §2-4, §6 are structural-only.
**Source chat:** https://claude.ai/chat/a0dc742b-cb62-4322-8cc1-3434114ddc43
**Original versions in source:** v1.0 (March 2026), v1.1, v1.2 (after restructure)

> NOTE: The original document was produced as a Word file with custom title page (dark teal/green palette, Cardo font). This is a markdown reconstruction. Where versions diverge, this document reflects v1.2 (the restructured version where ecology stops being privileged).

---

# PROJECT CONTEXTUS

## Theory Document

**A Framework for Ecosystem Comprehension**

**Version 1.2 | March 2026**
**Helpful Engineering**
**James Paget Butler**

*Classification: Open Source | Licence: TBD*

---

## 1. The Problem: Drowning in Data, Starving for Understanding

Ecology has a data problem, but it is not the problem most people assume. The challenge is not a shortage of measurements. Between satellite constellations, GPS collar networks, distributed sensor grids, citizen science platforms, and decades of field observations, we have more ecological data than any human or team of humans could review in a lifetime. The Yellowstone ecosystem alone generates terabytes of observational data annually across dozens of agencies, universities, and independent researchers.

The real problem is comprehension.

[v1.2 reframing applies here: the problem statement generalises beyond ecology. The structural limitation — pairwise joins, fixed time windows, single-domain search — is the same in any field where insight depends on the relationships between datasets rather than within them.]

### 1.1 The Comprehension Problem

Tools optimise for what they can compute, not for what researchers need to understand. Relational databases compute pairwise joins efficiently; they cannot represent the N-ary structure of an ecological hyperedge — wolf pack × geothermal feature × vegetation bloom × season — without flattening it into multiple two-table joins that fragment the very pattern the researcher is trying to see.

### 1.2 The Relationship Gap

[Substantial section on how analytical tools strip context to find causation, but in nature the outcome *is* the interaction pattern. The Womersley et al. (2025) whale shark paper is the motivating example: GAM with two variables is publishable but the researchers themselves know reality is higher-dimensional. The whale shark neonate doesn't care about Chl-a or DO independently — it cares about "am I fed and not eaten," which is a property of the conjunction of multiple gradients, predator physiology, bathymetry, and currents simultaneously.]

---

## 2. The Hypergraph as Substrate

[Reconstruction limited. Key points recovered:]

- MuninnDB hypergraph store; no engine changes needed
- 14 node types covering observations, locales, organisms, events, hypotheses, agents, provenance
- 9 hyperedge types covering observation links, agent inferences, hypothesis paths, provenance bindings
- HE_AGENT_INFERENCE edge with mandatory "D" (Derived) provenance tags — described as "the immune system against hallucination loops"
- Mandatory T/P/D/I provenance envelopes on every hyperedge
- The hypothesis IS a path through the graph; evidence is the density and confidence of edges

---

## 3. The Locale as Spatiotemporal Addressing

[Reconstruction limited. Key points recovered:]

- QBP Locale (Λ): quaternion-valued addressing unit
- Scalar part = proper time
- Vector part = spatial displacement
- An organism's trajectory through space-time is a path through Locale space
- Rendering math for the visualisation layer uses Locale directly rather than inventing a separate coordinate system; this keeps Contextus mathematically aligned with QBP

---

## 4. Visualisation Theory

[Section 4 was perception-native and domain-agnostic. Key elements recovered:]

### 4.4 Perceptual Constraints

- **Information density cap:** Maximum three active data layers at any time. Adding a fourth requires explicitly dismissing one. Hard constraint to prevent cognitive overload.
- **Animation speed:** Default playback with smooth transitions. Range adjustable. Transitions use ease-in-out curves, never abrupt jumps.
- **Colourblind safety:** All semantic colour encodings pass Okabe-Ito or equivalent validation. Shape and pattern are used as redundant channels alongside colour.

The three-layer visualisation cap is described in chat as "deliberately aggressive — the empathy constraint made architectural."

---

## 5. Case Studies: The Same Problem in Three Domains

> **Reconstruction status: NEAR-COMPLETE for §5.1-5.2 from retrieved chat content.**

Contextus is domain-agnostic. To demonstrate this, we present three case studies drawn from ecology, marine biology, and astrophysics. Each exhibits the same structural problem: insight latent in the relationships between datasets, invisible to tools that impose pairwise joins, fixed time windows, or single-domain search.

### 5.1 Yellowstone: The Trophic Cascade

The reintroduction of grey wolves to Yellowstone in 1995–96 is widely cited as a trophic cascade: wolves changed elk behaviour, which allowed willow and aspen recovery, which stabilised stream banks, which changed river morphology. The narrative is compelling but contested. The data to test it exists across multiple agencies and formats — wolf GPS tracks, elk surveys, vegetation transects, stream gauges, climate records — but no existing tool integrates them into a single queryable structure. Each study examines one or two links in the chain. Nobody has tested the full chain as a path through a unified relationship graph.

Contextus models the cascade as a hypergraph path and tests whether the temporal and spatial correlations along that path are consistent with a causal chain, an incidental correlation, or a mix. The hypothesis is a path through the graph. The evidence is the density and confidence of the edges.

A second Yellowstone hypothesis concerns the reported increase in grizzly bear visibility during the COVID-19 pandemic, when human visitation dropped dramatically. Testing this requires modelling confounding variables — mast crop years, hunting allocations, survey methodology changes — alongside human traffic data. The hypergraph makes confounders explicit as nodes and edges, preventing the researcher from inadvertently ignoring them.

### 5.2 Whale Sharks: The Confluence Problem

Womersley et al. (2025) studied neonate whale shark habitat, finding that neither sea surface temperature nor current velocity alone predicted neonate locations. What predicted them was the interaction of surface chlorophyll-a and dissolved oxygen at depth — the boundary condition between productivity and hypoxia.

The researchers used a generalised additive model (GAM) that could handle two-variable interactions, and they found a significant result. But the GAM is a tractable approximation of a higher-dimensional reality. The neonate's habitat selection is N-ary: it integrates Chl-a × DO × bathymetry × predator-physiology × current structure simultaneously. The GAM collapses this to pairwise. Contextus represents it as a single hyperedge.

### 5.3 GRB 250702B and the Ten-Second Window Problem

[Reconstruction limited. Key points recovered:]

The longest gamma-ray burst on record (~7 hours, detected July 2 2025) sits outside conventional GRB analysis pipelines because those pipelines use ten-second windows. Cross-correlation with LIGO O4 gravitational wave data — at the appropriate temporal grain — would test the QBP prediction of pulse-for-pulse temporal correlation. No existing tool supports the search at the right scale. Contextus's no-pre-specified-windows principle is exactly what's needed: the data reveals the structure; the tool does not impose it.

### Common Thread

[Closing line of §5, recovered verbatim:]

> "Contextus does not make researchers smarter. It stops the tools from making researchers blind."

---

## 6. Ethics and Data Sensitivity

[Reconstruction limited. Key points recovered:]

- Three-tier sensitivity model: PUBLIC (aggregated only), RESEARCHER (individual-level with temporal lag, requires institutional affiliation + DUA), STEWARD (real-time full-resolution, requires domain authority + approved protocol + audit logging)
- Sensitivity enforced at the MuninnDB query layer, not the API layer — every node carries a `sensitivity_tier` field
- Vulnerability disclosure protocol: poaching corridors, pollution pathways, etc. are auto-classified STEWARD with 90-day embargo

---

## 7. Design Principles

[Recovered with high confidence — these were the closing summary of v1.2.]

1. **Comprehension before generation.** The system optimises for what researchers need to understand, not for what tools can compute.
2. **Empathy as primary constraint.** Human perceptual and cognitive limits define the visualisation envelope. The information density cap is empathy made architectural.
3. **Hypothesis as path.** The space of askable questions determines the space of discoverable insights.
4. **Agents as attention, not authority.** AI directs the researcher's gaze; the researcher makes the judgments.
5. **No pre-specified windows.** Agents search across domains and timescales without assuming the correct correlation structure. The data reveals the structure; the tool does not impose it.
6. **Provenance is non-negotiable.** Every datum, inference, and hypothesis must be traceable to its source with explicit confidence.
7. **Perception-native visualisation.** Use motion, colour, and narrative to work with human perceptual strengths, not against them.
8. **Open tools, governed data.** The framework is open source. Sensitive data is access-tiered.
9. **Built on what exists.** MuninnDB, NATS, Podman, FastMCP, BMA-BRIDGE. Contextus is a domain layer, not a new platform.

> Note: Design Principle 10 ("Signals are evidence of patterns, not conclusions") was added in Theory v1.3 in chat 2. v1.2 had nine principles.

---

*End of v1.2. See Spec v1.0 for implementation details.*
*Subsequent revision (v1.3) extended the theory with §3.6 Insight Signal — see `02-bridge-and-signals/contextus-theory-v1.3-extract.md`.*
