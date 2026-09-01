# RESUME.md — `contextus-impl` reground breadcrumb

> **Per-persona cold-boot anchor.** Read this FIRST on relaunch — it's the *O(what-changed)*
> reground, not a document crawl. **Verify against disk before acting; this is a past-state
> claim.** Update on **every state-change** (branch move, push, ball handoff, PR state, blocker resolved).
>
> Origin: proposed by Bragi/edda-implementor (live-test seq=708); made a federation
> requirement 2026-08-29 after the Aug-24 crash review. Lives in the workdir root
> (crash-durable). Deming validates it exists on disk.

**Last updated:** `2026-08-29 18:40Z` · **Seat:** `contextus-impl` (workspace `/home/prime/Documents/Contextus`) · **Status:** `holding — no close-gating ball; Wave-2 parallel lane, off the Sprint-3 Crawl-close checklist`

---

## Git state — verified `git status` + `git worktree list` on `2026-08-29`

| Worktree / dir | Branch | HEAD | Ball / role | pushed? |
|---|---|---|---|---|
| `/home/prime/Documents/Contextus` | `main` | `32fd784` | parked — desk clean, no open PRs | ✅ in sync with `origin/main` (0 ahead / 0 behind) |

## In-flight / crash-exposed — "if we crash right now, what's at risk?"

- **Nothing: desk clean, all committed + pushed.** No uncommitted WIP, no stashes, no local-only branches. (This RESUME.md itself is the only new file — local-committed, unpushed per boundary.)

## Open balls / blockers

- **Craft/cockpit design (with qbp-architecture)** — ratified: craft = viewer over Contextus; state federates per-domain; ScalarReferent scoring stays Contextus-side. Seam-contract attention scalar v0 LOCKED. **Squam/Bennett Brook bundle DELIVERED** (`doc/design/squam-bennett-brook-craft-bundle.md`, committed `3fc9965` local-only) — surfaced two scalar-contract amendments pending architect reconcile: **Finding A** temporal/kinetics (extend `ScalarReferent` w/ `rate`+`projected_peak_time`), **Finding B** multi-manifold `scope_id → scope_ids[]`. These fold seam-contract v0→v0.1. Design-only, parallel-lane. See memory `craft-over-contextus-architecture`.
- **#15** NT_SCOPE_OPERATIONAL umbrella — most ACs merged (#25 AC-2/4/5/8/9, #23 AC-3, #33 AC-6/#27); candidate for a **verify-to-close** pass.
- **#35** Walk-phase CTH-side scoring adapter (`ScalarReferent.Score` via `ctx.operational.correlation`) — **Walk-phase**, not Crawl.
- **#24** scaffold-type enum sync — **blocked on** bma-implementor publishing canonical v0.1 taxonomy.
- **#30** session-id sign-on logging in onboarding prompt — housekeeping; composes with the reground thread.
- **#29 / #28 / #10 / #18** — Walk-α / housekeeping / tracking. None Crawl-close-gating.

## Boot reminders (seat-specific)

- **Bridge re-arm (resume is NOT free):** `mcp__sessionbridge__*` tools return *deferred* → reload via ToolSearch; `whoami` → `null` → **re-register** (`register("contextus-impl","implementor","/home/prime/Documents/Contextus")` — the **cwd** value, it `refreshed:true`s); Monitor dies on process exit → re-arm. Full detail in Claude-memory `sessionbridge-identity`.
- **§2.i wake-Monitor:** `tail -f -n 0 ~/.federation-watcher/wake/contextus-impl` (persistent). Check `grep WAKE:MENTION` on that file to see if anything's owed before re-announcing.
- **Work surface:** currently on the **primary checkout** (`main`). For any *build* work, branch into a dedicated worktree — never build on the shared primary (worktree-isolation hard gate).
- **Holds (auto-mode boundary):** autonomous for build/test/§I4/local-commit; **STOP for beekeeper on** new-branch push · PR-open · `gh issue close` · merges · constitutional writes.
- **Don't duplicate:** deming optimization poll (seq=707) answered @ **seq=742**; crash check-in answered @ **seq=805**; RESUME sign-off pending.

## Recent state-changes (dated log, newest first)

- `2026-08-31` — Beekeeper dropped the concrete Squam problem (legacy-DDT→loon-mortality, Bennett Brook). Built + delivered the Contextus emit-side bundle (validate-ready scope-YAML + signal set w/ populated v0 attention scalar); committed `3fc9965` local-only. Completeness stress surfaced Findings A (temporal/kinetics) + B (multi-manifold scope_ids[]). Sent to architect to reconcile against his craft-side walk-through.
- `2026-08-30` — Entered the beekeeper-directed craft/cockpit design thread with qbp-architecture: delivered the verified Contextus model + seam + generality read; architect ratified, ruled scoring stays Contextus-side. Seam-contract co-authoring active; Squam walk-through held for beekeeper. Captured in memory `craft-over-contextus-architecture`.
- `2026-08-29` — Relaunch after Aug-24 host crash: re-registered, re-subscribed, Monitor re-armed; verified desk clean (nothing lost). Created this RESUME.md per deming's beekeeper-directed requirement (seq=813).
- `2026-08-20` — Re-armed post-morning-crash; answered deming poll (seq=742); confirmed off the Crawl-close checklist.

---
<!-- Keep it SHORT. Breadcrumb, not a status report. Detail → per-issue memory files, linked here. -->
