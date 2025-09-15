// Package asterix binary utilities for ASTERIX message decoding
package asterix

import (
	"encoding/binary"
	"fmt"
)

// ReadUint16BE reads a 16-bit big-endian unsigned integer with optional bit mask
func ReadUint16BE(data []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(data[offset:])
}

// ReadUint16BEMasked reads a 16-bit big-endian unsigned integer and applies a bit mask
func ReadUint16BEMasked(data []byte, offset int, mask uint16) uint16 {
	return binary.BigEndian.Uint16(data[offset:]) & mask
}

// ReadUint24BE reads a 24-bit big-endian unsigned integer (returned as uint32)
func ReadUint24BE(data []byte, offset int) uint32 {
	return uint32(data[offset])<<16 | uint32(data[offset+1])<<8 | uint32(data[offset+2])
}

// ReadInt24BE reads a 24-bit big-endian signed integer with two's complement conversion
func ReadInt24BE(data []byte, offset int) int32 {
	raw := int32(data[offset])<<16 | int32(data[offset+1])<<8 | int32(data[offset+2])
	return TwosComplement24(raw)
}

// TwosComplement24 converts a 24-bit value to proper two's complement int32
func TwosComplement24(x int32) int32 {
	if x&0x800000 != 0 {
		return x - 0x1000000
	}
	return x
}

// SignExtend14To16 performs 14-bit to 16-bit sign extension
func SignExtend14To16(value uint16) int16 {
	if value&0x2000 != 0 { // if bit 14 is set (negative)
		return int16(value | 0xC000) // extend sign
	}
	return int16(value)
}

// ExtractBitsFromByte extracts a bit field from a single byte
// bitPos: bit position (1-8, where 8 is MSB, 1 is LSB) - ASTERIX convention
// numBits: number of consecutive bits to extract
func ExtractBitsFromByte(b byte, bitPos, numBits int) uint8 {
	shift := bitPos - numBits
	mask := uint8((1 << numBits) - 1)
	return (b >> shift) & mask
}

// ReadBit reads a single bit from a byte (ASTERIX bit numbering: 8=MSB, 1=LSB)
func ReadBit(b byte, bitPos int) bool {
	return (b>>(bitPos-1))&0x01 == 1
}

// FormatTransponderCode formats a 13-bit transponder code as octal string
func FormatTransponderCode(code uint16) string {
	return fmt.Sprintf("%04o", code&Mask13Bits)
}

// FormatTransponderCodeFromBytes formats transponder code directly from 2-byte data
func FormatTransponderCodeFromBytes(data []byte, offset int) string {
	code := ReadUint16BEMasked(data, offset, 0x1FFF)
	return FormatTransponderCode(code)
}

// ASTERIX bit field extraction helpers for common patterns

// Extract2Bits extracts a 2-bit field from a byte
func Extract2Bits(b byte, bitPos int) uint8 {
	return ExtractBitsFromByte(b, bitPos, 2)
}

// Extract3Bits extracts a 3-bit field from a byte  
func Extract3Bits(b byte, bitPos int) uint8 {
	return ExtractBitsFromByte(b, bitPos, 3)
}

// Extract4Bits extracts a 4-bit field from a byte
func Extract4Bits(b byte, bitPos int) uint8 {
	return ExtractBitsFromByte(b, bitPos, 4)
}

// Extract5Bits extracts a 5-bit field from a byte
func Extract5Bits(b byte, bitPos int) uint8 {
	return ExtractBitsFromByte(b, bitPos, 5)
}

// ReadFlightLevel decodes ASTERIX flight level from 2 bytes (I048/090)
// Returns flight level with V/G flags
func ReadFlightLevel(data []byte, offset int) (flightLevel float64, validated bool, garbled bool) {
	b1 := data[offset]
	b2 := data[offset+1]

	v := ReadBit(b1, 8) // bit 16 - V flag
	g := ReadBit(b1, 7) // bit 15 - G flag

	// Bits 14–1: 14-bit signed integer
	raw := (uint16(b1&0x3F) << 8) | uint16(b2) // clear top 2 bits (V & G)
	signedVal := SignExtend14To16(raw)

	return float64(signedVal) * 0.25, !v, g // 1/4 FL resolution, V=0 means validated
}

// ReadMode3ACode decodes ASTERIX Mode 3/A code with flags (I048/070)
func ReadMode3ACode(data []byte, offset int) (code string, validated bool, garbled bool, local bool) {
	b1 := data[offset]
	
	v := !ReadBit(b1, 8) // V bit - 0=validated
	g := ReadBit(b1, 7)  // G bit - garbled
	l := ReadBit(b1, 6)  // L bit - local
	
	codeVal := ReadUint16BEMasked(data, offset, 0x1FFF)
	
	return FormatTransponderCode(codeVal), v, g, l
}

// Common ASTERIX constants
const (
	// Bit masks for common field sizes
	Mask1Bit  = 0x01
	Mask2Bits = 0x03
	Mask3Bits = 0x07
	Mask4Bits = 0x0F
	Mask5Bits = 0x1F
	Mask6Bits = 0x3F
	Mask7Bits = 0x7F
	
	// ASTERIX specific masks
	Mask12Bits = 0x0FFF // Mode codes
	Mask13Bits = 0x1FFF // Transponder codes 
	Mask14Bits = 0x3FFF // Flight levels
)