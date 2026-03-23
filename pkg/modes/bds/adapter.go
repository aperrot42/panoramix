package bds

import "fmt"

// DecodedRegisters holds the decoded BDS registers from a message's I048/250 data.
// Each field is set if the corresponding BDS register was present and decoded successfully.
type DecodedRegisters struct {
	BDS10 *BDS10Decoded `json:"bds_1_0,omitempty"`
	BDS17 *BDS17Decoded `json:"bds_1_7,omitempty"`
	BDS20 *BDS20Decoded `json:"bds_2_0,omitempty"`
	BDS30 *BDS30Decoded `json:"bds_3_0,omitempty"`
	BDS40 *BDS40Decoded `json:"bds_4_0,omitempty"`
	BDS44 *BDS44Decoded `json:"bds_4_4,omitempty"`
	BDS50 *BDS50Decoded `json:"bds_5_0,omitempty"`
	BDS60 *BDS60Decoded `json:"bds_6_0,omitempty"`
}

// DecodeRegister decodes a single raw BDS register and sets the corresponding
// typed field on dr. Unknown or malformed registers are silently skipped.
func (dr *DecodedRegisters) DecodeRegister(bds1, bds2 uint8, data []byte) error {
	if len(data) < 7 {
		return fmt.Errorf("BDS data too short: %d bytes, need 7", len(data))
	}

	bdsCode := (bds1 << 4) | bds2
	decoder, ok := decoders[bdsCode]
	if !ok {
		return fmt.Errorf("unknown BDS code %d,%d", bds1, bds2)
	}

	decoded, err := decoder.Decode(data)
	if err != nil {
		return err
	}

	switch v := decoded.(type) {
	case BDS10Decoded:
		dr.BDS10 = &v
	case BDS17Decoded:
		dr.BDS17 = &v
	case BDS20Decoded:
		dr.BDS20 = &v
	case BDS30Decoded:
		dr.BDS30 = &v
	case BDS40Decoded:
		dr.BDS40 = &v
	case BDS44Decoded:
		dr.BDS44 = &v
	case BDS50Decoded:
		dr.BDS50 = &v
	case BDS60Decoded:
		dr.BDS60 = &v
	}

	return nil
}

// RawRegister is a BDS register address and its raw 7-byte payload.
type RawRegister struct {
	Code byte
	Data []byte
}

// DecodeAll decodes all raw registers and sets the corresponding typed fields.
func (dr *DecodedRegisters) DecodeAll(registers []RawRegister) {
	for _, reg := range registers {
		bds1 := (reg.Code >> 4) & 0x0F
		bds2 := reg.Code & 0x0F
		dr.DecodeRegister(bds1, bds2, reg.Data)
	}
}
