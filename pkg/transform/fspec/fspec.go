package fspec

import (
	"github.com/aperrot42/panoramix/pkg/asterix"
)

// ComuteAvailableFields creates a computed field containing all available FSPEC fields for a message
func ComputeAvailableFields(msg *asterix.AsterixMessage) []string {
	if msg.Items == nil {
		msg.Items = make(map[string]any)
	}

	// Extract available field names from the decoded Items keys
	var availableFields []string
	for fieldName := range msg.Items {
		availableFields = append(availableFields, fieldName)
	}

	return availableFields
}
