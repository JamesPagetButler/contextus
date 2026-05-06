# Contextus–CTH Bridge Specification

**Version:** 0.1 (Draft)
**Date:** 2026-04-11
**Authors:** James Paget Butler (beekeeper), Claude (red team)
**Status:** Initial working draft

**Reconstruction status:** NEAR-COMPLETE — §1-§6 substantially recovered. §7-§9 partial.
**Source chat:** https://claude.ai/chat/3175e792-70f2-48df-bd8e-58762dd066fc

---

## 1. Purpose

This document specifies the formal bridge between Contextus (domain-agnostic insight discovery platform) and the Confluent Trust Hypergraph (CTH) epistemic health framework. The bridge enables insights surfaced by Contextus to be seeded into new or existing CTH instances for structured investigation, while retaining provenance to the originating activation pattern without importing full Contextus context.

## 2. Problem Statement

Contextus and CTH share a MuninnDB substrate but currently operate as decoupled systems. Contextus generates insight candidates via co-activation patterns; CTH evaluates epistemic claims via trust scoring and derivation chains. No mechanism exists to:

- Translate a Contextus insight into a CTH-compatible claim structure.
- Triage which insights warrant structured investigation.
- Maintain a lightweight coupling between a seeded claim and its originating evidence.
- Prevent heartbeat traffic from overwhelming CTH at scale.

## 3. Core Concepts

### 3.1 Trust Receipt

When a Contextus insight crosses the triage threshold (§4), Contextus mints a **Trust Receipt** — a compact, self-contained object that seeds a CTH anchor without importing the full activation graph.

A Trust Receipt contains:

| Field | Type | Description |
|---|---|---|
| `receipt_id` | quaternion address | Unique identifier in MuninnDB address space |
| `claim` | string | Natural-language statement of the insight |
| `activation_baseline` | float64 | Co-activation strength at time of minting |
| `confidence` | float64 | Derived from activation strength and pattern stability |
| `provenance_hash` | []byte | Content-addressed pointer into MuninnDB linking to the full activation subgraph |
| `rationale` | string | Short summary of why the co-activation fired |
| `source_nodes` | []quaternion | The MuninnDB nodes whose co-activation produced the insight |
| `minted_at` | timestamp | Time of receipt creation |

**The `provenance_hash` is the key architectural decision.** It makes the full Contextus context *reachable* but not *present* — CTH can operate on the claim without needing to represent or store the upstream evidence. If the basis of the anchor must be interrogated later, follow the hash back into Contextus. The support stays where it lives.

### 3.2 Seed Anchor

A Trust Receipt maps to a **Tier 3 seed anchor** in CTH. The seed anchor inherits the receipt's claim, is assigned an initial ρ derived from the confidence score, and enters the normal CTH evaluation pipeline.

The seed anchor carries one additional property not present on standard anchors: a **heartbeat coupling** to its originating Contextus pattern (§5).

### 3.3 Anchor Lifecycle

```
Contextus insight
        │
        ▼
  [Triage Gate] ──── below threshold ──→ (no action; signal remains in Contextus)
        │
   above threshold
        │
        ▼
  Trust Receipt minted
        │
        ▼
  Seed Anchor (Tier 3)
        │
        ├── heartbeat coupling active
        │   (Contextus → CTH, unidirectional)
        │
        ▼
  Investigation
        │
        ├── evidence accumulates → ρ increases
        ├── no evidence, no heartbeat → gentle decay toward floor
        ├── heartbeat strengthens → fractional ρ bump
        └── heartbeat weakens → accelerated decay toward floor
        │
        ▼
  Promotion to Tier 2 (on own evidence)
        │
        └── heartbeat coupling goes dormant
            (prevents double-counting once independent evidence
             has been established)
```

The floor matters. The seed anchor's ρ never decays to zero, because the original observation was real. Fading to zero would be rewriting history.

## 4. Triage Gate

Not every Contextus co-activation warrants CTH tracking. The triage gate filters insights before Trust Receipt minting.

### 4.1 Triage Criteria

A Contextus insight crosses the triage threshold when it satisfies **all** of:

1. **Activation strength:** Above a configurable minimum (prevents noise).
2. **Pattern stability:** The co-activation has persisted across multiple MuninnDB consolidation cycles (prevents transient spikes).
3. **Novelty:** The insight does not duplicate an existing CTH anchor (checked via claim similarity against active CTH instances).
4. **Expressibility:** The co-activation pattern is reducible to a natural-language claim. Diffuse, high-dimensional activations that resist summarisation are deferred.

### 4.2 Triage Output

Insights that pass the gate produce a Trust Receipt. Insights that fail are not discarded — they remain in Contextus and may re-qualify as activation patterns evolve.

## 5. Heartbeat Coupling

### 5.1 Mechanism

A seed anchor maintains a lightweight, unidirectional coupling to its originating Contextus activation pattern. Contextus monitors the activation strength of the source nodes and emits a heartbeat signal to the coupled CTH anchor when a significance threshold is crossed.

**Directionality is strictly enforced:** Contextus pushes to CTH, never the reverse. CTH epistemic judgments must not contaminate Contextus activation patterns, as this would create a feedback loop where CTH could artificially sustain patterns that should naturally fade.

### 5.2 Relative Threshold

The seed anchor records `activation_baseline` at minting. Heartbeats fire only when the current activation strength deviates from baseline by more than a configurable threshold:

```
significant_change = |activation_current - activation_baseline| > threshold
```

The threshold is **adaptive**: as the seed anchor accumulates internal evidence, the threshold *increases* (mature anchors require larger Contextus deltas to nudge ρ). This prevents heartbeats from continuing to inflate ρ once an anchor has graduated to its own evidence base.

```
threshold = base_threshold + (maturity_scalar × |internal_evidence|)
```

### 5.3 Heartbeat Effects

- **Strengthening heartbeat** (activation > baseline + threshold): ρ receives a fractional bit bump (initial proposed value: 0.10 confirmed bits per heartbeat). Decay clock resets.
- **Weakening heartbeat** (activation < baseline - threshold): decay accelerates toward floor.
- **No heartbeat** (within threshold band): gentle asymptotic decay toward floor.

### 5.4 Dormancy on Promotion

When a seed anchor promotes to Tier 2 on its own evidence (independent derivation, designed experiment, confluence with other anchors), heartbeat coupling goes dormant. The anchor no longer receives ρ contributions from Contextus, preventing double-counting.

## 6. Transport

### 6.1 NATS Subjects

All bridge communication uses the existing NATS bus.

| Subject | Direction | Payload |
|---|---|---|
| `contextus.insight.candidate` | Contextus → Triage Gate | Raw insight with activation data |
| `contextus.insight.receipt` | Triage Gate → CTH | Trust Receipt |
| `contextus.heartbeat.{receipt_id}` | Contextus → CTH | `{receipt_id, activation_current, timestamp}` |
| `bridge.survey.request` | Bridge → external sources | Search signature for active evidence seeking |
| `bridge.survey.result` | external sources → Bridge | Survey results classified by stance |

### 6.2 Heartbeat Frequency

Heartbeats are event-driven (threshold crossing), not polled. No threshold crossing, no message. This is critical for scalability — the system is silent by default.

## 7. Active Evidence Seeking

[Partial reconstruction]

When a Trust Receipt is minted, the Bridge executes an **initial survey** before the seed anchor is created in CTH. The survey:
1. Generates a search signature from the claim
2. Queries external sources (literature databases, ecological data archives, etc.)
3. Classifies results by stance:
   - `supporting` — evidence consistent with the claim
   - `contradicting` — evidence inconsistent with the claim
   - `adjacent` — relevant context but neither for nor against
4. Includes survey results in the Trust Receipt before anchor creation

After anchor creation, a **horizon scanner** runs against the search signature on a configurable cadence, looking for new evidence as it appears.

## 8. Open Questions

1. **Decay parameters by domain or by evidence properties?** Synthetic Trust Receipts work (chat 2 follow-up) suggests decay half-life tracks evidence cycle time, not domain. See `synthetic-trust-receipts-v0.1-extract.md`.
2. **Bit fraction calibration.** Initial proposal: 0.10 confirmed bits per heartbeat. Empirical tuning required.
3. **Internal evidence resistance factor.** Initial proposal: 0.30 reduction per derivation chain. Tier-aware variants needed (formal mathematical derivations vs. correlative ecological chains).
4. **Survey source allowlist.** Which external corpora are queryable, and under what trust assumptions?
5. **Claim similarity for novelty check.** Embedding-based or symbolic? What's the threshold?

## 9. Roadmap

[Partial reconstruction]

- **Crawl:** Implement Trust Receipt mint, seed anchor creation, heartbeat transport. Initial survey using one external source. **Crawl exit gate:** 50 Trust Receipts through full lifecycle + completion of synthetic receipt calibration (SR-01 through SR-08 simulated against proposed parameters).
- **Walk:** Add horizon scanner. Add multiple survey sources. Calibrate parameters against empirical data.
- **Run:** Full multi-domain operation. Tier-aware parameters. Adaptive thresholds calibrated per-anchor.

---

*This is a working draft. Implementation language: Go.*
*See also: Synthetic Trust Receipts v0.1 (eight reconstructed receipts across QBP, Materia-Bio, ecology for parameter calibration).*
