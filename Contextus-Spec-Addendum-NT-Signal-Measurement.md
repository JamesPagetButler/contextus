# PROJECT CONTEXTUS

## Specification Addendum — NT_SIGNAL Measurement Payload

**Machine-liftable measurement schema + quantity-kind registry — the wire-format half of the NT_SIGNAL ↔ Edda `Measurement[PhysicalQuantity]` co-design**

Version 0.1 | June 2026

Helpful Engineering — Contextus project

Author: contextus-impl
Co-designed with: edda-implementor (Bragi) — live-test seq=317–343 + seq=504 design conversation
Co-Authored-By: James Paget Butler (Beekeeper)

**Extends:** Contextus Specification v1.3 (`Contextus-Spec-v1.3.md`) §4.4 (Signal Emission)
**Sibling addenda:** Research-Aid-Tenancy (PR #17, merged); NT_SCOPE_OPERATIONAL (PR #22, merged)
**Cross-repo counterpart:** Edda `Measurement[Q: PhysicalQuantity]` type (Stage 1; type-shape slice E3 ships this sprint)
**Tracking issue:** Contextus issue #31 (T1 of T1–T5)
**Status:** §I4 review surface. Federation-additive only.

---

## Changelog

| Version | Date | Changes |
|---|---|---|
| 0.1 | 2026-06-04 | Initial draft. Defines the NT_SIGNAL `measurement` payload contract (six fields, all machine-liftable), the witnessed/declared field classification, confidence semantics with CTH-side per-anchor-class floors, the versioned quantity-kind registry (Model B), version-skew evaluation semantics, the deprecation + split-governance lifecycle, and the three-rule registry-pin coordination model. Out of scope: JSON Schema fragment (T2), Go types (T2), registry seed file (T3), loader decode + `within_uncertainty` (T4), live consumer round-trip (T5). |

---

## 0. Scope

This addendum extends Contextus Spec v1.3 §4.4 with one additive surface: a typed,
machine-liftable schema for the `measurement` payload of `NT_SIGNAL` hyperedges, plus the
versioned quantity-kind registry that makes the payload's `quantity_kind` field stable
enough for a compiler to import.

The addendum specifies:

- The motivating failure and the design principle it forces (§1)
- The six-field `measurement` payload contract (§2)
- The witnessed/declared classification of payload fields (§3)
- Confidence semantics and the CTH-side per-anchor-class admission floor (§4)
- The quantity-kind registry: structure, ownership, versioning (§5)
- Version-skew evaluation semantics (§6)
- The deprecation lifecycle and split governance (§7)
- Registry-pin coordination for multi-party kinds — the three-rule model (§8)
- The Edda-side counterpart and the agreement-test contract (§9)
- Concern separation: scout / bridge / CTH / Edda / Contextus (§10)
- Backwards compatibility (§11)
- Sequencing (§12)
- Open §I4 questions and named reviewers (§13)

The addendum does NOT specify:

- `schema/nt-signal-measurement.schema.json` — T2
- `pkg/types/measurement.go` (`Measurement`, `Uncertainty`) — T2
- `registry/quantity-kinds.yaml` v0.1 seed — T3
- Loader decode, embedded-schema structural-equality test, `within_uncertainty` — T4
- The live consumer round-trip (scoutd emit-path or Edda E6 fallback) — T5
- CTH admission mechanics beyond the confidence-floor read (judge collective, multisig
  `cap(cth_admit)`, proposal/admission timestamps) — Verdandi Authority Theory v0.2
  Addendum A territory (`inter/theory/`)
- Edda language semantics (`PhysicalQuantity` primitive, §3.9 schema-versioned types) —
  Edda Stage 1 design issue

Relationship to existing surfaces: **federation-additive** in the shape of PR #17 and
PR #22. No existing v1.3 section is modified. Signals that do not carry a `measurement`
payload remain fully valid — this addendum types a payload that was previously free-form;
it does not require one.

---

## 1. Motivation — the extraction failure

During the 2026-05-29 scout-pipeline emulation (live-test seq=317), a numerical comparison
was delegated to a language model. The model reported a measured value of 0.68 ± 0.14 as
"not near 2/3" — when the interval comfortably contains 0.667. A categorical prose judgment
was returned where arithmetic was needed, and a correct, escalation-worthy physics finding
(`PRED-peak-sound-speed-Q`: untested → consistent) was very nearly lost at the extraction
step — upstream of every authority check, every significance computation, every admission
gate.

The structural lesson, ratified across both scout design threads:

> **The schema contract precedes the typed computation.** If the measurement crosses the
> bridge as the prose string "0.68 ± 0.14", no typed matching on the evaluation side can
> recover it — the information was already lost at emission. No downstream type discipline
> can recover information lost at the wire format.

This addendum is the wire-format half of the fix. The Edda-side half — the
`Measurement[Q: PhysicalQuantity]` in-language type whose comparison primitive
`within_uncertainty()` makes the "not near 2/3" error structurally impossible — is one
type with two encodings: Contextus owns the wire format, Edda owns the in-language type,
and they agree by construction (§9).

---

## 2. The `measurement` payload contract

An `NT_SIGNAL` hyperedge MAY carry a `measurement` object. When present, it MUST conform to:

| Field | Type | Required | Semantics |
|---|---|---|---|
| `quantity_kind` | string (registry-validated) | yes | The physical quantity kind, e.g. `SoundSpeedPeak`. MUST be a registered kind in the registry version named by `registry_version` (or carried in `extended_quantity_kind` instead — see §5.4) |
| `value` | number | yes | The measured value. Numeric, never a string. |
| `uncertainty` | object `{plus: number, minus: number}` | yes | Instrument precision, asymmetric. `plus`/`minus` are non-negative magnitudes (σ above / σ below). Symmetric uncertainty sets both to the same value. |
| `confidence` | number ∈ [0,1] | yes | Extraction fidelity — the lift's self-assessment of how reliably the value+uncertainty were extracted from the source (§4) |
| `unit_system` | string | yes | Explicit unit system, e.g. `"si"`, `"natural"`, `"dimensionless"`. Never implied. |
| `registry_version` | integer ≥ 1 | yes | The quantity-kind registry version the signal was emitted under (§6) |

All six fields are mandatory when `measurement` is present. There is no prose form: a
signal whose measurement cannot be expressed in this shape carries no `measurement` object
and is not eligible for typed significance evaluation.

The JSON Schema (T2) applies `additionalProperties: false` at the `measurement` level —
unknown fields are rejected at emission, consistent with every Contextus schema surface
since PR #11.

### 2.1 Worked example

The motivating finding, as a conformant payload:

```json
{
  "measurement": {
    "quantity_kind": "SoundSpeedPeak",
    "value": 0.68,
    "uncertainty": { "plus": 0.14, "minus": 0.13 },
    "confidence": 0.7,
    "unit_system": "dimensionless",
    "registry_version": 1
  }
}
```

The evaluation-side comparison is then arithmetic over typed values —
`within_uncertainty(measured, predicted)` — not an interpretation.

---

## 3. Witnessed vs declared fields

Per the witnessed/declared provenance axis (Verdandi Authority Theory v0.2 Addendum A,
§A.5; established live-test seq=331): security-load-bearing quantities are
substrate-witnessed; producer-supplied fields are advisory and asymmetric.

| Field | Class | Consequence |
|---|---|---|
| `value`, `uncertainty` | **witnessed-by-source** — transcribed from the published artifact (the paper's stated measurement) | Verifiable against the source URI; falsifying them is misrepresentation of a citable artifact |
| `quantity_kind`, `unit_system`, `registry_version` | **schema-validated** — checked mechanically at emission against the registry | Invalid values are rejected at the bridge; no trust required |
| `confidence` | **declared** — the lift's self-assessment | Asymmetric: can only RAISE the admission bar, never lower it (§4) |

The asymmetry rule for declared fields is load-bearing: a declared-high confidence must
not reduce collective scrutiny; a declared-low confidence must block automated status
change. Declared provenance can make admission harder, never easier.

---

## 4. Confidence semantics

`confidence` is **extraction fidelity**, distinct from `uncertainty` (instrument
precision):

- A fully automated extraction from a structured data table: `confidence` ≈ 1.0
- A model-lifted value from prose text: `confidence` ≈ 0.7
- A value inferred from a figure or indirect statement: lower still

### 4.1 Where the floor is enforced — CTH at evaluation (Option B)

The confidence floor is applied by **CTH at evaluation time**, per anchor class — NOT by
the bridge at emission. Rationale (ratified live-test seq=336/342):

1. The bridge lacks anchor-specific context: confidence 0.4 may be below the
   physics-prediction admission floor yet acceptable for a low-stakes administrative
   anchor. CTH, as the admission gatekeeper, can make that distinction.
2. Low-confidence signals still enter the hypergraph — visible, queryable, useful as
   deliberation metadata — they simply cannot trigger automated status changes.

The significance gate has the shape:

```
if measurement.confidence >= class_floor
   AND within_uncertainty(measurement, anchor.prediction)
then Consistent else Inconclusive
```

where `class_floor` is **witnessed-config**: set by the anchor-class owner as a governance
act, read by the evaluator at evaluation time. The scout cannot set its own admission
floor. WCET note (per seq=342): the gate's worst-case cost is certified over the operation
shape — `bounded(config_read) + compare + within_uncertainty` — with `class_floor` as a
bounded-read runtime parameter; the compiler never needs its value.

**Crawl residence (Q1 ruled by cth-implementor, live-test seq=511):** CTH v0.3 has no
reified anchor-class record (classes are implicit in ID prefixes), so the Crawl artifact
is a top-level optional `class_floors` map in the canonical CTH inventory — object-valued
entries, witnessed by inventory change-control, bounded-read (the evaluator already holds
the inventory), **fail-closed default**: a class absent from the map permits no automated
status change until its owner sets a floor. Setting the floor is a governance act, the
same declared-asymmetry direction as §3. Additive-optional ⟹ minor semver on the CTH side.

The bridge performs schema validation only (shape, types, registry membership). It does
not filter on confidence.

---

## 5. The quantity-kind registry

### 5.1 Model

**Model B — versioned federation registry** (ratified live-test seq=330/332): a canonical,
versioned list of quantity kinds that Contextus maintains. Registered kinds give
compile-time type identity (Edda Stage 1 imports a pinned registry version: a QBP
`SoundSpeedPeak` and a CTH `SoundSpeedPeak` are the same type because they reference the
same registry entry). Tenant-specific extensions ride an escape hatch (§5.4) opaque to
cross-tenant matching.

This follows the `HardwareClass` constant pattern (NT_SCOPE_OPERATIONAL addendum §4): a
stable v0.1 set, extensible by registered amendment.

### 5.2 Registry entry shape

Each kind carries:

| Field | Semantics |
|---|---|
| `kind` | Canonical name, e.g. `SoundSpeedPeak` |
| `owner_tenant` | The tenant that introduced the kind; owns its split/deprecation decisions (§7) |
| `unit_system` | The kind's canonical unit system |
| `description` | Human-readable definition |
| `introduced_in` | Registry version the kind first appeared in |
| `deprecated_in` | Registry version the kind was deprecated in (absent = active) |
| `maps_to` | Successor declaration on deprecation: single kind (exclusive) or list (ambiguous) — §7 |

### 5.3 Versioning invariant

**Kinds are never deleted, only deprecated.** The registry grows monotonically. Any pinned
registry version remains permanently valid; old signals never become unevaluable. This
invariant is what makes the version-skew semantics (§6) and the Edda static-enum import
(§9) sound.

### 5.4 The extension escape hatch

`extended_quantity_kind` (string, used in place of `quantity_kind`) carries
tenant-specific kinds outside the registry. Extended kinds are opaque to cross-tenant
significance matching — they are evaluable only within the emitting tenancy. A tenant
whose extended kind acquires cross-tenant relevance registers it (additive registry
amendment), at which point it gains a canonical entry and compile-time identity.

### 5.5 v0.1 seed (T3)

The registry ships seeded with the kinds the two live tenants need:

- QBP: `SoundSpeedPeak`, `GravitationalWaveStrain`
- BMA telemetry (aligned with ctx-adapter-system §5 observation kinds): `CPUTemperature`,
  `MemoryPressure`

---

## 6. Version-skew semantics

An NT_SIGNAL is an immutable fact about a point in time: a scout measured a value under
the registry version that was live. The bridge cannot retroactively re-emit under a later
version. Therefore (ratified live-test seq=336/342):

1. Every measurement-bearing signal carries `registry_version` (§2).
2. The bridge performs **version-indexed dispatch**: it holds all historical registry
   versions (small YAML artifacts; cost negligible) and routes each signal to its
   emission-version schema. No computation — a lookup.
3. CTH evaluates the signal under its emission version, where the signal's kinds are
   guaranteed valid by the monotonicity invariant (§5.3).

A deprecated kind remains fully valid for the evaluation of signals emitted under any
version where it was active. Deprecation constrains *new emission* only (§7).

---

## 7. Deprecation lifecycle and split governance

### 7.1 Deprecation

When a kind is deprecated at version N+1, the emission-time schema for current-version
signals rejects it; version-≤N signals continue to evaluate normally. Two deprecation
forms:

- **Plain deprecation** — kind retired, no successor. `maps_to` absent.
- **Split** — kind superseded by finer-grained kinds. `maps_to` declares the successor(s).

### 7.2 Exclusive vs ambiguous splits

| Split type | `maps_to` | Governance | Evaluation of old signals |
|---|---|---|---|
| **Exclusive** — deprecated kind unambiguously maps to exactly one successor | single kind | Registered amendment by `owner_tenant`; no §I4 cycle | Forward-mapped: evaluated under the mapped successor against the current anchor |
| **Ambiguous** — the kind label alone cannot disambiguate between successors | list (union) | **Requires a Contextus spec amendment with §I4 cycle**; affected tenants (any tenant whose CTH anchors reference the kind) are named reviewers | Evaluated against ALL successor anchors; a match against any is consistent (generous interpretation — accepted false-positive risk is the cost of not losing in-flight signals) |

Ambiguous splits are semantically load-bearing — they change what in-flight signals
*mean* — which is why they carry the only mandatory multi-party governance step in this
addendum (§8 rule 3).

### 7.3 Split-decision routing

`owner_tenant` makes governance routing mechanical: Contextus reads the field at
split-request time and routes the §I4 cycle to the affected tenants. The Edda-side mirror
(design-only this sprint): Edda §3.9 migration declarations carry the same
exclusive/ambiguous flag — `compatible_with_mapping` vs `incompatible_requires_review` —
so compiler-side and registry-side governance stay in lockstep.

---

## 8. Registry-pin coordination — the three-rule model

For multi-party kinds (emitter in one tenancy, evaluator in another, type-consumer in a
third), no global "all parties ≥ version N+1" barrier is needed. Three rules replace it
(authored by edda-implementor, live-test seq=504; adopted verbatim):

1. **Evaluator-leads, emitter-follows.** A kind is usable in cross-tenant matching only
   when the *evaluator* (bridge/CTH) holds the kind's introducing version. Emitters may
   adopt eagerly; evaluation is the gate. Deployment ordering: registry updates reach the
   evaluation side before emitters adopt the new kind. This is a deployment convention,
   not a language mechanism — a consumer's pin may lag the registry head safely forever
   (kinds never deleted); it must simply not lead the evaluator.

2. **Defer-not-reject at the bridge for forward-version signals.** A signal carrying
   `registry_version` greater than the evaluator's max-held version is well-formed under a
   version the registry will provably hold (monotonicity) — so the bridge queues it and
   flags an evaluator-upgrade need rather than rejecting valid data. The signal is never
   wrong; the evaluator is behind.

3. **The ambiguous-split §I4 cycle is the only true multi-party barrier** — applied
   exactly where meaning changes and nowhere else. Cross-tenant kind identity is carried
   by the registry version-chain (same entry, connected by migration declarations), so
   exclusive splits and additions need no coordination beyond rule 1.

---

## 9. The Edda-side counterpart — one type, two encodings

The Edda Stage 1 type:

```edda
Measurement[Q: PhysicalQuantity] = {
    value:       Q,
    uncertainty: { plus, minus },   -- witnessed: from the paper
    confidence:  Real               -- declared: the lift's self-assessment
}
```

`PhysicalQuantity` is a registered-kind enum imported from a **pinned registry version**;
`PhysicalQuantity.Extended(string)` carries the §5.4 escape hatch (runtime-validated,
opaque cross-tenant). The registry IS a versioned schema in Edda's §3.9 sense; a registry
bump is a `compatible_with_widening` migration.

### 9.1 The agreement contract (D4 seam)

"One type, two encodings" is an **executable claim, enforced in CI on both sides of the
bridge**:

- **Contextus side (T4):** embedded-schema structural-equality test — the schema constant
  embedded in the validator and `schema/nt-signal-measurement.schema.json` are
  JSON-decoded, stripped of cosmetic `description` keys, and compared with deep equality
  (the PR #20 `TestSchemaFileMatchesEmbedded` pattern).
- **Edda side (E3):** a structural-equality agreement test between the Edda `Measurement`
  type shape and the T2 schema — the same decode-strip-compare pattern, asserting the
  field-for-field correspondence.

Either side drifting breaks a CI gate, not a reader's trust.

### 9.2 Compiler treatment of deprecated kinds (ratified seq=342/343)

Mirroring the bridge's emission-time reject: new Edda source emitting a deprecated kind is
a **compile-time error** (with migration hint for exclusive splits; escalation message for
ambiguous splits; stark removal notice for plain deprecation). Evaluation-context
references to deprecated kinds are never rejected. A per-site
`#[allow(deprecated_quantity)]` annotation downgrades the error to a warning for
migration-in-progress code; production CI bans the annotation (`no_deprecated_quantity`
lint), making the escape hatch development-only.

---

## 10. Concern separation

| Concern | Owner | This addendum's claim on it |
|---|---|---|
| Wire format (`measurement` schema, registry) | **Contextus** | Specified here; implemented T2–T4 |
| Signal emission (minting typed payloads from sources) | Scout implementations (qbp-systema scoutd first) | Must conform to §2; T5 proves the path live |
| Schema validation at ingestion | Contextus bridge | Shape + registry membership only; no confidence filtering (§4.1); defer-not-reject on forward versions (§8.2) |
| Significance evaluation, confidence floor, admission | **CTH** | Reads `class_floor` witnessed-config; evaluates under emission version (§6); admission mechanics out of scope here |
| In-language type, typed comparison, compile-time kind identity | **Edda** (Stage 1) | Counterpart shape in §9; language semantics out of scope here |
| Registry governance (splits, deprecations) | `owner_tenant` per kind, routed by Contextus | §7 |

The signals-not-conclusions discipline (Spec v1.3 §4.4) is unchanged: a measurement-bearing
signal asserts *what was measured*, never *what it means*. Significance remains an
evaluation-side judgment.

---

## 11. Backwards compatibility

- `measurement` is optional on NT_SIGNAL: every existing signal remains valid.
- No existing Contextus type, schema, or loader surface is modified; T2–T4 are additive
  files plus additive loader decode.
- The registry starts at version 1; no signal predates it with a `measurement` payload
  (the payload type is introduced by this addendum).
- Federation-additive contract preserved: Wyrd's direct import of `pkg/types` gains one
  new type (`Measurement`), no changes to existing types.

---

## 12. Sequencing

| Phase | Artifact | Gate |
|---|---|---|
| T1 (this doc) | Spec addendum | §I4 quorum |
| T2 | JSON Schema + `pkg/types/measurement.go` | T1 merged |
| T3 | `registry/quantity-kinds.yaml` v0.1 | T1 merged |
| T4 | Loader decode + structural-equality test + `within_uncertainty` | T2 + T3 |
| T5 | Live round-trip — scoutd emit → validate → decode → compare on a real arXiv alias-match (PO ruling seq=500: live consumer required; Edda E6 is the standing fallback) | T4 + consumer capacity confirm |
| Walk-α | `ParseCTHRef`-style strict validation library; registry tooling (split-request automation); Edda Stage 1 `PhysicalQuantity` import | out of Crawl scope |

---

## 13. Open questions for §I4 reviewers

1. **Q1 (`@cth-implementor`) — RULED (seq=511, folded into §4.1):** `class_floors` map in
   the canonical CTH inventory (no anchor-class record exists in v0.3), fail-closed
   default, additive-optional minor semver. Non-blocking note carried to T4: ambiguous-split
   `Consistent` verdicts carry a `via_ambiguous_split` marker in the evaluation record so
   generous-interpretation matches stay auditable.
2. **Q2 (`@qbp-implementor`)** — scoutd emit-path shape: does the alias-matcher have
   access to structured value+uncertainty at match time, or does T5 need a small extraction
   step between match and mint? Capacity confirm requested either way (T5 consumer choice).
3. **Q3 (`@edda-implementor`)** — registry artifact format: T3 proposes YAML
   (`registry/quantity-kinds.yaml`) for human-readable governance diffs, with the JSON
   Schema validating signals against the *kinds extracted from it*. Does Edda's Stage 1
   import path prefer consuming the YAML directly or a generated JSON projection? My lean:
   generated JSON projection committed alongside (one source, two encodings — same
   discipline as the type itself).
4. **Q4 (`@qbp-architecture`)** — registry residence: `registry/` in repo-contextus (this
   addendum's assumption) vs `inter/` as federation-canonical. My lean: repo-contextus —
   Contextus owns the wire format and the §I4 routing; `inter/` references it. Counter-case
   is that the registry is federation-shared state like the wisdoms file; ruling requested.

**Named reviewers:** `@edda-implementor` (co-designer), `@cth-implementor` (evaluation
side), `@qbp-implementor` (emitter side), `@qbp-architecture` (federation coherence),
`@beekeeper` (HVR).
