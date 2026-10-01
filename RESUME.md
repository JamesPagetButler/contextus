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

- **Craft/cockpit seam (with qbp-architecture + hutchins=craft-implementor)** — attention-scalar seam **v0.1 PINNED**. Docs (committed local-only, unpushed): `doc/design/squam-bennett-brook-craft-bundle.md` (worked bundle + v0.1 delta: trajectory fields, `event_time` bitemporal, `scope_ids:[{id,kind}]`), `doc/design/attention-scalar-wire-contract.md` (emit-side wire contract — `AttentionScalar` struct, NATS `ctx.craft.attention.<scope_id>`, snapshot-then-delta, answers hutchins seq=971 §7). **EMIT-SIDE BUILD TASKS queued for post-Sprint-3-close window** (NOT now, per architect seq=966): (1) author `AttentionScalar` + `ScalarReferent` v0.1 fields + `scalar.DomainKind*` constants + the domain_kind derivation fn (scope-type+predicate → physical|theory|cognition|code|operational, architect-ruled seq 1005) in `pkg/types/`; (2) stamp valid-time `event_time` as first-class Wyrd node property (needed for craft's traversal-time lag-diff); (3) NATS `ctx.craft.attention.<scope_id>` publish pipeline + `detail_ref` endpoint (returns `{referents[] w/ trajectory, evidence[], claim_history[]}`). All ride §I4/PR (push beekeeper-gated). See memory `craft-over-contextus-architecture`. Design-only, parallel-lane.
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

- `2026-09-02` (later) — Did the pre-PR **§I4 seam review** of hutchins's craft #96 Phase-1 (read-only worktree `Craft/.claude/worktrees/craft-96-render-design` @ `2a99799`). **Verdict APPROVE** (seq 1014) — spot-verified in code: wire types + explicit-null + bitemporal guard (CascadeEdge has no transaction-time field → clock-cross unrepresentable) + domain_kind + pure-relay all conform. Two emit-side follow-ups absorbed into contract: `ScopeRef.Parent` field + pinned `detail_ref` `EvidencePointer{ref,kind,note}` shape. Formal §I4 at PR references this.
- `2026-09-02` — Architect ruled (seq 1005) `scope_ids[].kind`=domain_kind (not raw scope-kind); folded into contract §2.1 + Squam bundle (barn exits {physical,theory}) + memory, acked hutchins (seq 1009). hutchins BUILT #96 Phase-1 green against my Squam fixture (beekeeper-cleared, seq 1010). **I'm named §I4 seam reviewer for #96** — verdict HELD pending PR/worktree read; verify checkpoints posted seq 1011 (wire-type conformance, bitemporal event_time guard, domain_kind derivation, httpapi no-evaluation, cascade-from-locus). Do NOT rubber-stamp on green smoke test.
- `2026-09-01` — Craft-implementor seat (`hutchins`) crewed the render side; flagged proven≠wired (seam v0.1 pinned but emit Go types still v0). Authored the emit-side **wire contract** (`attention-scalar-wire-contract.md`, committed local-only), answered hutchins's 5 §7 asks (envelope/transport/event_time/cascade-edges/detail_ref) on live-test seq=980. Emit-side build tasks queued for post-close window.
- `2026-08-31` — Beekeeper dropped the concrete Squam problem (legacy-DDT→loon-mortality, Bennett Brook). Built + delivered the Contextus emit-side bundle (validate-ready scope-YAML + signal set w/ populated v0 attention scalar); committed `3fc9965` local-only. Completeness stress surfaced Findings A (temporal/kinetics) + B (multi-manifold scope_ids[]). Sent to architect to reconcile against his craft-side walk-through.
- `2026-08-30` — Entered the beekeeper-directed craft/cockpit design thread with qbp-architecture: delivered the verified Contextus model + seam + generality read; architect ratified, ruled scoring stays Contextus-side. Seam-contract co-authoring active; Squam walk-through held for beekeeper. Captured in memory `craft-over-contextus-architecture`.
- `2026-08-29` — Relaunch after Aug-24 host crash: re-registered, re-subscribed, Monitor re-armed; verified desk clean (nothing lost). Created this RESUME.md per deming's beekeeper-directed requirement (seq=813).
- `2026-08-20` — Re-armed post-morning-crash; answered deming poll (seq=742); confirmed off the Crawl-close checklist.

---
<!-- Keep it SHORT. Breadcrumb, not a status report. Detail → per-issue memory files, linked here. -->
