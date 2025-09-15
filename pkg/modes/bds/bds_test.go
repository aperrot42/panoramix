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

			reg, err := Decode(tt.bds1, tt.bds2, data)
			if err != nil {
				t.Fatalf("Failed to decode BDS %X,%X: %v", tt.bds1, tt.bds2, err)
			}

			t.Logf("BDS %X,%X decoded successfully", tt.bds1, tt.bds2)
			
			// Check specific fields based on BDS type
			bdsCode := (tt.bds1 << 4) | tt.bds2
			switch bdsCode {
			case 0x40: // BDS 4,0
				if valid, ok := reg.Fields["selected_altitude_valid"].(bool); ok && valid {
					if alt, ok := reg.Fields["selected_altitude_ft"].(float64); ok {
						t.Logf("  Selected altitude: %.0f ft", alt)
					}
				}
			case 0x50: // BDS 5,0
				if valid, ok := reg.Fields["roll_angle_valid"].(bool); ok && valid {
					if roll, ok := reg.Fields["roll_angle_deg"].(float64); ok {
						t.Logf("  Roll angle: %.2f degrees", roll)
					}
				}
				if valid, ok := reg.Fields["ground_speed_valid"].(bool); ok && valid {
					if speed, ok := reg.Fields["ground_speed_kt"].(float64); ok {
						t.Logf("  Ground speed: %.0f knots", speed)
					}
				}
			case 0x60: // BDS 6,0
				if valid, ok := reg.Fields["magnetic_heading_valid"].(bool); ok && valid {
					if heading, ok := reg.Fields["magnetic_heading_deg"].(float64); ok {
						t.Logf("  Magnetic heading: %.2f degrees", heading)
					}
				}
				if valid, ok := reg.Fields["indicated_airspeed_valid"].(bool); ok && valid {
					if ias, ok := reg.Fields["indicated_airspeed_kt"].(float64); ok {
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

	reg, err := Decode(4, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 4,0: %v", err)
	}

	// Check if selected altitude field exists
	if valid, ok := reg.Fields["selected_altitude_valid"].(bool); ok && valid {
		if alt, ok := reg.Fields["selected_altitude_ft"].(float64); ok {
			t.Logf("Selected altitude: %.0f ft", alt)
		}
	}

	// Check if FMS altitude field exists
	if valid, ok := reg.Fields["fms_altitude_valid"].(bool); ok && valid {
		if alt, ok := reg.Fields["fms_altitude_ft"].(float64); ok {
			t.Logf("FMS altitude: %.0f ft", alt)
		}
	}

	// Check if barometric pressure field exists
	if valid, ok := reg.Fields["baro_pressure_valid"].(bool); ok && valid {
		if pressure, ok := reg.Fields["baro_pressure_mb"].(float64); ok {
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

	reg, err := Decode(5, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 5,0: %v", err)
	}

	// Check roll angle
	if valid, ok := reg.Fields["roll_angle_valid"].(bool); ok && valid {
		if roll, ok := reg.Fields["roll_angle_deg"].(float64); ok {
			t.Logf("Roll angle: %.2f degrees", roll)
			// Validate reasonable range
			if roll < -90 || roll > 90 {
				t.Errorf("Roll angle out of range: %.2f", roll)
			}
		}
	}

	// Check ground speed
	if valid, ok := reg.Fields["ground_speed_valid"].(bool); ok && valid {
		if speed, ok := reg.Fields["ground_speed_kt"].(float64); ok {
			t.Logf("Ground speed: %.0f knots", speed)
			// Validate reasonable range
			if speed < 0 || speed > 1000 {
				t.Errorf("Ground speed out of range: %.0f", speed)
			}
		}
	}

	// Check track angle
	if valid, ok := reg.Fields["true_track_angle_valid"].(bool); ok && valid {
		if track, ok := reg.Fields["true_track_angle_deg"].(float64); ok {
			t.Logf("True track angle: %.2f degrees", track)
			// Validate range 0-360
			if track < 0 || track > 360 {
				t.Errorf("Track angle out of range: %.2f", track)
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

	reg, err := Decode(6, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 6,0: %v", err)
	}

	// Check magnetic heading
	if valid, ok := reg.Fields["magnetic_heading_valid"].(bool); ok && valid {
		if heading, ok := reg.Fields["magnetic_heading_deg"].(float64); ok {
			t.Logf("Magnetic heading: %.2f degrees", heading)
			// Validate range 0-360
			if heading < 0 || heading > 360 {
				t.Errorf("Magnetic heading out of range: %.2f", heading)
			}
		}
	}

	// Check indicated airspeed
	if valid, ok := reg.Fields["indicated_airspeed_valid"].(bool); ok && valid {
		if ias, ok := reg.Fields["indicated_airspeed_kt"].(float64); ok {
			t.Logf("Indicated airspeed: %.0f knots", ias)
			// Validate reasonable range
			if ias < 0 || ias > 600 {
				t.Errorf("IAS out of range: %.0f", ias)
			}
		}
	}

	// Check Mach number
	if valid, ok := reg.Fields["mach_number_valid"].(bool); ok && valid {
		if mach, ok := reg.Fields["mach_number"].(float64); ok {
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

	reg, err := Decode(4, 4, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 4,4: %v", err)
	}

	// Check wind data
	if valid, ok := reg.Fields["wind_valid"].(bool); ok && valid {
		if windSpeed, ok := reg.Fields["wind_speed_kt"].(uint32); ok {
			t.Logf("Wind speed: %d knots", windSpeed)
		}
		if windDir, ok := reg.Fields["wind_direction_deg"].(float64); ok {
			t.Logf("Wind direction: %.1f degrees", windDir)
		}
	}

	// Check temperature
	if valid, ok := reg.Fields["static_air_temperature_valid"].(bool); ok && valid {
		if temp, ok := reg.Fields["static_air_temperature_c"].(float64); ok {
			t.Logf("Static air temperature: %.2f°C", temp)
			// Validate reasonable range
			if temp < -80 || temp > 60 {
				t.Errorf("Temperature out of range: %.2f°C", temp)
			}
		}
	}

	// Check turbulence
	if valid, ok := reg.Fields["turbulence_valid"].(bool); ok && valid {
		if turb, ok := reg.Fields["turbulence"].(string); ok {
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

	reg, err := Decode(1, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 1,0: %v", err)
	}

	// Check various capability flags
	capabilities := []string{
		"bds_20_cap", "bds_40_cap", "bds_50_cap", "bds_60_cap",
	}

	for _, cap := range capabilities {
		if val, ok := reg.Fields[cap].(bool); ok {
			t.Logf("%s: %v", cap, val)
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

	reg, err := Decode(9, 9, data) // BDS 9,9 doesn't exist
	if err != nil {
		t.Fatalf("Unexpected error for unknown BDS code: %v", err)
	}

	// Should return raw data with unknown flag
	if unknown, ok := reg.Fields["unknown"].(bool); !ok || !unknown {
		t.Error("Expected unknown flag to be set for unknown BDS code")
	}

	if raw, ok := reg.Fields["raw"].(string); !ok || raw == "" {
		t.Error("Expected raw data to be present for unknown BDS code")
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