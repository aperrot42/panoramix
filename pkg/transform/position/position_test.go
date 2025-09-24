package position

import (
	"testing"
	"time"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

func TestPositionExtractorIntegration(t *testing.T) {
	// Create position extractor
	extractor := NewPositionExtractor()

	// Add a radar position
	radarPos := RadarPosition{
		Latitude:  46.5197, // Example Swiss radar position
		Longitude: 6.6323,
		Height:    430.0,
		SIC:       33,
		SAC:       40,
	}
	extractor.AddRadarPosition(radarPos)

	// Create a mock ASTERIX CAT 048 message
	msg := &asterix.AsterixMessage{
		Category: 48,
		Sic:      33,
		Sac:      40,
		Items: map[string]interface{}{
			"I048/240": "EZS36LP",        // Aircraft ID
			"I048/220": uint32(0x4B1A2B), // Aircraft address
			"I048/161": uint16(2063),     // Track number
			"I048/070": map[string]interface{}{
				"mode3a_code": "1000",
			},
			"I048/090": struct {
				FlightLevel float64
				RawValue    int16
				Validated   bool
				Garbled     bool
			}{
				FlightLevel: 100.0,
				RawValue:    400, // 100 * 4 (since 1/4 FL resolution)
				Validated:   true,
				Garbled:     false,
			},
			"I048/040": map[string]interface{}{
				"rho_nm":    96.1640625,
				"theta_deg": 31.3330078125,
			},
			"I048/042": map[string]interface{}{
				"x_nm": 50.0, // X coordinate in nautical miles
				"y_nm": 82.0, // Y coordinate in nautical miles
			},
			"I048/200": map[string]interface{}{
				"ground_speed_kt": 249.48,
				"track_angle_deg": 217.4853515625,
			},
		},
	}

	// Extract aircraft observation
	timestamp := time.Now()
	obs, err := extractor.ExtractFromMessage(msg, timestamp)

	if err != nil {
		t.Fatalf("Failed to extract aircraft observation: %v", err)
	}

	if obs == nil {
		t.Fatalf("Expected 1 observation, got nil")
	}

	// Verify basic fields
	if obs.RadarSIC != 33 || obs.RadarSAC != 40 {
		t.Errorf("Expected SIC=33 SAC=40, got SIC=%d SAC=%d", obs.RadarSIC, obs.RadarSAC)
	}

	// Verify WGS84 position is computed
	if obs.WGS84Position.PostionSource != "POLAR" {
		t.Errorf("Expected WGS84 source 'POLAR', got '%s'", obs.WGS84Position.PostionSource)
	}

	// Verify WGS84 coordinates are reasonable (should be in Europe)
	if obs.WGS84Position.Latitude_deg < 45 || obs.WGS84Position.Latitude_deg > 50 {
		t.Errorf("Expected latitude between 45-50°, got %.6f°", obs.WGS84Position.Latitude_deg)
	}

	if obs.WGS84Position.Longitude_deg < 5 || obs.WGS84Position.Longitude_deg > 10 {
		t.Errorf("Expected longitude between 5-10°, got %.6f°", obs.WGS84Position.Longitude_deg)
	}

	// Verify altitude is computed from flight level
	expectedAltitude := 10000. // FL100 -> 1000ft
	if obs.WGS84Position.AltitudeFt != expectedAltitude {
		t.Errorf("Expected altitude %.1f m, got %.1f m", expectedAltitude, obs.WGS84Position.AltitudeFt)
	}

	t.Logf("  WGS84: %.6f°, %.6f° (%.0f m)",
		obs.WGS84Position.Latitude_deg,
		obs.WGS84Position.Longitude_deg,
		obs.WGS84Position.AltitudeFt)
}

func TestRadarRegistry(t *testing.T) {
	registry := NewRadarRegistry()

	// Add radar positions
	pos1 := RadarPosition{Latitude: 46.5, Longitude: 6.6, Height: 400, SIC: 33, SAC: 40}
	pos2 := RadarPosition{Latitude: 47.4, Longitude: 8.5, Height: 500, SIC: 34, SAC: 41}

	registry.AddRadarPosition(pos1)
	registry.AddRadarPosition(pos2)

	// Test retrieval
	retrieved, found := registry.GetRadarPosition(33, 40)
	if !found {
		t.Error("Should find radar position for SIC=33 SAC=40")
	}
	if retrieved.Latitude != 46.5 {
		t.Errorf("Expected latitude 46.5, got %.1f", retrieved.Latitude)
	}

	// Test non-existent radar
	_, found = registry.GetRadarPosition(99, 99)
	if found {
		t.Error("Should not find radar position for SIC=99 SAC=99")
	}
}
