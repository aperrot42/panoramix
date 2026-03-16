package asterix

import (
	"fmt"
)

func extractFSPEC(data []byte) (fspec []byte, rest []byte, err error) {
	// FSPEC is typically 1-4 bytes; find the end by scanning for FX==0
	for i, b := range data {
		if b&0x01 == 0 { // FX == 0 → last FSPEC byte
			return data[:i+1], data[i+1:], nil
		}
	}
	return nil, nil, fmt.Errorf("FSPEC not properly terminated")
}

// WalkFSPEC processes an FSPEC (Field Specification) and decodes the corresponding data fields.
// Per the ASTERIX standard, each FRN maps to exactly one unique data item within a UAP,
// so duplicate detection is not needed.
func WalkFSPEC(fspec []byte, payload []byte, fieldTable map[int]DataItem) (map[string]any, int, error) {
	decodedFields := make(map[string]any)
	payloadCursor := payload
	fieldReferenceNumber := 1

	// Process each FSPEC byte
	for _, fspecByte := range fspec {
		// Check bits 7 down to 1 (bit 0 is FX - Field Extension)
		for bitPosition := 7; bitPosition >= 1; bitPosition-- {
			if isBitSet(fspecByte, bitPosition) {
				if err := decodeField(fieldReferenceNumber, fieldTable, &payloadCursor, decodedFields); err != nil {
					return nil, 0, err
				}
			}
			fieldReferenceNumber++
		}
	}

	bytesConsumed := len(payload) - len(payloadCursor)
	return decodedFields, bytesConsumed, nil
}

// isBitSet checks if a specific bit position is set in a byte (1-indexed from right)
func isBitSet(b byte, bitPosition int) bool {
	return b&(1<<bitPosition) != 0
}

// decodeField decodes a single ASTERIX field and updates the results
func decodeField(frn int, fieldTable map[int]DataItem, payloadCursor *[]byte, decodedFields map[string]any) error {
	fieldEntry, exists := fieldTable[frn]
	if !exists {
		return nil
	}

	value, bytesConsumed, err := fieldEntry.Decoder(*payloadCursor)
	if err != nil {
		return fmt.Errorf("failed to decode field %s (FRN %d): %w", fieldEntry.Name, frn, err)
	}

	decodedFields[fieldEntry.Name] = value
	*payloadCursor = (*payloadCursor)[bytesConsumed:]
	return nil
}
