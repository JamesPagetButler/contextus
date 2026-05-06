# Contextus Theory — Proposed §3.6

**Insert after §3.5 (The Human-Agent Relationship)**

---

## 3.6 The Insight Signal: What Agents Produce

Sections 3.1 through 3.5 describe what agents do: they traverse the hypergraph, detect anomalies, test correlations, propose narratives, and track provenance. But they do not define what an agent *produces* when it finds something worth flagging. Without a formal definition of this output, the system has a process but not a product. This section defines that product: the Insight Signal.

### 3.6.1 Why a Formal Signal Matters

An agent that detects an anomaly must communicate that detection to the rest of the system — to other agents, to the visualisation layer, and ultimately to the human researcher. If this communication is informal (a log entry, a highlighted node, a notification), then insights are ephemeral. They exist only in the moment of display. They cannot be compared, tracked over time, combined with other insights, or evaluated for reliability. They are sparks, not records.

Contextus requires insights to be records. Not conclusions — the theory is explicit that interpretation remains with the human — but structured descriptions of what an agent observed, where and when it observed it, how surprising the observation was, and how confident the agent is in the observation's validity. A record can be stored, decayed, strengthened by subsequent evidence, and eventually promoted or dismissed. A spark cannot.

The Insight Signal is this record. It is the atomic unit of agent output. Every time an agent flags something — an anomaly, a correlation, a narrative thread — it emits an Insight Signal into the hypergraph. The signal is itself a node in the hypergraph, connected to the nodes and hyperedges that constitute the pattern it describes. Insights about the data become part of the data. This reflexivity is deliberate: it allows agents to discover insights about insights, building layered understanding over time.

### 3.6.2 Anatomy of an Insight Signal

An Insight Signal is composed of six elements, each derived from structures the theory has already defined.

**The Subgraph.** The specific set of nodes and hyperedges the agent flagged. This is the pattern itself — the flower the agent picked out from the forest. In a trophic cascade analysis, the subgraph might be the set of hyperedges connecting wolf presence, elk behavioural change, willow recovery, and stream bank stability at specific Locale coordinates. The subgraph is not the agent's interpretation of the pattern; it is the pattern. It is a pointer into the hypergraph, not a copy of the data.

**The Traversal Path.** The sequence of hypergraph steps the agent took to arrive at the subgraph. This is provenance for the pattern — it records how the agent found it, which matters because different traversal paths to the same subgraph carry different epistemic weight. A scout that stumbled onto the pattern during random traversal and a correlation agent that was directed to test a specific hypothesis may flag the same subgraph, but the directed discovery provides stronger evidence than the undirected one. The traversal path also enables the system to detect when multiple agents independently converge on the same pattern — a convergence signal that is itself significant.

**The Anomaly Characterisation.** What makes this subgraph interesting. The theory defines three agent types that produce three corresponding classes of anomaly:

- *Density anomaly* (from scout agents): the subgraph exhibits unexpected density — more or fewer connections, stronger or weaker edge weights, or faster or slower temporal dynamics than the surrounding hypergraph baseline. Something changed, or something is different from what was expected.

- *Correlation anomaly* (from correlation agents): the subgraph exhibits a statistical relationship — positive, negative, or conditional — that was not previously recorded or that contradicts a previously recorded relationship. The agent is not claiming causation; it is reporting co-variation.

- *Narrative anomaly* (from synthesis agents): the subgraph forms a coherent causal or temporal chain — a sequence of events or conditions that, taken together, tell a story. The agent is proposing that the individual edges form a path, and that the path has explanatory power. This is the weakest class of anomaly epistemically, but often the most valuable for directing human attention, because narrative is the format in which human insight most naturally operates.

Each class carries a score quantifying how surprising the observation is relative to the agent's baseline expectations. This score is not a p-value; it is an information-theoretic measure of how much the observation would change an informed observer's beliefs. A density anomaly in a well-characterised region of the hypergraph (many observations, stable patterns) is more surprising than the same anomaly in a sparsely observed region. The score accounts for this.

**The Locale Envelope.** The spatiotemporal bounds within which the pattern holds. Derived directly from the Locale framework (§2.2), the envelope defines where and when the insight is valid. A seasonal pattern has a Locale envelope with temporal periodicity. A geographically bounded pattern has a spatial envelope that may or may not shift over time. The envelope is not a bounding box; it is a region in Locale space whose shape is determined by the data.

The Locale envelope carries an additional property: its *grain* — the finest resolution at which the pattern is detectable. A trend visible in annual averages but not monthly data has coarse temporal grain. A pattern visible at the watershed level but not at individual stream reaches has coarse spatial grain. Grain matters because it constrains what evidence can confirm or deny the insight. Evidence at a finer grain than the pattern may be noise; evidence at a coarser grain may average out the signal.

**The Persistence Measure.** How many MuninnDB consolidation cycles the pattern has survived. A signal that fires once and never again may be an artifact — sensor noise, data ingestion error, or a genuinely transient event. A signal that persists across multiple consolidation cycles, during which Ebbinghaus decay would naturally attenuate unsupported connections, is more likely to reflect a real pattern in the data. Persistence is not proof, but it is a necessary condition for the signal to be taken seriously by downstream processes.

Persistence interacts with the landscape/exploration distinction from §2.3. A structural pattern (one that is embedded in the topology of the hypergraph) will persist indefinitely because it is part of the landscape. A contingent pattern (one that depends on specific conditions being active) will appear and disappear as those conditions come and go. Both are real patterns; they differ in kind, not in quality. The Insight Signal records which type it appears to be, so that downstream processes can apply appropriate expectations.

**The Confidence Estimate.** The agent's Bayesian confidence in the signal's validity, derived from the provenance chain (§3.2, provenance agents). Confidence incorporates data quality, source independence, temporal coverage, and known confounders. It is explicitly *not* confidence in the signal's importance or meaning — those are human judgments. It is confidence that the pattern described by the subgraph actually exists in the data and is not an artifact of incomplete observation, measurement error, or analytical mistake.

Confidence is reported as a distribution, not a point estimate. A signal with high mean confidence and narrow variance is well-established. A signal with moderate mean confidence and wide variance is plausible but uncertain. A signal with high mean but wide variance deserves particular attention: something is probably there, but the evidence is inconsistent. Each case suggests a different response from the human researcher.

### 3.6.3 Signal Lifecycle

An Insight Signal is not a static record. It evolves as new data enters the hypergraph, as agents re-traverse the relevant subgraph, and as the consolidation cycle runs.

**Emission.** An agent detects a pattern and emits a signal into the hypergraph. The signal is a new node, connected to the subgraph it describes. At emission, persistence is zero, confidence reflects only the initial observation, and the anomaly score is computed against the agent's current baseline.

**Strengthening.** Subsequent observations that are consistent with the pattern increase the signal's confidence and persistence. If multiple independent agents converge on the same subgraph — a scout flags it and a correlation agent independently confirms it — the signal gains convergence weight. This is the Hebbian principle applied at the insight level: signals that are repeatedly activated grow stronger.

**Decay.** In the absence of new supporting evidence, the signal's activation decays according to MuninnDB's Ebbinghaus schedule. The signal does not disappear; its activation approaches a floor determined by the strength of its initial evidence. A signal based on a single noisy observation decays to near-zero. A signal based on multiple high-confidence observations decays to a higher floor. The floor represents the permanent record that the observation occurred, even if it is no longer considered active.

**Mutation.** New data may change the shape of the pattern. The subgraph may expand (additional nodes or edges join the pattern), contract (some elements are no longer anomalous), or shift (the Locale envelope moves). When mutation occurs, the signal records the change as a version transition, preserving the history of how the pattern evolved. The signal's identity is continuous across mutations; it is the same insight, refined.

**Absorption.** When two signals describe overlapping subgraphs and a synthesis agent determines they are aspects of the same underlying pattern, the signals merge. The weaker signal is absorbed into the stronger one, its traversal path and provenance are incorporated, and the merged signal's confidence and anomaly score are recomputed. Absorption reduces inventory size without losing information.

**Promotion.** When a signal's confidence and persistence exceed a threshold, it becomes a candidate for external evaluation — a potential seed for structured investigation by systems outside Contextus. Promotion does not change the signal's status within the hypergraph; it creates an export event. The promoted signal retains its lifecycle within Contextus regardless of what happens downstream.

### 3.6.4 Relationship to Human Attention

The Insight Signal is not a conclusion, a recommendation, or an alert. It is a structured description of something an agent noticed. The visualisation layer (§4) renders active signals as perceptual cues — subtle colour shifts, motion changes, density variations — that guide the researcher's attention without demanding it. A signal with high anomaly score and high confidence renders more prominently. A signal with low persistence renders as a tentative ghost. A signal whose Locale envelope overlaps the researcher's current viewport renders; one outside the viewport does not.

The researcher engages with signals by investigating the subgraph they describe, not by reading the signal itself. The signal says "look here." The looking is the researcher's job. The understanding that results from looking is the insight. The signal is the finger pointing at the moon; it is not the moon.

### 3.6.5 Design Constraints on the Signal

Three constraints preserve the integrity of the Insight Signal as a design element:

**Signals are never conclusions.** An agent may report "these five nodes co-activate anomalously" but never "these five nodes co-activate because of X." Causal claims require human judgment. The signal structure has no field for causal explanation. This is a deliberate omission, not an oversight.

**Signals are always traceable.** Every element of the signal — the subgraph, the traversal path, the anomaly score, the confidence — is derived from specific data in the hypergraph and can be audited. A signal that cannot be traced to its source data is malformed and must be rejected by the provenance agent.

**Signals are never suppressed.** The system may deprioritise a signal for display (low confidence, low persistence), but it never deletes one. Suppressed signals have historically turned out to be important when context changes. The Ebbinghaus decay ensures that genuinely unsupported signals fade naturally; forced suppression is not necessary and could hide valid patterns.

---

*This section defines the Insight Signal as the formal output type for Contextus agents. Implementation details — data structures, NATS subjects, MuninnDB schema — are specified in the Contextus Technical Specification v1.0. The relationship between the Insight Signal and external evaluation frameworks is specified in the Contextus–CTH Bridge Specification.*
