// Package bds provides decoding for Mode S BDS (Comm-B Data Selector) registers
// as specified in ED-73F MOPS for Secondary Surveillance Radar Mode S Transponders
package bds

import (
	"fmt"
)

// Register represents a decoded BDS register
type Register struct {
	BDS1   uint8                  // First digit of BDS code
	BDS2   uint8                  // Second digit of BDS code
	Data   []byte                 // Raw 56-bit data
	Fields map[string]interface{} // Decoded fields
}

// BDS10Register represents BDS 1,0 - Data Link Capability Report
type BDS10Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS10Decoded `json:"decoded"`
}

type BDS10Decoded struct {
	// Add fields as they are defined in BDS10Decoder
}

// BDS17Register represents BDS 1,7 - Common Usage GICB Capability Report
type BDS17Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS17Decoded `json:"decoded"`
}

type BDS17Decoded struct {
	// Add fields as they are defined in BDS17Decoder
}

// BDS20Register represents BDS 2,0 - Aircraft Identification
type BDS20Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS20Decoded `json:"decoded"`
}

type BDS20Decoded struct {
	Callsign               string `json:"callsign"`
	AircraftIdentification string `json:"aircraft_identification"`
}

// BDS30Register represents BDS 3,0 - ACAS Active Resolution Advisory
type BDS30Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS30Decoded `json:"decoded"`
}

type BDS30Decoded struct {
	// Add fields as they are defined in BDS30Decoder
}

// BDS40Register represents BDS 4,0 - Selected Vertical Intention
type BDS40Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS40Decoded `json:"decoded"`
}

type BDS40Decoded struct {
	// Add fields as they are defined in BDS40Decoder
}

// BDS44Register represents BDS 4,4 - Meteorological Routine Air Report
type BDS44Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS44Decoded `json:"decoded"`
}

type BDS44Decoded struct {
	// Add fields as they are defined in BDS44Decoder
}

// BDS50Register represents BDS 5,0 - Track and Turn Report
type BDS50Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS50Decoded `json:"decoded"`
}

type BDS50Decoded struct {
	RollAngleDeg           float64 `json:"roll_angle_deg,omitempty"`
	RollAngleValid         bool    `json:"roll_angle_valid"`
	TrueTrackAngleDeg      float64 `json:"true_track_angle_deg,omitempty"`
	TrueTrackAngleValid    bool    `json:"true_track_angle_valid"`
	GroundSpeedKt          float64 `json:"ground_speed_kt,omitempty"`
	GroundSpeedValid       bool    `json:"ground_speed_valid"`
	TrackAngleRateDegS     float64 `json:"track_angle_rate_deg_s,omitempty"`
	TrackAngleRateValid    bool    `json:"track_angle_rate_valid"`
	TrueAirspeedKt         float64 `json:"true_airspeed_kt,omitempty"`
	TrueAirspeedValid      bool    `json:"true_airspeed_valid"`
}

// BDS60Register represents BDS 6,0 - Heading and Speed Report
type BDS60Register struct {
	BDS1        uint8  `json:"bds1"`
	BDS2        uint8  `json:"bds2"`
	BDSDataRaw  string `json:"bds_data_raw"`
	Decoded     BDS60Decoded `json:"decoded"`
}

type BDS60Decoded struct {
	MagneticHeadingDeg            float64 `json:"magnetic_heading_deg,omitempty"`
	MagneticHeadingValid          bool    `json:"magnetic_heading_valid"`
	IndicatedAirspeedKt           float64 `json:"indicated_airspeed_kt,omitempty"`
	IndicatedAirspeedValid        bool    `json:"indicated_airspeed_valid"`
	MachNumber                    float64 `json:"mach_number,omitempty"`
	MachNumberValid               bool    `json:"mach_number_valid"`
	BaroAltitudeRateFpm           float64 `json:"baro_altitude_rate_fpm,omitempty"`
	BaroAltitudeRateValid         bool    `json:"baro_altitude_rate_valid"`
	InertialVerticalVelocityFpm   float64 `json:"inertial_vertical_velocity_fpm,omitempty"`
	InertialVerticalVelocityValid bool    `json:"inertial_vertical_velocity_valid"`
}

// Decoder is the interface for BDS register decoders
type Decoder interface {
	Decode(data []byte) (map[string]interface{}, error)
	BDSCode() (uint8, uint8) // Returns BDS1, BDS2
}

// Registry of BDS decoders
var decoders = map[uint8]Decoder{
	0x10: &BDS10Decoder{}, // BDS 1,0 - Data Link Capability Report
	0x17: &BDS17Decoder{}, // BDS 1,7 - Common Usage GICB Capability Report
	0x20: &BDS20Decoder{}, // BDS 2,0 - Aircraft Identification
	0x30: &BDS30Decoder{}, // BDS 3,0 - ACAS Active Resolution Advisory
	0x40: &BDS40Decoder{}, // BDS 4,0 - Selected Vertical Intention
	0x44: &BDS44Decoder{}, // BDS 4,4 - Meteorological Routine Air Report
	0x50: &BDS50Decoder{}, // BDS 5,0 - Track and Turn Report
	0x60: &BDS60Decoder{}, // BDS 6,0 - Heading and Speed Report
}

// Decode decodes a BDS register based on its code
func Decode(bds1, bds2 uint8, data []byte) (*Register, error) {
	if len(data) < 7 {
		return nil, fmt.Errorf("BDS data too short: %d bytes, need 7", len(data))
	}

	reg := &Register{
		BDS1:   bds1,
		BDS2:   bds2,
		Data:   data[:7],
		Fields: make(map[string]interface{}),
	}

	// Combine BDS1 and BDS2 into single code for lookup
	bdsCode := (bds1 << 4) | bds2

	decoder, ok := decoders[bdsCode]
	if !ok {
		// Unknown BDS code, return raw data
		reg.Fields["raw"] = fmt.Sprintf("%014X", data[:7])
		reg.Fields["unknown"] = true
		return reg, nil
	}

	fields, err := decoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("BDS %X,%X decode error: %w", bds1, bds2, err)
	}

	reg.Fields = fields
	return reg, nil
}

// Helper functions for bit extraction

// ExtractBits extracts bits from a byte array
func ExtractBits(data []byte, startBit, numBits int) uint32 {
	if numBits > 32 {
		return 0
	}

	result := uint32(0)
	for i := 0; i < numBits; i++ {
		bitPos := startBit + i
		bytePos := bitPos / 8
		bitInByte := 7 - (bitPos % 8)

		if bytePos < len(data) {
			if (data[bytePos] & (1 << bitInByte)) != 0 {
				result |= 1 << (numBits - 1 - i)
			}
		}
	}
	return result
}

// ExtractSignedBits extracts signed bits (two's complement)
func ExtractSignedBits(data []byte, startBit, numBits int) int32 {
	value := ExtractBits(data, startBit, numBits)

	// Check if sign bit is set
	if value&(1<<(numBits-1)) != 0 {
		// Extend sign bit
		mask := ^((1 << numBits) - 1)
		return int32(value | uint32(mask))
	}
	return int32(value)
}
