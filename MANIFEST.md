# Contextus — Ecosystem Insight Discovery — Document Manifest

## Current Documents

### Root (active, current versions)

| File | Version | Description |
|---|---|---|
| `Contextus-Theory-v1.5.md` | v1.5 | **Current theory.** Formalises the *structural anomaly* (fourth class, §3.6.2) for findings that describe the structure of the search/exploration process; adds §3.6.6 (Synthesis as Persistence Boundary); reframes the Colorado River failure as worked example (§7.6). Closes issue #3. |
| `Contextus-Theory-v1.4.md` | v1.4 | Superseded by v1.5 — kept at root for the duration of the v1.5 review window per ADR-003 §I4; will move to `Archive/` after v1.5 PR merges. |
| `Contextus-Spec-v1.3.md` | v1.3 | **Current spec.** Closes the four v1.2 open architectural questions: §4.6 Scope Nodes (`NT_SCOPE_PHYSICAL`, `NT_SCOPE_CONCEPTUAL`, `HE_SCOPE_MEMBERSHIP`); §5.4 Evidence Pointer Discipline (tier-conditional fields, cap-per-tier eviction); §4.4 Synthesis-as-persistence-boundary clause; §11.1 `EvidencePointer` + `AnomalyStructural`; §11.4 scope-node Go types. Resolves Wyrd issue [#6](https://github.com/JamesPagetButler/wyrd/issues/6) `SignalSource` enum. Reviewed under qbp-compute-unit ADR-003 §I4. |
| `contextus-wyrd-integration-architecture-2026-05-05.md` | — | Architecture-instance integration doc; resolves Wyrd issue #6; inputs to Spec v1.3. |
| `doc/contextus-impl-onboarding-prompt.md` | — | Bootstrap prompt for fresh contextus-impl sessions. |

### Archive (historical versions and source documents)

| File | Version | Notes |
|---|---|---|
| `Contextus-Spec-v1.2.md` | v1.2 | **Superseded by v1.3** — kept at root for the duration of the v1.3 review window per ADR-003 §I4; will move to `Archive/` after v1.3 PR merges. |
| `Archive/contextus-theory-v1.3.md` | v1.3 | Superseded by v1.4 |
| `Archive/contextus-spec-v1.1.md` | v1.1 | Superseded by v1.2 |
| `Archive/Contextus-Theory-Addendum-v1.3.md` | Addendum | Colorado River test + surveillance mode — content now integrated into Theory v1.4 |
| `Archive/contextus_theory_addendum_v1.3.docx` | Addendum | Source .docx for the addendum |
| `Archive/contextus-theory-section-3.6-insight-signal.md` | — | Standalone §3.6 draft — content integrated into Theory v1.3+ |
| `Archive/contextus-cth-bridge-spec-v0.1.md` | v0.1 | CTH Bridge Specification (standalone — not superseded) |

## Needs Retrieval from Chat

### docs/
- [ ] eDNA Contextus Prompt
  - **Source:** https://claude.ai/chat/6f0eae0c-ce00-4ca1-924b-7ec420f6496b
  - **Search terms:** `*edna*contextus*prompt*`
- [ ] eDNA Materia-Bio Prompt
  - **Source:** https://claude.ai/chat/6f0eae0c-ce00-4ca1-924b-7ec420f6496b
  - **Search terms:** `*edna*materia*prompt*`

## Not Yet Written

- Contextus–CTH Bridge Synthetic Receipts v0.1 (parameter calibration test cases)
- Contextus Governance Document
