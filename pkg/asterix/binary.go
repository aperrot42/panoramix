// Package asterix binary utilities for ASTERIX message decoding
package asterix

import "fmt"

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

// FormatTransponderCodeFromBytes formats transponder code directly from 2-byte data
func FormatTransponderCodeFromBytes(data []byte, offset int) string {
	code := (uint16(data[offset])<<8 | uint16(data[offset+1])) & 0x1FFF
	return fmt.Sprintf("%04o", code)
}

// Common ASTERIX constants
const (
	// ASTERIX specific masks
	Mask13Bits = 0x1FFF // Transponder codes
	Mask14Bits = 0x3FFF // Flight levels
)
