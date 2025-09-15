package bds

// BDS17Decoder decodes BDS 1,7 - Common Usage GICB Capability Report
type BDS17Decoder struct{}

func (d *BDS17Decoder) BDSCode() (uint8, uint8) {
	return 1, 7
}

func (d *BDS17Decoder) Decode(data []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// Bits 1-5: Reserved
	
	// Bit 6: BDS 0,5 capability
	fields["bds_05_cap"] = ExtractBits(data, 5, 1) == 1

	// Bit 7: BDS 0,6 capability
	fields["bds_06_cap"] = ExtractBits(data, 6, 1) == 1

	// Bit 8: BDS 0,7 capability
	fields["bds_07_cap"] = ExtractBits(data, 7, 1) == 1

	// Bit 9: BDS 0,8 capability
	fields["bds_08_cap"] = ExtractBits(data, 8, 1) == 1

	// Bit 10: BDS 0,9 capability
	fields["bds_09_cap"] = ExtractBits(data, 9, 1) == 1

	// Bit 11: BDS 0,A capability
	fields["bds_0A_cap"] = ExtractBits(data, 10, 1) == 1

	// Bits 12-25: Reserved

	// Bit 26: BDS 2,0 capability
	fields["bds_20_cap"] = ExtractBits(data, 25, 1) == 1

	// Bit 27: BDS 2,1 capability
	fields["bds_21_cap"] = ExtractBits(data, 26, 1) == 1

	// Bit 28: BDS 4,0 capability
	fields["bds_40_cap"] = ExtractBits(data, 27, 1) == 1

	// Bit 29: BDS 4,1 capability
	fields["bds_41_cap"] = ExtractBits(data, 28, 1) == 1

	// Bit 30: BDS 4,2 capability
	fields["bds_42_cap"] = ExtractBits(data, 29, 1) == 1

	// Bit 31: BDS 4,3 capability
	fields["bds_43_cap"] = ExtractBits(data, 30, 1) == 1

	// Bit 32: BDS 4,4 capability
	fields["bds_44_cap"] = ExtractBits(data, 31, 1) == 1

	// Bit 33: BDS 4,5 capability
	fields["bds_45_cap"] = ExtractBits(data, 32, 1) == 1

	// Bit 34: BDS 4,8 capability
	fields["bds_48_cap"] = ExtractBits(data, 33, 1) == 1

	// Bit 35: BDS 5,0 capability
	fields["bds_50_cap"] = ExtractBits(data, 34, 1) == 1

	// Bit 36: BDS 5,1 capability
	fields["bds_51_cap"] = ExtractBits(data, 35, 1) == 1

	// Bit 37: BDS 5,2 capability
	fields["bds_52_cap"] = ExtractBits(data, 36, 1) == 1

	// Bit 38: BDS 5,3 capability
	fields["bds_53_cap"] = ExtractBits(data, 37, 1) == 1

	// Bit 39: BDS 5,4 capability
	fields["bds_54_cap"] = ExtractBits(data, 38, 1) == 1

	// Bit 40: BDS 5,5 capability
	fields["bds_55_cap"] = ExtractBits(data, 39, 1) == 1

	// Bit 41: BDS 5,6 capability
	fields["bds_56_cap"] = ExtractBits(data, 40, 1) == 1

	// Bit 42: BDS 5,F capability
	fields["bds_5F_cap"] = ExtractBits(data, 41, 1) == 1

	// Bit 43: BDS 6,0 capability
	fields["bds_60_cap"] = ExtractBits(data, 42, 1) == 1

	// Bits 44-56: Reserved
	
	return fields, nil
}