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

// Gray2Binary converts Gray code to binary
func Gray2Binary(gray uint32) uint32 {
	binary := gray
	for gray >>= 1; gray != 0; gray >>= 1 {
		binary ^= gray
	}
	return binary
}

// NauticalMilesToMeters converts nautical miles to meters
func NauticalMilesToMeters(nm float64) float64 {
	return nm * 1852.0
}

// FeetToMeters converts feet to meters
func FeetToMeters(ft float64) float64 {
	return ft * 0.3048
}

// KnotsToMetersPerSecond converts knots to m/s
func KnotsToMetersPerSecond(knots float64) float64 {
	return knots * 0.514444
}