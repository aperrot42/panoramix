package bds

import (
	"strings"
)

// BDS20Decoder decodes BDS 2,0 - Aircraft Identification
type BDS20Decoder struct{}

func (d *BDS20Decoder) BDSCode() (uint8, uint8) {
	return 2, 0
}

func (d *BDS20Decoder) Decode(data []byte) (any, error) {
	decoded := BDS20Decoded{}

	// Bits 9-56: Aircraft identification (8 characters, 6 bits each)
	chars := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bitStart := 8 + i*6
		charCode := ExtractBits(data, bitStart, 6)
		chars[i] = decodeAISCharacter(charCode)
	}

	callsign := strings.TrimRight(string(chars), " ")
	decoded.Callsign = callsign
	decoded.AircraftIdentification = callsign

	return decoded, nil
}

// decodeAISCharacter decodes a 6-bit character code according to ICAO Annex 10
func decodeAISCharacter(code uint32) byte {
	switch {
	case code >= 1 && code <= 26:
		return byte('A' + code - 1)
	case code >= 48 && code <= 57:
		return byte('0' + code - 48)
	case code == 32:
		return ' '
	case code == 45:
		return '-'
	default:
		return ' ' // undefined or padding
	}
}
