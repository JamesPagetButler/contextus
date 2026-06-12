package measurement

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

// repoRoot returns the absolute path to the repository root, derived from the
// test file's source location so tests run correctly under `go test ./...` from
// any working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// file is .../internal/contextus/measurement/validate_test.go
	// Walk up 4 levels: measurement → contextus → internal → repo root
	dir := file
	for range 4 {
		dir = filepath.Dir(dir)
	}
	return dir
}

// ---------------------------------------------------------------------------
// T4 — TestSchemaFileMatchesEmbedded
// D4 seam structural-equality contract: the embedded schemaJSON constant and
// schema/nt-signal-measurement.schema.json on disk must decode to the same
// normalised JSON structure (description keys stripped).
//
// Spec Addendum §9.1: "Either side drifting breaks a CI gate, not a reader's
// trust." This is that gate.
// ---------------------------------------------------------------------------

// stripDescriptions recursively removes "description" keys from a decoded
// JSON value. Used to permit cosmetic description differences between the
// embedded constant and the on-disk file while still catching structural drift.
func stripDescriptions(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			if k == "description" {
				continue
			}
			out[k] = stripDescriptions(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, val := range t {
			out[i] = stripDescriptions(val)
		}
		return out
	default:
		return v
	}
}

// normaliseJSON decodes JSON bytes and re-marshals with sorted keys (via
// map[string]interface{}) so that two semantically-identical JSON documents
// compare equal even if written with different key ordering.
func normaliseJSON(t *testing.T, label string, b []byte) []byte {
	t.Helper()
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("normaliseJSON: decode %s: %v", label, err)
	}
	stripped := stripDescriptions(v)
	out, err := json.Marshal(stripped)
	if err != nil {
		t.Fatalf("normaliseJSON: re-marshal %s: %v", label, err)
	}
	return out
}

func TestSchemaFileMatchesEmbedded(t *testing.T) {
	// Read the on-disk schema file.
	schemaPath := filepath.Join(repoRoot(t), "schema", "nt-signal-measurement.schema.json")
	onDisk, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read on-disk schema %s: %v", schemaPath, err)
	}

	embeddedNorm := normaliseJSON(t, "embedded schemaJSON constant", []byte(schemaJSON))
	onDiskNorm := normaliseJSON(t, "on-disk schema file", onDisk)

	if !bytes.Equal(embeddedNorm, onDiskNorm) {
		t.Errorf("structural drift between embedded schemaJSON constant and\n  %s\n\nUpdate one of them.\n\nEmbedded (normalised):\n%s\n\nOn-disk (normalised):\n%s",
			schemaPath, string(embeddedNorm), string(onDiskNorm))
	}
}

// ---------------------------------------------------------------------------
// T4 — Validate round-trip
// ---------------------------------------------------------------------------

func TestValidate_HappyPath(t *testing.T) {
	// The motivating example from Spec Addendum §2.1.
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)
	if err := Validate(context.Background(), payload); err != nil {
		t.Fatalf("Validate: unexpected error: %v", err)
	}
}

func TestValidate_MissingField(t *testing.T) {
	// Omitting "value" must fail.
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)
	err := Validate(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Validate: want ErrMeasurementInvalid, got: %v", err)
	}
}

func TestValidate_AdditionalProperty(t *testing.T) {
	// additionalProperties: false — unknown field must fail.
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 1,
		"extra_field": "forbidden"
	}`)
	err := Validate(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Validate: want ErrMeasurementInvalid on extra field, got: %v", err)
	}
}

func TestValidate_NegativeUncertainty(t *testing.T) {
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": -0.1, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)
	err := Validate(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Validate: want ErrMeasurementInvalid on negative uncertainty.plus, got: %v", err)
	}
}

func TestValidate_ConfidenceOutOfRange(t *testing.T) {
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 1.5,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)
	err := Validate(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Validate: want ErrMeasurementInvalid on confidence > 1, got: %v", err)
	}
}

func TestValidate_ZeroRegistryVersion(t *testing.T) {
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 0
	}`)
	err := Validate(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Validate: want ErrMeasurementInvalid on registry_version 0, got: %v", err)
	}
}

func TestValidate_ExtendedKind(t *testing.T) {
	// extended_quantity_kind is valid in place of quantity_kind, but the schema
	// requires quantity_kind — so a payload with ONLY extended_quantity_kind and
	// no quantity_kind should fail schema validation (quantity_kind is required).
	payload := []byte(`{
		"extended_quantity_kind": "tenant:qbp:custom-kind",
		"value": 0.42,
		"uncertainty": { "plus": 0.05, "minus": 0.05 },
		"confidence": 0.8,
		"unit_system": "si",
		"registry_version": 1
	}`)
	// quantity_kind is in "required" — this must fail schema validation.
	err := Validate(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Validate: want ErrMeasurementInvalid (missing required quantity_kind), got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// T4 — Decode round-trip
// ---------------------------------------------------------------------------

func TestDecode_HappyPath(t *testing.T) {
	payload := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)
	m, err := Decode(context.Background(), payload)
	if err != nil {
		t.Fatalf("Decode: unexpected error: %v", err)
	}
	if m.QuantityKind != "SoundSpeedPeak" {
		t.Errorf("QuantityKind: want %q, got %q", "SoundSpeedPeak", m.QuantityKind)
	}
	if m.Value != 0.68 {
		t.Errorf("Value: want 0.68, got %v", m.Value)
	}
	if m.Uncertainty.Plus != 0.14 {
		t.Errorf("Uncertainty.Plus: want 0.14, got %v", m.Uncertainty.Plus)
	}
	if m.Uncertainty.Minus != 0.13 {
		t.Errorf("Uncertainty.Minus: want 0.13, got %v", m.Uncertainty.Minus)
	}
	if m.Confidence != 0.7 {
		t.Errorf("Confidence: want 0.7, got %v", m.Confidence)
	}
	if m.UnitSystem != "dimensionless" {
		t.Errorf("UnitSystem: want %q, got %q", "dimensionless", m.UnitSystem)
	}
	if m.RegistryVersion != 1 {
		t.Errorf("RegistryVersion: want 1, got %v", m.RegistryVersion)
	}
}

func TestDecode_Invalid(t *testing.T) {
	payload := []byte(`{"quantity_kind": "SoundSpeedPeak"}`) // missing required fields
	_, err := Decode(context.Background(), payload)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("Decode: want ErrMeasurementInvalid, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// T4 — Registry round-trip
// ---------------------------------------------------------------------------

func registryFixturePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "registry", name)
}

func TestLoadRegistryJSON_HappyPath(t *testing.T) {
	f, err := os.Open(registryFixturePath(t, "quantity-kinds.json"))
	if err != nil {
		t.Fatalf("open registry JSON: %v", err)
	}
	defer func() { _ = f.Close() }()

	reg, err := LoadRegistryJSON(context.Background(), f)
	if err != nil {
		t.Fatalf("LoadRegistryJSON: unexpected error: %v", err)
	}

	if reg.Version != 1 {
		t.Errorf("Version: want 1, got %d", reg.Version)
	}

	// All four seed kinds must be present.
	seedKinds := []string{"SoundSpeedPeak", "GravitationalWaveStrain", "CPUTemperature", "MemoryPressure"}
	kinds := reg.Kinds()
	kindSet := make(map[string]bool, len(kinds))
	for _, k := range kinds {
		kindSet[k.Kind] = true
	}
	for _, name := range seedKinds {
		if !kindSet[name] {
			t.Errorf("seed kind %q missing from registry", name)
		}
	}
}

func TestLoadRegistryYAML_HappyPath(t *testing.T) {
	f, err := os.Open(registryFixturePath(t, "quantity-kinds.yaml"))
	if err != nil {
		t.Fatalf("open registry YAML: %v", err)
	}
	defer func() { _ = f.Close() }()

	reg, err := LoadRegistryYAML(context.Background(), f)
	if err != nil {
		t.Fatalf("LoadRegistryYAML: unexpected error: %v", err)
	}
	if reg.Version != 1 {
		t.Errorf("Version: want 1, got %d", reg.Version)
	}
}

func TestRegistry_ValidateKind_HappyPath(t *testing.T) {
	reg := loadTestRegistry(t)

	m := types.Measurement{
		QuantityKind:    "SoundSpeedPeak",
		Value:           0.68,
		Uncertainty:     types.Uncertainty{Plus: 0.14, Minus: 0.13},
		Confidence:      0.7,
		UnitSystem:      "dimensionless",
		RegistryVersion: 1,
	}
	if err := reg.ValidateKind(context.Background(), m); err != nil {
		t.Fatalf("ValidateKind: unexpected error: %v", err)
	}
}

func TestRegistry_ValidateKind_UnknownKind(t *testing.T) {
	reg := loadTestRegistry(t)

	m := types.Measurement{
		QuantityKind:    "NoSuchKind",
		Value:           1.0,
		Uncertainty:     types.Uncertainty{Plus: 0.1, Minus: 0.1},
		Confidence:      0.9,
		UnitSystem:      "si",
		RegistryVersion: 1,
	}
	err := reg.ValidateKind(context.Background(), m)
	if !errors.Is(err, ErrQuantityKindUnknown) {
		t.Fatalf("ValidateKind: want ErrQuantityKindUnknown, got: %v", err)
	}
}

func TestRegistry_ValidateKind_ForwardVersion(t *testing.T) {
	// Signal claims registry_version > registry's version. Bridge must defer-not-reject.
	reg := loadTestRegistry(t)

	m := types.Measurement{
		QuantityKind:    "SoundSpeedPeak",
		Value:           0.68,
		Uncertainty:     types.Uncertainty{Plus: 0.14, Minus: 0.13},
		Confidence:      0.7,
		UnitSystem:      "dimensionless",
		RegistryVersion: 999, // far ahead of registry v1
	}
	err := reg.ValidateKind(context.Background(), m)
	if !errors.Is(err, ErrRegistryVersionMismatch) {
		t.Fatalf("ValidateKind: want ErrRegistryVersionMismatch, got: %v", err)
	}
}

func TestRegistry_ValidateKind_ExtendedKindBypassesRegistry(t *testing.T) {
	reg := loadTestRegistry(t)

	m := types.Measurement{
		ExtendedQuantityKind: "tenant:qbp:custom-kind",
		Value:                0.42,
		Uncertainty:          types.Uncertainty{Plus: 0.05, Minus: 0.05},
		Confidence:           0.8,
		UnitSystem:           "si",
		RegistryVersion:      1,
	}
	if err := reg.ValidateKind(context.Background(), m); err != nil {
		t.Fatalf("ValidateKind: extended kind should bypass registry check, got: %v", err)
	}
}

func TestRegistry_ValidateKind_BothKindFieldsError(t *testing.T) {
	reg := loadTestRegistry(t)

	m := types.Measurement{
		QuantityKind:         "SoundSpeedPeak",
		ExtendedQuantityKind: "tenant:qbp:custom-kind",
		Value:                0.68,
		Uncertainty:          types.Uncertainty{Plus: 0.14, Minus: 0.13},
		Confidence:           0.7,
		UnitSystem:           "dimensionless",
		RegistryVersion:      1,
	}
	err := reg.ValidateKind(context.Background(), m)
	if !errors.Is(err, ErrMeasurementInvalid) {
		t.Fatalf("ValidateKind: want ErrMeasurementInvalid on both kind fields set, got: %v", err)
	}
}

func loadTestRegistry(t *testing.T) *Registry {
	t.Helper()
	f, err := os.Open(registryFixturePath(t, "quantity-kinds.json"))
	if err != nil {
		t.Fatalf("open registry JSON: %v", err)
	}
	defer func() { _ = f.Close() }()
	reg, err := LoadRegistryJSON(context.Background(), f)
	if err != nil {
		t.Fatalf("LoadRegistryJSON: %v", err)
	}
	return reg
}

// ---------------------------------------------------------------------------
// T4 — WithinUncertainty typed comparison
// The motivating fix from Spec Addendum §1: the "0.68 is not near 2/3" failure
// is structurally impossible after this test gate.
// ---------------------------------------------------------------------------

func TestWithinUncertainty_MotivatingCase(t *testing.T) {
	// seq=317 live-test failure case: 0.68 ± 0.14 must contain 2/3 = 0.6667...
	m := types.Measurement{
		Value:       0.68,
		Uncertainty: types.Uncertainty{Plus: 0.14, Minus: 0.13},
	}
	predicted := 2.0 / 3.0 // 0.6666...

	// [0.68 - 0.13, 0.68 + 0.14] = [0.55, 0.82] — 0.667 is inside.
	if !types.WithinUncertainty(m, predicted) {
		t.Errorf("WithinUncertainty(0.68 ± (0.14/0.13), 2/3): expected true, got false — the seq=317 failure would recur")
	}
}

func TestWithinUncertainty_Table(t *testing.T) {
	cases := []struct {
		name      string
		value     float64
		plus      float64
		minus     float64
		predicted float64
		want      bool
	}{
		{
			name:  "predicted equals value",
			value: 1.0, plus: 0.1, minus: 0.1,
			predicted: 1.0, want: true,
		},
		{
			name:  "predicted at upper bound (inclusive)",
			value: 1.0, plus: 0.1, minus: 0.1,
			predicted: 1.1, want: true,
		},
		{
			name:  "predicted at lower bound (inclusive)",
			value: 1.0, plus: 0.1, minus: 0.1,
			predicted: 0.9, want: true,
		},
		{
			name:  "predicted just above upper bound",
			value: 1.0, plus: 0.1, minus: 0.1,
			predicted: 1.10001, want: false,
		},
		{
			name:  "predicted just below lower bound",
			value: 1.0, plus: 0.1, minus: 0.1,
			predicted: 0.89999, want: false,
		},
		{
			name:  "asymmetric uncertainty — within",
			value: 0.68, plus: 0.14, minus: 0.13,
			predicted: 0.55, want: true,
		},
		{
			name:  "asymmetric uncertainty — outside lower",
			value: 0.68, plus: 0.14, minus: 0.10,
			predicted: 0.55, want: false, // [0.58, 0.82] — 0.55 is outside
		},
		{
			name:  "zero uncertainty — exact match only",
			value: 1.0, plus: 0.0, minus: 0.0,
			predicted: 1.0, want: true,
		},
		{
			name:  "zero uncertainty — no match",
			value: 1.0, plus: 0.0, minus: 0.0,
			predicted: 1.000001, want: false,
		},
		{
			name:  "NaN value returns false",
			value: math.NaN(), plus: 0.1, minus: 0.1,
			predicted: 1.0, want: false,
		},
		{
			name:  "NaN predicted returns false",
			value: 1.0, plus: 0.1, minus: 0.1,
			predicted: math.NaN(), want: false,
		},
		{
			name:  "NaN uncertainty.plus returns false",
			value: 1.0, plus: math.NaN(), minus: 0.1,
			predicted: 1.0, want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := types.Measurement{
				Value:       c.value,
				Uncertainty: types.Uncertainty{Plus: c.plus, Minus: c.minus},
			}
			got := types.WithinUncertainty(m, c.predicted)
			if got != c.want {
				t.Errorf("WithinUncertainty: want %v, got %v (value=%v, plus=%v, minus=%v, predicted=%v)",
					c.want, got, c.value, c.plus, c.minus, c.predicted)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// T5 — Live round-trip: decode → registry validate → within_uncertainty
// This is the D4 seam cell: scoutd emit → validate → decode → compare on a
// real measurement value. The live arXiv consumer (Q2 from @qbp-implementor)
// is pending; this test proves the decode-registry-compare pipeline is wired
// end-to-end with a realistic physics payload.
// ---------------------------------------------------------------------------

func TestLiveRoundTrip_SoundSpeedPeakWithinTwoThirds(t *testing.T) {
	// The motivating measurement payload, encoded as emitted by a scout.
	emittedJSON := []byte(`{
		"quantity_kind": "SoundSpeedPeak",
		"value": 0.68,
		"uncertainty": { "plus": 0.14, "minus": 0.13 },
		"confidence": 0.7,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)

	ctx := context.Background()

	// Step 1: schema validation (bridge ingestion gate).
	if err := Validate(ctx, emittedJSON); err != nil {
		t.Fatalf("Step 1 Validate: %v", err)
	}

	// Step 2: decode into typed struct.
	m, err := Decode(ctx, emittedJSON)
	if err != nil {
		t.Fatalf("Step 2 Decode: %v", err)
	}

	// Step 3: registry validation (quantity_kind membership check).
	reg := loadTestRegistry(t)
	if err := reg.ValidateKind(ctx, m); err != nil {
		t.Fatalf("Step 3 ValidateKind: %v", err)
	}

	// Step 4: typed significance comparison (the primitive that makes the
	// seq=317 "not near 2/3" error structurally impossible).
	predicted := 2.0 / 3.0
	if !types.WithinUncertainty(m, predicted) {
		t.Errorf("Step 4 WithinUncertainty: 0.68 ± (0.14/0.13) does not contain 2/3 — seq=317 failure would recur")
	}

	// Step 5: verify the confidence gate shape (§4.1 WCET note).
	// A realistic class_floor for physics predictions; bridge does not filter
	// on confidence, but the gate shape is tested here for documentation.
	const physicsClassFloor = 0.5 // realistic floor; CTH sets this per anchor class
	if m.Confidence < physicsClassFloor {
		t.Logf("Note: confidence %.2f below class_floor %.2f — CTH would block automated status change", m.Confidence, physicsClassFloor)
	} else {
		// CTH admission gate: confidence >= class_floor AND within_uncertainty
		t.Logf("D4 seam: confidence=%.2f >= floor=%.2f, within_uncertainty=true → Consistent", m.Confidence, physicsClassFloor)
	}
}

// TestLiveRoundTrip_BMAMemoryPressure proves the BMA telemetry kind path.
func TestLiveRoundTrip_BMAMemoryPressure(t *testing.T) {
	emittedJSON := []byte(`{
		"quantity_kind": "MemoryPressure",
		"value": 0.72,
		"uncertainty": { "plus": 0.02, "minus": 0.02 },
		"confidence": 0.95,
		"unit_system": "dimensionless",
		"registry_version": 1
	}`)

	ctx := context.Background()

	if err := Validate(ctx, emittedJSON); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	m, err := Decode(ctx, emittedJSON)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	reg := loadTestRegistry(t)
	if err := reg.ValidateKind(ctx, m); err != nil {
		t.Fatalf("ValidateKind: %v", err)
	}

	// BMA autonomic uses 0.80 as the high-pressure threshold; check 0.72 is
	// NOT within the [0.80, 0.80] predicted value (boundary detection).
	if types.WithinUncertainty(m, 0.80) {
		// 0.72 ± 0.02 = [0.70, 0.74]; 0.80 is outside — correct
		t.Errorf("0.72 ± 0.02 should NOT contain 0.80 (outside uncertainty interval)")
	}
	// But 0.72 is within_uncertainty of itself.
	if !types.WithinUncertainty(m, 0.72) {
		t.Errorf("0.72 ± 0.02 should contain 0.72")
	}
}

// TestJSONYAMLRegistryConsistency verifies YAML and JSON sources decode to
// the same version and kind set. YAML governs; JSON feeds — they must agree.
func TestJSONYAMLRegistryConsistency(t *testing.T) {
	ctx := context.Background()

	jsonFile, err := os.Open(registryFixturePath(t, "quantity-kinds.json"))
	if err != nil {
		t.Fatalf("open JSON registry: %v", err)
	}
	defer func() { _ = jsonFile.Close() }()
	jsonReg, err := LoadRegistryJSON(ctx, jsonFile)
	if err != nil {
		t.Fatalf("LoadRegistryJSON: %v", err)
	}

	yamlFile, err := os.Open(registryFixturePath(t, "quantity-kinds.yaml"))
	if err != nil {
		t.Fatalf("open YAML registry: %v", err)
	}
	defer func() { _ = yamlFile.Close() }()
	yamlReg, err := LoadRegistryYAML(ctx, yamlFile)
	if err != nil {
		t.Fatalf("LoadRegistryYAML: %v", err)
	}

	if jsonReg.Version != yamlReg.Version {
		t.Errorf("registry_version mismatch: JSON=%d YAML=%d", jsonReg.Version, yamlReg.Version)
	}

	jsonKinds := make(map[string]bool)
	for _, k := range jsonReg.Kinds() {
		jsonKinds[k.Kind] = true
	}
	yamlKinds := make(map[string]bool)
	for _, k := range yamlReg.Kinds() {
		yamlKinds[k.Kind] = true
	}

	for name := range yamlKinds {
		if !jsonKinds[name] {
			t.Errorf("kind %q present in YAML but missing from JSON projection — run registry-gen to regenerate", name)
		}
	}
	for name := range jsonKinds {
		if !yamlKinds[name] {
			t.Errorf("kind %q present in JSON but missing from YAML — JSON projection is stale", name)
		}
	}
}

// TestWireValueConstants locks the wire values for LocatorKind-style constants
// that affect JSON serialisation.
func TestMeasurementJSONRoundtrip(t *testing.T) {
	original := types.Measurement{
		QuantityKind:    "SoundSpeedPeak",
		Value:           0.68,
		Uncertainty:     types.Uncertainty{Plus: 0.14, Minus: 0.13},
		Confidence:      0.7,
		UnitSystem:      "dimensionless",
		RegistryVersion: 1,
	}

	b, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	// The JSON must decode to the same struct.
	var decoded types.Measurement
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if decoded.QuantityKind != original.QuantityKind {
		t.Errorf("QuantityKind: want %q, got %q", original.QuantityKind, decoded.QuantityKind)
	}
	if decoded.Value != original.Value {
		t.Errorf("Value: want %v, got %v", original.Value, decoded.Value)
	}
	if decoded.Uncertainty.Plus != original.Uncertainty.Plus {
		t.Errorf("Uncertainty.Plus: want %v, got %v", original.Uncertainty.Plus, decoded.Uncertainty.Plus)
	}
	if decoded.Uncertainty.Minus != original.Uncertainty.Minus {
		t.Errorf("Uncertainty.Minus: want %v, got %v", original.Uncertainty.Minus, decoded.Uncertainty.Minus)
	}

	// Re-validate the round-tripped JSON against the schema.
	if err := Validate(context.Background(), b); err != nil {
		t.Errorf("Validate on re-serialised struct: %v", err)
	}

	// Check the JSON keys are snake_case as expected.
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	for _, key := range []string{"quantity_kind", "value", "uncertainty", "confidence", "unit_system", "registry_version"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q missing from serialised Measurement", key)
		}
	}
	if _, ok := raw["extended_quantity_kind"]; ok {
		// omitempty — should not appear when empty.
		t.Errorf("extended_quantity_kind should be omitted when empty (omitempty)")
	}
}
