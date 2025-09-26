package bds

// BDS30Decoder decodes BDS 3,0 - ACAS Active Resolution Advisory
type BDS30Decoder struct{}

func (d *BDS30Decoder) BDSCode() (uint8, uint8) {
	return 3, 0
}

func (d *BDS30Decoder) Decode(data []byte) (interface{}, error) {
	decoded := BDS30Decoded{}

	// This is an ACAS resolution advisory message
	// The structure follows RTCA DO-185B

	// Bit 1: Threat type indicator
	decoded.ThreatTypeIndicator = ExtractBits(data, 0, 1) == 1

	// Bits 2-8: Active Resolution Advisory (ARA)
	ara := ExtractBits(data, 1, 7)
	decoded.ActiveResolutionAdvisory = ara

	// Decode individual ARA bits
	decoded.AraCorrectiveRa = (ara & 0x40) != 0
	decoded.AraDownwardSense = (ara & 0x20) != 0
	decoded.AraIncreasedRate = (ara & 0x10) != 0
	decoded.AraSenseReversal = (ara & 0x08) != 0
	decoded.AraAltitudeCrossing = (ara & 0x04) != 0
	decoded.AraPositiveRa = (ara & 0x02) != 0
	decoded.AraVerticalSpeedLimit = (ara & 0x01) != 0

	// Bits 9-13: Resolution Advisory Complement (RAC)
	decoded.ResolutionAdvisoryComplement = ExtractBits(data, 8, 5)

	// Bit 14: RA terminated
	decoded.RaTerminated = ExtractBits(data, 13, 1) == 1

	// Bit 15: Multiple threat encounter
	decoded.MultipleThreatEncounter = ExtractBits(data, 14, 1) == 1

	// Bits 16-19: Threat type
	threatType := ExtractBits(data, 15, 4)
	decoded.ThreatTypeRaw = threatType
	switch threatType {
	case 0:
		decoded.ThreatType = "No threat"
	case 1:
		decoded.ThreatType = "Traffic advisory"
	case 2:
		decoded.ThreatType = "Resolution advisory"
	default:
		decoded.ThreatType = "Reserved"
	}

	// Bits 20-26: Threat identity data, bits 26-20 of Mode S address
	decoded.ThreatIdentityData = ExtractBits(data, 19, 7)

	// Bits 27-30: Reserved for ACAS III

	// Bits 31-39: Threat identity data, bits 19-11 of Mode S address
	decoded.ThreatIdentityDataMid = ExtractBits(data, 30, 9)

	// Bits 40-50: Threat identity data, bits 10-0 of Mode S address
	decoded.ThreatIdentityDataLow = ExtractBits(data, 39, 11)

	// Combine threat identity to get full Mode S address if available
	if decoded.ThreatTypeIndicator {
		threatAddr := (decoded.ThreatIdentityData << 17) |
			(decoded.ThreatIdentityDataMid << 8) |
			decoded.ThreatIdentityDataLow
		decoded.ThreatModeSAddress = threatAddr
	}

	// Bits 51-56: Reserved

	return decoded, nil
}
