// Package types defines the canonical Go types for Contextus, mirroring
// Contextus-Spec-v1.3.md §11. This package is the source of truth for type
// shape; consumer repos (wyrd, bma) import from here directly.
//
// Federation contract (per Wyrd PR #40 / wyrd-implementor #toddle-design seq=15,
// confirmed by contextus-impl PR #40 review): this package lives at the
// non-internal `pkg/types/` path specifically so that cross-repo consumers
// can import it. Go's internal-package rule blocks any path under
// `internal/`. Do not move this package back under `internal/`.
//
// Changes to public type shape must be additive-only (new fields with
// json:"omitempty"; never rename or retype existing fields) to avoid
// breaking consumer builds. See PR #11 §7 (Spec v1.4 design) for the
// `OntologyURI` + `OntologyURIs` precedent.
//
// Wyrd marshals these types into model.Node.Payload bytes opaquely;
// Wyrd does not enforce shape. The retention layer enforces tier-conditional
// field population (§5.4.3).
//
// Spec cross-references:
//   - §11.1 InsightSignal, EvidencePointer, AnomalyKind, AgentClass, ClaimVersion,
//     HyperedgeRef, LocaleBounds, GrainSpec, EvidenceTierCap.
//   - §11.3 EdgeScoutFlag, CorpusDiversityReport, BridgeIntervention.
//   - §11.4 ScopePhysical, ScopeConceptual, ScopeMembership.
//   - §11 cross-reference: SedenionResult is canonically defined in
//     qbp-compute-unit/doc/wyrd-integration.md and is NOT redefined here.
package types
