// Package tenancy provides tenant scope-node configuration loading for the
// Contextus federation layer. See Spec v1.3 §4.6 (Scope Nodes) and §11.4 (Go
// type definitions). The counterpart Wyrd-side loader is wyrd/store.LoadScopeConfig
// (Wyrd PR #40 design surface / Wyrd issue #33).
package tenancy

import "errors"

// Sentinel errors for LoadScopeConfig. Names are intentionally identical to
// the wyrd/store.Err* sentinels defined in Wyrd PR #40 §6 so that consumers
// can swap tenancy.Err* for store.Err* with a one-line import change once the
// Wyrd impl PR lands. errors.Is unwrapping works across the boundary because
// both sides follow the same fmt.Errorf("...: %w", Err*) wrapping convention.
var (
	// ErrScopeConfigParse is returned when the YAML or JSON cannot be decoded.
	ErrScopeConfigParse = errors.New("tenancy: scope-config parse error")

	// ErrScopeConfigInvalid is returned when the decoded config fails JSON
	// Schema 2020-12 validation or a semantic constraint (e.g. empty id,
	// salience out of range, bounds min > max, path traversal).
	ErrScopeConfigInvalid = errors.New("tenancy: scope-config validation failed")

	// ErrScopeLoadConflict is returned when the same scope_id appears more than
	// once within the config (intra-file duplicate). v0.2 may add an upsert
	// variant; v0.1 errors on collision per Wyrd PR #40 §3 open-question resolution.
	ErrScopeLoadConflict = errors.New("tenancy: scope-node ID already present in config")
)
