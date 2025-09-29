package bds

import (
	"encoding/hex"
	"testing"
)

func TestRealBDSData(t *testing.T) {
	// Real BDS data extracted from CAT048 test messages
	tests := []struct {
		name    string
		bds1    uint8
		bds2    uint8
		hexData string
	}{
		// From Message 1
		{
			name:    "Message_1_BDS_1",
			bds1:    0x05,
			bds2:    0x00,
			hexData: "FFD4F7396004E0",
			// BDS 5,0 - Track and Turn Report
		},
		{
			name:    "Message_1_BDS_2",
			bds1:    0x06,
			bds2:    0x00,
			hexData: "A609F930A00C00",
			// BDS 6,0 - Heading and Speed Report
		},
		{
			name:    "Message_1_BDS_3",
			bds1:    0x04,
			bds2:    0x00,
			hexData: "C8480030A40000",
			// BDS 4,0 - Selected Vertical Intention
		},
		// From Message 2
		{
			name:    "Message_2_BDS_1",
			bds1:    0x05,
			bds2:    0x00,
			hexData: "F01A9D12FD2C52",
			// BDS 5,0 - Track and Turn Report
		},
		{
			name:    "Message_2_BDS_2",
			bds1:    0x06,
			bds2:    0x00,
			hexData: "D2A92F0FA1F43B",
			// BDS 6,0 - Heading and Speed Report
		},
		{
			name:    "Message_2_BDS_3",
			bds1:    0x04,
			bds2:    0x00,
			hexData: "97700030A40000",
			// BDS 4,0 - Selected Vertical Intention
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hexData)
			if err != nil {
				t.Fatalf("Failed to decode hex: %v", err)
			}

			decoded, err := Decode(tt.bds1, tt.bds2, data)
			if err != nil {
				t.Fatalf("Failed to decode BDS %X,%X: %v", tt.bds1, tt.bds2, err)
			}

			t.Logf("BDS %X,%X decoded successfully", tt.bds1, tt.bds2)

			// Check specific fields based on BDS type
			bdsCode := (tt.bds1 << 4) | tt.bds2
			switch bdsCode {
			case 0x40: // BDS 4,0
				if rawMap, ok := decoded.(map[string]interface{}); ok {
					if alt, ok := rawMap["selected_altitude_ft"].(float64); ok {
						t.Logf("  Selected altitude: %.0f ft", alt)
					}
				}
			case 0x50: // BDS 5,0
				if bds50, ok := decoded.(BDS50Decoded); ok {
					if bds50.RollAngleValid {
						t.Logf("  Roll angle: %.2f degrees", bds50.RollAngleDeg)
					}
					if bds50.GroundSpeedValid {
						t.Logf("  Ground speed: %.0f knots", bds50.GroundSpeedKt)
					}
				}
			case 0x60: // BDS 6,0
				if rawMap, ok := decoded.(map[string]interface{}); ok {
					if heading, ok := rawMap["magnetic_heading_deg"].(float64); ok {
						t.Logf("  Magnetic heading: %.2f degrees", heading)
					}
					if ias, ok := rawMap["indicated_airspeed_kt"].(float64); ok {
						t.Logf("  Indicated airspeed: %.0f knots", ias)
					}
				}
			}
		})
	}
}

func TestBDS40_SelectedVerticalIntention(t *testing.T) {
	// Test data for BDS 4,0
	hexData := "A8282880380300" // Sample BDS 4,0 data
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(4, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 4,0: %v", err)
	}

	// For BDS 4,0 (unknown), it returns a raw map
	if rawMap, ok := decoded.(map[string]interface{}); ok {
		if alt, ok := rawMap["selected_altitude_ft"].(float64); ok {
			t.Logf("Selected altitude: %.0f ft", alt)
		}
		if alt, ok := rawMap["fms_altitude_ft"].(float64); ok {
			t.Logf("FMS altitude: %.0f ft", alt)
		}
		if pressure, ok := rawMap["baro_pressure_mb"].(float64); ok {
			t.Logf("Barometric pressure: %.1f mb", pressure)
		}
	}
}

func TestBDS50_TrackAndTurnReport(t *testing.T) {
	// Test data for BDS 5,0 with various fields set
	hexData := "9F3840B7E90200" // Sample BDS 5,0 data
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(5, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 5,0: %v", err)
	}

	// Check if we got a typed BDS50Decoded struct
	if bds50, ok := decoded.(BDS50Decoded); ok {
		// Check roll angle
		if bds50.RollAngleValid {
			t.Logf("Roll angle: %.2f degrees", bds50.RollAngleDeg)
			// Validate reasonable range
			if bds50.RollAngleDeg < -90 || bds50.RollAngleDeg > 90 {
				t.Errorf("Roll angle out of range: %.2f", bds50.RollAngleDeg)
			}
		}

		// Check ground speed
		if bds50.GroundSpeedValid {
			t.Logf("Ground speed: %.0f knots", bds50.GroundSpeedKt)
			// Validate reasonable range
			if bds50.GroundSpeedKt < 0 || bds50.GroundSpeedKt > 1000 {
				t.Errorf("Ground speed out of range: %.0f", bds50.GroundSpeedKt)
			}
		}

		// Check track angle
		if bds50.TrueTrackAngleValid {
			t.Logf("True track angle: %.2f degrees", bds50.TrueTrackAngleDeg)
			// Validate range 0-360
			if bds50.TrueTrackAngleDeg < 0 || bds50.TrueTrackAngleDeg > 360 {
				t.Errorf("Track angle out of range: %.2f", bds50.TrueTrackAngleDeg)
			}
		}
	}
}

func TestBDS60_HeadingAndSpeedReport(t *testing.T) {
	// Test data for BDS 6,0
	hexData := "990940B80C0100" // Sample BDS 6,0 data
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(6, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 6,0: %v", err)
	}

	// For BDS 6,0 (unknown), it returns a raw map
	if rawMap, ok := decoded.(map[string]interface{}); ok {
		if heading, ok := rawMap["magnetic_heading_deg"].(float64); ok {
			t.Logf("Magnetic heading: %.2f degrees", heading)
			// Validate range 0-360
			if heading < 0 || heading > 360 {
				t.Errorf("Magnetic heading out of range: %.2f", heading)
			}
		}
		if ias, ok := rawMap["indicated_airspeed_kt"].(float64); ok {
			t.Logf("Indicated airspeed: %.0f knots", ias)
			// Validate reasonable range
			if ias < 0 || ias > 600 {
				t.Errorf("IAS out of range: %.0f", ias)
			}
		}
		if mach, ok := rawMap["mach_number"].(float64); ok {
			t.Logf("Mach number: %.3f", mach)
			// Validate reasonable range
			if mach < 0 || mach > 2.0 {
				t.Errorf("Mach number out of range: %.3f", mach)
			}
		}
	}
}

func TestBDS44_MeteorologicalReport(t *testing.T) {
	// Test data for BDS 4,4
	hexData := "8D40000000000000" // Sample BDS 4,4 data
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	// Adjust the data to have 7 bytes
	if len(data) > 7 {
		data = data[:7]
	}

	decoded, err := Decode(4, 4, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 4,4: %v", err)
	}

	// For BDS 4,4 (unknown), it returns a raw map
	if rawMap, ok := decoded.(map[string]interface{}); ok {
		if windSpeed, ok := rawMap["wind_speed_kt"].(uint32); ok {
			t.Logf("Wind speed: %d knots", windSpeed)
		}
		if windDir, ok := rawMap["wind_direction_deg"].(float64); ok {
			t.Logf("Wind direction: %.1f degrees", windDir)
		}
		if temp, ok := rawMap["static_air_temperature_c"].(float64); ok {
			t.Logf("Static air temperature: %.2f°C", temp)
			// Validate reasonable range
			if temp < -80 || temp > 60 {
				t.Errorf("Temperature out of range: %.2f°C", temp)
			}
		}
		if turb, ok := rawMap["turbulence"].(string); ok {
			t.Logf("Turbulence: %s", turb)
		}
	}
}

func TestBDS10_DataLinkCapability(t *testing.T) {
	// Test data for BDS 1,0
	hexData := "0000F21F820700" // Sample BDS 1,0 data
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(1, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 1,0: %v", err)
	}

	// For BDS 1,0 (unknown), it returns a raw map
	if rawMap, ok := decoded.(map[string]interface{}); ok {
		capabilities := []string{
			"bds_20_cap", "bds_40_cap", "bds_50_cap", "bds_60_cap",
		}

		for _, cap := range capabilities {
			if val, ok := rawMap[cap].(bool); ok {
				t.Logf("%s: %v", cap, val)
			}
		}
	}
}

func TestUnknownBDSCode(t *testing.T) {
	// Test with an unknown BDS code
	hexData := "1234567890ABCD"
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(9, 9, data) // BDS 9,9 doesn't exist
	if err != nil {
		t.Fatalf("Unexpected error for unknown BDS code: %v", err)
	}

	// Should return raw data with unknown flag
	if rawMap, ok := decoded.(map[string]interface{}); ok {
		if unknown, ok := rawMap["unknown"].(bool); !ok || !unknown {
			t.Error("Expected unknown flag to be set for unknown BDS code")
		}
		if raw, ok := rawMap["raw"].(string); !ok || raw == "" {
			t.Error("Expected raw data to be present for unknown BDS code")
		}
	}
}

func TestExtractBits(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		startBit int
		numBits  int
		expected uint32
	}{
		{
			name:     "Extract first byte",
			data:     []byte{0xFF, 0x00, 0x00},
			startBit: 0,
			numBits:  8,
			expected: 0xFF,
		},
		{
			name:     "Extract across byte boundary",
			data:     []byte{0x0F, 0xF0, 0x00},
			startBit: 4,
			numBits:  8,
			expected: 0xFF,
		},
		{
			name:     "Extract single bit",
			data:     []byte{0x80, 0x00, 0x00},
			startBit: 0,
			numBits:  1,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractBits(tt.data, tt.startBit, tt.numBits)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExtractSignedBits(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		startBit int
		numBits  int
		expected int32
	}{
		{
			name:     "Positive value",
			data:     []byte{0x7F, 0xFF, 0x00},
			startBit: 0,
			numBits:  8,
			expected: 127,
		},
		{
			name:     "Negative value",
			data:     []byte{0x80, 0x00, 0x00},
			startBit: 0,
			numBits:  8,
			expected: -128,
		},
		{
			name:     "Negative value across boundary",
			data:     []byte{0x0F, 0x80, 0x00},
			startBit: 4,
			numBits:  8,
			expected: -8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSignedBits(tt.data, tt.startBit, tt.numBits)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
