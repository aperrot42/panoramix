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
		cat048, ok := ar.Message.Record.(*asterix.Cat048Message)
		if !ok || cat048.BDSRegister == nil {
			return ar, true, nil
		}

		for key, reg := range cat048.BDSRegister.Registers {
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
			cat048.BDSRegister.Registers[key] = reg
		}

		return ar, true, nil
	}
}