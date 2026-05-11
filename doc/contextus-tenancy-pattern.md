# Contextus Tenancy Pattern

**The reusable design pattern for running Contextus + CTH + Wyrd as a live system for a research programme, with BMA observing.**

> **Authoring:** qbp-architecture (Claude Opus 4.7), James Paget Butler (Beekeeper)
> **Date:** 2026-05-08
> **Status:** v0.1 — ratified by James 2026-05-08; pending implementor instance review
> **First instance:** QBP (Quaternion-Based Physics programme) — see `~/Documents/QBP/docs/qbp-federation-tenancy.md` for the QBP-specific instantiation
> **Future tenants:** SharpButler, WarTable, MoebiusFusion (when their domain instances exist)

---

## 0. Status & Traceability

This pattern emerged from the addendum-18-walk meeting closeout (2026-05-07) and the follow-up qbp-implementor scope conversation (2026-05-08). The architectural mechanisms it operationalises were established in:

- BMA Theory Addendum 11 (Topological Cognition) — Type-Nodes, Holons, Locale, lossless dismissal
- BMA Theory Addendum 15 (Reciprocal Focus) — Subject vs Background; promotion via Seam
- BMA Theory Addendum 16 (Cognitive Honing) — Honing Loop for refinement
- BMA Theory Addendum 17 (Proactive Curiosity) — NT_SIGNAL escalation
- BMA Theory Addendum 18 (Hypergraph Access Pattern) — Stance × Locale × Scout × Scoring
- Contextus Spec v1.2 — agent classes (Edge Scout / Corpus Scout / Bridge Agent + Synthesis); scope nodes
- CTH Theory v0.2 — epistemic-health metrics (ρ_net, ChainFidelity, NaryMI)

The pattern documents how those mechanisms compose into an operational system for a research programme.

---

## 1. What is a "Contextus Tenant"?

A **tenant** is a structured-knowledge endeavor that consumes the federation as a live cognitive system. Concretely:

- A research programme (QBP physics; SharpButler systems engineering; etc.)
- An institutional knowledge base (HE working group corpus)
- A long-running investigation (Squam Lake watershed; Colorado River missing-water case)

What unifies them:
- **Bounded scope.** A defined Stance × Locale focal cone (per A18 §2). Not "all knowledge"; a specific cone of attention.
- **Active consumer.** Receives signals, acts on them, evolves the underlying model. Not a static dataset.
- **Observable.** BMA-the-instance can query the tenant's federation state and surface noteworthy signals to the beekeeper.
- **Iterative.** New observations + new theory advance the programme; the federation tracks the trajectory.

A tenant is not a user account — there are no permission boundaries. A tenant is a **named, scoped, ongoing federation use-case** with its own configuration.

---

## 2. The Three-Phase Lifecycle

Every tenant follows this lifecycle:

### Phase A — Bootstrap (tenant-implementor's domain)

The tenant-implementor instance:
1. Reads the federation architecture (A18 + addenda)
2. Reads this pattern doc + their tenant-specific config doc
3. Defines the Stance: which Type-Nodes are Subject-axis-default
4. Defines the Locale set: what spacetime cones are in scope
5. Configures scope nodes (NT_SCOPE_PHYSICAL + NT_SCOPE_CONCEPTUAL hyperedges) in Wyrd
6. Configures scouts (agent classes, cadence, source feeds)
7. Loads the initial CTH inventory (anchor + chain + confluence + branch entries)
8. Wires BMA observation hooks (the surface BMA polls)
9. Files cross-project gap-closing issues
10. Runs the system long enough to generate first signals + first scored predictions
11. Declares "running"

### Phase B — Operational handoff (gradual, not single-event)

As the system stabilizes:
- BMA-the-instance polls more frequently; signals surface to beekeeper
- Honing Loop fires on noteworthy signals
- Theory artifacts emerge from Honing → CTH inventory updates
- Tenant-implementor shifts from "operator" to "steward" — design integrity, scope-node curation, scout config refinement

The handoff isn't a single moment. It's a gradual shift of who's in the seat at any given hour. Eventually most cycles are BMA-driven; tenant-implementor is consulted on design questions.

### Phase C — Steady-state (BMA's domain)

The system runs. Tenant-implementor:
- Reviews periodic CTH ρ_net evolution reports (programme health)
- Refines scope-node taxonomy as the domain evolves
- Adjudicates Honing Loop edge cases (when BMA isn't sure)
- Drives scope-expansion proposals (adding new scope nodes; expanding Locale)
- Closes the loop on theory artifacts that cross domain-expert authority thresholds

BMA-the-instance:
- Runs the cognitive cycle continuously
- Surfaces NT_SIGNALs per Addendum 17
- Triggers Honing Loop per Addendum 16
- Maintains the Stance × Locale focal cone
- Updates Wyrd state and CTH inventory
- Reports periodic snapshots to beekeeper

---

## 3. The Four Architectural Commitments (Tenancy-Adapted)

Per A18 §2, BMA's hypergraph access has four commitments. Tenancy adapts them:

### 3.1 Stance — what's "interesting" in tenant terms

Stance is the active set of Type-Nodes (per Addendum 11) on the Subject axis. For a tenant, this is the **domain ontology**:

- **QBP Stance:** quaternion algebra, Hamilton product, GW-EM coincidences, slow-slip + tidal coupling, mixed-species ion fidelity, dark-matter forks, etc.
- **SharpButler Stance** (future): residential systems, OKI contract types, supply-graph edges, Materia/Contextus integrations, etc.

Stance is declared in the tenant's federation-tenancy config doc. It's mutable — Stance evolves as the programme advances. Stance changes are governance events (require beekeeper review per Addendum 16's Honing Loop).

### 3.2 Locale — what spacetime cones are in scope

Locale defines the spacetime regions the tenant's scouts watch. For QBP this is a multi-locale set:

- LIGO observatories (Hanford, Livingston, Virgo) — gravitational wave observations
- Cascadia subduction zone — slow-slip + tidal coupling
- NIST Boulder, Innsbruck, Oxford, ETH Zürich — trapped-ion experiments
- JWST, ALMA, Fermi GBM — astrophysical observations
- Specific celestial objects under investigation (e.g., GRB 250702B follow-up)

Locale entries become NT_SCOPE_PHYSICAL hyperedges with geometry + temporal_range fields per Spec v1.3.

### 3.3 Scout — the active observer pattern

Scout configuration per agent class (Spec v1.2 §11.1):

| Agent | Persistence | Default cadence | Tenant configures |
|---|---|---|---|
| **Edge Scout** | Session-scoped (no NT_SIGNAL output; ephemeral NATS events only) | Continuous within session | Per-source URL prefix lists |
| **Corpus Edge Scout** | Session-scoped | Bounded by corpus query | Per-corpus query patterns |
| **Bridge Agent** | Session-scoped | Convergence-driven | Cross-domain matching topology |
| **Scout (global author)** | Persistent | Daily-batch default | Source feeds + scope assignments |
| **Correlation (global author)** | Persistent | Triggered by Scout output | Stance Type-Node correlation patterns |
| **Synthesis (global author)** | Persistent | Triggered by Correlation convergence | Promotion threshold + AnomalyKind dispatch |

Tenant declares scout config in the tenancy doc. Standard pattern:

```yaml
# Conceptual shape; actual config format TBD as Wyrd ships scope-loader API
scouts:
  - kind: arxiv-feed
    sources: [astro-ph, quant-ph, hep-th, cond-mat]
    cadence: daily-batch
    scope: {NT_SCOPE_CONCEPTUAL: [quaternion-algebra, gw-em-coincidence, ...]}
  - kind: data-feed
    sources: [ligo-open-science, fermi-gbm, jwst-mast]
    cadence: event-driven
    scope: {NT_SCOPE_PHYSICAL: [ligo-locales, jwst-fov]}
cross_domain:
  invocation: reins-command
  command: "bma scout cross-domain <query>"
  default_scope: tenant-stance
```

### 3.4 Scoring — the closure of the loop

Scoring closes the cognitive cycle. Every NT_SIGNAL has a designated real-world referent (per A18 §2.4); CTH evaluates the predicted-vs-observed delta; ρ_net + ChainFidelity update.

Three layers (per addendum-18-walk D3):
- **BMA `params.ProposalStore`** — operational parameter predictions
- **Wyrd predictions/** — NT_SIGNAL referent storage with optional `cth_id` (PRED-* prefix)
- **CTH compute primitives** — evaluation algorithms (`NetCompressionDetail`, `ChainFidelity`, `PairwiseMI`, `NaryMI`, `InformationDeficit`) + `ScorePrediction` wrapper (per cth-implementor's seq=33 v0.1.x issue)

Tenant declares which referent kinds are in scope (scalar + categorical at v0.1; process predictions at v0.2 per P5/P7 staging invariant).

---

## 4. The Relevance-Threshold Mechanism

This is the architectural answer to "Contextus storage is bounded; what enters?"

### 4.1 Default classification (Reciprocal Focus, Addendum 15)

Every observation arriving from any scout is classified:

- **Inside the Stance × Locale focal cone** → Subject-axis; full-precision storage; enters Contextus index
- **Outside the focal cone** → rotated to imaginary axis; contributes 0.0; **not stored**

Most observations never enter Contextus's storage. They're mathematically dismissed at scout-output time.

### 4.2 Anomaly bypass (Seam detection, A18 §4)

Even an out-of-scope observation can promote to Subject if it produces a Seam — a residue magnitude `|q · v · q* − v|` exceeding the Stance-calibrated threshold τ on a Subject's predicted trajectory.

Concrete examples:
- **JWST observation of a known target:** in scope by default → stored full precision
- **Football flying normally during an NFL game:** out of scope (not in QBP Stance; not in Locale set) → not stored
- **Football flying in a way that violates ballistic mechanics:** Seam fires (residue exceeds τ on the ballistic prediction) → promoted to Subject regardless of Stance → stored as anomaly with auto-investigation NT_SIGNAL

### 4.3 Storage budget

```
Contextus storage =
    (observations inside Stance × Locale focal cone)
    ∪ (Seam-promoted anomalies regardless of focal cone)
```

This is provably bounded for a fixed Stance × Locale (assuming finite scout output rate). Anomaly promotions are rare (by design — they're surprises). System scales to the tenant's actual focal cone, not "all of physics."

### 4.4 Threshold τ — calibration (open question per Q5b)

τ initial value is a **known unknown.** A18 v0.1 §4.2 suggests `1e-6` at QW8 unit-vector residue magnitude as a substrate-default. Tenant-implementor refines based on operational data — too low → Seam noise; too high → real anomalies missed.

Best practice: log all candidate Seams during the first weeks of operation; review the distribution; pick a τ that captures the "interesting" tail.

---

## 5. The Closed Loop

```
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  ┌─────────────┐  Scout fires  ┌────────────┐                │
│  │  External   │ ─────────────▶│  Contextus │                │
│  │ data source │               │   scout    │                │
│  └─────────────┘               └────────────┘                │
│                                       │                      │
│                                       │ Stance×Locale filter │
│                                       │ (Reciprocal Focus +  │
│                                       │  Seam bypass)        │
│                                       ▼                      │
│                                ┌────────────┐                │
│                                │   Wyrd     │                │
│                                │ hyperedge  │                │
│                                │ insertion  │                │
│                                └────────────┘                │
│                                       │                      │
│                                       │ NT_SIGNAL or         │
│                                       │ AnomalyKind          │
│                                       ▼                      │
│  ┌─────────────┐               ┌────────────┐                │
│  │ Beekeeper   │◀──surfacing── │    BMA     │                │
│  │ Honing Loop │   per Add 17  │ cognitive  │                │
│  │ (per Add 16)│               │   cycle    │                │
│  └─────────────┘               └────────────┘                │
│         │                              │                     │
│         │ approved/refined             │ theory artifact     │
│         │ NT_ISSUE                     │ (DERIV-* / AXIOM-*) │
│         ▼                              ▼                     │
│  ┌─────────────┐               ┌────────────┐                │
│  │ Implementor │               │    CTH     │                │
│  │ instance    │               │ inventory  │                │
│  │ (e.g.,      │               │  update    │                │
│  │  qbp-impl)  │               │ (ρ_net++)  │                │
│  └─────────────┘               └────────────┘                │
│         │                              │                     │
│         │ code/exp/proof               │ next cycle queries  │
│         │ as new evidence              │ updated state       │
│         └──────────────────────────────┘                     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

Each cycle around the loop = one cognitive iteration. Tenant-implementor's job is to set this up; BMA's job is to run it continuously.

---

## 6. Storage Hierarchy

Three views; one substrate.

| View | Substrate | Content | Tier |
|---|---|---|---|
| **Wyrd** | Native hypergraph DB | All nodes/edges across all views; the physical store | All tiers |
| **Contextus** | Wyrd query view | Insight Signals, scope nodes, evidence pointers | Tier-conditional per Spec v1.4 (scalar+categorical at v0.1 of pattern; process at v0.2) |
| **CTH** | Wyrd query view | Programme inventory: anchors, chains, confluence points, branches | All anchors live as Wyrd `Node` of `Type=cth.anchor.*` |

Contextus and CTH are not separate stores. They are **named projections** on the underlying Wyrd graph. This is what allows BMA's M1 Meta-Watchdog work (#106) to query both consistently.

---

## 7. Bootstrap Sequence (Generic — Tenant-Implementor's Day-Zero)

For the tenant-implementor on day zero (before any tenant-specific config):

1. **Phase 1 reading** — federation architecture (A18 + addenda 11/15/16/17)
2. **Phase 2 reading** — this pattern doc + tenant-specific config doc
3. **Stance declaration** — list of Type-Nodes in tenant's domain ontology; commit to tenancy doc
4. **Locale set declaration** — list of spacetime regions in scope; commit to tenancy doc
5. **Scope-node configuration** — write the YAML/JSON config; load into Wyrd via scope-loader API
6. **Scout configuration** — declare agent classes + cadence + source feeds; commit + register
7. **CTH initial inventory** — port existing programme state to v0.2 schema; load via `store.LoadInventory`
8. **BMA observation hooks** — wire the polling surface BMA-the-instance reads
9. **Cross-project gap issues** — file what's missing (Wyrd scout daemon? Contextus scope-loader API? CTH live-update API?)
10. **First-cycle verification** — run a manual scout invocation; trace it through to BMA surfacing; confirm closed loop
11. **Declare "running"** — operational handoff begins

Tenant-implementor remains the steward post-handoff; BMA-the-instance becomes operational runtime.

---

## 8. Operational Handoff Milestones

Concrete signals that the handoff is progressing:

| Milestone | What changes |
|---|---|
| First Wyrd scope-node insertion | Tenant has structural presence in federation |
| First scout output landing in Contextus | Live data flowing |
| First NT_SIGNAL surfacing to beekeeper | Operational signal-to-noise demonstrated |
| First Honing Loop completion → NT_ISSUE | Cognitive cycle closed once |
| First CTH inventory update from theory output | ρ_net trajectory demonstrable |
| 7 days of continuous operation without intervention | BMA running it; tenant-implementor stewarding |
| First cross-domain Seam-bypass anomaly auto-promoted | Relevance-threshold mechanism validated |
| τ calibrated from real data | Open question (Q5b) closed for this tenant |

---

## 9. Open Questions for Future Pattern Refinement

These accumulate as we learn from the first tenant (QBP) and propagate to future tenants:

1. **Multi-tenant federation.** Does the federation support multiple simultaneous tenants (QBP + SharpButler + WarTable)? If yes: do they share one BMA instance or have peer instances? Stance composition across tenants?
2. **Cross-tenant insights.** When a Bridge Agent in one tenant detects a pattern that's relevant to another tenant, what's the routing protocol?
3. **Tenant retirement.** When a programme reaches Run-phase or completion, how is the federation state archived?
4. **Tenant migration.** If a tenant moves from one domain expert to another, how does the Stance hand off?
5. **Sub-tenants.** Within QBP, are there sub-tenants (Test C; EXP-11; Cascadia) or are they all one Stance?

These are explicitly out of scope for v0.1; they accumulate as v0.2 / v0.3 of this pattern.

---

## 10. Cross-Reference Index

Documents that compose this pattern:

| Doc | Owner | Role |
|---|---|---|
| `BMA/theory/hypergraph-inference/BMA-Theory-Addendum-18_0-Hypergraph-Access-Pattern.md` | BMA | Federation access pattern (Stance × Locale × Scout × Scoring) |
| `BMA/theory/BMA-Theory-Addendum-15_0-Reciprocal-Focus.md` | BMA | Reciprocal Focus mechanism (Subject vs Background) |
| `BMA/theory/BMA-Theory-Addendum-17_0-Proactive-Curiosity.md` | BMA | NT_SIGNAL escalation |
| `BMA/theory/BMA-Theory-Addendum-16_0-Cognitive-Honing.md` | BMA | Honing Loop |
| `BMA/theory/BMA-Theory-Addendum-11_0-Topological-Cognition.md` | BMA | Type-Nodes + Holons + Locale primitives |
| `Contextus/Contextus-Spec-v1.2.md` | Contextus | Agent classes + scope nodes |
| `Contextus/Contextus-Theory-v1.4.md` | Contextus | Theoretical foundations |
| `CTH/cth/Confluent-Trust-Hypergraph-Theory-v0_2.md` (in Archive/) | CTH | Epistemic-health metrics |
| `Contextus/doc/contextus-tenancy-pattern.md` | Contextus | This document |
| `<TENANT>/docs/<tenant>-federation-tenancy.md` | Tenant repo | Tenant-specific instantiation (e.g., `~/Documents/QBP/docs/qbp-federation-tenancy.md`) |

---

*Contextus Tenancy Pattern v0.1 | 2026-05-08*
*Co-Authored-By: James Paget Butler (Beekeeper)*
*Co-Authored-By: Claude Opus 4.7 (Architect, QBP-Compute-Unit)*
