package types

import "time"

// Addr is a placeholder for q8.Addr (the QBP Locale quaternion address type).
// Imported as a type alias once github.com/helpfulengineering/q8 is available
// as a Go module dependency. For now, treat as opaque 16-byte content-address.
type Addr [16]byte

// LocaleBounds defines the spatiotemporal region where a pattern holds.
// Spec v1.3 §11.1.
type LocaleBounds struct {
	SpatialMin  Addr      `json:"spatial_min"`
	SpatialMax  Addr      `json:"spatial_max"`
	TemporalMin time.Time `json:"temporal_min"`
	TemporalMax time.Time `json:"temporal_max"`
}

// GrainSpec defines the finest resolution at which a pattern is detectable.
// Spec v1.3 §11.1.
type GrainSpec struct {
	SpatialGrainM float64       `json:"spatial_grain_m"`
	TemporalGrain time.Duration `json:"temporal_grain"`
}
