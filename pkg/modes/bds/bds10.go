package bds

// BDS10Decoder decodes BDS 1,0 - Data Link Capability Report
type BDS10Decoder struct{}

func (d *BDS10Decoder) BDSCode() (uint8, uint8) {
	return 1, 0
}

func (d *BDS10Decoder) Decode(data []byte) (interface{}, error) {
	decoded := BDS10Decoded{}

	// Bits 1-16: Reserved for ACAS
	decoded.AcasReserved = ExtractBits(data, 0, 16)

	// Bit 17: BDS 1,0 bit 16 = 1
	decoded.BDS10Cf = ExtractBits(data, 16, 1) == 1

	// Bit 18: BDS 1,7 capability
	decoded.BDS17Cap = ExtractBits(data, 17, 1) == 1

	// Bits 19-23: Reserved

	// Bit 24: Comm-B broadcast message 1 capability
	decoded.CommBBroadcast1Cap = ExtractBits(data, 23, 1) == 1

	// Bit 25: BDS 2,0 capability
	decoded.BDS20Cap = ExtractBits(data, 24, 1) == 1

	// Bit 26: BDS 2,1 capability
	decoded.BDS21Cap = ExtractBits(data, 25, 1) == 1

	// Bit 27: BDS 4,0 capability
	decoded.BDS40Cap = ExtractBits(data, 26, 1) == 1

	// Bit 28: BDS 4,1 capability
	decoded.BDS41Cap = ExtractBits(data, 27, 1) == 1

	// Bit 29: BDS 4,2 capability
	decoded.BDS42Cap = ExtractBits(data, 28, 1) == 1

	// Bit 30: BDS 4,3 capability
	decoded.BDS43Cap = ExtractBits(data, 29, 1) == 1

	// Bit 31: BDS 4,4 capability
	decoded.BDS44Cap = ExtractBits(data, 30, 1) == 1

	// Bit 32: BDS 4,5 capability
	decoded.BDS45Cap = ExtractBits(data, 31, 1) == 1

	// Bit 33: BDS 4,8 capability
	decoded.BDS48Cap = ExtractBits(data, 32, 1) == 1

	// Bit 34: BDS 5,0 capability
	decoded.BDS50Cap = ExtractBits(data, 33, 1) == 1

	// Bit 35: BDS 5,1 capability
	decoded.BDS51Cap = ExtractBits(data, 34, 1) == 1

	// Bit 36: BDS 5,2 capability
	decoded.BDS52Cap = ExtractBits(data, 35, 1) == 1

	// Bit 37: BDS 5,3 capability
	decoded.BDS53Cap = ExtractBits(data, 36, 1) == 1

	// Bit 38: BDS 5,4 capability
	decoded.BDS54Cap = ExtractBits(data, 37, 1) == 1

	// Bit 39: BDS 5,5 capability
	decoded.BDS55Cap = ExtractBits(data, 38, 1) == 1

	// Bit 40: BDS 5,6 capability
	decoded.BDS56Cap = ExtractBits(data, 39, 1) == 1

	// Bit 41: BDS 5,F capability
	decoded.BDS5FCap = ExtractBits(data, 40, 1) == 1

	// Bit 42: BDS 6,0 capability
	decoded.BDS60Cap = ExtractBits(data, 41, 1) == 1

	// Bits 43-56: Reserved

	return decoded, nil
}
