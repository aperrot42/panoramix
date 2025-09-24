package asterix

import (
	"fmt"
)

func extractFSPEC(data []byte) (fspec []byte, rest []byte, err error) {
	for i, b := range data {
		fspec = append(fspec, b)
		if b&0x01 == 0 { // FX == 0 → last FSPEC byte
			return fspec, data[i+1:], nil
		}
	}
	return nil, nil, fmt.Errorf("FSPEC not properly terminated")
}

// WalkFSPEC processes an FSPEC (Field Specification) and decodes the corresponding data fields
func WalkFSPEC(fspec []byte, payload []byte, fieldTable map[int]DataItem) (map[string]any, int, error) {
	decodedFields := make(map[string]any)
	payloadCursor := payload
	fieldReferenceNumber := 1
	alreadyDecodedFields := make(map[string]bool)

	// Process each FSPEC byte
	for _, fspecByte := range fspec {
		// Check bits 7 down to 1 (bit 0 is FX - Field Extension)
		for bitPosition := 7; bitPosition >= 1; bitPosition-- {
			if isBitSet(fspecByte, bitPosition) {
				// Skip if we already decoded a field with this name (handles duplicates)
				if fieldEntry, exists := fieldTable[fieldReferenceNumber]; exists {
					if alreadyDecodedFields[fieldEntry.Name] {
						fieldReferenceNumber++
						continue
					}
				}

				// Decode the field if it exists in our field table
				if err := decodeField(fieldReferenceNumber, fieldTable, &payloadCursor, decodedFields, alreadyDecodedFields); err != nil {
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
func decodeField(frn int, fieldTable map[int]DataItem, payloadCursor *[]byte, decodedFields map[string]any, alreadyDecoded map[string]bool) error {
	fieldEntry, exists := fieldTable[frn]
	if !exists {
		// Unknown field - stop processing (this is normal for optional fields)
		return nil
	}

	// Decode the field
	value, bytesConsumed, err := fieldEntry.Decoder(*payloadCursor)
	if err != nil {
		return fmt.Errorf("failed to decode field %s (FRN %d): %w", fieldEntry.Name, frn, err)
	}

	// Store the decoded value
	decodedFields[fieldEntry.Name] = value
	*payloadCursor = (*payloadCursor)[bytesConsumed:]

	// Mark this field name as decoded to prevent duplicates
	alreadyDecoded[fieldEntry.Name] = true

	return nil
}
