package fspec

import (
	"github.com/aperrot42/panoramix/pkg/asterix"
)

// ComputeAvailableFields creates a computed field containing all available FSPEC fields for a message
func ComputeAvailableFields(msg *asterix.AsterixMessage) []string {
	if msg.Items == nil {
		msg.Items = make(map[string]any)
	}

	// Extract available field names from the decoded Items keys
	var availableFields []string
	for fieldName := range msg.Items {
		// Skip the computed field itself to avoid recursion
		if fieldName != "computed" {
			availableFields = append(availableFields, fieldName)
		}
	}

	return availableFields
}

// AddAvailableFields adds computed FSPEC available fields to the message
func AddAvailableFields(msg *asterix.AsterixMessage) {
	if msg.Items == nil {
		msg.Items = make(map[string]any)
	}

	// Get or create computed field
	var computed map[string]any
	if existing, ok := msg.Items["computed"]; ok {
		computed = existing.(map[string]any)
	} else {
		computed = make(map[string]any)
		msg.Items["computed"] = computed
	}

	// Add available fields to computed
	availableFields := ComputeAvailableFields(msg)
	computed["fspec_available_fields"] = availableFields
}
