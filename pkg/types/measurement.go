package types

import "math"

// Measurement is the machine-liftable payload an NT_SIGNAL MAY carry when it
// asserts a numerical physical measurement. When present, it is the
// wire-format half of the NT_SIGNAL ↔ Edda Measurement[PhysicalQuantity]
// co-design: Contextus owns this shape; Edda owns the in-language type; they
// agree by construction (D4 seam, Spec Addendum NT-Signal-Measurement §9).
//
// All six fields (QuantityKind, Value, Uncertainty, Confidence, UnitSystem,
// RegistryVersion) are mandatory when Measurement is present. ExtendedQuantityKind
// is the §5.4 escape hatch for tenant-specific kinds; it is mutually exclusive
// with QuantityKind in the JSON Schema (additionalProperties: false).
//
// The JSON Schema (schema/nt-signal-measurement.schema.json) is the authoritative
// structural contract. TestSchemaFileMatchesEmbedded in
// internal/contextus/measurement/validate_test.go enforces that the schema constant
// embedded in the validator stays in sync with the on-disk schema file.
//
// Spec cross-reference: Addendum NT-Signal-Measurement §2, §3, §5.4.
type Measurement struct {
	// QuantityKind is the registered physical quantity kind, e.g. "SoundSpeedPeak".
	// MUST be a kind present in RegistryVersion of the canonical registry
	// (registry/quantity-kinds.yaml → generated registry/quantity-kinds.json).
	// Mutually exclusive with ExtendedQuantityKind.
	QuantityKind string `json:"quantity_kind,omitempty"`

	// ExtendedQuantityKind carries a tenant-specific kind outside the registry
	// (§5.4 escape hatch). Opaque to cross-tenant significance matching.
	// Mutually exclusive with QuantityKind.
	ExtendedQuantityKind string `json:"extended_quantity_kind,omitempty"`

	// Value is the measured value. Numeric — never a string.
	Value float64 `json:"value"`

	// Uncertainty is the instrument precision, asymmetric. Witnessed from the
	// source artifact (§3 — witnessed-by-source class).
	Uncertainty Uncertainty `json:"uncertainty"`

	// Confidence is the lift's self-assessment of extraction fidelity (§4).
	// Declared class: can only raise the admission bar, never lower it.
	// CTH applies per-anchor-class floors at evaluation; the bridge does not
	// filter on confidence.
	Confidence float64 `json:"confidence"`

	// UnitSystem is the explicit unit system. Never implied.
	// Typical values: "si", "natural", "dimensionless".
	UnitSystem string `json:"unit_system"`

	// RegistryVersion is the quantity-kind registry version the signal was
	// emitted under. Immutable after emission — signals carry a permanent
	// record of the registry state they were minted against (§6).
	RegistryVersion int `json:"registry_version"`
}

// Uncertainty is the asymmetric instrument precision for a Measurement.
// Plus is σ above the value; Minus is σ below. Both are non-negative magnitudes.
// For symmetric uncertainty, set Plus == Minus.
//
// Spec cross-reference: Addendum NT-Signal-Measurement §2.
type Uncertainty struct {
	// Plus is the σ above the measured value. Non-negative.
	Plus float64 `json:"plus"`
	// Minus is the σ below the measured value. Non-negative.
	Minus float64 `json:"minus"`
}

// WithinUncertainty reports whether the predicted value falls within the
// measured value's uncertainty interval (inclusive).
//
// The interval is [value - uncertainty.minus, value + uncertainty.plus].
// This is the arithmetic primitive that makes the seq=317 "not near 2/3"
// error structurally impossible: the check is deterministic, typed, and
// never delegated to a language model.
//
// A predicted value of NaN or an interval containing NaN returns false.
//
// Spec cross-reference: Addendum NT-Signal-Measurement §1 (motivation) +
// §9.1 (D4 seam agreement contract); Spec v1.3 §4.1 (confidence semantics,
// CTH admission gate shape).
func WithinUncertainty(m Measurement, predicted float64) bool {
	if math.IsNaN(m.Value) || math.IsNaN(predicted) {
		return false
	}
	if math.IsNaN(m.Uncertainty.Plus) || math.IsNaN(m.Uncertainty.Minus) {
		return false
	}
	low := m.Value - m.Uncertainty.Minus
	high := m.Value + m.Uncertainty.Plus
	return predicted >= low && predicted <= high
}
