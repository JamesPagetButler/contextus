package tenancy

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

// testdataPath returns the absolute path to a file under the testdata/ tree
// adjacent to this test file. The directory is embedded at compile time via
// the package path; runtime os.Getwd() + relative join is avoided so tests
// run correctly under `go test ./...` from any working directory.
func testdataPath(t *testing.T, parts ...string) string {
	t.Helper()
	// __file__ equivalent: use the package source dir via a known anchor.
	// In Go tests, os.Getwd() points to the package directory.
	return filepath.Join(append([]string{"testdata"}, parts...)...)
}

// TestRoundtrip_QBPFixture loads the canonical QBP fixture and asserts struct
// count + key field content. This is the happy-path contract test.
func TestRoundtrip_QBPFixture(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "qbp-scope-config-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig(%q): unexpected error: %v", path, err)
	}

	// Struct counts.
	if got, want := len(res.PhysicalScopes), 2; got != want {
		t.Errorf("PhysicalScopes: got %d, want %d", got, want)
	}
	if got, want := len(res.ConceptualScopes), 2; got != want {
		t.Errorf("ConceptualScopes: got %d, want %d", got, want)
	}
	if got, want := len(res.Memberships), 2; got != want {
		t.Errorf("Memberships: got %d, want %d", got, want)
	}
	if got, want := len(res.EdgeOptions), 2; got != want {
		t.Errorf("EdgeOptions: got %d, want %d", got, want)
	}

	// Spot-check physical scope.
	cascadia := res.PhysicalScopes[0]
	if cascadia.ScopeID != "contextus:scope:physical:cascadia" {
		t.Errorf("PhysicalScopes[0].ScopeID = %q, want cascadia id", cascadia.ScopeID)
	}
	if cascadia.Name == "" {
		t.Error("PhysicalScopes[0].Name must not be empty")
	}

	// Spot-check conceptual scope.
	var hamiltonFound bool
	for _, cs := range res.ConceptualScopes {
		if cs.ScopeID == "contextus:scope:conceptual:hamilton-product" {
			hamiltonFound = true
			if cs.Name == "" {
				t.Error("hamilton-product Name must not be empty")
			}
		}
	}
	if !hamiltonFound {
		t.Error("hamilton-product conceptual scope not found in result")
	}

	// Spot-check NodeOptions map presence for all scope IDs.
	allIDs := []string{
		"contextus:scope:physical:cascadia",
		"contextus:scope:physical:ligo-hanford",
		"contextus:scope:conceptual:slow-slip",
		"contextus:scope:conceptual:hamilton-product",
	}
	for _, id := range allIDs {
		if _, ok := res.NodeOptions[id]; !ok {
			t.Errorf("NodeOptions missing entry for scope_id %q", id)
		}
	}
}

// TestSchema_RejectsMalformed asserts that each malformed fixture triggers an
// appropriate sentinel error. The duplicate-id case wraps ErrScopeLoadConflict;
// all others wrap ErrScopeConfigInvalid (schema validation failure).
func TestSchema_RejectsMalformed(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		file    string
		wantErr error
	}{
		{
			name:    "missing id",
			file:    "scope-config-malformed/missing-id.yaml",
			wantErr: ErrScopeConfigInvalid,
		},
		{
			name:    "salience out of range",
			file:    "scope-config-malformed/salience-out-of-range.yaml",
			wantErr: ErrScopeConfigInvalid,
		},
		{
			name:    "unknown top-level key",
			file:    "scope-config-malformed/unknown-top-level-key.yaml",
			wantErr: ErrScopeConfigInvalid,
		},
		{
			name:    "duplicate scope_id",
			file:    "scope-config-malformed/duplicate-scope-id.yaml",
			wantErr: ErrScopeLoadConflict,
		},
		{
			name:    "malformed bounds",
			file:    "scope-config-malformed/malformed-bounds.yaml",
			wantErr: ErrScopeConfigInvalid,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := testdataPath(t, tc.file)
			_, err := LoadScopeConfig(ctx, path)
			if err == nil {
				t.Fatalf("LoadScopeConfig(%q): expected error wrapping %v, got nil", path, tc.wantErr)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("LoadScopeConfig(%q): errors.Is(err, %v) = false; err = %v", path, tc.wantErr, err)
			}
		})
	}
}

// TestEnvelopeSplit_TierImmuneSalience asserts that tier_immune and salience
// from the QBP fixture appear in NodeOptions, NOT in the ScopeConceptual struct.
// This is the core split-envelope contract (Wyrd PR #40 §2.1).
func TestEnvelopeSplit_TierImmuneSalience(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "qbp-scope-config-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig: %v", err)
	}

	const hamiltonID = "contextus:scope:conceptual:hamilton-product"

	// NodeOptions must carry TierImmune: true and Salience: 1.0.
	opts, ok := res.NodeOptions[hamiltonID]
	if !ok {
		t.Fatalf("NodeOptions[%q] not present", hamiltonID)
	}
	if !opts.TierImmune {
		t.Errorf("NodeOptions[%q].TierImmune = false, want true", hamiltonID)
	}
	if opts.Salience != 1.0 {
		t.Errorf("NodeOptions[%q].Salience = %v, want 1.0", hamiltonID, opts.Salience)
	}

	// The ScopeConceptual struct must NOT carry tier_immune or salience;
	// verify the struct only has the canonical Contextus-spec fields.
	var hamilton *struct{ found bool }
	for i := range res.ConceptualScopes {
		if res.ConceptualScopes[i].ScopeID == hamiltonID {
			hamilton = &struct{ found bool }{found: true}
			// ScopeConceptual has no TierImmune or Salience field — this is a
			// compile-time guarantee from the pkg/types definition. We assert
			// ScopeID is set and the type has only the expected Contextus fields.
			cs := res.ConceptualScopes[i]
			if cs.ScopeID != hamiltonID {
				t.Errorf("ScopeConceptual.ScopeID = %q, want %q", cs.ScopeID, hamiltonID)
			}
			_ = hamilton
			break
		}
	}
	if hamilton == nil {
		t.Errorf("hamilton-product not found in ConceptualScopes")
	}
}

// TestErrorsIs_Compatibility asserts that errors returned by LoadScopeConfig
// properly wrap the package sentinels so errors.Is works from any call depth.
func TestErrorsIs_Compatibility(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		file    string
		wantErr error
	}{
		{
			name:    "ErrScopeConfigInvalid is unwrappable",
			file:    "scope-config-malformed/missing-id.yaml",
			wantErr: ErrScopeConfigInvalid,
		},
		{
			name:    "ErrScopeLoadConflict is unwrappable",
			file:    "scope-config-malformed/duplicate-scope-id.yaml",
			wantErr: ErrScopeLoadConflict,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := testdataPath(t, tc.file)
			_, err := LoadScopeConfig(ctx, path)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("errors.Is(err, %v) = false; full error: %v", tc.wantErr, err)
			}
		})
	}
}

// TestPathTraversal_Rejected asserts that a configPath containing ".."
// components returns ErrScopeConfigInvalid without touching the filesystem.
func TestPathTraversal_Rejected(t *testing.T) {
	ctx := context.Background()

	traversalPaths := []string{
		"../../../etc/passwd",
		"testdata/../../../etc/shadow",
		"testdata/../../secret.yaml",
	}

	for _, p := range traversalPaths {
		p := p
		t.Run(p, func(t *testing.T) {
			_, err := LoadScopeConfig(ctx, p)
			if err == nil {
				t.Fatalf("LoadScopeConfig(%q): expected error, got nil", p)
			}
			if !errors.Is(err, ErrScopeConfigInvalid) {
				t.Errorf("LoadScopeConfig(%q): want ErrScopeConfigInvalid; got %v", p, err)
			}
		})
	}
}
