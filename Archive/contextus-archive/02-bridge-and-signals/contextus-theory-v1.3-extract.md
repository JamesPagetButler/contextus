# Contextus Theory Document v1.3 — Reconstruction

**Reconstruction status:** PARTIAL — §3.6 (the v1.2 → v1.3 addition) recovered NEAR-COMPLETE. Other sections inherited from v1.2, see `01-contextus-theory-spec/contextus-theory-v1.0-extract.md`.
**Source chat:** https://claude.ai/chat/3175e792-70f2-48df-bd8e-58762dd066fc

---

# PROJECT CONTEXTUS

## Theory Document v1.3

*April 2026 — adds §3.6 Insight Signal and Design Principle 10*

---

## Sections 1-3.5

> Inherited from v1.2 — see `01-contextus-theory-spec/contextus-theory-v1.0-extract.md` for §1 (the comprehension problem), §2 (hypergraph as substrate), §3.1-3.5 (Locale).

---

## 3.6 The Insight Signal

> NEW IN v1.3.

The Insight Signal is the atomic unit of agent output in Contextus. It is the formal record of an anomalous, correlated, or narratively coherent pattern detected in the hypergraph.

### 3.6.1 Why a Formal Signal

Without a first-class signal type, agent outputs would be ad-hoc messages: scout agents would emit "I noticed X" notifications, correlation agents would emit "X and Y co-vary" messages, synthesis agents would emit "this could be explained by Z" messages. Each format would differ; each would carry different metadata; each would be hard to compose, query, or evaluate at scale.

The Insight Signal unifies these. Every agent — scout, correlation, synthesis — emits the same structured object. The object describes what was noticed, where in the hypergraph it lives, how anomalous it is, and how confident the agent is. The signal does not say what the pattern *means*; meaning is the human's job, or the Bridge's job when promoting to CTH.

### 3.6.2 Anatomy of a Signal

A Signal contains:

| Element | Role |
|---|---|
| Subgraph | The hyperedges and nodes implicated in the pattern |
| Traversal path | The sequence of nodes the agent walked to find the pattern |
| Anomaly type | density / correlation / narrative |
| Anomaly score | How surprising is this pattern, given history? |
| Confidence | How certain is the agent that the pattern is real? |
| Locale envelope | The spatial/temporal bounds of the pattern |
| Persistence | How many consolidation cycles has the pattern survived? |
| Provenance | Which data streams contributed; which agents observed |

A signal is an InsightSignal node in MuninnDB, connected to the hyperedges it describes via SIG_DESCRIBES edges. It is itself part of the hypergraph it describes — a reflexive structure where insights about data are part of the data.

### 3.6.3 Lifecycle

**Decay.** In the absence of new supporting evidence, the signal's activation decays according to MuninnDB's Ebbinghaus schedule. The signal does not disappear; its activation approaches a floor determined by the strength of its initial evidence. A signal based on a single noisy observation decays to near-zero. A signal based on multiple high-confidence observations decays to a higher floor. The floor represents the permanent record that the observation occurred, even if it is no longer considered active.

**Mutation.** New data may change the shape of the pattern. The subgraph may expand (additional nodes or edges join the pattern), contract (some elements are no longer anomalous), or shift (the Locale envelope moves). When mutation occurs, the signal records the change as a version transition, preserving the history of how the pattern evolved. The signal's identity is continuous across mutations; it is the same insight, refined.

**Absorption.** When two signals describe overlapping subgraphs and a synthesis agent determines they are aspects of the same underlying pattern, the signals merge. The weaker signal is absorbed into the stronger one, its traversal path and provenance are incorporated, and the merged signal's confidence and anomaly score are recomputed. Absorption reduces inventory size without losing information.

**Promotion.** When a signal's confidence and persistence exceed a threshold, it becomes a candidate for external evaluation — a potential seed for structured investigation by systems outside Contextus. Promotion does not change the signal's status within the hypergraph; it creates an export event. The promoted signal retains its lifecycle within Contextus regardless of what happens downstream.

### 3.6.4 Relationship to Human Attention

The Insight Signal is not a conclusion, a recommendation, or an alert. It is a structured description of something an agent noticed. The visualisation layer (§4) renders active signals as perceptual cues — subtle colour shifts, motion changes, density variations — that guide the researcher's attention without demanding it. A signal with high anomaly score and high confidence renders more prominently. A signal with low persistence renders as a tentative ghost. A signal whose Locale envelope overlaps the researcher's current viewport renders; one outside the viewport does not.

The researcher engages with signals by investigating the subgraph they describe, not by reading the signal itself. The signal says "look here." The looking is the researcher's job. The understanding that results from looking is the insight. The signal is the finger pointing at the moon; it is not the moon.

### 3.6.5 Design Constraints on the Signal

Three constraints preserve the integrity of the Insight Signal as a design element:

**Signals are never conclusions.** An agent may report "these five nodes co-activate anomalously" but never "these five nodes co-activate because of X." Causal interpretation is reserved for humans and for the Bridge boundary. The signal describes what was noticed; it does not explain why.

**Signals are always traceable.** Every signal carries full provenance back to the source data, the agent that produced it, the model version, and the traversal path. This is non-negotiable. A signal without provenance is not a signal; it is a hallucination.

**Signals are never suppressed.** Even contradictory signals remain in the hypergraph. The system does not pre-filter or hide signals based on their relationship to existing beliefs. The human (or the Bridge, with CTH judgment) decides what to do with each signal. Suppression is a category of curation, and curation belongs to the user, not the system.

---

## Sections 4-6

> Inherited from v1.2.

---

## 7. Design Principles

[v1.2's nine principles plus Principle 10 added in v1.3]

1. **Comprehension before generation.** The system optimises for what researchers need to understand, not for what tools can compute.
2. **Empathy as primary constraint.** Human perceptual and cognitive limits define the visualisation envelope. The information density cap is empathy made architectural.
3. **Hypothesis as path.** The space of askable questions determines the space of discoverable insights.
4. **Agents as attention, not authority.** AI directs the researcher's gaze; the researcher makes the judgments.
5. **No pre-specified windows.** Agents search across domains and timescales without assuming the correct correlation structure. The data reveals the structure; the tool does not impose it.
6. **Provenance is non-negotiable.** Every datum, inference, and hypothesis must be traceable to its source with explicit confidence.
7. **Perception-native visualisation.** Use motion, colour, and narrative to work with human perceptual strengths, not against them.
8. **Open tools, governed data.** The framework is open source. Sensitive data is access-tiered.
9. **Built on what exists.** MuninnDB, NATS, Podman, FastMCP, BMA-BRIDGE. Contextus is a domain layer, not a new platform.
10. **Signals are evidence of patterns, not conclusions.** *(NEW in v1.3)* The Insight Signal records what was observed. The human decides what it means. External evaluation frameworks (CTH) judge whether claims derived from signals hold up. These are separate concerns.

---

*End of Theory Document v1.3.*
*See: Contextus Specification v1.1 for implementation details.*
*See: Contextus–CTH Bridge Specification v0.1 for external evaluation integration.*
