package position

import (
	"testing"
)

func TestLoadRadarConfig(t *testing.T) {
	// Test loading the DOLS radar configuration
	registry, err := LoadRadarConfig("../../../radar_config.yaml")
	if err != nil {
		t.Fatalf("Failed to load radar config: %v", err)
	}

	// Test retrieving DOLS radar position (SAC=40, SIC=33)
	pos, found := registry.GetRadarPosition(33, 40)
	if !found {
		t.Fatal("DOLS radar position not found in registry")
	}

	// Verify coordinates
	expectedLat := 46.42567939
	expectedLon := 6.09995219
	expectedAlt := 1734.13

	if pos.Latitude != expectedLat {
		t.Errorf("Expected latitude %f, got %f", expectedLat, pos.Latitude)
	}
	if pos.Longitude != expectedLon {
		t.Errorf("Expected longitude %f, got %f", expectedLon, pos.Longitude)
	}
	if pos.Height != expectedAlt {
		t.Errorf("Expected altitude %f, got %f", expectedAlt, pos.Height)
	}
	if pos.SAC != 40 {
		t.Errorf("Expected SAC 40, got %d", pos.SAC)
	}
	if pos.SIC != 33 {
		t.Errorf("Expected SIC 33, got %d", pos.SIC)
	}
}

func TestNewPositionExtractorFromConfig(t *testing.T) {
	extractor, err := NewPositionExtractorFromConfig("../../../radar_config.yaml")
	if err != nil {
		t.Fatalf("Failed to create PositionExtractor from config: %v", err)
	}

	// Test that the DOLS radar is available
	pos, found := extractor.radarRegistry.GetRadarPosition(33, 40)
	if !found {
		t.Fatal("DOLS radar position not found in extractor registry")
	}

	if pos.Latitude != 46.42567939 {
		t.Errorf("Expected DOLS latitude 46.42567939, got %f", pos.Latitude)
	}
}
