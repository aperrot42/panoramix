package bds

import (
	"fmt"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

// DecodedRegister holds the decoded result for a single BDS register.
type DecodedRegister struct {
	BDSCode byte   `json:"bds_code"`
	Decoded any    `json:"decoded,omitempty"`
	Error   string `json:"error,omitempty"`
}

// DecodeRegisters decodes all BDS registers from raw ASTERIX I048/250 data.
func DecodeRegisters(data *asterix.BDSRegisterData) map[string]DecodedRegister {
	if data == nil {
		return nil
	}
	result := make(map[string]DecodedRegister, len(data.Registers))
	for _, reg := range data.Registers {
		dr := DecodedRegister{BDSCode: reg.BDSCode}
		bds1 := (reg.BDSCode >> 4) & 0x0F
		bds2 := reg.BDSCode & 0x0F
		decoded, err := Decode(bds1, bds2, reg.RawData)
		if err != nil {
			dr.Error = err.Error()
		} else {
			dr.Decoded = decoded
		}
		result[fmt.Sprintf("%d_%d", bds1, bds2)] = dr
	}
	return result
}