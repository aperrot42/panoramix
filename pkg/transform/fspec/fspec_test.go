package fspec

import (
	"testing"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

func TestAddAvailableFields(t *testing.T) {
	// Create a sample ASTERIX message with some decoded fields
	msg := &asterix.AsterixMessage{
		Category: 48,
		Sic:      1,
		Sac:      2,
		Items: map[string]any{
			"I048/010": map[string]uint8{"SIC": 1, "SAC": 2},
			"I048/130": "some radar characteristics",
			"I048/090": "flight level data",
		},
		FSPEC: []byte{0x80, 0x40}, // Sample FSPEC
	}

	// Apply the transformation
	AddAvailableFields(msg)

	// Check that computed field was added
	computed, ok := msg.Items["computed"]
	if !ok {
		t.Fatal("computed field not added")
	}

	computedMap, ok := computed.(map[string]any)
	if !ok {
		t.Fatal("computed field is not a map")
	}

	// Check that fspec_available_fields was added
	availableFields, ok := computedMap["fspec_available_fields"]
	if !ok {
		t.Fatal("fspec_available_fields not added to computed")
	}

	fields, ok := availableFields.([]string)
	if !ok {
		t.Fatal("fspec_available_fields is not a []string")
	}

	// Check that the fields match what we expect
	expectedFields := []string{"I048/010", "I048/130", "I048/090"}
	if len(fields) != len(expectedFields) {
		t.Fatalf("expected %d fields, got %d", len(expectedFields), len(fields))
	}

	// Convert to map for easier checking (order doesn't matter)
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}

	for _, expected := range expectedFields {
		if !fieldMap[expected] {
			t.Errorf("expected field %s not found in available fields", expected)
		}
	}
}

func TestAddAvailableFields_EmptyMessage(t *testing.T) {
	msg := &asterix.AsterixMessage{
		Category: 48,
		Items:    make(map[string]any),
	}

	AddAvailableFields(msg)

	computed := msg.Items["computed"].(map[string]any)
	fields := computed["fspec_available_fields"].([]string)

	if len(fields) != 0 {
		t.Errorf("expected empty fields list, got %v", fields)
	}
}

func TestAddAvailableFields_ExistingComputed(t *testing.T) {
	// Test that we don't overwrite existing computed fields
	msg := &asterix.AsterixMessage{
		Category: 48,
		Items: map[string]any{
			"I048/010": "data",
			"computed": map[string]any{
				"existing_field": "existing_value",
			},
		},
	}

	AddAvailableFields(msg)

	computed := msg.Items["computed"].(map[string]any)

	// Check existing field is preserved
	if computed["existing_field"] != "existing_value" {
		t.Error("existing computed field was overwritten")
	}

	// Check new field was added
	fields := computed["fspec_available_fields"].([]string)
	if len(fields) != 1 || fields[0] != "I048/010" {
		t.Errorf("expected [I048/010], got %v", fields)
	}
}
