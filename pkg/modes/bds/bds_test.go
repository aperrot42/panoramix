package bds

import (
	"encoding/hex"
	"testing"
)

func TestRealBDSData(t *testing.T) {
	tests := []struct {
		name    string
		bds1    uint8
		bds2    uint8
		hexData string
	}{
		{name: "BDS_5_0_msg1", bds1: 0x05, bds2: 0x00, hexData: "FFD4F7396004E0"},
		{name: "BDS_6_0_msg1", bds1: 0x06, bds2: 0x00, hexData: "A609F930A00C00"},
		{name: "BDS_4_0_msg1", bds1: 0x04, bds2: 0x00, hexData: "C8480030A40000"},
		{name: "BDS_5_0_msg2", bds1: 0x05, bds2: 0x00, hexData: "F01A9D12FD2C52"},
		{name: "BDS_6_0_msg2", bds1: 0x06, bds2: 0x00, hexData: "D2A92F0FA1F43B"},
		{name: "BDS_4_0_msg2", bds1: 0x04, bds2: 0x00, hexData: "97700030A40000"},
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

			bdsCode := (tt.bds1 << 4) | tt.bds2
			switch bdsCode {
			case 0x40:
				bds40, ok := decoded.(BDS40Decoded)
				if !ok {
					t.Fatalf("expected BDS40Decoded, got %T", decoded)
				}
				if bds40.SelectedAltitudeValid {
					t.Logf("  Selected altitude: %.0f ft", bds40.SelectedAltitudeFt)
				}
			case 0x50:
				bds50, ok := decoded.(BDS50Decoded)
				if !ok {
					t.Fatalf("expected BDS50Decoded, got %T", decoded)
				}
				if bds50.RollAngleValid {
					t.Logf("  Roll angle: %.2f degrees", bds50.RollAngleDeg)
				}
				if bds50.GroundSpeedValid {
					t.Logf("  Ground speed: %.0f knots", bds50.GroundSpeedKt)
				}
			case 0x60:
				bds60, ok := decoded.(BDS60Decoded)
				if !ok {
					t.Fatalf("expected BDS60Decoded, got %T", decoded)
				}
				if bds60.MagneticHeadingValid {
					t.Logf("  Magnetic heading: %.2f degrees", bds60.MagneticHeadingDeg)
				}
				if bds60.IndicatedAirspeedValid {
					t.Logf("  Indicated airspeed: %.0f knots", bds60.IndicatedAirspeedKt)
				}
			}
		})
	}
}

func TestBDS40_SelectedVerticalIntention(t *testing.T) {
	data, err := hex.DecodeString("A8282880380300")
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(4, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 4,0: %v", err)
	}

	bds40, ok := decoded.(BDS40Decoded)
	if !ok {
		t.Fatalf("expected BDS40Decoded, got %T", decoded)
	}
	if bds40.SelectedAltitudeValid {
		t.Logf("Selected altitude: %.0f ft", bds40.SelectedAltitudeFt)
	}
	if bds40.FmsAltitudeValid {
		t.Logf("FMS altitude: %.0f ft", bds40.FmsAltitudeFt)
	}
	if bds40.BaroPressureValid {
		t.Logf("Barometric pressure: %.1f mb", bds40.BaroPressureMb)
	}
}

func TestBDS50_TrackAndTurnReport(t *testing.T) {
	data, err := hex.DecodeString("9F3840B7E90200")
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(5, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 5,0: %v", err)
	}

	bds50, ok := decoded.(BDS50Decoded)
	if !ok {
		t.Fatalf("expected BDS50Decoded, got %T", decoded)
	}
	if bds50.RollAngleValid {
		t.Logf("Roll angle: %.2f degrees", bds50.RollAngleDeg)
		if bds50.RollAngleDeg < -90 || bds50.RollAngleDeg > 90 {
			t.Errorf("Roll angle out of range: %.2f", bds50.RollAngleDeg)
		}
	}
	if bds50.GroundSpeedValid {
		t.Logf("Ground speed: %.0f knots", bds50.GroundSpeedKt)
		if bds50.GroundSpeedKt < 0 || bds50.GroundSpeedKt > 1000 {
			t.Errorf("Ground speed out of range: %.0f", bds50.GroundSpeedKt)
		}
	}
	if bds50.TrueTrackAngleValid {
		t.Logf("True track angle: %.2f degrees", bds50.TrueTrackAngleDeg)
		if bds50.TrueTrackAngleDeg < 0 || bds50.TrueTrackAngleDeg > 360 {
			t.Errorf("Track angle out of range: %.2f", bds50.TrueTrackAngleDeg)
		}
	}
}

func TestBDS60_HeadingAndSpeedReport(t *testing.T) {
	data, err := hex.DecodeString("990940B80C0100")
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(6, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 6,0: %v", err)
	}

	bds60, ok := decoded.(BDS60Decoded)
	if !ok {
		t.Fatalf("expected BDS60Decoded, got %T", decoded)
	}
	if bds60.MagneticHeadingValid {
		t.Logf("Magnetic heading: %.2f degrees", bds60.MagneticHeadingDeg)
		if bds60.MagneticHeadingDeg < 0 || bds60.MagneticHeadingDeg > 360 {
			t.Errorf("Magnetic heading out of range: %.2f", bds60.MagneticHeadingDeg)
		}
	}
	if bds60.IndicatedAirspeedValid {
		t.Logf("Indicated airspeed: %.0f knots", bds60.IndicatedAirspeedKt)
		if bds60.IndicatedAirspeedKt < 0 || bds60.IndicatedAirspeedKt > 600 {
			t.Errorf("IAS out of range: %.0f", bds60.IndicatedAirspeedKt)
		}
	}
	if bds60.MachNumberValid {
		t.Logf("Mach number: %.3f", bds60.MachNumber)
		if bds60.MachNumber < 0 || bds60.MachNumber > 2.0 {
			t.Errorf("Mach number out of range: %.3f", bds60.MachNumber)
		}
	}
}

func TestBDS44_MeteorologicalReport(t *testing.T) {
	data, err := hex.DecodeString("8D40000000000000")
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}
	data = data[:7]

	decoded, err := Decode(4, 4, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 4,4: %v", err)
	}

	bds44, ok := decoded.(BDS44Decoded)
	if !ok {
		t.Fatalf("expected BDS44Decoded, got %T", decoded)
	}
	if bds44.WindValid {
		t.Logf("Wind speed: %d knots, direction: %.1f degrees", bds44.WindSpeedKt, bds44.WindDirectionDeg)
	}
	if bds44.StaticAirTemperatureValid {
		t.Logf("Static air temperature: %.2f°C", bds44.StaticAirTemperatureC)
		if bds44.StaticAirTemperatureC < -80 || bds44.StaticAirTemperatureC > 60 {
			t.Errorf("Temperature out of range: %.2f°C", bds44.StaticAirTemperatureC)
		}
	}
	if bds44.TurbulenceValid {
		t.Logf("Turbulence: %s", bds44.Turbulence)
	}
}

func TestBDS10_DataLinkCapability(t *testing.T) {
	data, err := hex.DecodeString("0000F21F820700")
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	decoded, err := Decode(1, 0, data)
	if err != nil {
		t.Fatalf("Failed to decode BDS 1,0: %v", err)
	}

	bds10, ok := decoded.(BDS10Decoded)
	if !ok {
		t.Fatalf("expected BDS10Decoded, got %T", decoded)
	}
	t.Logf("BDS 2,0 capable: %v", bds10.BDS20Cap)
	t.Logf("BDS 4,0 capable: %v", bds10.BDS40Cap)
	t.Logf("BDS 5,0 capable: %v", bds10.BDS50Cap)
	t.Logf("BDS 6,0 capable: %v", bds10.BDS60Cap)
}

func TestUnknownBDSCode(t *testing.T) {
	data, err := hex.DecodeString("1234567890ABCD")
	if err != nil {
		t.Fatalf("Failed to decode hex: %v", err)
	}

	_, err = Decode(9, 9, data) // BDS 9,9 doesn't exist
	if err == nil {
		t.Fatal("Expected error for unknown BDS code, got nil")
	}
	t.Logf("Got expected error: %v", err)
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