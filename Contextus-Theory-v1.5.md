# PROJECT CONTEXTUS

## Theory Document

*Picking Out the Flowers Amongst the Forest*

Version 1.5 | May 2026

Helpful Engineering

James Paget Butler

Classification: Open Source | Licence: TBD

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 1.0 | March 2026 | Initial theory document (ecology-focused) |
| 1.1 | March 2026 | Restructured as domain-agnostic framework |
| 1.2 | March 2026 | Three case studies (Yellowstone, whale sharks, GRB 250702B); ethical framework |
| 1.3 | April 2026 | Added §3.6 (Insight Signal); bridge integration references |
| 1.4 | April 2026 | Added §7 (Colorado River Test — worked failure); added §8 (Surveillance Mode, edge scouting, corpus scouting, bridge agent); updated Design Principles to 14 items (§9) |
| 1.5 | May 2026 | Formalised the *structural anomaly* — fourth anomaly class (§3.6.2) for findings that describe the structure of the search/exploration process rather than the structure of the data. Source mapping to Synthesis-promoted session-scoped agent output (Edge Scout boundary blindspots; Corpus Edge Scout diversity convergence). Added §3.6.6 (Synthesis as Persistence Boundary) cross-referencing Spec §4.4. Reframed the Colorado River failure (§7.6) as a worked example of a structural anomaly that an Edge Scout could have flagged. Provisional `AnomalyStructural` constant in Spec v1.3 §11.1 is no longer provisional after this revision. |

---

## 1. The Core Problem: Insight Lost in the Weave

The defining challenge of modern research is not a shortage of data. It is the inability to perceive the patterns that exist across datasets. Every scientific discipline now generates more measurements than any individual or team can review. The data is there. The insight is there. But the insight lives in the relationships between measurements — in the weave of the data, not in any single thread — and the tools we use to analyse data are designed to examine threads, not weaves.

This is not a problem of compute power or storage capacity. It is a problem of structure. The way we organise, query, and visualise data determines which questions we can ask. And the questions we can ask determine the insights we can find. When the tooling constrains the question space, it constrains the discovery space. Insights that require relating data across domains, across timescales, or across more than two variables simultaneously are structurally invisible to tools built on pairwise joins, fixed correlation windows, and single-domain pipelines.

Contextus is a framework for finding those invisible insights. It is not a database, not a dashboard, and not an AI that produces conclusions. It is a system designed to help a human being perceive patterns in interwoven data — to pick out the flowers amongst the forest.

### 1.1 The Structural Constraint

Most analytical tools are built on relational databases, which model the world as entities and pairwise relationships between them. This is a powerful abstraction for many purposes, but it imposes a fundamental limitation: every relationship must connect exactly two things. When the phenomenon under study involves three, five, or twenty things interacting simultaneously, the relational model forces the researcher to decompose it into pairs. The decomposition is not neutral. It fragments the very pattern the researcher is trying to see.

Similarly, most search and correlation tools operate within pre-specified windows — a time lag, a spatial radius, a domain boundary. These windows are necessary engineering choices, but they encode assumptions about where the interesting patterns live. When the interesting pattern lives outside the window, the tool will never find it, no matter how much data it processes or how fast it runs.

The result is a systematic bias toward the expected. Tools built on pairwise joins find pairwise relationships. Tools built on ten-second windows find ten-second correlations. Tools built within a single domain find single-domain patterns. The unexpected — the N-ary relationship, the decade-long lag, the cross-domain correlation — falls through the structural gaps.

### 1.2 The Human in the Loop

The goal of Contextus is not to automate insight. It is to restore the human researcher's ability to perceive patterns that the tooling has rendered invisible. Human cognition has extraordinary strengths: we detect motion against backgrounds, perceive subtle colour gradients, recognise narrative arcs, and intuit spatial relationships. These are not incidental skills; they are the product of millions of years of evolutionary pressure. A system that works with these strengths can make a researcher orders of magnitude more effective.

Conversely, we are poor at holding large numbers of variables in working memory, correlating across long time lags, and searching exhaustively across domains. A useful system compensates for these weaknesses by doing the exhaustive search, maintaining the full variable set, and presenting the results in forms that human perception can process naturally. The system searches. The human sees. The insight emerges in the seeing.

This commitment to the human as the locus of understanding is not a philosophical nicety. It is the primary design constraint. Every architectural decision in Contextus is evaluated against one question: does this help a prepared mind perceive a pattern it could not perceive before?

---

## 2. The Theoretical Framework

### 2.1 The Hypergraph: Representing Interwoven Relationships

A hypergraph is a generalisation of a graph in which an edge can connect any number of nodes simultaneously. Where a standard graph edge connects exactly two nodes, a hyperedge connects an arbitrary set. This mathematical structure directly addresses the pairwise limitation of relational databases: an event involving five entities, three conditions, and two timescales is represented as a single hyperedge connecting all of them, not as a combinatorial explosion of pairwise joins.

In Contextus, every observation, relationship, hypothesis, and inference is a structure within the hypergraph. Nodes represent entities (an individual organism, a geographic feature, a sensor, a data source, a temporal period). Hyperedges represent relationships between entities, and carry their own metadata: provenance, confidence, temporal bounds, observational context, and the method by which the relationship was established.

The hypergraph is not a metaphor. It is the literal data structure. A research question becomes a graph traversal. A hypothesis becomes a path. The evidence for a hypothesis is the density and confidence of the edges along that path. The absence of evidence is the absence of edges. The system does not answer questions; it makes the structure of evidence navigable.

### 2.2 The Locale: Unified Spatiotemporal Addressing

Every observation occurs somewhere and somewhen. Contextus adopts a unified spatiotemporal addressing scheme drawn from quaternion algebra: the Locale. A Locale is a quaternion-valued unit in which the scalar component encodes time and the vector component encodes spatial position. This is not an arbitrary mathematical choice; it reflects the physical inseparability of space and time in any system where entities move.

An entity's trajectory through a domain over time is a continuous path through Locale space: a world-line. Two entities whose world-lines converge in Locale space had a potential interaction. A feature whose Locale is spatially fixed but temporally periodic creates a recurring attractor. The mathematics of proximity, intersection, and divergence in Locale space are the mathematics of interaction, and they work identically whether the entities are wolves, whale sharks, or photons.

### 2.3 Landscape and Exploration

The hypergraph has two distinct aspects. The topology — which nodes exist and how they are connected — is the landscape. It describes what relationships are possible. The actual flow of entities, energy, and information through that topology is the exploration. It describes what relationships are active at any given time.

This distinction matters because different questions interrogate different aspects. Questions about structure ("what are the possible interaction pathways?") are queries against the landscape. Questions about dynamics ("what is happening now?") are queries against the exploration. Questions about change ("what shifted after this intervention?") are queries about how the exploration pattern changed in response to a landscape modification. At low disturbance, entities explore a narrow subset of available relationships. As disturbance increases, higher-energy interaction modes activate. Understanding which relationships are structural (persistent under perturbation) and which are contingent (activated or suppressed by disturbance) is a central analytical goal.

### 2.4 Empathy as Architecture

The deepest design commitment in Contextus is empathetic: the system exists to help a human being understand something complex. A feature that increases data throughput but makes the display harder to read is a net negative. A simplification that loses nuance but makes a pattern visible is a net positive, provided the simplification is transparent and the full data remains accessible.

This is not a soft requirement. It is the primary design constraint. The system's value is zero if it cannot improve a competent researcher's comprehension of the domain they are studying. Practically, this means:

- **Progressive disclosure.** The default view shows the minimum information needed to orient the researcher. Detail is available on demand, never forced.

- **Perceptual honesty.** Visualisations must not exaggerate, distort, or aestheticise data in ways that create false impressions. Beauty is welcome; deception is not.

- **Cognitive pacing.** Animations and transitions match human perceptual processing speed. Fast enough to feel responsive; slow enough to be legible.

- **Narrative support.** Human insight often takes narrative form. The system supports narrative exploration of the hypergraph, not just statistical querying.

---

## 3. Agent-Mediated Insight Discovery

### 3.1 The Attention Problem

Even with a well-structured hypergraph and empathetic visualisation, any real dataset contains more observations than a researcher can attend to. The system needs a way to direct human attention toward the parts of the data most likely to yield insight. This is where AI agents enter — not as analysts who produce conclusions, but as attentional guides who surface anomalies, suggest relationships, and maintain awareness of the full dataset while the human focuses on specific questions.

### 3.2 Agents as Data, Not Code

The agent architecture draws from the War Table model: agents are defined by doctrine documents (natural language descriptions of analytical perspective and priorities) and parameter sidecars (operational configuration). They are spawned dynamically based on the current situation, not pre-compiled into the system. This makes the agent ensemble adaptive — the system can deploy different analytical perspectives depending on what the data is showing, without code changes.

- **Scout agents** continuously traverse the hypergraph looking for anomalies: unexpected density changes, broken correlations, temporal patterns that deviate from norms. They do not interpret; they flag.

- **Correlation agents** test specific hypothetical relationships by querying the hypergraph for supporting and contradicting evidence, reporting the balance and flagging confounders.

- **Synthesis agents** combine outputs from multiple scouts and correlators to propose narrative explanations. Their outputs are always flagged as hypotheses, never conclusions.

- **Provenance agents** track data lineage and confidence. Every claim must be traceable to specific observations with explicit confidence. This is the immune system against hallucination.

- **Edge scouts** monitor the boundary between the researcher's currently explored subgraph and the unexplored remainder of the hypergraph. They operate concurrently with researcher queries, flagging significant unexplored connections from nodes the researcher has visited. (§8.2)

- **Corpus edge scouts** monitor the diversity of search results returned to the researcher or to querying agents, detecting domain convergence and injecting broadening queries before results are delivered. (§8.5)

- **Bridge agents** maintain a concurrent view of surveillance flags and active search trajectories, detecting when the two are converging on the same physical process from different domain perspectives — the signature of a local minimum in insight space. (§8.6)

### 3.3 Cross-Domain Search Without Pre-Specified Windows

The most important capability agents provide is searching across domain boundaries without pre-specified correlation windows. A scout agent does not know that two signals "should" be correlated within ten seconds rather than ten hours. It has no prior about the correct time window. It traverses the hypergraph looking for co-occurrence patterns at whatever timescale the data presents. It does not need to be told the lag; it discovers the lag from the data.

This is the direct architectural answer to the structural constraint described in Section 1. The fixed time window, the pairwise correlation, the domain boundary — none of these exist in the hypergraph. The agent searches the structure that actually exists in the data, not the structure that a previous researcher assumed would exist.

### 3.4 Metadata Construction from Raw Data

Much raw data arrives without rich metadata: a GPS fix with a timestamp, a sensor reading, an image. Agents construct metadata that transforms raw observations into knowledge:

- **Behavioural inference:** A sequence of position fixes can be segmented into behavioural states (resting, travelling, hunting) based on speed, turning angle, and spatial pattern.

- **Interaction detection:** When two entities' world-lines converge in Locale space below a threshold, the agent creates a potential-interaction hyperedge.

- **Environmental context linking:** A bare coordinate is enriched with contextual data — terrain type, elevation, nearest features, current conditions — drawn from other datasets in the hypergraph.

- **Temporal pattern detection:** Long time series are monitored for regime shifts, periodicity changes, and trend breaks invisible in short-window analysis.

### 3.5 The Human-Agent Relationship

The agents are not autonomous analysts. They are more like a team of very attentive research assistants: they notice things, they flag things, they test specific hypotheses when asked, and they maintain awareness of parts of the dataset the researcher is not currently examining. But the intellectual judgment — what matters, what it means, and what to do about it — remains with the human. The system is most valuable when it amplifies existing expertise rather than attempting to replace it.

### 3.6 The Insight Signal: What Agents Produce

Sections 3.1 through 3.5 describe what agents do: they traverse the hypergraph, detect anomalies, test correlations, propose narratives, and track provenance. But they do not define what an agent *produces* when it finds something worth flagging. Without a formal definition of this output, the system has a process but not a product. This section defines that product: the Insight Signal.

#### 3.6.1 Why a Formal Signal Matters

An agent that detects an anomaly must communicate that detection to the rest of the system — to other agents, to the visualisation layer, and ultimately to the human researcher. If this communication is informal (a log entry, a highlighted node, a notification), then insights are ephemeral. They exist only in the moment of display. They cannot be compared, tracked over time, combined with other insights, or evaluated for reliability. They are sparks, not records.

Contextus requires insights to be records. Not conclusions — the theory is explicit that interpretation remains with the human — but structured descriptions of what an agent observed, where and when it observed it, how surprising the observation was, and how confident the agent is in the observation's validity. A record can be stored, decayed, strengthened by subsequent evidence, and eventually promoted or dismissed. A spark cannot.

The Insight Signal is this record. It is the atomic unit of agent output. Every time an agent flags something — an anomaly, a correlation, a narrative thread — it emits an Insight Signal into the hypergraph. The signal is itself a node in the hypergraph, connected to the nodes and hyperedges that constitute the pattern it describes. Insights about the data become part of the data. This reflexivity is deliberate: it allows agents to discover insights about insights, building layered understanding over time.

#### 3.6.2 Anatomy of an Insight Signal

An Insight Signal is composed of six elements, each derived from structures the theory has already defined.

**The Subgraph.** The specific set of nodes and hyperedges the agent flagged. This is the pattern itself — the flower the agent picked out from the forest. The subgraph is not the agent's interpretation of the pattern; it is the pattern. It is a pointer into the hypergraph, not a copy of the data.

**The Traversal Path.** The sequence of hypergraph steps the agent took to arrive at the subgraph. This is provenance for the pattern — it records how the agent found it, which matters because different traversal paths to the same subgraph carry different epistemic weight. A scout that stumbled onto the pattern during random traversal and a correlation agent that was directed to test a specific hypothesis may flag the same subgraph, but the directed discovery provides stronger evidence than the undirected one. The traversal path also enables the system to detect when multiple agents independently converge on the same pattern — a convergence signal that is itself significant.

**The Anomaly Characterisation.** What makes this subgraph interesting. The theory defines four classes of anomaly. The first three are direct outputs of the global agent types named in §3.2 (Scout, Correlation, Synthesis). The fourth is a derived class — an output of Synthesis when Synthesis is acting as the *persistence boundary* for findings raised by session-scoped agents (Edge Scout, Corpus Edge Scout, Bridge Agent; see §3.6.6 below).

- *Density anomaly* (from scout agents): the subgraph exhibits unexpected density — more or fewer connections, stronger or weaker edge weights, or faster or slower temporal dynamics than the surrounding hypergraph baseline. Something changed, or something is different from what was expected.

- *Correlation anomaly* (from correlation agents): the subgraph exhibits a statistical relationship — positive, negative, or conditional — that was not previously recorded or that contradicts a previously recorded relationship. The agent is not claiming causation; it is reporting co-variation.

- *Narrative anomaly* (from synthesis agents): the subgraph forms a coherent causal or temporal chain — a sequence of events or conditions that, taken together, tell a story. Epistemically the weakest of the data-shaped anomalies, but often the most valuable for directing human attention, because narrative is the format in which human insight most naturally operates.

- *Structural anomaly* (from synthesis agents acting as persistence boundary): the subgraph describes the *structure of the search or exploration process* itself, not the structure of the underlying data. A boundary connection a researcher's queries have systematically avoided; a search-result distribution that has converged into a single domain; a topological connection across a knowledge-domain boundary that has no implied causal narrative. Structural anomalies are the persistent record of what surveillance mode (§8) noticed: they let a fleeting observation about the process become a durable finding without crossing the "Contextus detects, never evaluates" line. They are epistemically distinct from the first three because they are observations *about how the system was looking*, not *about what the system found*.

Each class carries a score quantifying how surprising the observation is relative to the agent's baseline expectations. This score is not a p-value; it is an information-theoretic measure of how much the observation would change an informed observer's beliefs.

The four classes are not fully orthogonal. A bridge-agent convergence promoted as a Synthesis-authored signal often reads as both narrative (it tells a cross-domain story) and structural (it bridges a knowledge-domain boundary). When this happens, the agent picks the class that better serves the consumer: narrative when the cross-domain story is the load-bearing part; structural when the *fact of bridging* is the load-bearing part. The choice is recorded with the signal so consumers know the agent's framing without overriding it.

**The Locale Envelope.** The spatiotemporal bounds within which the pattern holds. Derived directly from the Locale framework (§2.2), the envelope defines where and when the insight is valid. The envelope carries an additional property: its *grain* — the finest resolution at which the pattern is detectable. A trend visible in annual averages but not monthly data has coarse temporal grain. Grain matters because it constrains what evidence can confirm or deny the insight.

**The Persistence Measure.** How many MuninnDB consolidation cycles the pattern has survived. A signal that fires once and never again may be an artifact. A signal that persists across multiple consolidation cycles, during which Ebbinghaus decay would naturally attenuate unsupported connections, is more likely to reflect a real pattern.

Persistence interacts with the landscape/exploration distinction from §2.3. A structural pattern (one embedded in the topology of the hypergraph) will persist indefinitely. A contingent pattern (one that depends on specific conditions being active) will appear and disappear as those conditions come and go. Both are real patterns; they differ in kind, not in quality.

**The Confidence Estimate.** The agent's Bayesian confidence in the signal's validity, derived from the provenance chain (§3.2). Confidence incorporates data quality, source independence, temporal coverage, and known confounders. It is explicitly *not* confidence in the signal's importance or meaning — those are human judgments. Confidence is reported as a distribution, not a point estimate.

#### 3.6.3 Signal Lifecycle

An Insight Signal is not a static record. It evolves as new data enters the hypergraph.

**Emission.** An agent detects a pattern and emits a signal. At emission, persistence is zero and confidence reflects only the initial observation.

**Strengthening.** Subsequent observations consistent with the pattern increase confidence and persistence. If multiple independent agents converge on the same subgraph, the signal gains convergence weight. This is the Hebbian principle applied at the insight level.

**Decay.** In the absence of new supporting evidence, the signal's activation decays according to MuninnDB's Ebbinghaus schedule. The signal does not disappear; its activation approaches a floor determined by the strength of its initial evidence.

**Mutation.** New data may change the shape of the pattern. When mutation occurs, the signal records the change as a version transition, preserving the history of how the pattern evolved.

**Absorption.** When two signals describe overlapping subgraphs and a synthesis agent determines they are aspects of the same underlying pattern, the signals merge. The weaker signal is absorbed into the stronger one.

**Promotion.** When a signal's confidence and persistence exceed a threshold, it becomes a candidate for external evaluation — a potential seed for structured investigation by systems outside Contextus. Promotion creates an export event. The promoted signal retains its lifecycle within Contextus regardless of what happens downstream. The relationship between the Insight Signal and external evaluation frameworks is specified in the Contextus–CTH Bridge Specification.

#### 3.6.4 Relationship to Human Attention

The Insight Signal is not a conclusion, a recommendation, or an alert. It is a structured description of something an agent noticed. The visualisation layer (§4) renders active signals as perceptual cues — subtle colour shifts, motion changes, density variations — that guide the researcher's attention without demanding it.

The researcher engages with signals by investigating the subgraph they describe, not by reading the signal itself. The signal says "look here." The looking is the researcher's job. The understanding that results from looking is the insight. The signal is the finger pointing at the moon; it is not the moon.

#### 3.6.5 Design Constraints on the Signal

Three constraints preserve the integrity of the Insight Signal:

**Signals are never conclusions.** An agent may report "these five nodes co-activate anomalously" but never "these five nodes co-activate because of X." Causal claims require human judgment. The signal structure has no field for causal explanation. This is a deliberate omission, not an oversight.

**Signals are always traceable.** Every element of the signal — the subgraph, the traversal path, the anomaly score, the confidence — is derived from specific data in the hypergraph and can be audited.

**Signals are never suppressed.** The system may deprioritise a signal for display (low confidence, low persistence), but it never deletes one. The Ebbinghaus decay ensures that genuinely unsupported signals fade naturally; forced suppression is not necessary and could hide valid patterns.

#### 3.6.6 Synthesis as Persistence Boundary

The session-scoped agents introduced in §8 — Edge Scout, Corpus Edge Scout, Bridge Agent — produce ephemeral findings. They flag boundary connections, search-distribution convergences, and cross-domain local minima while a session is active, and their output is rendered as visualisation cues in real time. When the session ends, those findings end. They do not become Insight Signals; they leave no record in the hypergraph.

This is a deliberate constraint, not an oversight. Spec §4.5 captures it as "session state is not persisted." The reasoning is twofold:

- **Storage discipline.** Surveillance-mode flags are noisy by design — Edge Scout flags every significant unexplored connection, not only the consequential ones. Persisting all of them would inflate the hypergraph with low-quality records and accelerate the trip toward the §9 Storage Sentinel thresholds. The cost-benefit is wrong.

- **Epistemic discipline.** A flag is the system saying *"look here"* — an attentional cue, not a claim. A persistent record in the hypergraph carries weight: other agents will encounter it, future researchers will inherit it, downstream consumers will see it as part of the corpus's knowledge. Promoting an attentional cue to a hypergraph record requires a deliberate act of judgment. That judgment belongs to a Synthesis agent.

This generates a structural pattern: the **persistence boundary**. Session-scoped findings flow on NATS subjects. A Synthesis agent subscribes. When a finding's significance, confidence, or pattern stability crosses a configured threshold, Synthesis mints an Insight Signal — an `agent_type = synthesis` signal whose anomaly class is *structural* (for boundary blindspots and diversity convergences) or *narrative* (for Bridge Agent convergences that imply a cross-domain causal story). The session-scoped agents themselves never write to the hypergraph.

The architectural consequence is that there is exactly one gate — the Synthesis agent's persistence-boundary subscription — at which ephemeral surveillance findings can become durable. The boundary is auditable (every minted signal records its source subject and finding ID in `traversal_path` metadata), tunable (thresholds are configuration, not hardcoded), and reversible (a signal demoted by Ebbinghaus decay returns to the inactive tail without explicit suppression). Spec v1.3 §4.4 contains the operational specification; this section is the theoretical justification.

The persistence boundary explains why the structural anomaly class exists. Without it, the system would have either to forget surveillance findings entirely or to admit them indiscriminately. Neither is acceptable. Synthesis-as-persistence-boundary is the third option: surveillance findings persist *only when an agent has decided they merit it*, and the resulting record carries an anomaly class that names what they are — observations about the search process, not the data.

---

## 4. Visualisation Theory: Working With Human Perception

### 4.1 Motion as Information Channel

Human visual systems are extraordinarily sensitive to motion. We detect moving objects faster and with less conscious effort than static ones. Contextus uses motion as a primary information channel:

- **World-line playback:** Entity trajectories rendered as moving points along their world-lines. Motion itself communicates velocity, directedness, and state.

- **Breathing heat maps:** Variables rendered as slowly pulsing gradients rather than static overlays. Pulse amplitude proportional to rate of change — stable regions pulse gently, rapidly changing regions pulse visibly.

- **Flow fields:** Aggregate movement patterns rendered as particle flows, similar to weather visualisation wind maps.

### 4.2 Colour as Semantic Layer

Colour is used semantically, not decoratively. The palette is designed for long viewing sessions, colourblind accessibility, and intuitive domain mapping:

- **Temporal gradient:** A slow shift through a colour spectrum provides continuous temporal context without requiring the researcher to check a timeline.

- **Confidence encoding:** Opacity maps to confidence. Low-confidence inferences appear as ghosts — visible but clearly tentative. Below a threshold, elements render as dotted outlines.

- **Anomaly highlighting:** Agent-flagged anomalies (Insight Signals with high anomaly scores) are indicated by subtle colour shifts, not jarring alerts. Peripheral vision detects the shift; focused attention investigates.

- **Edge scout boundary:** Nodes at the boundary of the researcher's explored subgraph that have significant unexplored connections are indicated by a distinct peripheral hue shift — the system's way of saying "there is more here." The shift scales with the statistical significance of the unexplored connections. (§8.2)

### 4.3 Narrative Navigation

The most powerful visualisation mode is narrative: following a causal chain through the hypergraph as if reading a story. The researcher selects a starting point and the system renders connected events as a navigable timeline, with branching paths and dead ends. This is not a pre-computed animation but an interactive traversal rendered in real time.

### 4.4 Perceptual Constraints

- **Information density cap:** Maximum three active data layers at any time. Adding a fourth requires explicitly dismissing one. This is a hard constraint to prevent cognitive overload.

- **Animation speed:** Default playback with smooth transitions. Range adjustable. Transitions use ease-in-out curves, never abrupt jumps.

- **Colourblind safety:** All semantic colour encodings pass Okabe-Ito or equivalent validation. Shape and pattern are used as redundant channels alongside colour.

---

## 5. Case Studies: The Same Problem in Three Domains

Contextus is domain-agnostic. To demonstrate this, we present three case studies drawn from ecology, marine biology, and astrophysics. Each exhibits the same structural problem: insight latent in the relationships between datasets, invisible to tools that impose pairwise joins, fixed time windows, or single-domain search.

### 5.1 Yellowstone: The Trophic Cascade

The reintroduction of grey wolves to Yellowstone in 1995–96 is widely cited as a trophic cascade: wolves changed elk behaviour, which allowed willow and aspen recovery, which stabilised stream banks, which changed river morphology. The narrative is compelling but contested. The data to test it exists across multiple agencies and formats — wolf GPS tracks, elk surveys, vegetation transects, stream gauges, climate records — but no existing tool integrates them into a single queryable structure. Each study examines one or two links in the chain. Nobody has tested the full chain as a path through a unified relationship graph.

Contextus models the cascade as a hypergraph path and tests whether the temporal and spatial correlations along that path are consistent with a causal chain, an incidental correlation, or a mix.

A second Yellowstone hypothesis concerns the reported increase in grizzly bear visibility during the COVID-19 pandemic, when human visitation dropped dramatically. Testing this requires modelling confounding variables — mast crop years, hunting allocations, survey methodology changes — alongside human traffic data.

### 5.2 Whale Sharks: The Confluence Problem

Womersley et al. (2025) studied neonate whale shark habitat, finding that neither sea surface temperature nor current velocity alone predicted neonate locations. What predicted them was the interaction of surface chlorophyll-a and dissolved oxygen at depth — the boundary condition between productivity and hypoxia.

The researchers used a generalised additive model (GAM) that could handle two-variable interactions, and they found a significant result. But they explicitly acknowledged that the real interaction space was higher-dimensional — bathymetry, currents, predator physiology, and reporting bias all play roles — and their statistical tools could not model the full conjunction.

A neonate whale shark does not care about chlorophyll-a. It does not care about dissolved oxygen. It cares about "am I fed and not eaten," which is a property of the simultaneous conjunction of multiple environmental gradients. A hypergraph represents this conjunction as a single hyperedge. A GAM cannot.

### 5.3 Black Holes: The Hidden Correlation

In March 2026, an analysis of GRB 250702B — a seven-hour gamma-ray burst with three distinct pulses — produced a prediction about temporal correlation between gravitational wave and electromagnetic signals from black hole consumption events. The prediction was specific: pulse-level synchrony between the two channels, arising from a shared boundary process.

The prediction was immediately testable on existing archival data. The LIGO O4 run had accumulated approximately 250 gravitational wave candidates. Fermi had catalogued ultra-long gamma-ray bursts over the same period. Both datasets were public. But the most comprehensive archival cross-correlation (Wang et al., 2022) had searched with a ten-second time lag — because the standard theory predicts seconds-scale correlations for merger events. If consumption events produce hour-long correlated structure, a ten-second window will never find it.

The data existed. The question had not been asked. The tool could not ask it.

### 5.4 The Common Thread

In all three cases:

- **The data exists** across multiple, well-curated, publicly accessible datasets.

- **The insight lives in the relationship between datasets,** not within any single one.

- **Existing tools impose structural constraints** (pairwise joins, fixed time windows, domain boundaries) that are correct for the questions they were designed to answer but that actively suppress the novel pattern.

- **A different way of asking the question** immediately reveals that existing data could answer a question nobody had posed.

Contextus does not make researchers smarter. It stops the tools from making researchers blind.

---

## 6. Ethical Framework: Open Source and Discretion

### 6.1 The Sensitivity Problem

Cross-domain insight tools are powerful, and power requires responsibility. High-resolution animal location data is a poaching risk. Medical data cross-referenced with genomics raises privacy concerns. Vulnerability discovery in any domain — ecological, infrastructural, social — can be exploited by bad actors. Contextus must implement sensitivity tiering that is structural, not policy-based: the system architecture itself prevents inappropriate access.

### 6.2 Tiered Access

- **Public tier:** Aggregated data only. Population estimates, regional trends, coarse-resolution maps. No individual tracks. No precise coordinates for sensitive features.

- **Researcher tier:** Individual-level data with temporal lag. Precise coordinates for non-sensitive features. Requires institutional affiliation and data use agreement.

- **Steward tier:** Real-time, full-resolution data. Limited to domain authorities and approved protocols. Audit-logged.

### 6.3 Open Source Boundaries

The Contextus codebase, schema, and analytical tools are open source. Domain-specific datasets are not. This distinction is critical: anyone can deploy Contextus for their own domain, but sensitive data is governed by the access tiers.

### 6.4 Vulnerability Disclosure

When an agent or researcher identifies a domain vulnerability — a poaching corridor, a pollution pathway, a system weakness — the finding is automatically classified at the steward sensitivity level. Designated authorities are notified via a secure channel. An embargo period applies before reclassification. The provenance agent logs the full disclosure chain.

---

## 7. A Worked Failure: The Colorado River Test

In April 2026, we ran a deliberate test of the Contextus thesis using the Colorado River system as a case study. The test was designed to determine whether an analytical process guided by the Contextus principles would discover a novel cross-domain insight that had recently been published in the scientific literature. The test failed in an instructive way. This section documents the failure and what it teaches about the system's architecture.

### 7.1 The Setup

The Colorado River is one of the most data-rich, multi-domain, and politically contested natural systems on Earth. It serves 40 million people across seven U.S. states and Mexico, and is governed by a century-old legal framework (the 1922 Colorado River Compact and subsequent "Law of the River") that over-allocated the river from day one. Since 2000, the system has been in persistent deficit, with total storage falling from near capacity to approximately 37% by April 2026.

The test proceeded in stages. First, a comprehensive baseline was assembled covering the river's physical geography, legal framework, infrastructure, ecology, and current status. Then, the question was posed: "Is there missing water? Where is it going?"

### 7.2 What the Analysis Found

The analysis identified five mechanisms draining the system that sit outside the management framework's field of view:

- **Groundwater mining:** GRACE satellite data shows 65% of the basin's total water storage loss since 2002 is groundwater depletion, not reservoir decline. Groundwater is unregulated by the interstate compact.
- **Dust on snow:** Anthropogenic dust loading darkens mountain snowpacks, accelerating melt by up to six weeks and reducing annual runoff by approximately 5%. No snowmelt forecast model accounts for dust.
- **Wildfire watershed destruction:** Burned soils lose infiltration capacity, converting slow percolation into flash runoff and sediment loading. Fire seasons are lengthening as the drought deepens.
- **Warming-driven evapotranspiration:** Each degree Celsius of warming reduces discharge by approximately 9.3%. Warmer forests in the headwaters consume more water, reducing groundwater recharge.
- **Irrigation method paradox:** Government-subsidised centre-pivot sprinklers, marketed as conservation, produce zero return flow. System-level efficiency may have decreased even as field-level efficiency increased.

These five mechanisms were identified through systematic cross-domain search and represent genuine analytical value. The feedback loops between them — drought exposes soil, dust accelerates melt, less water means more drought — are real and largely absent from the management framework's models.

### 7.3 What the Analysis Missed

A paper published simultaneously (Hogan et al., 2026, *Geophysical Research Letters*) had identified the dominant mechanism: vegetation interception of snowmelt. The finding is that warmer, drier springs cause plants to consume snowmelt before it reaches rivers, and this single mechanism explains approximately 70% of the gap between snowpack-predicted and actual streamflow since 2000. A companion Princeton study found the "drought paradox": plants maintain or increase transpiration during drought by switching from soil moisture to groundwater, competing directly with the river for the same water.

The analysis had the components of this insight. It cited headwater groundwater loss from increased forest water use. It discussed warming-driven evapotranspiration. But it treated vegetation as a passive loss channel — a coefficient in the energy balance — rather than as an active biological agent competing for water. It modelled biology as physics rather than as ecology. The critical mechanism — plants switching water sources in response to drought, intercepting snowmelt before it becomes streamflow — was structurally invisible to an analysis framed as "where is the water going?" because the water was never "going" anywhere. It was being intercepted before it arrived.

### 7.4 Root Cause Analysis

Three structural reasons explain the failure:

#### 7.4.1 Search by Mechanism, Not by Anomaly

The analysis searched for known loss mechanisms: things that remove water from the system. It did not search for unexplained discrepancies in the data. The snowpack-to-streamflow forecasting gap — good snowpack but less water than predicted, consistently since 2000 — is an anomaly sitting in publicly available data. A scout agent monitoring the hypergraph would have flagged this broken correlation without being asked. The analysis failed because it searched for causes rather than for unexplained patterns.

#### 7.4.2 Biology Modelled as Physics

Vegetation was treated as a passive variable — a loss term in the energy balance that scales with temperature. The Hogan paper shows vegetation is a strategic actor that adapts its behaviour to conditions: when soil moisture drops, plants switch to snowmelt; when snowmelt is unavailable, they tap groundwater. This is ecological behaviour, not physical evaporation. In a properly constructed hypergraph, plants would be nodes with behavioural edges — entities that change their interaction patterns in response to environmental state — not just coefficients.

#### 7.4.3 Inherited Domain Boundaries

The search terms used to investigate the problem were domain-specific: "water budget," "evaporation losses," "groundwater depletion." These are hydrology terms. The insight lives at the intersection of hydrology, plant physiology, and climate science. The analysis reproduced exactly the siloed thinking that Contextus is designed to overcome. It found what hydrologists find, because it searched the way hydrologists search.

### 7.5 The Role of User Framing

A fourth cause is more subtle and has significant implications for system design. The question that initiated the analysis was "Is there missing water? Where is it going?" This framing treats the problem as a loss accounting exercise. It naturally drives toward subtractive explanations: mechanisms that remove water from the system. The Hogan paper's insight is not about water being removed. It is about water being intercepted before it enters the system. The accounting frame cannot find interception, because in an accounting frame the water was never "there" to go missing.

A better question would have been "Why do snowpack measurements no longer predict streamflow?" This frames the problem as a broken correlation rather than a missing quantity, and points directly at the interception mechanism. But a researcher approaching the Colorado River from the existing management context — the Law of the River, the 24-Month Study, the allocation disputes — will naturally frame questions in allocation terms. Their expertise defines their framing, and their framing constrains their discovery space.

This is not a criticism of the researcher. It is a structural observation about how expert knowledge interacts with analytical tools. The researcher's expertise is simultaneously their greatest asset and their greatest liability: they know where to look, which means they also know where not to look. A system designed to augment expert comprehension must account for the fact that expert framing will systematically bias the question space.

### 7.6 The Failure as a Worked Example of a Structural Anomaly

The Colorado River failure (§7.3) is the canonical worked example for the *structural anomaly* class formalised in §3.6.2. It is worth tracing the connection explicitly because the failure motivated the class.

In the analysis as it actually ran, the researcher's queries touched hydrology nodes (streamflow, reservoirs, snowpack), climate nodes (temperature, precipitation), and infrastructure nodes (dams, diversions). They did not touch ecology nodes — vegetation phenology, plant-physiology, soil biology — even though those nodes had hundreds of strong connections to the queried set. An Edge Scout running concurrently with those queries would have flagged the unexplored boundary repeatedly: "Your active subgraph includes snowpack and streamflow nodes. These have 847 connections to vegetation phenology nodes that you have not explored. Three of these connections show statistically significant changes since 2000."

In the system as it ran, those flags would have been visualisation cues only. They would have appeared as peripheral hue shifts on the boundary nodes and faded when the session ended. The researcher might or might not have noticed; if they did, they might or might not have followed the flags into the ecology subgraph. Either way, no record would persist.

In the system Spec v1.3 specifies and Theory v1.5 formalises, when an Edge Scout flag's significance crosses the persistence-boundary threshold (§3.6.6), a Synthesis agent subscribed to `ctx.edge.boundary.*` mints an Insight Signal of *structural* anomaly class. The signal records: the boundary node, the unexplored domain, the significance of the unexplored connections, the session in which the flag fired. It is not a claim that the unexplored domain *is* the answer; it is a record that the system noticed the researcher's query pattern was not exploring it.

The epistemic value of the signal is this: when the Hogan et al. paper subsequently appeared and the gap between snowpack-predicted and actual streamflow turned out to be vegetation interception, the system's record of the structural anomaly would have provided three things:

1. A retrospective audit trail showing that the boundary had been visible during the analysis, even though it had not been pursued.
2. A forward-looking attentional cue — any subsequent researcher querying the same hydrology subgraph would inherit the prior structural anomaly as a peripheral hue shift, raising the chance that their session would explore the boundary even if their framing was the same as the original researcher's.
3. A diagnostic signal at the system level: a region of the hypergraph where structural anomalies accumulate without conversion into other anomaly types is a region where the framing bias is repeatedly suppressing genuine cross-domain insight. That region itself becomes a Scout target.

None of this requires the system to have known that vegetation interception was the answer. It only requires the system to have noticed that the researcher was not asking about vegetation. The structural anomaly class is exactly that noticing, persisted.

This is also why the structural class is epistemically distinct from density, correlation, and narrative. Those three describe properties of the data the agent examined. A structural anomaly describes a property of the *examination*, not the data. Treating it as a fourth class — rather than as a flavour of one of the existing three — keeps the distinction explicit so that consumers (other agents, downstream systems, human researchers reviewing audit trails) can tell the two kinds of finding apart at a glance.

---

## 8. Surveillance Mode: Counteracting User Framing Bias

### 8.1 Two Modes of Operation

The Colorado River test reveals that Contextus requires two distinct modes of operation, and that neither is sufficient alone.

- **Query mode:** The researcher asks a question. Correlation agents search the hypergraph for answers. This mode is powerful but vulnerable to framing bias. The quality of answers is bounded by the quality of the question, and the question inherits the researcher's domain assumptions.
- **Surveillance mode:** Scout agents continuously monitor the hypergraph for anomalies, broken correlations, and regime changes, without being asked. No human question frames the search. The agent notices that a historically reliable relationship has degraded and flags it. This mode is immune to framing bias because there is no frame.

The plant-snowmelt insight lives in surveillance mode. It is a broken correlation (snowpack no longer predicts streamflow), not a missing quantity (water disappearing from the budget). A scout agent would have flagged the degrading correlation coefficient between April snowpack and summer streamflow across the Upper Basin without any human prompt. The researcher then investigates the flag, and the cross-domain traversal through vegetation, precipitation, and groundwater nodes leads to the insight.

### 8.2 Edge Scouting: The Framing Bias Countermeasure

Surveillance mode alone is necessary but not sufficient. The Colorado River test shows that even when a researcher is actively querying the system, their framing will bias the search toward their existing mental model of the problem. The system needs a mechanism that actively works along the edge of the researcher's framing, looking for patterns that the current question set is systematically missing.

We call this **edge scouting**. When a researcher submits a query or a series of queries, edge scouts analyse the query pattern to identify what the researcher is not asking. They do this by comparing the subgraph being explored (the nodes and edges the researcher's queries are touching) against the full connectivity of those nodes. If a node in the researcher's active subgraph has significant connections to domains the researcher has not queried, the edge scout flags these unexplored connections.

In the Colorado River test, the researcher's queries touched hydrology nodes (streamflow, reservoirs, snowpack), climate nodes (temperature, precipitation), and infrastructure nodes (dams, diversions). These nodes have strong connections to ecology nodes (vegetation, species, soil biology) that were never queried. An edge scout would have flagged: "Your active subgraph includes snowpack and streamflow nodes. These have 847 connections to vegetation phenology nodes that you have not explored. Three of these connections show statistically significant changes since 2000."

The edge scout does not answer a question. It does not interpret data. It simply observes that the researcher's exploration pattern has a blind spot, and makes the blind spot visible. The researcher decides whether to follow the flag. The system's role is to ensure that the researcher's expertise-driven framing does not inadvertently exclude the domains where the insight lives.

### 8.3 The Empathy Dimension

Edge scouting must be implemented with care. A system that constantly interrupts the researcher with "you haven't looked at X" will be turned off. The flags must be calibrated to the researcher's current focus, presented unobtrusively (the visualisation equivalent of a peripheral colour shift, not a modal alert), and ranked by the statistical significance of the unexplored connections. The system should feel like a quiet colleague who occasionally says "have you considered?" — not like an audit trail of everything the researcher hasn't done.

The design constraint is familiar: the system exists to help the researcher think better, not to think for them. Edge scouting extends this principle to the meta-level: it helps the researcher *question* better, by making visible the boundaries of their current questioning pattern. It is attention management applied to the question space itself.

### 8.4 Design Requirement

The Colorado River test generates a specific design requirement for the Contextus specification:

> Every researcher session must have at least one active edge scout that monitors the boundary between the explored and unexplored subgraph. The edge scout operates independently of the researcher's queries. It does not require configuration or domain knowledge from the researcher. Its output is a ranked list of unexplored connections from the active subgraph, filtered by statistical significance and recency of change. The output is rendered in the visualisation layer as subtle indicators on nodes at the boundary of the explored region, visible in peripheral vision but not demanding focused attention.

This is the architectural counterweight to expert framing bias. The researcher's expertise drives the query. The edge scout ensures the query does not become a cage.

### 8.5 Corpus Edge Scouting: The Distributed Confirmation Bias Problem

Edge scouting on the researcher's query boundary is necessary but not sufficient. The Colorado River test revealed a second-order bias: a feedback loop between the researcher and the search system. The researcher asks hydrology questions. The search engine returns hydrology papers. Those papers frame the problem in hydrology terms. That framing reinforces the researcher's initial frame, which generates more hydrology queries. The researcher and the search system confirm each other's assumptions in a closed loop. Neither party is wrong individually — every result is relevant, every query is reasonable — but the system converges on a basin of attraction that excludes cross-domain insights.

This is confirmation bias distributed across two agents. It is worse than individual confirmation bias because it feels like validation. The analysis finds real mechanisms, gets substantive results, builds a coherent picture. The quality of individual results masks the narrowing of the search space.

Corpus edge scouting monitors the boundary of what the search system is returning, not just what the researcher is asking. It evaluates: Are the search results clustering in a single domain? Are the sources from a single institutional perspective? Is the vocabulary converging? Are there adjacent domains with relevant data that the current search terms cannot reach? Before results are returned to the researcher or querying agent, a diversity check evaluates domain distribution. If convergence is detected, the system autonomously injects broadening queries designed to reach adjacent domains that the current vocabulary cannot access — not replacing the researcher's queries but supplementing them.

This applies equally to any searching agent: an LLM performing literature review, a BMA agent traversing the hypergraph, or a human using a search engine. Automated agents may be more vulnerable than humans because they process results faster and build self-consistent narratives before the narrowing becomes apparent. Corpus edge scouting is the countermeasure at the search interface itself.

### 8.6 The Local Minimum Problem and the Bridge Agent

The Colorado River test exhibited a local minimum in the insight space. The analysis converged on groundwater extraction by wells and pumping — a coherent, well-evidenced, quantified explanation for the missing water. This explanation was physically adjacent to the actual dominant mechanism (plants extracting groundwater and intercepting snowmelt) but separated from it by a domain boundary. Both are "things taking water out of the ground." The physical distance is zero. The institutional distance is enormous.

The local minimum was stable because the groundwater depletion story was satisfying: it explained the missing water, had clear policy implications, connected to other mechanisms in a tidy feedback loop, and was quantified by satellite data. There was no gradient pulling toward the adjacent insight. The loss function was flat between "wells pumping groundwater" and "plants pumping groundwater" because the search system could not see across the domain boundary to know there was a better explanation next door.

This generates a requirement for a third agent type beyond scouts and search agents: the **bridge agent**. The bridge agent maintains a higher-altitude view of three things simultaneously: what surveillance mode has flagged (anomalies, broken correlations), what query mode is actively investigating (the researcher's current question set and search results), and the topological distance between them in the hypergraph. Its job is to detect when a surveillance finding and an active search are close in physical or process space but separated by a domain boundary in knowledge space — which is the signature of a local minimum.

In the Colorado River case, the bridge agent would see: the search agent is investigating groundwater depletion (hydrology/engineering nodes). The scout has flagged a broken snowpack-streamflow correlation (hydrology/climate nodes). Both involve water leaving the subsurface. But the search agent's explanation (pumping) does not explain the scout's anomaly (forecasting gap). The bridge agent notices the gap and asks: what else extracts groundwater in these watersheds? That question crosses the domain boundary into plant physiology without the researcher or the search agent needing to know that domain exists.

The bridge agent must operate concurrently with the search, not retrospectively. If it reviews scout flags only after the search is complete, the researcher has already converged on the local minimum and formed a coherent narrative. The bridge needs to continuously compare the search trajectory against outstanding scout flags, and intervene when it detects that a flag and a search are converging on the same physical process from different domain perspectives. It provides the energy to escape the basin of attraction by injecting information from outside the current search's convergence zone.

The bias in this system is not located in any single component. It is an emergent property of the interaction between the researcher's expertise, the question framing, the search system's relevance ranking, and the corpus's institutional structure. No single fix addresses it. The architecture requires edge scouting on the researcher (query boundary monitoring), corpus edge scouting on the search results, and a bridge agent connecting surveillance findings to active searches across domain boundaries. Together, these three mechanisms provide the escape velocity from local minima in the insight space.

---

## 9. Design Principles

1. **Comprehension over computation.** The system exists to help a human understand a complex domain. Every feature must serve this goal.
2. **Relationships over entities.** The hypergraph models how things connect, not just what things are.
3. **Questions over answers.** The system's most important function is enabling researchers to ask questions that existing tools structurally prevent. The space of askable questions determines the space of discoverable insights.
4. **Agents as attention, not authority.** AI directs the researcher's gaze; the researcher makes the judgments.
5. **No pre-specified windows.** Agents search across domains and timescales without assuming the correct correlation structure. The data reveals the structure; the tool does not impose it.
6. **Surveillance before query.** Scout agents monitor for anomalies and broken correlations continuously, without being asked. The most important insights are the ones no one thought to look for.
7. **Edge scouting against framing bias.** When a researcher queries the system, edge scouts monitor the boundary of the explored subgraph and flag significant unexplored connections. The researcher's expertise must not become a cage.
8. **Corpus diversity monitoring.** Search results must be monitored for domain convergence. When results cluster in a single domain or institutional vocabulary, the system injects broadening queries into adjacent domains. Distributed confirmation bias between researcher and search system is an emergent property that must be architecturally prevented.
9. **Bridge agents against local minima.** A bridge agent operates concurrently with active searches, connecting surveillance findings to search trajectories across domain boundaries. When a scout flag and an active search converge on the same physical process from different domains, the bridge agent intervenes to provide escape velocity from local minima in the insight space.
10. **Biology is behaviour, not background.** Living systems are active agents that adapt to conditions, not passive coefficients. Vegetation, organisms, and ecosystems must be modelled as entities with behavioural edges, not as loss terms in physical equations.
11. **Provenance is non-negotiable.** Every datum, inference, and hypothesis must be traceable to its source with explicit confidence.
12. **Perception-native visualisation.** Use motion, colour, and narrative to work with human perceptual strengths, not against them.
13. **Open tools, governed data.** The framework is open source. Sensitive data is access-tiered.
14. **Built on what exists.** MuninnDB, NATS, Podman, FastMCP, BMA-BRIDGE. Contextus is a domain layer, not a new platform.

---

*End of Theory Document v1.5. See: Contextus Specification v1.3 for implementation details (the `AnomalyStructural` constant in §11.1 is no longer provisional after this revision). See: Contextus–CTH Bridge Specification v0.1 for external evaluation integration. Closes contextus issue [#3](https://github.com/JamesPagetButler/contextus/issues/3).*
