// Package types defines the canonical Go types for Contextus, mirroring
// Contextus-Spec-v1.3.md §11. This package is the source of truth for type
// shape; consumer repos (wyrd/contextus, bma) import from here.
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
