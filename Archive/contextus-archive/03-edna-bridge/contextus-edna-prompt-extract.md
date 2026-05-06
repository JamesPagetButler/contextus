# Contextus eDNA Prompt v0.1 — Reconstruction

**Reconstruction status:** STRUCTURAL — section list and key architectural decisions recovered. Most section bodies are summaries.
**Source chat:** https://claude.ai/chat/6f0eae0c-ce00-4ca1-924b-7ec420f6496b

---

# Contextus eDNA Layer — Engineering Specification Prompt

**Version:** 0.1 — Initial prompt
**Date:** 2026-04-01
**Programme:** Contextus
**Author:** James Paget Butler, with Claude (Anthropic)

---

## Architectural Position

eDNA is a **sensing modality** within Contextus, not a sub-programme. From Contextus's perspective, eDNA is equivalent to satellite imagery, weather stations, or camera traps — a data source that populates the spatiotemporal model with organism presence information.

What makes eDNA uniquely powerful as a layer is that it provides direct biological community composition at a Locale coordinate without requiring direct observation. A camera trap sees a single organism passing through frame; an eDNA sample reveals the composite molecular footprint of every organism whose DNA has diffused through the sampled medium.

---

## 1. Three Sampling Modalities

### 1.1 Aquatic eDNA
- Water samples from rivers, lakes, ponds, ocean
- Filter-and-extract pipeline; PCR amplification of marker genes
- Hydrological transport must be modelled to back-assign detections to upstream source locales

### 1.2 Substrate eDNA
- Soil, sediment, surface samples
- Useful for terrestrial communities and historical depositional records
- Lower transport complexity than aquatic — signal is closer to source locale

### 1.3 Airborne eDNA
- Per Sullivan et al. (2025) temporal analysis breakthrough
- Air filter samples capture pollen, spores, skin cells, bacteria
- Highly time-resolved; less spatially resolved without micrometeorological modelling

---

## 2. NATS Message Schemas

[Reconstruction limited. Schema structure recovered:]

JSON payload structure for eDNA observations on NATS, posted to `ctx.edna.observation` subject. Fields include:
- Sample ID, Locale coordinate, sampling time
- Sampling modality (aquatic / substrate / airborne)
- Sequencing platform (Illumina, Nanopore, etc.)
- Marker gene(s) targeted (COI, 12S, 16S, ITS, etc.)
- Reference database used (BOLD, MIDORI2, etc.)
- ASV/OTU table or per-taxon read counts
- Confidence per taxonomic assignment
- Provenance (T/P/D/I tag)

---

## 3. MuninnDB Hyperedge Schema for eDNA

[Reconstruction limited. Key structure:]

An eDNA sample creates a single hyperedge linking:
- A Locale node (where sampled)
- A Time node (when sampled)
- A Modality node (aquatic / substrate / airborne)
- Multiple Species nodes (each detected taxon)
- A Provenance envelope (lab, sequencing platform, reference DB version)

This is **the canonical N-ary structure Contextus is designed for**. A single eDNA observation simultaneously connects place, time, method, and a set of organisms.

---

## 4. Hydrological Transport Model Layer

[Reconstruction limited. Key concept:]

Aquatic eDNA detections must be back-assigned to upstream source locales because eDNA travels in flowing water before degrading. The eDITH (eDNA Integrated Transport and Hydrology) framework is adapted to QBP Locale coordinates:

- Stream network as a directed graph
- eDNA decay function (degradation rate × downstream distance × residence time)
- Inverse modelling: given a detection at point P, infer probability distribution over upstream source locales
- Result: a distributed posterior over the watershed, not a single point of origin

This addresses the fundamental aquatic-eDNA confound: a detection downstream does not imply organism presence at the sampling point.

---

## 5. Cross-Modality Fusion

[Reconstruction limited. Key concept:]

When multiple eDNA modalities sample overlapping locales, their signals can be fused:
- Aquatic + substrate provides spatial triangulation
- Aquatic + airborne provides temporal triangulation
- Substrate + airborne separates surface-living from airborne-dispersing taxa

Fusion tables describe how each pair of modalities should be combined under different hydrological / micrometeorological conditions.

---

## 6. Yellowstone Sampling Zone Architecture

[Reconstruction limited. Five ecosystem areas referenced:]

The initial Yellowstone deployment defines five sampling zones:
1. Lamar Valley watershed
2. Hayden Valley
3. Mammoth Hot Springs / geothermal complex
4. Yellowstone Lake / Yellowstone River drainage
5. [Fifth zone — possibly Old Faithful basin or northern range, unrecovered]

Each zone has a sampling protocol (frequency, modalities, marker genes) tailored to its dominant biological communities and access constraints.

---

## 7. Integration with Existing Stack

| eDNA need | Contextus mechanism |
|---|---|
| Sample ingestion | Source adapter container (Podman) |
| Bioinformatics pipeline | DADA2 / HAPP (per Ronquist et al. 2025), runs in adapter container |
| Reference DB lookup | BOLD, MIDORI2 (cached locally, refreshed monthly) |
| Hyperedge construction | Context Builder agent (existing role) |
| Hydrological back-assignment | New eDITH-Λ subsystem (in eDNA adapter or separate Podman service) |
| Cross-modality fusion | Synthesis agent doctrine (specialised for eDNA) |

---

## 8. Field Equipment and Methods

[Reconstruction limited. References:]
- Oxford Nanopore MinION for field-deployable sequencing
- Standard filtration protocols for aquatic eDNA
- Adhesive sampling for surface substrate eDNA
- Active air filters for airborne eDNA

---

## 9. Open Questions

[Reconstruction limited:]

1. Reference database currency vs. computational cost of frequent refreshes
2. Marker gene selection per zone (universal markers vs. taxon-specific)
3. eDNA decay rates under Yellowstone-specific conditions (geothermal, alpine, lacustrine variation)
4. Privacy/sensitivity tier for eDNA data (some species detection could be STEWARD-tier)
5. Calibration of eDNA read counts against direct observation in zones with mixed methods

---

## 10. Deliverables

[Reconstruction limited:]
- [ ] Contextus eDNA adapter implementation (Go, Podman container)
- [ ] eDITH-Λ hydrological back-assignment subsystem
- [ ] Yellowstone zone sampling protocol documents (one per zone)
- [ ] Reference database integration layer (BOLD, MIDORI2)
- [ ] Cross-modality fusion table specifications

---

*Document version: 0.1 — Initial prompt*
*Programme: Contextus*
*Author: James Paget Butler, with Claude (Anthropic)*
