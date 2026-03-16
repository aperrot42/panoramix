package pipeline

import (
	"fmt"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

// BDSEnricher decodes raw BDS register data in I048/250.
// It satisfies asterix.BDSDecoder via structural typing.
type BDSEnricher interface {
	DecodeBDS(bdsRegisterAddress byte, data []byte) (any, error)
}

// NewBDSEnrichFilter returns a filter that enriches I048/250 registers
// using the provided BDS decoder.
func NewBDSEnrichFilter(decoder BDSEnricher) func(AsterixResult) (AsterixResult, bool, error) {
	return func(ar AsterixResult) (AsterixResult, bool, error) {
		if ar.Message.Category != 48 {
			return ar, true, nil
		}
		raw, ok := ar.Message.Items["I048/250"]
		if !ok {
			return ar, true, nil
		}
		bdsData, ok := raw.(asterix.BDSRegisterData)
		if !ok {
			return ar, true, nil
		}

		for key, reg := range bdsData.Registers {
			if reg.RawData == nil {
				continue
			}
			decoded, err := decoder.DecodeBDS(reg.BDSCode, reg.RawData)
			if err != nil {
				reg.Error = err.Error()
			} else {
				bds1 := (reg.BDSCode >> 4) & 0x0F
				bds2 := reg.BDSCode & 0x0F
				reg.Decoded = map[string]any{
					fmt.Sprintf("%d_%d", bds1, bds2): decoded,
				}
			}
			bdsData.Registers[key] = reg
		}

		ar.Message.Items["I048/250"] = bdsData
		return ar, true, nil
	}
}
