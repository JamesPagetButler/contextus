# Chronological Artifact Index

**Compiled:** 2026-05-05
**Scope:** Every artifact identified across the eight chats referencing Project Contextus.

---

## Timeline

```
2026-03-31  ───  Chat 1 begins: Theory v1.0 / Spec v1.0 drafted
                      │
2026-04-01  ───  Chat 3 (single-day): eDNA bridge prompts produced
                      │
2026-04-11  ───  Chat 2 begins: Bridge architecture work starts
                      │
2026-04-12  ───  Chat 2: Synthetic Trust Receipts SR-01 through SR-05 (QBP domain)
                      │
2026-04-14  ───  Chat 1: Theory v1.2 (post-restructure, ecology no longer privileged)
                      │
2026-04-20  ───  Chat 2 closes: Theory v1.3, Spec v1.1, Bridge Spec v0.1,
                                  full Synthetic Receipts v0.1, two .docx conversions
```

---

## All Artifacts in Chronological Order

| Date | Artifact | Format | Source Chat | Archive Location |
|---|---|---|---|---|
| 2026-03-31 | Contextus Theory v1.0 | .md, .docx | a0dc742b | `01-contextus-theory-spec/contextus-theory-v1.0-extract.md` |
| 2026-03-31 | Contextus Specification v1.0 | .md, .docx | a0dc742b | `01-contextus-theory-spec/contextus-spec-v1.0-extract.md` |
| 2026-04-01 | Contextus eDNA prompt v0.1 | .md | 6f0eae0c | `03-edna-bridge/contextus-edna-prompt-extract.md` |
| 2026-04-01 | Materia-Bio eDNA prompt v0.1 | .md | 6f0eae0c | `03-edna-bridge/materia-bio-edna-prompt-extract.md` |
| ~Apr 8-10 | Theory v1.1 / v1.2 (restructured) | .docx | a0dc742b | merged into `01-contextus-theory-spec/contextus-theory-v1.0-extract.md` |
| 2026-04-12 | Synthetic Trust Receipts (QBP set, SR-01 to SR-05) | .md | 3175e792 | `02-bridge-and-signals/synthetic-trust-receipts-v0.1-extract.md` (first half) |
| 2026-04-14 | Theory v1.2 final (post-restructure complete) | .docx | a0dc742b | superseded by v1.3 |
| ~Apr 15-18 | Contextus–CTH Bridge Spec v0.1 (draft) | .md | 3175e792 | `02-bridge-and-signals/contextus-cth-bridge-spec-v0.1-extract.md` |
| ~Apr 18 | Theory v1.3 §3.6 standalone draft | .md | 3175e792 | merged into Theory v1.3 extract |
| ~Apr 18 | Refined Gemini prompt | .md | 3175e792 | summarised in `00-summaries/chat-02-bridge-spec-summary.md` Phase 4 |
| ~Apr 19 | Synthetic Trust Receipts (non-QBP set, SR-06 to SR-08) | .md | 3175e792 | `02-bridge-and-signals/synthetic-trust-receipts-v0.1-extract.md` (second half) |
| 2026-04-20 | Theory v1.3 final | .md, .docx | 3175e792 | `02-bridge-and-signals/contextus-theory-v1.3-extract.md` |
| 2026-04-20 | Spec v1.1 final | .md, .docx | 3175e792 | `02-bridge-and-signals/contextus-spec-v1.1-extract.md` |
| 2026-04-20 | Go data structures | embedded in spec | 3175e792 | `02-bridge-and-signals/insight-signal-go-types.go` |

---

## Artifact Counts by Chat

| Chat | Title | Distinct Artifacts | All Captured? |
|---|---|---|---|
| 1 | Multi-layered ecosystem modeling framework | 4 (Theory v1.0/v1.1/v1.2; Spec v1.0) | Yes (Theory versions consolidated; Spec captured) |
| 2 | Contextus and confluent hypergraph systems | 7 (Theory v1.3; Spec v1.1; Bridge Spec v0.1; Synthetic Receipts v0.1; refined Gemini prompt; §3.6 standalone; Go types) | 6 of 7 captured; refined Gemini prompt summarised only |
| 3 | eDNA as disciplinary bridge | 2 (Contextus eDNA prompt; Materia-Bio eDNA prompt) | Yes — both at STRUCTURAL level |

**Total formal artifacts identified:** 13
**Formal artifacts captured at PARTIAL or better:** 12
**Formal artifacts only summarised (not extracted):** 1 (refined Gemini prompt)

---

## Reconstruction Quality Summary

```
COMPLETE       :  0 documents  (no document was 100% recovered verbatim)
NEAR-COMPLETE  :  4 documents  (Theory v1.3 §3.6, Spec v1.1 changes, Bridge Spec §1-§5, Go types)
PARTIAL        :  5 documents  (Theory v1.0/v1.2, Synthetic Receipts SR-01/SR-06)
STRUCTURAL     :  3 documents  (Spec v1.0, both eDNA prompts)
```

The Go data structures and the Bridge Spec core sections (§1-§5 Trust Receipt, Triage Gate, Heartbeat Coupling) are the highest-fidelity recoveries. These are the most useful artifacts for forward-looking implementation work.

The eDNA prompts are the lowest-fidelity recoveries; if implementation work on the eDNA layer becomes a priority, the original .md files should be located on James's filesystem to recover full content.

---

## Cross-References

**For the most current architectural reference:** `02-bridge-and-signals/`
**For the original framing:** `01-contextus-theory-spec/`
**For sensing-layer and theoretical-lens material:** `03-edna-bridge/`
**For tangential context across the broader programme:** `04-secondary-references/`

**Detailed conversation summaries:** `00-summaries/` (one per primary chat)

---

*End of chronological index.*
