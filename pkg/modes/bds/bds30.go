package bds

// BDS30Decoder decodes BDS 3,0 - ACAS Active Resolution Advisory
type BDS30Decoder struct{}

func (d *BDS30Decoder) BDSCode() (uint8, uint8) {
	return 3, 0
}

func (d *BDS30Decoder) Decode(data []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// This is an ACAS resolution advisory message
	// The structure follows RTCA DO-185B

	// Bit 1: Threat type indicator
	fields["threat_type_indicator"] = ExtractBits(data, 0, 1) == 1

	// Bits 2-8: Active Resolution Advisory (ARA)
	ara := ExtractBits(data, 1, 7)
	fields["active_resolution_advisory"] = ara

	// Decode individual ARA bits
	fields["ara_corrective_ra"] = (ara & 0x40) != 0
	fields["ara_downward_sense"] = (ara & 0x20) != 0
	fields["ara_increased_rate"] = (ara & 0x10) != 0
	fields["ara_sense_reversal"] = (ara & 0x08) != 0
	fields["ara_altitude_crossing"] = (ara & 0x04) != 0
	fields["ara_positive_ra"] = (ara & 0x02) != 0
	fields["ara_vertical_speed_limit"] = (ara & 0x01) != 0

	// Bits 9-13: Resolution Advisory Complement (RAC)
	rac := ExtractBits(data, 8, 5)
	fields["resolution_advisory_complement"] = rac

	// Bit 14: RA terminated
	fields["ra_terminated"] = ExtractBits(data, 13, 1) == 1

	// Bit 15: Multiple threat encounter
	fields["multiple_threat_encounter"] = ExtractBits(data, 14, 1) == 1

	// Bits 16-19: Threat type
	threatType := ExtractBits(data, 15, 4)
	fields["threat_type_raw"] = threatType
	switch threatType {
	case 0:
		fields["threat_type"] = "No threat"
	case 1:
		fields["threat_type"] = "Traffic advisory"
	case 2:
		fields["threat_type"] = "Resolution advisory"
	default:
		fields["threat_type"] = "Reserved"
	}

	// Bits 20-26: Threat identity data, bits 26-20 of Mode S address
	fields["threat_identity_data"] = ExtractBits(data, 19, 7)

	// Bits 27-30: Reserved for ACAS III

	// Bits 31-39: Threat identity data, bits 19-11 of Mode S address
	fields["threat_identity_data_mid"] = ExtractBits(data, 30, 9)

	// Bits 40-50: Threat identity data, bits 10-0 of Mode S address
	fields["threat_identity_data_low"] = ExtractBits(data, 39, 11)

	// Combine threat identity to get full Mode S address if available
	if fields["threat_type_indicator"].(bool) {
		threatAddr := (fields["threat_identity_data"].(uint32) << 17) |
			(fields["threat_identity_data_mid"].(uint32) << 8) |
			fields["threat_identity_data_low"].(uint32)
		fields["threat_mode_s_address"] = threatAddr
	}

	// Bits 51-56: Reserved

	return fields, nil
}
