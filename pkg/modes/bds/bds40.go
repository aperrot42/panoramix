package bds

// BDS 4,0 - Selected Vertical Intention

// BDS40Decoder decodes BDS 4,0 - Selected Vertical Intention
type BDS40Decoder struct{}

func (d *BDS40Decoder) BDSCode() (uint8, uint8) {
	return 4, 0
}

func (d *BDS40Decoder) Decode(data []byte) (interface{}, error) {
	decoded := BDS40Decoded{}

	// Bit 1: Status of MCP/FCU selected altitude
	if ExtractBits(data, 0, 1) == 1 {
		// Bits 2-13: MCP/FCU selected altitude (12 bits)
		altCode := ExtractBits(data, 1, 12)
		altitude := float64(altCode) * 16.0 // 16 ft resolution
		decoded.SelectedAltitudeFt = altitude
		decoded.SelectedAltitudeValid = true
	} else {
		decoded.SelectedAltitudeValid = false
	}

	// Bit 14: Status of FMS selected altitude
	if ExtractBits(data, 13, 1) == 1 {
		// Bits 15-26: FMS selected altitude (12 bits)
		altCode := ExtractBits(data, 14, 12)
		altitude := float64(altCode) * 16.0 // 16 ft resolution
		decoded.FmsAltitudeFt = altitude
		decoded.FmsAltitudeValid = true
	} else {
		decoded.FmsAltitudeValid = false
	}

	// Bit 27: Status of barometric pressure setting
	if ExtractBits(data, 26, 1) == 1 {
		// Bits 28-39: Barometric pressure setting (12 bits)
		pressCode := ExtractBits(data, 27, 12)
		pressure := float64(pressCode)*0.1 + 800.0 // 0.1 mb resolution, 800 mb offset
		decoded.BaroPressureMb = pressure
		decoded.BaroPressureValid = true
	} else {
		decoded.BaroPressureValid = false
	}

	// Bits 40-47: Reserved (8 bits)

	// Bit 48: Status of MCP/FCU mode
	if ExtractBits(data, 47, 1) == 1 {
		decoded.McpFcuModeValid = true
		// Bit 49: VNAV mode
		decoded.VnavMode = ExtractBits(data, 48, 1) == 1
		// Bit 50: ALT HOLD mode
		decoded.AltHoldMode = ExtractBits(data, 49, 1) == 1
		// Bit 51: Approach mode
		decoded.ApproachMode = ExtractBits(data, 50, 1) == 1
	} else {
		decoded.McpFcuModeValid = false
	}

	// Bits 52-53: Reserved

	// Bit 54: Status of target altitude source
	if ExtractBits(data, 53, 1) == 1 {
		// Bits 55-56: Target altitude source
		source := ExtractBits(data, 54, 2)
		decoded.TargetAltitudeSourceRaw = uint32(source)
		switch source {
		case 0:
			decoded.TargetAltitudeSource = "Unknown"
		case 1:
			decoded.TargetAltitudeSource = "Aircraft altitude"
		case 2:
			decoded.TargetAltitudeSource = "FCU/MCP selected altitude"
		case 3:
			decoded.TargetAltitudeSource = "FMS selected altitude"
		}
		decoded.TargetAltitudeSourceValid = true
	} else {
		decoded.TargetAltitudeSourceValid = false
	}

	return decoded, nil
}
