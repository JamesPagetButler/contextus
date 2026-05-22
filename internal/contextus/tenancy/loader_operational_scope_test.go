package tenancy

import (
	"context"
	"errors"
	"testing"

	"github.com/JamesPagetButler/contextus/pkg/types"
)

// TestOperationalScope_RoundTrip loads the BMA-prime operational-scope fixture
// (§6.1 worked example: host umbrella + 4 hardware-class child scopes) and
// asserts that LoadResult.OperationalScopes is populated correctly with all
// canonical payload fields per Contextus-Spec-Addendum-NT-Scope-Operational §3.
//
// AC-8: round-trip YAML→ScopeOperational struct.
func TestOperationalScope_RoundTrip(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "bma-prime-operational-scopes-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig(%q): unexpected error: %v", path, err)
	}

	if got, want := len(res.OperationalScopes), 5; got != want {
		t.Fatalf("OperationalScopes count: got %d, want %d", got, want)
	}

	// Index by ScopeID for field-level assertions.
	byID := make(map[string]types.ScopeOperational, len(res.OperationalScopes))
	for _, s := range res.OperationalScopes {
		byID[s.ScopeID] = s
	}

	// Host umbrella (hardware_class = "" per spec §5 hierarchical-aggregation pattern).
	umbrella, ok := byID["contextus:scope:operational:host-bma-prime"]
	if !ok {
		t.Fatalf("missing host umbrella scope")
	}
	if umbrella.HostID != "bma-prime" {
		t.Errorf("umbrella.HostID = %q, want %q", umbrella.HostID, "bma-prime")
	}
	if umbrella.HardwareClass != "" {
		t.Errorf("umbrella.HardwareClass = %q, want empty (host-level umbrella)", umbrella.HardwareClass)
	}
	if umbrella.ParentScopeID != "" {
		t.Errorf("umbrella.ParentScopeID = %q, want empty (top-level)", umbrella.ParentScopeID)
	}
	if umbrella.Name == "" {
		t.Error("umbrella.Name must not be empty")
	}
	if len(umbrella.Tags) == 0 || umbrella.Tags[0] != types.HardwareClassRuntimeBMAInstance {
		t.Errorf("umbrella.Tags = %v, want first element %q", umbrella.Tags, types.HardwareClassRuntimeBMAInstance)
	}

	// Per-subsystem children — each carries the exact v0.1 hardware-class
	// taxonomy value from spec §4.
	childExpectations := map[string]string{
		"contextus:scope:operational:host-bma-prime-cpu":    types.HardwareClassHardwareCPU,
		"contextus:scope:operational:host-bma-prime-disk":   types.HardwareClassHardwareDisk,
		"contextus:scope:operational:host-bma-prime-gpu":    types.HardwareClassHardwareGPU,
		"contextus:scope:operational:host-bma-prime-memory": types.HardwareClassHardwareMemory,
	}
	for id, wantClass := range childExpectations {
		child, ok := byID[id]
		if !ok {
			t.Errorf("missing child operational scope %q", id)
			continue
		}
		if child.HostID != "bma-prime" {
			t.Errorf("%s.HostID = %q, want %q", id, child.HostID, "bma-prime")
		}
		if child.HardwareClass != wantClass {
			t.Errorf("%s.HardwareClass = %q, want %q", id, child.HardwareClass, wantClass)
		}
		if child.ParentScopeID != "contextus:scope:operational:host-bma-prime" {
			t.Errorf("%s.ParentScopeID = %q, want host-bma-prime umbrella", id, child.ParentScopeID)
		}
	}
}

// TestOperationalScope_EnvelopeSplit asserts the Wyrd PR #40 §2.1 envelope
// contract: tier_immune + salience appear in LoadResult.NodeOptions (which
// thread to model.Node.TierImmune / model.Node.Salience), NOT on the
// ScopeOperational struct payload. The struct's field-set is a compile-time
// guarantee (the type has no TierImmune / Salience fields).
//
// AC-8: envelope-split for tier_immune.
func TestOperationalScope_EnvelopeSplit(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "bma-prime-operational-scopes-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig: %v", err)
	}

	const umbrellaID = "contextus:scope:operational:host-bma-prime"

	opts, ok := res.NodeOptions[umbrellaID]
	if !ok {
		t.Fatalf("NodeOptions[%q] not present", umbrellaID)
	}
	if !opts.TierImmune {
		t.Errorf("NodeOptions[%q].TierImmune = false, want true (foundational host identity)", umbrellaID)
	}
	if opts.Salience != 1.0 {
		t.Errorf("NodeOptions[%q].Salience = %v, want 1.0", umbrellaID, opts.Salience)
	}

	// Sibling children: tier_immune true, salience 0.9 from fixture.
	const cpuID = "contextus:scope:operational:host-bma-prime-cpu"
	cpuOpts, ok := res.NodeOptions[cpuID]
	if !ok {
		t.Fatalf("NodeOptions[%q] not present", cpuID)
	}
	if !cpuOpts.TierImmune {
		t.Errorf("NodeOptions[%q].TierImmune = false, want true", cpuID)
	}
	if cpuOpts.Salience != 0.9 {
		t.Errorf("NodeOptions[%q].Salience = %v, want 0.9", cpuID, cpuOpts.Salience)
	}

	// Compile-time guarantee: the ScopeOperational struct does NOT carry
	// TierImmune or Salience. We exercise the canonical fields and rely on
	// the type definition (pkg/types/scope.go) to enforce field-set purity.
	for _, s := range res.OperationalScopes {
		if s.ScopeID == "" || s.Name == "" || s.HostID == "" {
			t.Errorf("ScopeOperational payload incomplete: %+v", s)
		}
	}
}

// TestOperationalScope_BackwardsCompat asserts the federation-additive contract
// per spec §10: a v1.3 (or v1.3+tenant_profile) scope-config without an
// operational_scopes key loads with LoadResult.OperationalScopes == nil.
// This is the regression guard against accidentally requiring operational_scopes.
//
// AC-8: backwards-compatibility check (existing QBP fixture loads unchanged).
func TestOperationalScope_BackwardsCompat(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "qbp-scope-config-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig(%q): unexpected error: %v", path, err)
	}

	if res.OperationalScopes != nil {
		t.Errorf("OperationalScopes = %v, want nil (v1.3 fixture has no operational_scopes key)", res.OperationalScopes)
	}
	// And the existing scope arrays remain populated — no accidental shuffling.
	if got, want := len(res.PhysicalScopes), 2; got != want {
		t.Errorf("PhysicalScopes regression: got %d, want %d", got, want)
	}
	if got, want := len(res.ConceptualScopes), 2; got != want {
		t.Errorf("ConceptualScopes regression: got %d, want %d", got, want)
	}
}

// TestOperationalScope_MalformedRejected asserts that an operational_scope
// entry with a hardware_class value outside the v0.1 enum (§4 closed set:
// runtime.bma-instance, hardware.{cpu,disk,gpu,memory,network}, plus empty
// string for host-umbrella scopes) fails JSON Schema validation and yields
// ErrScopeConfigInvalid. The example fixture uses "hardware.fpga" — a v0.2+
// expansion candidate that is NOT in the v0.1 enum.
//
// AC-8: malformed-hardware-class rejection.
func TestOperationalScope_MalformedRejected(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "scope-config-malformed/operational-scope-unknown-hardware-class.yaml")

	_, err := LoadScopeConfig(ctx, path)
	if err == nil {
		t.Fatalf("LoadScopeConfig(%q): expected ErrScopeConfigInvalid, got nil", path)
	}
	if !errors.Is(err, ErrScopeConfigInvalid) {
		t.Errorf("LoadScopeConfig(%q): errors.Is(err, ErrScopeConfigInvalid) = false; err = %v", path, err)
	}
}

// TestOperationalScope_DuplicateIDRejected asserts that the federation-unique
// scope_id discipline is enforced ACROSS scope sibling types: an
// operational_scope.id that collides with a physical_scope.id (or
// conceptual_scope.id) yields ErrScopeLoadConflict, not a silent overwrite.
//
// AC-8: cross-scope-type duplicate-id rejection.
func TestOperationalScope_DuplicateIDRejected(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "scope-config-malformed/operational-scope-duplicate-id-cross-type.yaml")

	_, err := LoadScopeConfig(ctx, path)
	if err == nil {
		t.Fatalf("LoadScopeConfig(%q): expected ErrScopeLoadConflict, got nil", path)
	}
	if !errors.Is(err, ErrScopeLoadConflict) {
		t.Errorf("LoadScopeConfig(%q): errors.Is(err, ErrScopeLoadConflict) = false; err = %v", path, err)
	}
}

// TestOperationalScope_HierarchicalParent asserts the spec §5 aggregation
// pattern: an operational_scope can declare parent_scope_id pointing at a
// host-level umbrella scope; the loaded struct preserves the ParentScopeID
// field. Wildcards remain unsupported at v0.1 — aggregation across children
// is a Synthesis-layer concern, not a predicate-layer concern.
//
// AC-8: hierarchical parent-pointer round-trip.
func TestOperationalScope_HierarchicalParent(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "bma-prime-operational-scopes-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig: %v", err)
	}

	const (
		umbrellaID = "contextus:scope:operational:host-bma-prime"
		gpuID      = "contextus:scope:operational:host-bma-prime-gpu"
	)

	var umbrella, gpu *types.ScopeOperational
	for i := range res.OperationalScopes {
		s := &res.OperationalScopes[i]
		switch s.ScopeID {
		case umbrellaID:
			umbrella = s
		case gpuID:
			gpu = s
		}
	}
	if umbrella == nil {
		t.Fatalf("umbrella scope %q not loaded", umbrellaID)
	}
	if gpu == nil {
		t.Fatalf("child scope %q not loaded", gpuID)
	}
	if umbrella.ParentScopeID != "" {
		t.Errorf("umbrella.ParentScopeID = %q, want empty (host-level top of hierarchy)", umbrella.ParentScopeID)
	}
	if gpu.ParentScopeID != umbrellaID {
		t.Errorf("gpu.ParentScopeID = %q, want %q", gpu.ParentScopeID, umbrellaID)
	}
	if gpu.HardwareClass != types.HardwareClassHardwareGPU {
		t.Errorf("gpu.HardwareClass = %q, want %q", gpu.HardwareClass, types.HardwareClassHardwareGPU)
	}
}
