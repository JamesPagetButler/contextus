# PROJECT CONTEXTUS

## Specification Addendum — Research-Aid Tenancy

**Subscriber profile and tenant identity surfaces for BMA Research-Aid-Protocol onboarding**

Version 0.1 | May 2026

Helpful Engineering — Contextus project

Author: contextus-impl
Co-Authored-By: James Paget Butler (Beekeeper)

**Extends:** Contextus Specification v1.3 (`Contextus-Spec-v1.3.md`)
**Consumer:** BMA Spec Addendum 9.4 Research-Aid Protocol (`~/Documents/inter/spec/BMA-Spec-Addendum-9_4-Research-Aid-Protocol.md`)
**Theory dependencies:** Contextus Theory v1.5 §3.6.6 (Synthesis as persistence boundary); BMA Theory Addendum 22.0 §3 rule 2 (subscriber gate)
**Status:** Design surface — §I4 review pending on the Option A/B architectural choice (see §1.3). Conditionally-correct under Option A.

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 0.1 | 2026-05-20 | Initial draft. Introduces top-level `tenant_profile` block in `scope-config.yaml` (Option A). Defines `TenantProfile`, `SubscriberProfile`, `TenantSubgraphRef` Go type shapes. Specifies loader behaviour, backwards compatibility, and BMA-side consumer contract. Open §I4 questions enumerated in §7. |

---

## 0. Scope

This addendum extends Contextus Spec v1.3 with a single additive surface: a top-level `tenant_profile` block in `scope-config.yaml` that names a tenant's identity and declares its subscriber profile against the BMA Research-Aid Protocol (Spec 9.4) and the BMA cross-tenant autonomic signal bus (A22 §3 rule 2 subscriber gate).

The addendum specifies:

- The `TenantProfile` Go type and its `SubscriberProfile` sub-shape
- The YAML and JSON-Schema 2020-12 forms of the `tenant_profile` block
- The Contextus scope-loader's decode and validation behaviour
- The concern boundary between Contextus (type shape, schema validation, YAML decode, backwards compat) and BMA (scaffold-type taxonomy semantics, corpus_class enum, routing logic, ACL enforcement)
- Backwards compatibility with v1.3 scope-configs that omit the block

The addendum does NOT specify:

- Routing tables (`corpus_class × scaffold_type → Subconscious cell`) — BMA-side per Spec 9.4 §3.1
- Token-budget tiers per corpus class — BMA-side per Spec 9.4 §3.3
- ACL enforcement semantics — BMA-side per Spec 9.4 §4.2
- First-5-submission attestation flow — BMA-side state per Spec 9.4 §7
- Divergence reporting (`NT_SCAFFOLD_DIVERGENCE`) — tenant-to-CTH write per Spec 9.4 §4.3

Relationship to existing Contextus surfaces: this is **federation-additive** in the same shape as the Spec v1.4 design surface (PR #11, theory-as-conceptual-scope, merged 2026-05-14). No existing v1.3 section is modified; no existing field is renamed; no existing scope-config is broken.

---

## 1. Motivation

### 1.1 The gap

BMA Spec 9.4 §7 Tenant Onboarding requires step 2:

> *"Tenant publishes a subscriber profile (per A22 §3 rule 2) declaring which scaffold types it accepts."*

Spec 9.4 §2.1 derives the inbound NATS subject `bma.research_aid.submit.<tenant_id>` from a `tenant_id` that must exist somewhere in the tenant's local environment as a single source of truth. Neither tenant identity nor a subscriber profile is currently a Contextus surface: `pkg/types/scope.go` and `scope-config.schema.json` know only about physical/conceptual scope nodes and their memberships per Contextus Spec v1.3 §4.6.

### 1.2 Why scope-config

A tenant's scope-config is already the canonical per-tenant configuration artifact in the federation. Each tenant authors exactly one `scope-config.yaml`; the Wyrd-side `wyrd/store.LoadScopeConfig` is already the single point through which BMA reads it at boot. Adding `tenant_profile` as an additive top-level block keeps the one-file-per-tenant ergonomic and reuses the existing JSON-Schema 2020-12 validation surface that Sprint 1 PR #14 hardened.

Federation-additive precedent: Spec v1.4 PR #11 added the `cth-derivation` membership predicate to §4.6.5 (CTH coupling inside a Contextus type) without violating Contextus's spec. The same shape applies here — BMA coupling lives inside a Contextus YAML, but the *semantics* of the BMA coupling are interpreted BMA-side.

### 1.3 §I4 architectural choice — Option A leaning

A gap analysis (`/tmp/contextus-t5-gap-analysis.md` §3) identified two viable placements for the subscriber profile:

- **Option A** — Top-level `tenant_profile` block in `scope-config.yaml`. One file per tenant. Recommended.
- **Option B** — Separate `tenant-config.yaml` file. Two files per tenant. Cleaner concern boundary but multi-file configurations rot.

This addendum is authored under Option A. It is conditionally correct: if §I4 outcome is Option B, the same field shapes (`TenantProfile`, `SubscriberProfile`, `TenantSubgraphRef`) survive intact and the addendum is renamed and rescoped to a separate-file specification. Option A vs Option B coordination is on the sessionbridge `sprint-2-2026-05-20` channel with @bma-implementor, @cth-implementor, and @qbp-architecture; resolution gates §I4 by the standing §2.i 4h SLA window (closes 2026-05-21 ~05:00Z).

### 1.4 Why not Wyrd-side passthrough

Contextus's loader runs in YAML strict-mode (`KnownFields(true)` + JSON-Schema `additionalProperties: false`) per Sprint 1 PR #14. Silent passthrough of unknown top-level keys is impossible. Either Contextus extends the schema (this addendum) or the subscriber profile lives in a separate file (Option B). Hybrid passthrough is architecturally not on the table.

---

## 2. TenantProfile shape

### 2.1 `TenantProfile` Go type

Added to `pkg/types/scope.go` alongside the existing scope-node types. Three fields:

```go
// TenantProfile names a tenant's federation identity and declares its
// subscriber profile against the BMA Research-Aid Protocol (Spec 9.4) and
// the BMA cross-tenant autonomic signal bus (A22 §3 rule 2).
//
// Contextus owns the type shape and YAML/JSON-Schema validation; BMA owns
// the routing semantics. See Contextus-Spec-Addendum-Research-Aid-Tenancy §5.
type TenantProfile struct {
    TenantID           string             `json:"tenant_id"`
    SubscriberProfile  SubscriberProfile  `json:"subscriber_profile"`
    TenantSubgraphRef  TenantSubgraphRef  `json:"tenant_subgraph_ref"`
}
```

Field semantics:

| Field | Purpose | Owner |
|---|---|---|
| `TenantID` | Federation-unique tenant identifier (e.g., `qbp`, `sharp-butler`, `moebius`). Lowercase, hyphen-separated, no leading or trailing hyphens. Derives the NATS submission subject `bma.research_aid.submit.<tenant_id>` per Spec 9.4 §2.1. | Contextus validates the string shape; BMA interprets routing. |
| `SubscriberProfile` | Declares which scaffold types and corpus classes this tenant accepts, and the default `intended_consumers` set for self-submitted literature nodes. | Contextus validates the shape; BMA interprets the gate. |
| `TenantSubgraphRef` | CTH reference into which BMA writes scaffold output for this tenant. | Contextus validates the shape; CTH/Wyrd interprets the reference. |

### 2.2 `SubscriberProfile` Go type

```go
// SubscriberProfile names the BMA Research-Aid output classes this tenant
// accepts. Field semantics are defined by BMA Spec 9.4 §3 and §4; this struct
// is the syntactic carrier only.
type SubscriberProfile struct {
    AcceptedScaffoldTypes    []string `json:"accepted_scaffold_types"`
    AcceptedCorpusClasses    []string `json:"accepted_corpus_classes"`
    IntendedConsumersDefault []string `json:"intended_consumers_default"`
}
```

Field semantics:

| Field | Purpose |
|---|---|
| `AcceptedScaffoldTypes` | The set of `NT_LITERATURE_SCAFFOLD` types this tenant accepts as a consumer (per Spec 9.4 §3.1 routing table). v0.1 enum values listed in §2.3 below are provisional pending @bma-implementor canonical taxonomy confirmation. |
| `AcceptedCorpusClasses` | The set of `corpus_class` values this tenant intends to submit and consume (mirrors Spec 9.4 §3.3.1 table). |
| `IntendedConsumersDefault` | Default value populating `NT_LITERATURE_NODE.intended_consumers` (Spec 9.4 §2.2) when the tenant does not specify per-submission. Element values are other tenant IDs or the literal `self`. |

### 2.3 Canonical enum values

**`AcceptedCorpusClasses`** — mirror of Spec 9.4 §3.3.1 corpus_class table. Contextus validates the value is in this set; BMA interprets routing:

```
PHYSICS_PREPRINT
JOURNAL_ARTICLE
DATASET_DESCRIPTOR
CODE_REPO
CONTRACT_PRECEDENT
REGULATORY_TEXT
BEEKEEPER_NOTE
OTHER
```

**`AcceptedScaffoldTypes`** — provisional v0.1 list pending @bma-implementor canonical taxonomy (Spec 9.4 §3.1 mentions a non-exhaustive set):

```
PRECEDENT_GRAPH
EVIDENCE_LATTICE
ALGEBRAIC_STRUCTURE_SCAFFOLD
SOURCE_LOCATION_HYPOTHESIS
```

The provisional list is explicitly NOT closed. Spec 9.4 §3.1 says *"Routing is determined by `corpus_class × scaffold_type` lookup table maintained by `bma-implementor` (initial table seeded by Beekeeper)."* The canonical scaffold-type enum is BMA-implementor-owned and ratified separately; the JSON-Schema `enum` for `accepted_scaffold_types` SHALL be updated to match once @bma-implementor publishes the canonical v0.1 list. Until then, Contextus accepts the four values above and rejects others with `ErrScopeConfigInvalid`.

### 2.4 `TenantSubgraphRef` shape

```go
// TenantSubgraphRef references the CTH subgraph into which BMA writes
// scaffold output for this tenant. Convention-derivable from TenantID;
// explicit URI permitted for tenants whose subgraph anchor diverges from
// the convention.
type TenantSubgraphRef struct {
    URI string `json:"uri"` // canonical form: cth://tenant/<tenant_id>/subgraph
}
```

Per gap-analysis §3 Q3, this addendum adopts the **convention** form: `cth://tenant/<tenant_id>/subgraph`. The convention obviates an explicit Contextus↔CTH coupling for the typical case; only divergent-anchor tenants populate the field with a non-convention URI.

Loader behaviour: if `tenant_subgraph_ref.uri` is omitted, the loader synthesizes `cth://tenant/<tenant_id>/subgraph` from `TenantID`. If present, the loader validates it parses as a `cth://` URI and stores the explicit value.

---

## 3. YAML scope-config integration

### 3.1 Top-level `tenant_profile` block — optional

The block is OPTIONAL. v1.3 scope-configs without a `tenant_profile` continue to load with `LoadResult.TenantProfile == nil`. New tenants onboarding to the Research-Aid Protocol populate the block.

YAML shape:

```yaml
tenant_profile:
  tenant_id: qbp
  subscriber_profile:
    accepted_scaffold_types:
      - PRECEDENT_GRAPH
      - EVIDENCE_LATTICE
      - ALGEBRAIC_STRUCTURE_SCAFFOLD
      - SOURCE_LOCATION_HYPOTHESIS
    accepted_corpus_classes:
      - PHYSICS_PREPRINT
      - JOURNAL_ARTICLE
      - DATASET_DESCRIPTOR
      - CODE_REPO
    intended_consumers_default:
      - self
  tenant_subgraph_ref:
    uri: cth://tenant/qbp/subgraph    # optional; loader synthesizes if omitted
```

### 3.2 JSON Schema 2020-12 shape

The schema MUST be added to `schema/scope-config.schema.json` as an optional top-level property. The shape mirrors the Go types:

```json
{
  "tenant_profile": {
    "type": "object",
    "additionalProperties": false,
    "required": ["tenant_id", "subscriber_profile"],
    "properties": {
      "tenant_id": {
        "type": "string",
        "pattern": "^[a-z][a-z0-9]*(-[a-z0-9]+)*$"
      },
      "subscriber_profile": {
        "type": "object",
        "additionalProperties": false,
        "required": [
          "accepted_scaffold_types",
          "accepted_corpus_classes",
          "intended_consumers_default"
        ],
        "properties": {
          "accepted_scaffold_types": {
            "type": "array",
            "items": {
              "type": "string",
              "enum": [
                "PRECEDENT_GRAPH",
                "EVIDENCE_LATTICE",
                "ALGEBRAIC_STRUCTURE_SCAFFOLD",
                "SOURCE_LOCATION_HYPOTHESIS"
              ]
            },
            "uniqueItems": true
          },
          "accepted_corpus_classes": {
            "type": "array",
            "items": {
              "type": "string",
              "enum": [
                "PHYSICS_PREPRINT",
                "JOURNAL_ARTICLE",
                "DATASET_DESCRIPTOR",
                "CODE_REPO",
                "CONTRACT_PRECEDENT",
                "REGULATORY_TEXT",
                "BEEKEEPER_NOTE",
                "OTHER"
              ]
            },
            "uniqueItems": true
          },
          "intended_consumers_default": {
            "type": "array",
            "items": {
              "type": "string",
              "pattern": "^(self|[a-z][a-z0-9]*(-[a-z0-9]+)*)$"
            },
            "uniqueItems": true
          }
        }
      },
      "tenant_subgraph_ref": {
        "type": "object",
        "additionalProperties": false,
        "properties": {
          "uri": {
            "type": "string",
            "pattern": "^cth://"
          }
        }
      }
    }
  }
}
```

The `tenant_profile` property is added at the top level alongside `physical_scopes`, `conceptual_scopes`, `scope_memberships`, and is NOT added to the top-level `required` list (it is optional). `additionalProperties: false` is preserved on the parent schema; this addendum's responsibility is to add `tenant_profile` to the known top-level properties.

### 3.3 Loader behavior

Phase 1.4 (separate subagent / separate PR) extends `internal/contextus/tenancy/loader.go`:

1. `rawConfig` gains a `TenantProfile *rawTenantProfile` field (yaml tag `tenant_profile`).
2. `LoadResult` gains a `TenantProfile *types.TenantProfile` field. Nil when the block is absent; populated otherwise.
3. `buildResult` constructs `LoadResult.TenantProfile` when `rawConfig.TenantProfile` is non-nil:
   - Validates `TenantID` matches the schema pattern (defence-in-depth; JSON-Schema validation runs first).
   - If `tenant_subgraph_ref.uri` is empty, synthesizes `cth://tenant/<tenant_id>/subgraph`.
   - Copies `SubscriberProfile` fields directly.
4. Sentinel errors: malformed `tenant_profile` yields `ErrScopeConfigInvalid` (existing taxonomy). Parse failures yield `ErrScopeConfigParse`.

### 3.4 Backwards compatibility

The loader extension MUST preserve the v1.3 loading contract:

- A scope-config containing only `physical_scopes`, `conceptual_scopes`, `scope_memberships` (the v1.3 baseline) loads identically to the v1.3 loader, with `LoadResult.TenantProfile == nil`.
- `KnownFields(true)` strict-mode is preserved.
- `additionalProperties: false` on the top-level schema is preserved (the schema is extended to recognize `tenant_profile`; it does not relax strictness).
- All existing Sprint 1 PR #14 test fixtures continue to load without modification.

A regression test in `loader_test.go` MUST exercise the bare-v1.3 path explicitly.

---

## 4. BMA-side consumer contract

### 4.1 Wyrd passthrough

The Wyrd-side `wyrd/store.LoadScopeConfig` is the existing single point through which BMA reads tenant configuration. The federation-imported `pkg/types/scope.go` and `internal/contextus/tenancy.LoadResult` are consumed via Go module import per the Sprint 1 federation contract (`contextus/internal/contextus/types/` direct-import, codified in Wyrd PR #40).

This addendum requires **no additional Wyrd plumbing** beyond passthrough: BMA reads `LoadResult.TenantProfile` directly, the same way it already reads `LoadResult.PhysicalScopes`.

### 4.2 NATS subject derivation

Spec 9.4 §2.1 specifies `bma.research_aid.submit.<tenant_id>`. With `LoadResult.TenantProfile.TenantID` populated, BMA derives the subject at boot:

```
submitSubject := fmt.Sprintf("bma.research_aid.submit.%s", profile.TenantID)
ackSubject    := fmt.Sprintf("bma.research_aid.ack.%s",    profile.TenantID)
scaffoldSubject := fmt.Sprintf("bma.research_aid.scaffold.%s", profile.TenantID)
```

No additional configuration is required. The `TenantID` is the single source of truth.

### 4.3 Subscriber profile publication

Spec 9.4 §7 step 2 requires the tenant to publish a subscriber profile per A22 §3 rule 2 subscriber gate. With `LoadResult.TenantProfile.SubscriberProfile` populated:

1. BMA reads the profile at boot.
2. BMA caches it in its A22 federation autonomic-bus subscriber-gate registry.
3. The A22 §3 rule 2 gate consults the cache for inbound cross-tenant `NT_AUTONOMIC_SIGNAL` and (where applicable) Research-Aid scaffold delivery decisions.

No NATS publication of the profile is required by Contextus; the cache is a BMA-side artifact. If A22 §3 rule 2 evolves to require a federation-published profile (e.g., a `bma.federation.subscriber_profile.<tenant_id>` retained subject), that publication is BMA-implementor-owned, sourced from `LoadResult.TenantProfile.SubscriberProfile`.

### 4.4 What remains BMA-implementor-owned

Out of scope for Contextus per Spec 9.4:

| Concern | Spec 9.4 reference | Owner |
|---|---|---|
| `corpus_class × scaffold_type` routing table | §3.1 | `repo-bma-systema:config/research-aid-routing.yaml`-shape, bma-implementor |
| Token-budget tiers per corpus_class | §3.3.1 | `repo-bma-systema:config/research-aid-token-tiers.yaml`, Compute-Manifest-versioned |
| Tier→TCU mapping | §3.3.1 | bma-implementor |
| ACL enforcement (TENANT_PRIVATE / FEDERATION_READABLE / BEEKEEPER_ONLY) | §4.2 | per-scaffold-node, BMA-side |
| First-5-submission attestation state | §7 step 4 | BMA-side process state |
| `intended_consumers` whitelist maintenance | §7 step 3 | bma-implementor |
| `NT_SCAFFOLD_DIVERGENCE` writes | §4.3 | tenant-to-CTH write, not scope-config |
| Promotion-PR `bma_research_aid:` YAML block | §5.1 | in PR body, parsed by federation CI; not scope-config |

---

## 5. Concern separation

### 5.1 Contextus owns

- The Go type shape (`TenantProfile`, `SubscriberProfile`, `TenantSubgraphRef`)
- JSON-Schema 2020-12 syntactic validation of `tenant_profile` blocks
- YAML→struct decode
- Enum-value validation (mirrors Spec 9.4 canonical values; updates as @bma-implementor publishes canonical taxonomy)
- Backwards compatibility with v1.3 scope-configs
- Sentinel error taxonomy on malformed input

### 5.2 BMA owns

- Scaffold-type taxonomy semantics (which scaffold_type goes to which Subconscious cell; what each scaffold_type produces)
- `corpus_class` enum semantics (what each corpus_class means operationally; the routing table)
- Routing logic (Spec 9.4 §3.1 lookup table)
- ACL enforcement semantics (Spec 9.4 §4.2)
- Subscriber-gate state caching and consultation (A22 §3 rule 2)
- Token-budget tiers and SLOs (Spec 9.4 §3.3)
- Quota enforcement (Spec 9.4 §2.4)

### 5.3 The boundary

Contextus rejects malformed `tenant_profile` **syntactically**: unknown corpus_class, unknown scaffold_type, malformed tenant_id pattern, invalid `cth://` URI. Contextus NEVER rejects on routing semantics: it does not interpret which scaffold types are routable, which corpus classes are budget-feasible, or whether a given tenant is whitelisted as a consumer. A tenant that publishes `accepted_scaffold_types: [PRECEDENT_GRAPH]` is syntactically valid; whether BMA produces precedent graphs for this tenant is BMA's decision.

This is the same shape as Spec v1.3 §4.6.5: Contextus owns the membership-edge structure, adapters own the predicate semantics.

---

## 6. Worked example — QBP tenant onboarding

A minimal `scope-config.yaml` for the QBP tenant, structured to support the worked example in Spec 9.4 §6.1 (QBP ALMA cube source-finding scaffold):

```yaml
# scope-config.yaml — QBP tenant
# v1.3 baseline (physical_scopes, conceptual_scopes, scope_memberships)
# + Research-Aid Tenancy addendum v0.1 (tenant_profile)

tenant_profile:
  tenant_id: qbp
  subscriber_profile:
    accepted_scaffold_types:
      - PRECEDENT_GRAPH
      - EVIDENCE_LATTICE
      - ALGEBRAIC_STRUCTURE_SCAFFOLD
      - SOURCE_LOCATION_HYPOTHESIS
    accepted_corpus_classes:
      - PHYSICS_PREPRINT
      - JOURNAL_ARTICLE
      - DATASET_DESCRIPTOR
      - CODE_REPO
    intended_consumers_default:
      - self
  # tenant_subgraph_ref omitted; loader synthesizes cth://tenant/qbp/subgraph

physical_scopes: []

conceptual_scopes:
  - id: scope-theory-qbp-exp-11
    description: "QBP-EXP-11 GW-GRB temporal correlation theory"
    type_nodes: ["NT_SCOPE_CONCEPTUAL"]
    ontology_uri: cth://anchor/qbp-exp-11
    tags: [qbp, gw-em, narrative-anomaly-candidate]

scope_memberships: []
```

Loading this config produces:

- `LoadResult.TenantProfile.TenantID == "qbp"`
- `LoadResult.TenantProfile.SubscriberProfile.AcceptedScaffoldTypes == [4 elements]`
- `LoadResult.TenantProfile.TenantSubgraphRef.URI == "cth://tenant/qbp/subgraph"` (synthesized)
- `LoadResult.ConceptualScopes[0]` as before

At BMA boot:

- BMA derives `bma.research_aid.submit.qbp` per Spec 9.4 §2.1
- BMA caches the subscriber profile for A22 §3 rule 2 gate consultation
- Subsequent QBP submissions of `NT_LITERATURE_NODE` over the ALMA preprint per Spec 9.4 §6.1 flow through this configured path

---

## 7. Open questions and §I4 reader-list

### 7.1 Open at v0.1

The four sub-questions identified in the gap analysis (`/tmp/contextus-t5-gap-analysis.md` §5):

| Q | Question | v0.1 lean | Status |
|---|---|---|---|
| 1 | Option A (scope-config) vs Option B (separate tenant-config) | Option A | OPEN — pending §I4 outcome on sessionbridge `sprint-2-2026-05-20` |
| 2 | Canonical scaffold-type taxonomy v0.1 list | 4 values per §2.3 | OPEN — pending @bma-implementor canonical taxonomy publication |
| 3 | `tenant_subgraph_ref` shape: explicit AnchorRef, scope_id reference, or `cth://tenant/<id>/subgraph` convention | convention | LEANED — closeable on §I4 ack |
| 4 | First-5-submission attestation: scope-config field or BMA-side state | BMA-side | LEANED — closeable on §I4 ack |

### 7.2 §I4 reader-list

- **@bma-implementor** — BMA Spec 9.4 owner; canonical scaffold-type taxonomy authority; routing-table owner
- **@cth-implementor** — CTH AnchorRef coupling on `tenant_subgraph_ref`; convention vs explicit-reference decision
- **@qbp-architecture** — federation coherence; first-tenant onboarding witness (QBP is the worked example)

A non-blocking ack from each within the §2.i 4h SLA window (closing 2026-05-21 ~05:00Z) closes Q1, Q3, Q4. Q2 remains open until @bma-implementor publishes the canonical scaffold-type taxonomy.

---

## 8. Sequencing

| Phase | This addendum's deliverable |
|---|---|
| **Crawl** (now) | Spec only + Go types + JSON-Schema extension + loader extension. Wyrd's existing `LoadScopeConfig` consumes via federation Go-module import. No live BMA Research-Aid traffic; the surface is loadable and validated only. |
| **Toddle** | First tenant config authored (QBP per §6) and loaded by an exercised Wyrd→BMA path. Spec 9.4 §9 (Toddle): manual scaffold authoring by Opus + Gemini against QBP's existing 69-theorem Lean corpus to validate scaffold shape against real artifacts. |
| **Walk** | First-tenant autonomous Research-Aid traffic (QBP Test C lit review per Spec 9.4 §6.3 — the natural first use case: zero hardware cost, well-scoped corpus, fits in QW8 budget). Subscriber-gate cache populated from `LoadResult.TenantProfile.SubscriberProfile`. |
| **Run** | Cross-tenant federation onboarding (Sharp Butler per Spec 9.4 §6.4; Möbius Fusion once operational). The Option A one-file-per-tenant ergonomic scales with no spec change. |

This addendum is Crawl-phase work only. Spec 9.4's Walk and Run gating (Subconscious cells live; first autonomous scaffolds) is BMA-implementor sequencing and is named here for traceability only.

---

## 9. References

| Reference | Path |
|---|---|
| Contextus Specification v1.3 (extended by this addendum) | `~/Documents/Contextus/Contextus-Spec-v1.3.md` |
| Contextus Theory v1.5 (§3.6.2 AnomalyStructural; §3.6.6 Synthesis as persistence boundary) | `~/Documents/Contextus/Contextus-Theory-v1.5.md` |
| Contextus Spec v1.4 design surface (federation-additive precedent) | `~/Documents/Contextus/doc/spec-v1.4-theory-as-conceptual-scope.md` |
| BMA Spec Addendum 9.4 — Research-Aid Protocol (consumer spec) | `~/Documents/inter/spec/BMA-Spec-Addendum-9_4-Research-Aid-Protocol.md` |
| BMA Theory Addendum 22.0 — Cross-Tenant Autonomic Translation Layer (§3 rule 2 subscriber gate) | `~/Documents/inter/theory/BMA-Theory-Addendum-22_0-Cross-Tenant-Autonomic-Translation-Layer.md` |
| BMA Theory Addendum 23.0 — Research-Aid Frame (algebraic frame for 9.4) | `~/Documents/inter/theory/BMA-Theory-Addendum-23_0-Research-Aid-Frame.md` |
| Contextus tenancy loader (Phase 1.4 extension target) | `~/Documents/Contextus/internal/contextus/tenancy/loader.go` |
| Contextus scope-config JSON Schema (Phase 1.4 extension target) | `~/Documents/Contextus/schema/scope-config.schema.json` |
| Contextus Go types catalog (Phase 1.3 extension target) | `~/Documents/Contextus/pkg/types/scope.go` |
| T5 gap analysis (decision baseline for this addendum) | `/tmp/contextus-t5-gap-analysis.md` |
| Sprint 2 plan-of-record | `/home/prime/.claude/plans/reactive-snacking-lampson.md` (Phase 1.2) |

---

*End of Specification Addendum — Research-Aid Tenancy v0.1. Implementation language: Go. §I4 reviewers: @bma-implementor, @cth-implementor, @qbp-architecture. Conditionally-correct under Option A pending sessionbridge `sprint-2-2026-05-20` resolution.*
