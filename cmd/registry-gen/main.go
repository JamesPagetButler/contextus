// Command registry-gen generates the canonical JSON projection of the
// quantity-kind registry from its YAML governance source.
//
// Usage:
//
//	go run ./cmd/registry-gen ./registry/quantity-kinds.yaml
//
// The output is written to stdout; redirect to registry/quantity-kinds.json.
// YAML governs; JSON feeds. Edda Stage 1 pins by version and verifies by
// content hash. Spec Addendum NT-Signal-Measurement §5.6 (Q3 ruling, seq=512).
//
// Output contract:
//   - Deterministic (keys sorted lexicographically within each object)
//   - Self-describing: carries registry_version + content_hash (SHA-256 over
//     the JSON bytes with content_hash set to its zero value "")
//   - generated_by field records provenance
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

// rawRegistry is the decode target for quantity-kinds.yaml.
type rawRegistry struct {
	RegistryVersion int       `yaml:"registry_version"`
	Kinds           []rawKind `yaml:"kinds"`
}

// rawKind mirrors one entry in the YAML kinds list.
type rawKind struct {
	Kind         string `yaml:"kind"`
	OwnerTenant  string `yaml:"owner_tenant"`
	UnitSystem   string `yaml:"unit_system"`
	Description  string `yaml:"description"`
	IntroducedIn int    `yaml:"introduced_in"`
	DeprecatedIn int    `yaml:"deprecated_in,omitempty"`
	MapsTo       any    `yaml:"maps_to,omitempty"` // string or []string
}

// RegistryJSON is the canonical JSON projection shape.
type RegistryJSON struct {
	Schema          string     `json:"$schema"`
	RegistryVersion int        `json:"registry_version"`
	ContentHash     string     `json:"content_hash"` // SHA-256 over bytes with content_hash=""
	GeneratedAt     string     `json:"generated_at"` // RFC3339 UTC
	GeneratedBy     string     `json:"generated_by"`
	Kinds           []KindJSON `json:"kinds"`
}

// KindJSON is one entry in the JSON projection.
type KindJSON struct {
	Kind         string `json:"kind"`
	OwnerTenant  string `json:"owner_tenant"`
	UnitSystem   string `json:"unit_system"`
	Description  string `json:"description"`
	IntroducedIn int    `json:"introduced_in"`
	DeprecatedIn int    `json:"deprecated_in,omitempty"`
	MapsTo       any    `json:"maps_to,omitempty"` // string or []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: registry-gen <quantity-kinds.yaml>\n")
		os.Exit(1)
	}

	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "registry-gen: read %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	var reg rawRegistry
	if err := yaml.Unmarshal(raw, &reg); err != nil {
		fmt.Fprintf(os.Stderr, "registry-gen: parse YAML: %v\n", err)
		os.Exit(1)
	}

	if reg.RegistryVersion < 1 {
		fmt.Fprintf(os.Stderr, "registry-gen: registry_version must be >= 1\n")
		os.Exit(1)
	}

	// Sort kinds deterministically: by kind name, then owner_tenant as tie-break.
	// rawKind and KindJSON have the same field layout; convert directly.
	kinds := make([]KindJSON, len(reg.Kinds))
	for i, k := range reg.Kinds {
		kinds[i] = KindJSON(k)
	}
	sort.Slice(kinds, func(i, j int) bool {
		if kinds[i].Kind != kinds[j].Kind {
			return kinds[i].Kind < kinds[j].Kind
		}
		return kinds[i].OwnerTenant < kinds[j].OwnerTenant
	})

	out := RegistryJSON{
		Schema:          "https://github.com/JamesPagetButler/contextus/registry/quantity-kinds-schema.json",
		RegistryVersion: reg.RegistryVersion,
		ContentHash:     "", // computed below
		GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
		GeneratedBy:     "cmd/registry-gen",
		Kinds:           kinds,
	}

	// Deterministic JSON: marshal with sorted keys (encoding/json sorts struct
	// fields by declaration order; we control declaration order above).
	intermediateBytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "registry-gen: marshal intermediate JSON: %v\n", err)
		os.Exit(1)
	}

	// Content hash is SHA-256 of the JSON bytes with content_hash set to "".
	h := sha256.Sum256(intermediateBytes)
	out.ContentHash = hex.EncodeToString(h[:])

	finalBytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "registry-gen: marshal final JSON: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(finalBytes); err != nil {
		fmt.Fprintf(os.Stderr, "registry-gen: write output: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stdout.WriteString("\n"); err != nil {
		fmt.Fprintf(os.Stderr, "registry-gen: write newline: %v\n", err)
		os.Exit(1)
	}
}
