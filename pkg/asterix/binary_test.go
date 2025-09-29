package asterix

import (
	"testing"
)

func TestReadUint24BE(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		offset   int
		expected uint32
	}{
		{"Simple case", []byte{0x12, 0x34, 0x56}, 0, 0x123456},
		{"With offset", []byte{0x00, 0x12, 0x34, 0x56}, 1, 0x123456},
		{"Zero value", []byte{0x00, 0x00, 0x00}, 0, 0x000000},
		{"Max value", []byte{0xFF, 0xFF, 0xFF}, 0, 0xFFFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadUint24BE(tt.data, tt.offset)
			if result != tt.expected {
				t.Errorf("Expected %06x, got %06x", tt.expected, result)
			}
		})
	}
}

func TestReadInt24BE(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		offset   int
		expected int32
	}{
		{"Positive value", []byte{0x12, 0x34, 0x56}, 0, 0x123456},
		{"Negative value", []byte{0x80, 0x00, 0x00}, 0, -8388608}, // Sign extend 0x800000
		{"Max positive", []byte{0x7F, 0xFF, 0xFF}, 0, 8388607},
		{"Zero", []byte{0x00, 0x00, 0x00}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadInt24BE(tt.data, tt.offset)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestTwosComplement24(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected int32
	}{
		{"Positive value", 0x123456, 0x123456},
		{"Negative value", 0x800000, -8388608},
		{"Max positive", 0x7FFFFF, 8388607},
		{"Zero", 0x000000, 0},
		{"Edge case", 0x800001, -8388607},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TwosComplement24(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestSignExtend14To16(t *testing.T) {
	tests := []struct {
		name     string
		input    uint16
		expected int16
	}{
		{"Positive value", 0x1234, 0x1234},
		{"Negative value", 0x2000, -8192}, // Bit 14 set = negative
		{"Max positive", 0x1FFF, 8191},    // Bit 14 clear = positive
		{"Zero", 0x0000, 0},
		{"Edge negative", 0x2001, -8191}, // Smallest negative + 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SignExtend14To16(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExtractBitsFromByte(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		bitPos   int
		numBits  int
		expected uint8
	}{
		{"Extract 1 bit from MSB", 0b10101010, 8, 1, 1},
		{"Extract 1 bit from LSB", 0b10101010, 1, 1, 0},
		{"Extract 3 bits", 0b10101010, 8, 3, 0b101},
		{"Extract middle bits", 0b10101010, 6, 2, 0b10}, // Bits 6-5 of 10101010 = 10
		{"Extract all bits", 0b10101010, 8, 8, 0b10101010},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractBitsFromByte(tt.b, tt.bitPos, tt.numBits)
			if result != tt.expected {
				t.Errorf("Expected %08b, got %08b", tt.expected, result)
			}
		})
	}
}

func TestReadBit(t *testing.T) {
	b := byte(0b10101010) // 170 decimal

	tests := []struct {
		bitPos   int
		expected bool
	}{
		{8, true}, // MSB
		{7, false},
		{6, true},
		{5, false},
		{4, true},
		{3, false},
		{2, true},
		{1, false}, // LSB
	}

	for _, tt := range tests {
		t.Run("Bit position "+string(rune(tt.bitPos+'0')), func(t *testing.T) {
			result := ReadBit(b, tt.bitPos)
			if result != tt.expected {
				t.Errorf("Bit %d: expected %t, got %t", tt.bitPos, tt.expected, result)
			}
		})
	}
}

// Test compatibility with legacy functions
func TestLegacyCompatibility(t *testing.T) {
	testData := []byte{0x12, 0x34, 0x56}

	// Test that legacy functions still work
	result24 := int24(testData)
	expected24 := int32(0x123456)
	if result24 != expected24 {
		t.Errorf("int24(): expected %d, got %d", expected24, result24)
	}

	resultU24 := uint24(testData)
	expectedU24 := uint32(0x123456)
	if resultU24 != expectedU24 {
		t.Errorf("uint24(): expected %d, got %d", expectedU24, resultU24)
	}

	// Test two's complement function
	negativeTwos := twosComplement24(0x800000)
	expectedTwos := int32(-8388608)
	if negativeTwos != expectedTwos {
		t.Errorf("twosComplement24(): expected %d, got %d", expectedTwos, negativeTwos)
	}
}
