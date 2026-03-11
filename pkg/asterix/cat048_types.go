package asterix

// BDSRegisterData represents I048/250 Mode S MB Data
type BDSRegisterData struct {
	Repetition uint8                  `json:"repetition"`
	Registers  map[string]BDSRegister `json:"registers"`
}

// BDSRegister represents a single BDS register entry
type BDSRegister struct {
	BDSCode    byte   `json:"bds_code"`         // raw register address byte (e.g. 0x50 for BDS 5,0)
	BDSDataRaw string `json:"bds_data_raw"`     // hex-encoded 7-byte payload
	Decoded    any    `json:"decoded,omitempty"` // decoded register data, if a BDSDecoder was provided
	Error      string `json:"error,omitempty"`   // decoding error message, if any
}
