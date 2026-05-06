# Contextus Theory — Addendum v1.3
## The Colorado River Test: User Bias, Search Bias, and the Bridge Agent

**Addendum to Theory v1.2** | April 2026
**Author:** James Paget Butler | Helpful Engineering

*Inserts as Sections 8–10 of the Theory Document, before Design Principles.*

---

## 8. A Worked Failure: The Colorado River Test

In April 2026, we ran a deliberate test of the Contextus thesis using the Colorado River system as a case study. The test was designed to determine whether an analytical process guided by the Contextus principles would discover a novel cross-domain insight that had recently been published in the scientific literature. The test failed in an instructive way. This section documents the failure and what it teaches about the system's architecture.

### 8.1 The Setup

The Colorado River is one of the most data-rich, multi-domain, and politically contested natural systems on Earth. It serves 40 million people across seven U.S. states and Mexico, and is governed by a century-old legal framework (the 1922 Colorado River Compact and subsequent "Law of the River") that over-allocated the river from day one. Since 2000, the system has been in persistent deficit, with total storage falling from near capacity to approximately 37% by April 2026.

The test proceeded in stages. First, a comprehensive baseline was assembled covering the river's physical geography, legal framework, infrastructure, ecology, and current status. Then, the question was posed: "Is there missing water? Where is it going?"

### 8.2 What the Analysis Found

The analysis identified five mechanisms draining the system that sit outside the management framework's field of view:

- **Groundwater mining:** GRACE satellite data shows 65% of the basin's total water storage loss since 2002 is groundwater depletion, not reservoir decline. Groundwater is unregulated by the interstate compact.
- **Dust on snow:** Anthropogenic dust loading darkens mountain snowpacks, accelerating melt by up to six weeks and reducing annual runoff by approximately 5%. No snowmelt forecast model accounts for dust.
- **Wildfire watershed destruction:** Burned soils lose infiltration capacity, converting slow percolation into flash runoff and sediment loading. Fire seasons are lengthening as the drought deepens.
- **Warming-driven evapotranspiration:** Each degree Celsius of warming reduces discharge by approximately 9.3%. Warmer forests in the headwaters consume more water, reducing groundwater recharge.
- **Irrigation method paradox:** Government-subsidised centre-pivot sprinklers, marketed as conservation, produce zero return flow. System-level efficiency may have decreased even as field-level efficiency increased.

These five mechanisms were identified through systematic cross-domain search and represent genuine analytical value. The feedback loops between them — drought exposes soil, dust accelerates melt, less water means more drought — are real and largely absent from the management framework's models.

### 8.3 What the Analysis Missed

A paper published simultaneously (Hogan et al., 2026, *Geophysical Research Letters*) had identified the dominant mechanism: vegetation interception of snowmelt. The finding is that warmer, drier springs cause plants to consume snowmelt before it reaches rivers, and this single mechanism explains approximately 70% of the gap between snowpack-predicted and actual streamflow since 2000. A companion Princeton study found the "drought paradox": plants maintain or increase transpiration during drought by switching from soil moisture to groundwater, competing directly with the river for the same water.

The analysis had the components of this insight. It cited headwater groundwater loss from increased forest water use. It discussed warming-driven evapotranspiration. But it treated vegetation as a passive loss channel — a coefficient in the energy balance — rather than as an active biological agent competing for water. It modelled biology as physics rather than as ecology. The critical mechanism — plants switching water sources in response to drought, intercepting snowmelt before it becomes streamflow — was structurally invisible to an analysis framed as "where is the water going?" because the water was never "going" anywhere. It was being intercepted before it arrived.

### 8.4 Root Cause Analysis

Three structural reasons explain the failure:

#### 8.4.1 Search by Mechanism, Not by Anomaly

The analysis searched for known loss mechanisms: things that remove water from the system. It did not search for unexplained discrepancies in the data. The snowpack-to-streamflow forecasting gap — good snowpack but less water than predicted, consistently since 2000 — is an anomaly sitting in publicly available data. A scout agent monitoring the hypergraph would have flagged this broken correlation without being asked. The analysis failed because it searched for causes rather than for unexplained patterns.

#### 8.4.2 Biology Modelled as Physics

Vegetation was treated as a passive variable — a loss term in the energy balance that scales with temperature. The Hogan paper shows vegetation is a strategic actor that adapts its behaviour to conditions: when soil moisture drops, plants switch to snowmelt; when snowmelt is unavailable, they tap groundwater. This is ecological behaviour, not physical evaporation. In a properly constructed hypergraph, plants would be nodes with behavioural edges — entities that change their interaction patterns in response to environmental state — not just coefficients.

#### 8.4.3 Inherited Domain Boundaries

The search terms used to investigate the problem were domain-specific: "water budget," "evaporation losses," "groundwater depletion." These are hydrology terms. The insight lives at the intersection of hydrology, plant physiology, and climate science. The analysis reproduced exactly the siloed thinking that Contextus is designed to overcome. It found what hydrologists find, because it searched the way hydrologists search.

### 8.5 The Role of User Framing

A fourth cause is more subtle and has significant implications for system design. The question that initiated the analysis was "Is there missing water? Where is it going?" This framing treats the problem as a loss accounting exercise. It naturally drives toward subtractive explanations: mechanisms that remove water from the system. The Hogan paper's insight is not about water being removed. It is about water being intercepted before it enters the system. The accounting frame cannot find interception, because in an accounting frame the water was never "there" to go missing.

A better question would have been "Why do snowpack measurements no longer predict streamflow?" This frames the problem as a broken correlation rather than a missing quantity, and points directly at the interception mechanism. But a researcher approaching the Colorado River from the existing management context — the Law of the River, the 24-Month Study, the allocation disputes — will naturally frame questions in allocation terms. Their expertise defines their framing, and their framing constrains their discovery space.

This is not a criticism of the researcher. It is a structural observation about how expert knowledge interacts with analytical tools. The researcher's expertise is simultaneously their greatest asset and their greatest liability: they know where to look, which means they also know where not to look. A system designed to augment expert comprehension must account for the fact that expert framing will systematically bias the question space.

---

## 9. Surveillance Mode: Counteracting User Framing Bias

### 9.1 Two Modes of Operation

The Colorado River test reveals that Contextus requires two distinct modes of operation, and that neither is sufficient alone.

- **Query mode:** The researcher asks a question. Correlation agents search the hypergraph for answers. This mode is powerful but vulnerable to framing bias. The quality of answers is bounded by the quality of the question, and the question inherits the researcher's domain assumptions.
- **Surveillance mode:** Scout agents continuously monitor the hypergraph for anomalies, broken correlations, and regime changes, without being asked. No human question frames the search. The agent notices that a historically reliable relationship has degraded and flags it. This mode is immune to framing bias because there is no frame.

The plant-snowmelt insight lives in surveillance mode. It is a broken correlation (snowpack no longer predicts streamflow), not a missing quantity (water disappearing from the budget). A scout agent would have flagged the degrading correlation coefficient between April snowpack and summer streamflow across the Upper Basin without any human prompt. The researcher then investigates the flag, and the cross-domain traversal through vegetation, precipitation, and groundwater nodes leads to the insight.

### 9.2 Edge Scouting: The Framing Bias Countermeasure

Surveillance mode alone is necessary but not sufficient. The Colorado River test shows that even when a researcher is actively querying the system, their framing will bias the search toward their existing mental model of the problem. The system needs a mechanism that actively works along the edge of the researcher's framing, looking for patterns that the current question set is systematically missing.

We call this **edge scouting**. When a researcher submits a query or a series of queries, edge scouts analyse the query pattern to identify what the researcher is not asking. They do this by comparing the subgraph being explored (the nodes and edges the researcher's queries are touching) against the full connectivity of those nodes. If a node in the researcher's active subgraph has significant connections to domains the researcher has not queried, the edge scout flags these unexplored connections.

In the Colorado River test, the researcher's queries touched hydrology nodes (streamflow, reservoirs, snowpack), climate nodes (temperature, precipitation), and infrastructure nodes (dams, diversions). These nodes have strong connections to ecology nodes (vegetation, species, soil biology) that were never queried. An edge scout would have flagged: "Your active subgraph includes snowpack and streamflow nodes. These have 847 connections to vegetation phenology nodes that you have not explored. Three of these connections show statistically significant changes since 2000."

The edge scout does not answer a question. It does not interpret data. It simply observes that the researcher's exploration pattern has a blind spot, and makes the blind spot visible. The researcher decides whether to follow the flag. The system's role is to ensure that the researcher's expertise-driven framing does not inadvertently exclude the domains where the insight lives.

### 9.3 The Empathy Dimension

Edge scouting must be implemented with care. A system that constantly interrupts the researcher with "you haven't looked at X" will be turned off. The flags must be calibrated to the researcher's current focus, presented unobtrusively (the visualisation equivalent of a peripheral colour shift, not a modal alert), and ranked by the statistical significance of the unexplored connections. The system should feel like a quiet colleague who occasionally says "have you considered?" — not like an audit trail of everything the researcher hasn't done.

The design constraint is familiar: the system exists to help the researcher think better, not to think for them. Edge scouting extends this principle to the meta-level: it helps the researcher *question* better, by making visible the boundaries of their current questioning pattern. It is attention management applied to the question space itself.

### 9.4 Design Requirement

The Colorado River test generates a specific design requirement for the Contextus specification:

> Every researcher session must have at least one active edge scout that monitors the boundary between the explored and unexplored subgraph. The edge scout operates independently of the researcher's queries. It does not require configuration or domain knowledge from the researcher. Its output is a ranked list of unexplored connections from the active subgraph, filtered by statistical significance and recency of change. The output is rendered in the visualisation layer as subtle indicators on nodes at the boundary of the explored region, visible in peripheral vision but not demanding focused attention.

This is the architectural counterweight to expert framing bias. The researcher's expertise drives the query. The edge scout ensures the query does not become a cage.

### 9.5 Corpus Edge Scouting: The Distributed Confirmation Bias Problem

Edge scouting on the researcher's query boundary is necessary but not sufficient. The Colorado River test revealed a second-order bias: a feedback loop between the researcher and the search system. The researcher asks hydrology questions. The search engine returns hydrology papers. Those papers frame the problem in hydrology terms. That framing reinforces the researcher's initial frame, which generates more hydrology queries. The researcher and the search system confirm each other's assumptions in a closed loop. Neither party is wrong individually — every result is relevant, every query is reasonable — but the system converges on a basin of attraction that excludes cross-domain insights.

This is confirmation bias distributed across two agents. It is worse than individual confirmation bias because it feels like validation. The analysis finds real mechanisms, gets substantive results, builds a coherent picture. The quality of individual results masks the narrowing of the search space. The material surrounding the data — management reports, institutional analyses, policy documents — all frame the problem in the same domain vocabulary because that is the institutional context. The search corpus itself has a domain bias, and the search engine amplifies it by returning the most relevant results within that bias.

Corpus edge scouting monitors the boundary of what the search system is returning, not just what the researcher is asking. It evaluates: Are the search results clustering in a single domain? Are the sources from a single institutional perspective? Is the vocabulary converging? Are there adjacent domains with relevant data that the current search terms cannot reach? Before results are returned to the researcher or querying agent, a diversity check evaluates domain distribution. If convergence is detected, the system autonomously injects broadening queries designed to reach adjacent domains that the current vocabulary cannot access — not replacing the researcher's queries but supplementing them.

This applies equally to any searching agent: an LLM performing literature review, a BMA agent traversing the hypergraph, or a human using a search engine. The agent's initial framing generates queries, the queries return domain-clustered results, the results reinforce the framing, and the agent converges on a coherent but incomplete picture. Automated agents may be more vulnerable than humans because they process results faster and build self-consistent narratives before the narrowing becomes apparent. Corpus edge scouting is the countermeasure at the search interface itself.

### 9.6 The Local Minimum Problem and the Bridge Agent

The Colorado River test exhibited a local minimum in the insight space. The analysis converged on groundwater extraction by wells and pumping — a coherent, well-evidenced, quantified explanation for the missing water. This explanation was physically adjacent to the actual dominant mechanism (plants extracting groundwater and intercepting snowmelt) but separated from it by a domain boundary. Both are "things taking water out of the ground." The physical distance is zero. The institutional distance is enormous. One is measured by USGS monitoring wells and managed by state water engineers. The other is measured by eddy covariance towers and sap flow sensors and studied by plant physiologists.

The local minimum was stable because the groundwater depletion story was satisfying: it explained the missing water, had clear policy implications, connected to other mechanisms in a tidy feedback loop, and was quantified by satellite data. There was no gradient pulling toward the adjacent insight. The loss function was flat between "wells pumping groundwater" and "plants pumping groundwater" because the search system could not see across the domain boundary to know there was a better explanation next door.

This generates a requirement for a third agent type beyond scouts and search agents: the **bridge agent**. The bridge agent maintains a higher-altitude view of three things simultaneously: what surveillance mode has flagged (anomalies, broken correlations), what query mode is actively investigating (the researcher's current question set and search results), and the topological distance between them in the hypergraph. Its job is to detect when a surveillance finding and an active search are close in physical or process space but separated by a domain boundary in knowledge space — which is the signature of a local minimum.

In the Colorado River case, the bridge agent would see: the search agent is investigating groundwater depletion (hydrology/engineering nodes). The scout has flagged a broken snowpack-streamflow correlation (hydrology/climate nodes). Both involve water leaving the subsurface. But the search agent's explanation (pumping) does not explain the scout's anomaly (forecasting gap). The bridge agent notices the gap and asks: what else extracts groundwater in these watersheds? That question crosses the domain boundary into plant physiology without the researcher or the search agent needing to know that domain exists.

The bridge agent must operate concurrently with the search, not retrospectively. If it reviews scout flags only after the search is complete, the researcher has already converged on the local minimum and formed a coherent narrative. The bridge needs to continuously compare the search trajectory against outstanding scout flags, and intervene when it detects that a flag and a search are converging on the same physical process from different domain perspectives. It provides the energy to escape the basin of attraction by injecting information from outside the current search's convergence zone.

The bias in this system is not located in any single component. It is an emergent property of the interaction between the researcher's expertise, the question framing, the search system's relevance ranking, and the corpus's institutional structure. No single fix addresses it. The architecture requires edge scouting on the researcher (query boundary monitoring), edge scouting on the search results (corpus diversity monitoring), and a bridge agent connecting surveillance findings to active searches across domain boundaries. Together, these three mechanisms provide the escape velocity from local minima in the insight space.

---

## 10. Design Principles (v1.3 — Updated)

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
