# Chat 2 — Artifact Manifest

**Chat:** Contextus and confluent hypergraph systems
**URL:** https://claude.ai/chat/3175e792-70f2-48df-bd8e-58762dd066fc
**Period:** 2026-04-11 to 2026-04-20

---

## Artifacts Produced

| # | Artifact | Format | Reconstruction Status | File |
|---|---|---|---|---|
| 1 | Contextus Theory v1.3 | .md, .docx | PARTIAL — §3.6 NEAR-COMPLETE; rest inherits from v1.2 (see chat 1 archive) | `contextus-theory-v1.3-extract.md` |
| 2 | Contextus Spec v1.1 | .md, .docx | NEAR-COMPLETE for v1.0 → v1.1 changes; inherited content referenced | `contextus-spec-v1.1-extract.md` |
| 3 | Contextus–CTH Bridge Spec v0.1 | .md | NEAR-COMPLETE through §6; §7-§9 partial | `contextus-cth-bridge-spec-v0.1-extract.md` |
| 4 | Synthetic Trust Receipts v0.1 | .md | PARTIAL — SR-01 and SR-06 complete; SR-02, SR-03 framing complete; SR-04, SR-05, SR-07, SR-08 framing only | `synthetic-trust-receipts-v0.1-extract.md` |
| 5 | Refined Gemini prompt | .md | NOT SEPARATELY EXTRACTED — content captured in chat-02 summary | (see summary) |
| 6 | Theory v1.3 §3.6 standalone draft | .md | Captured within Theory v1.3 extract | (merged) |
| 7 | Go data structures | .go | RECOVERED with HIGH confidence | `insight-signal-go-types.go` |

---

## Reconstruction Confidence Notes

**HIGH confidence:**
- All field definitions in `InsightSignal`, `TrustReceipt`, `SeedAnchor`, `Heartbeat` Go structs
- Bridge Spec §1-§5 (Purpose, Problem Statement, Core Concepts, Triage Gate, Heartbeat Coupling)
- Theory v1.3 §3.6 (Insight Signal — full five subsections recovered verbatim)
- Spec v1.1 §8 (Bridge Integration), §9 (Proportional Retention), §11 (Go data structures)
- Synthetic Receipt SR-01 and SR-06 full reconstruction
- Three-observation logic for parameter generalisation (evidence cycle / independence / formality)

**MEDIUM confidence:**
- Bridge Spec §6 (NATS subjects) — subjects recovered, payload schemas approximate
- Bridge Spec §7 (Active Evidence Seeking) — process recovered, exact survey source list approximate
- Spec v1.1 §4.4 (Insight Signal emission pipeline) — eight steps recovered, but exact ordering of steps 6-7 may be inverted
- Synthetic Receipts SR-02, SR-03 (framing recovered, body summarised)

**LOW confidence:**
- Bridge Spec §8 (Open Questions) — list approximate, may be incomplete
- Bridge Spec §9 (Roadmap) — Crawl/Walk/Run phases recovered, but specific milestone counts are reconstructions
- Synthetic Receipts SR-04, SR-05, SR-07, SR-08 — framing only, no body content

---

## Documents NOT Recovered

The following artifacts were created in the chat but could not be reconstructed from search snippets:

- The **refined Gemini prompt** itself (referenced in summary; full text not in retrieved snippets).
  → Captured at the conceptual level in `00-summaries/chat-02-bridge-spec-summary.md` Phase 4.

- The exact **.docx formatting** of Theory v1.3 and Spec v1.1 (custom title pages, colour palette, page layout).
  → The .docx files were validated against the docx skill and presented to James in chat. They presumably exist on James's filesystem.

---

## To Recover Remaining Content

If a more complete reconstruction is needed:
1. Locate the .docx files James received (`contextus-theory-v1.3.docx`, `contextus-spec-v1.1.docx`)
2. Locate the original `.md` source files in `/home/claude/` from the past session (these would have been ephemeral)
3. Re-run targeted `conversation_search` queries on remaining sections (e.g. "Contextus refined Gemini prompt standing instructions")

---

## Total Artifact Count

7 distinct artifacts produced in chat 2. All 7 are at least partially captured in this archive.
