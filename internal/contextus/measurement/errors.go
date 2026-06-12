// Package measurement provides validation, decoding, and registry lookup for
// NT_SIGNAL measurement payloads.
//
// Spec cross-reference: Addendum NT-Signal-Measurement §2 (payload contract),
// §4 (confidence semantics), §5 (registry), §6 (version-skew), §9.1 (D4 seam
// structural-equality agreement contract).
package measurement

import "errors"

// Sentinel errors for measurement validation and registry operations.
var (
	// ErrMeasurementInvalid is returned when a JSON measurement payload fails
	// JSON Schema 2020-12 validation (schema/nt-signal-measurement.schema.json).
	ErrMeasurementInvalid = errors.New("measurement: payload failed schema validation")

	// ErrRegistryParse is returned when the registry JSON cannot be decoded.
	ErrRegistryParse = errors.New("measurement: registry parse error")

	// ErrRegistryInvalid is returned when the registry JSON is structurally
	// invalid (missing required fields, version < 1, etc.).
	ErrRegistryInvalid = errors.New("measurement: registry structurally invalid")

	// ErrQuantityKindUnknown is returned when a signal's quantity_kind is not
	// present in the registry version it claims. See §5.3 (monotonicity) —
	// a kind unknown to the named version was never registered, not deprecated.
	ErrQuantityKindUnknown = errors.New("measurement: quantity_kind not found in registry version")

	// ErrRegistryVersionMismatch is returned when the signal's registry_version
	// exceeds the max version held by the registry. The bridge defers-not-rejects
	// such signals per §8.2 (rule 2: defer-not-reject for forward-version signals).
	// Callers that implement the bridge SHOULD queue on this error rather than drop.
	ErrRegistryVersionMismatch = errors.New("measurement: signal registry_version exceeds registry max version")
)
