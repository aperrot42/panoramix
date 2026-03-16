package pipeline

import "github.com/aperrot42/panoramix/pkg/transform/position"

// NewPositionEnrichFilter returns a filter that computes WGS84 positions.
func NewPositionEnrichFilter(extractor *position.PositionExtractor) func(AsterixResult) (AsterixResult, bool, error) {
	return func(ar AsterixResult) (AsterixResult, bool, error) {
		pos, err := extractor.ExtractFromMessage(&ar.Message, ar.Record.Timestamp)
		if err == nil {
			if ar.Computed == nil {
				ar.Computed = make(map[string]any)
			}
			ar.Computed["position"] = pos
		}
		return ar, true, nil
	}
}
