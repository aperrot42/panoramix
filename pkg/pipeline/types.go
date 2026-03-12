package pipeline

import (
	"github.com/aperrot42/panoramix/pkg/asterix"
	"github.com/aperrot42/panoramix/pkg/internal_format"
)

// AsterixResult combines the original record metadata with the decoded message.
type AsterixResult struct {
	Record   internal_format.Record
	Message  asterix.AsterixMessage
	Computed map[string]any
}
