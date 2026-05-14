package synthesis

import (
	"context"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

// Persister is the abstraction over the Wyrd write path.
//
// The actual implementation will live in a separate adapter package once the
// Wyrd `contextus/` subpackage corrections (Wyrd issue #6) land — the Wyrd
// subpackage exposes FromSignal / PromoteToCTH which a real Persister wraps.
//
// Until then, callers wire a NoopPersister (logs the mint, no-ops the write)
// or test fakes. The Synthesis subscriber's behaviour is independent of the
// substrate, so swapping Persisters at deployment is the cleanest path.
type Persister interface {
	// MintSignal records a new InsightSignal. Returns the signal's persistent
	// SignalID (assigned by the substrate) or an error.
	//
	// Implementations must be idempotent on the (SourceSubject, SourceID)
	// pair carried in the signal's traversal_path metadata: if Synthesis
	// re-evaluates the same ephemeral source twice (e.g., due to a redelivered
	// NATS message), MintSignal should return the same persistent SignalID
	// rather than creating a duplicate.
	MintSignal(ctx context.Context, sig types.InsightSignal) (types.Addr, error)
}

// NoopPersister implements Persister with no side effects beyond returning a
// deterministic Addr per (SourceSubject, SourceID). Used in tests and during
// the gap before Wyrd's contextus/ subpackage lands.
//
// In v1.3 implementation rollout this is the production Persister too —
// Synthesis emits, the NoopPersister records that the mint happened, and
// real Wyrd writes pick up once the substrate is ready. This is a
// beekeeper-override implementation per ADR-003 §I4 (James 2026-05-06).
type NoopPersister struct {
	// Calls counts MintSignal invocations. Useful for tests.
	Calls int
	// LastSignal preserves the most recently minted signal for inspection.
	LastSignal types.InsightSignal
}

// MintSignal logs the mint by counting calls and capturing the signal.
// It returns a deterministic Addr derived from the signal's content so
// repeated calls with identical input produce identical addresses.
func (n *NoopPersister) MintSignal(ctx context.Context, sig types.InsightSignal) (types.Addr, error) {
	n.Calls++
	n.LastSignal = sig
	return sig.SignalID, nil
}
