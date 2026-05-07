package types

import "time"

// EvidencePointer is the load-bearing primitive for the where-not-what
// principle (Spec v1.3 §5.4): a signal carries pointers to evidence outside
// its flagged subgraph, not the evidence itself.
//
// Field population is tier-conditional and enforced by the retention layer
// at write time. The struct shape is uniform; tier policy zeros disallowed
// fields on demote and populates them lazily on promote.
//
// Tier rules (§5.4.3):
//
//	Skeleton/Distant: Locator + LocatorKind only (~110 bytes).
//	Peripheral:       + Hash, SizeBytes, LoadedAt, AccessHint.
//	Near:             same as Peripheral.
//	Core:             + Note.
//
// At Distant the AccessHint and Note fields are dropped to preserve the
// §9.1 byte budget (§5.4.5).
type EvidencePointer struct {
	Locator     string    `json:"locator"`                 // URL, file path, archive coordinate, NATS subject, DOI, telescope_obs_id, etc.
	LocatorKind string    `json:"locator_kind"`            // "https" | "file" | "archive" | "nats" | "doi" | "telescope_obs_id" | "summary" | ...
	Hash        []byte    `json:"hash,omitempty"`          // SHA-256 of dereferenced content (Peripheral+)
	SizeBytes   int64     `json:"size_bytes,omitempty"`    // size of dereferenced content (Peripheral+)
	LoadedAt    time.Time `json:"loaded_at,omitempty"`     // first-dereference timestamp (Peripheral+)
	AccessHint  string    `json:"access_hint,omitempty"`   // "cold" | "warm" | "hot" — adapter hint (Peripheral, Near; dropped at Distant)
	Note        string    `json:"note,omitempty"`          // human-readable provenance (Core only)
}

// Tier names used by the retention layer. String constants rather than an
// enum to keep the wire format readable across the §9.1 retention pipeline.
const (
	TierCore       = "core"
	TierNear       = "near"
	TierPeripheral = "peripheral"
	TierDistant    = "distant"
	TierSkeleton   = "skeleton"
)

// EvidenceTierCap returns the per-tier cap on the number of EvidencePointers
// that may live on a single signal. Eviction beyond the cap is by
// summary-pointer eviction (Spec v1.3 §5.4.4); the retention layer is the
// authoritative implementation of the eviction algorithm. This is the
// numeric contract only.
//
// Returns 0 for unknown tiers (caller should treat 0 as "no cap configured;
// reject writes" rather than "unbounded").
func EvidenceTierCap(tier string) int {
	switch tier {
	case TierCore:
		return 50
	case TierNear:
		return 20
	case TierPeripheral:
		return 5
	case TierDistant:
		return 5
	case TierSkeleton:
		return 1
	default:
		return 0
	}
}

// LocatorKindSummary identifies an EvidencePointer that aggregates
// previously-evicted pointers. The Locator field points at an external
// aggregate index (S3, file, NATS subject) that can be re-dereferenced if
// the signal is later promoted. The Note field (Core only) records
// "N evicted pointers from Strengthening cycles X–Y".
const LocatorKindSummary = "summary"
