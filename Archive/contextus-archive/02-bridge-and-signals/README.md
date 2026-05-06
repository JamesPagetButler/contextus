# 02 — Bridge, Signals, and Architectural Maturity (Chat 2)

Source: https://claude.ai/chat/3175e792-70f2-48df-bd8e-58762dd066fc
Period: 2026-04-11 to 2026-04-20
Detailed summary: [../00-summaries/chat-02-bridge-spec-summary.md](../00-summaries/chat-02-bridge-spec-summary.md)

## This is the architecturally most mature directory in the archive

If you are returning to Contextus work after a gap, **start here**.

## Files in this directory

| File | What it is |
|---|---|
| `contextus-theory-v1.3-extract.md` | Theory document with §3.6 Insight Signal added. Design Principle 10 added. |
| `contextus-spec-v1.1-extract.md` | Specification with InsightSignal schema, bridge integration, proportional retention, Go data structures. |
| `contextus-cth-bridge-spec-v0.1-extract.md` | The Bridge Spec — the missing translation layer between Contextus and CTH. |
| `synthetic-trust-receipts-v0.1-extract.md` | Eight synthetic receipts across QBP, eDNA, and Materia-Bio. Parameter calibration exercise. |
| `insight-signal-go-types.go` | Go type definitions for InsightSignal, TrustReceipt, SeedAnchor, Heartbeat, ThresholdConfig. |
| `ARTIFACT-MANIFEST.md` | Per-artifact reconstruction completeness notes |

## Architectural decisions captured here

1. **Strict separation of concerns:**
   - Contextus = pattern detection (Insight Signals)
   - Bridge = claim translation (Trust Receipts)
   - CTH = epistemic evaluation (trust scores, derivation chains)

2. **Trust Receipt with provenance hash** — full Contextus context is reachable but not present; CTH operates on the claim, not the upstream evidence.

3. **Heartbeat coupling, unidirectional** — Contextus → CTH only. Goes dormant once seed anchor promotes to Tier 2 on its own evidence. Floor never zero.

4. **Proportional retention five tiers** — Core / Active / Watching / Background / Skeleton, governed by Ebbinghaus decay, with a continuous throttle for storage capacity.

5. **Parameter generalisation** — three parameters previously flagged as "domain-configurable" actually derive from per-insight properties:
   - decay half-life ← evidence cycle time
   - heartbeat bit fraction ← evidence independence
   - internal evidence resistance ← evidence formality

## Crawl exit gate

Per the Bridge Spec roadmap: **50 Trust Receipts through full lifecycle + completion of synthetic receipt simulator runs (SR-01 through SR-08)** before scaling to Walk phase.

## Standing principles to preserve

- Signals are never conclusions
- Signals are always traceable
- Signals are never suppressed
- Heartbeat is unidirectional (Contextus → CTH)
- The Skeleton retention tier is the floor — fading to zero would erase history
