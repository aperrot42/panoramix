package pipeline

import "github.com/aperrot42/panoramix/pkg/transform/fspec"

// FSPECEnrichFilter adds computed FSPEC available fields to the result.
func FSPECEnrichFilter(ar AsterixResult) (AsterixResult, bool, error) {
	fields := fspec.ComputeAvailableFields(&ar.Message)
	if ar.Computed == nil {
		ar.Computed = make(map[string]any)
	}
	ar.Computed["fspec_available_fields"] = fields
	return ar, true, nil
}
