# Materia-Bio eDNA Prompt v0.1 — Reconstruction

**Reconstruction status:** STRUCTURAL — section list, key concepts (community-as-morphotype, cross-programme handshake), and §7-§9 recovered. Body text for §1-§6 partial.
**Source chat:** https://claude.ai/chat/6f0eae0c-ce00-4ca1-924b-7ec420f6496b

---

# Materia-Bio eDNA — Theoretical Lens Specification Prompt

**Version:** 0.1 — Initial prompt
**Date:** 2026-04-01
**Programme:** Materia-Bio
**Author:** James Paget Butler, with Claude (Anthropic)

---

## Architectural Position

For Materia-Bio, eDNA is a **theoretical lens**, not infrastructure. While Contextus treats eDNA as a sensing modality (engineering concern), Materia-Bio uses eDNA data to investigate a deep theoretical question:

> **Are the rules governing community assembly structurally analogous to the rules governing developmental regulatory networks within a single organism?**

Both involve diverse phenotypic expression emerging from a shared substrate under environmental constraint. Brassica oleracea expresses cabbage, kale, broccoli, kohlrabi, Brussels sprouts, and cauliflower from the same genome through regulatory differences. A Yellowstone watershed expresses extremophile microbial communities, salmonid fish assemblages, and large mammal trophic structures from the same regional species pool through environmental gradients.

This is the **community-as-morphotype analogy** — the centrepiece of Materia-Bio's eDNA contribution.

---

## 1. The Community-as-Morphotype Analogy

[Reconstruction limited. Key concept:]

Within an organism: regulatory differences (transcription factors, chromatin state, signalling thresholds) determine which morphotype expresses from a shared genome.

Within an ecosystem: environmental differences (temperature, nutrient availability, chemistry, predation pressure) determine which community composition expresses from a shared regional species pool.

Both can be modelled as a **state space** where shared substrate constrains the achievable phenotypes/communities, and environment selects which state is realised. The mathematical structure may be formally isomorphic — though this requires careful development.

---

## 2. Five Research Axes

### 2.1 eDNA Resolution Limits
Can metabarcoding distinguish between closely related species or subspecies within a radiation (cichlids) or domesticated species complex (Brassica, dogs)? At what phylogenetic distance does resolution fail for standard markers (COI, 12S, 16S, ITS)?

### 2.2 Geothermal Gradient as Continuous Regulatory Landscape (highlighted in chat as the most interesting axis)

Treating a hot spring runoff channel as a continuous regulatory landscape with punctuated community transitions. Directly testable with existing Yellowstone eDNA data:
- At each point along a thermal gradient, what is the community composition?
- Are transitions continuous (gradual species turnover) or punctuated (discrete attractor states)?
- Can the same regulatory-network mathematics that explains Brassica morphotype switches explain community state transitions?

### 2.3 Regulatory vs Coding Signal
When eDNA distinguishes morphotypes of the same species, where is the distinguishing signal — coding regions, regulatory regions, or epigenetic marks? This has direct implications for sequencing approach (targeted amplicon vs shotgun metagenomics).

### 2.4 Community Attractor States
Do Yellowstone eDNA community composition data show evidence of discrete attractor states (consistent with bistable switch hypothesis) or continuous variation (consistent with neutral assembly theory)?

### 2.5 eDNA as a Phenotype
Should eDNA shedding rate and profile be treated as a phenotypic trait under selection, rather than merely a measurement artefact? If so, eDNA enters Materia-Bio's trait analysis framework directly.

---

## 3. Cross-Programme Handshake Protocol (THE KEY ARCHITECTURAL DECISION)

This protocol keeps Contextus and Materia-Bio coupled without merging — exactly the architecture James specified ("eDNA is a layer in Contextus, not a new programme").

```
Contextus observes anomaly
        │
        ▼
"Anomalous community composition at Locale X, Time T:
  unexpected co-occurrence of taxa A, B, C"
        │
        ▼ (handshake message)
        │
Materia-Bio receives mechanism question
        │
        ▼
Materia-Bio generates hypothesis
        │
        ▼
"Hypothesis: A, B, C share a regulatory pattern
  responsive to gradient G; testable by..."
        │
        ▼ (handshake message)
        │
Contextus receives testable prediction
        │
        ▼
Contextus seeks corroborating data, refines anomaly
        │
        ▼
(loop until prediction confirmed, refuted, or refined)
```

The protocol is asymmetric:
- **Contextus → Materia-Bio:** observational anomalies as mechanism questions
- **Materia-Bio → Contextus:** mechanistic hypotheses as testable predictions

Neither programme owns the loop. The loop is the bridge.

---

## 4. Connection to Existing Model Organisms

Each research axis ties back to one or more model organisms:

| Axis | Brassica oleracea | Dogs | Pigeons | Cichlids |
|---|---|---|---|---|
| Resolution limits | Subspecies markers | Breed-level eDNA | Fancy varieties | Lake radiation |
| Regulatory vs coding | FLC2 master-hub regulatory data | IGF1 master-hub | Plumage genes | Jaw morphology |
| Community attractors | (parallel to community states) | Behavioural attractors | Flight/colour syndromes | Trophic morphs |
| Cross-scale isomorphism | Within-genome → between-genome | Breed → population | Variety → species | Morph → species |
| eDNA as phenotype | Pollen / leaf-litter eDNA | Saliva, hair, dander | Feather/dust eDNA | Mucus, faecal eDNA |

The IGF1 master-hub finding in dogs (parallel to FLC2 in Brassica) is identified as the strongest cross-kingdom prediction the framework makes.

---

## 5. Materia-Chem Integration

[Reconstruction limited. Key reference:]

eDNA degradation chemistry (DNase activity, oxidative damage, UV photolysis) is a Materia-Chem concern. The Materia-Chem eDNA degradation model becomes a shared computational substrate between Contextus (for transport modelling) and Materia-Bio (for understanding what survives to be measured).

---

## 6. Yellowstone Geothermal Gradient — The Flagship Experiment

[Reconstruction limited. Conceptual frame:]

Hot spring runoff channels in Yellowstone provide naturally-occurring continuous gradients in temperature, chemistry, and oxygen. Each channel is a continuous regulatory landscape. eDNA sampling along the gradient at high spatial resolution would test:

- Are community transitions continuous or punctuated?
- Do transition points correspond to threshold values in known geochemical gradients?
- Is the same bistable mathematics that governs FLC2 switching in Brassica observable in community state transitions?

This experiment uses existing Yellowstone access, well-characterised gradients, and standard eDNA methodologies. It is among the lowest-cost / highest-information experiments the programme could run.

---

## 7. Open Questions (recovered with HIGH confidence)

1. **Resolution limit:** Can eDNA metabarcoding distinguish between closely related species or subspecies within a radiation (cichlids) or a domesticated species complex (Brassica, dogs)? At what phylogenetic distance does resolution fail for standard markers (COI, 12S, 16S)?

2. **Regulatory vs coding signal:** When eDNA distinguishes morphotypes of the same species, is the distinguishing signal in coding regions, regulatory regions, or epigenetic marks? This has implications for which sequencing approach to use (targeted amplicon vs shotgun metagenomics).

3. **Community attractor states:** Do Yellowstone's eDNA community composition data show evidence of discrete attractor states (consistent with the bistable switch hypothesis) or continuous variation (consistent with neutral assembly)?

4. **Cross-scale isomorphism:** Is the mathematical structure of within-organism regulatory diversity (Brassica morphotypes) formally isomorphic to between-species community assembly (eDNA community composition)? This is a strong claim that requires careful mathematical development.

5. **eDNA as a phenotype:** Should eDNA shedding rate and profile be treated as a phenotypic trait under selection, rather than merely a measurement artefact? If so, it enters Materia-Bio's trait analysis framework directly.

---

## 8. Deliverables (recovered with HIGH confidence)

- [ ] Materia-Bio eDNA section for Theory v2.0
- [ ] Cross-programme handshake protocol specification (with Contextus)
- [ ] Model organism eDNA experimental design documents (one per organism)
- [ ] Community-as-morphotype theoretical framework paper (formal hypothesis)
- [ ] Materia-Chem eDNA degradation model specification
- [ ] Yellowstone geothermal gradient community assembly analysis plan

---

## 9. Attribution (recovered with HIGH confidence)

This work builds on Materia-Bio Theory v1.0 and Spec v1.0, and draws on published eDNA research including: Carraro & Carraro 2025 (eDNA occupancy modelling); Rowe et al. 2024 (Yellowstone geothermal microbial diversity); Sullivan et al. 2025 (airborne eDNA temporal analysis); HAPP pipeline (Ronquist et al. 2025). The community-as-morphotype theoretical framework is original to this programme and should be attributed to the Materia-Bio collaboration.

---

*Document version: 0.1 — Initial prompt*
*Programme: Materia-Bio*
*Author: James Paget Butler, with Claude (Anthropic)*
