package tenancy

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

// TestTenantProfile_RoundTrip loads the QBP tenant-profile fixture (mirror of
// Contextus-Spec-Addendum-Research-Aid-Tenancy §6 worked example) and asserts
// every field round-trips through the YAML → struct decode path.
func TestTenantProfile_RoundTrip(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "qbp-tenant-profile-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig(%q): unexpected error: %v", path, err)
	}

	if res.TenantProfile == nil {
		t.Fatalf("LoadResult.TenantProfile is nil; expected populated")
	}

	if got, want := res.TenantProfile.TenantID, "qbp"; got != want {
		t.Errorf("TenantID = %q, want %q", got, want)
	}

	wantScaffold := []string{
		"PRECEDENT_GRAPH",
		"EVIDENCE_LATTICE",
		"ALGEBRAIC_STRUCTURE_SCAFFOLD",
		"SOURCE_LOCATION_HYPOTHESIS",
	}
	if !reflect.DeepEqual(res.TenantProfile.SubscriberProfile.AcceptedScaffoldTypes, wantScaffold) {
		t.Errorf("AcceptedScaffoldTypes = %v, want %v",
			res.TenantProfile.SubscriberProfile.AcceptedScaffoldTypes, wantScaffold)
	}

	wantCorpus := []string{
		"PHYSICS_PREPRINT",
		"JOURNAL_ARTICLE",
		"DATASET_DESCRIPTOR",
		"CODE_REPO",
	}
	if !reflect.DeepEqual(res.TenantProfile.SubscriberProfile.AcceptedCorpusClasses, wantCorpus) {
		t.Errorf("AcceptedCorpusClasses = %v, want %v",
			res.TenantProfile.SubscriberProfile.AcceptedCorpusClasses, wantCorpus)
	}

	wantConsumers := []string{"self"}
	if !reflect.DeepEqual(res.TenantProfile.SubscriberProfile.IntendedConsumersDefault, wantConsumers) {
		t.Errorf("IntendedConsumersDefault = %v, want %v",
			res.TenantProfile.SubscriberProfile.IntendedConsumersDefault, wantConsumers)
	}

	if got, want := res.TenantProfile.TenantSubgraphRef.URI, "cth://tenant/qbp/subgraph"; got != want {
		t.Errorf("TenantSubgraphRef.URI = %q, want %q", got, want)
	}
}

// TestTenantProfile_BackwardsCompat is the critical regression guard for
// Contextus-Spec-Addendum-Research-Aid-Tenancy §3.4: a v1.3 scope-config
// without a tenant_profile block MUST load with LoadResult.TenantProfile == nil
// and all other fields intact.
func TestTenantProfile_BackwardsCompat(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "qbp-scope-config-v0_1.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig(%q): unexpected error: %v", path, err)
	}

	if res.TenantProfile != nil {
		t.Errorf("LoadResult.TenantProfile = %+v, want nil (v1.3 fixture has no tenant_profile block)",
			res.TenantProfile)
	}

	// Sanity check: the v1.3 baseline content still parsed correctly.
	if len(res.PhysicalScopes) == 0 {
		t.Error("PhysicalScopes is empty; v1.3 backwards-compat path corrupted other fields")
	}
	if len(res.ConceptualScopes) == 0 {
		t.Error("ConceptualScopes is empty; v1.3 backwards-compat path corrupted other fields")
	}
}

// TestTenantProfile_MalformedRejected asserts that a tenant_profile with an
// out-of-enum accepted_corpus_classes value is rejected with
// ErrScopeConfigInvalid (schema validation failure). This closes the
// PR #17 test-plan malformed-rejection checkbox.
func TestTenantProfile_MalformedRejected(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "scope-config-malformed/tenant-profile-unknown-corpus-class.yaml")

	_, err := LoadScopeConfig(ctx, path)
	if err == nil {
		t.Fatalf("LoadScopeConfig(%q): expected error wrapping ErrScopeConfigInvalid, got nil", path)
	}
	if !errors.Is(err, ErrScopeConfigInvalid) {
		t.Errorf("LoadScopeConfig(%q): errors.Is(err, ErrScopeConfigInvalid) = false; err = %v", path, err)
	}
}

// TestTenantProfile_SynthesizesSubgraphRef asserts that when tenant_subgraph_ref
// is omitted, the loader synthesises cth://tenant/<tenant_id>/subgraph per
// Contextus-Spec-Addendum-Research-Aid-Tenancy §2.4 + §3.3 step 3b.
func TestTenantProfile_SynthesizesSubgraphRef(t *testing.T) {
	ctx := context.Background()
	path := testdataPath(t, "qbp-tenant-profile-no-subgraph-ref.yaml")

	res, err := LoadScopeConfig(ctx, path)
	if err != nil {
		t.Fatalf("LoadScopeConfig(%q): unexpected error: %v", path, err)
	}

	if res.TenantProfile == nil {
		t.Fatalf("LoadResult.TenantProfile is nil; expected populated")
	}

	const want = "cth://tenant/qbp/subgraph"
	if got := res.TenantProfile.TenantSubgraphRef.URI; got != want {
		t.Errorf("TenantSubgraphRef.URI = %q, want %q (synthesised)", got, want)
	}
}
