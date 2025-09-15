package bds

// BDS 4,0 - Selected Vertical Intention

// BDS40Decoder decodes BDS 4,0 - Selected Vertical Intention
type BDS40Decoder struct{}

func (d *BDS40Decoder) BDSCode() (uint8, uint8) {
	return 4, 0
}

func (d *BDS40Decoder) Decode(data []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// Bit 1: Status of MCP/FCU selected altitude
	if ExtractBits(data, 0, 1) == 1 {
		// Bits 2-13: MCP/FCU selected altitude (12 bits)
		altCode := ExtractBits(data, 1, 12)
		altitude := float64(altCode) * 16.0 // 16 ft resolution
		fields["selected_altitude_ft"] = altitude
		fields["selected_altitude_m"] = FeetToMeters(altitude)
		fields["selected_altitude_valid"] = true
	} else {
		fields["selected_altitude_valid"] = false
	}

	// Bit 14: Status of FMS selected altitude
	if ExtractBits(data, 13, 1) == 1 {
		// Bits 15-26: FMS selected altitude (12 bits)
		altCode := ExtractBits(data, 14, 12)
		altitude := float64(altCode) * 16.0 // 16 ft resolution
		fields["fms_altitude_ft"] = altitude
		fields["fms_altitude_m"] = FeetToMeters(altitude)
		fields["fms_altitude_valid"] = true
	} else {
		fields["fms_altitude_valid"] = false
	}

	// Bit 27: Status of barometric pressure setting
	if ExtractBits(data, 26, 1) == 1 {
		// Bits 28-39: Barometric pressure setting (12 bits)
		pressCode := ExtractBits(data, 27, 12)
		pressure := float64(pressCode)*0.1 + 800.0 // 0.1 mb resolution, 800 mb offset
		fields["baro_pressure_mb"] = pressure
		fields["baro_pressure_inhg"] = pressure * 0.02953 // Convert to inHg
		fields["baro_pressure_valid"] = true
	} else {
		fields["baro_pressure_valid"] = false
	}

	// Bits 40-47: Reserved (8 bits)

	// Bit 48: Status of MCP/FCU mode
	if ExtractBits(data, 47, 1) == 1 {
		fields["mcp_fcu_mode_valid"] = true
		// Bit 49: VNAV mode
		fields["vnav_mode"] = ExtractBits(data, 48, 1) == 1
		// Bit 50: ALT HOLD mode
		fields["alt_hold_mode"] = ExtractBits(data, 49, 1) == 1
		// Bit 51: Approach mode
		fields["approach_mode"] = ExtractBits(data, 50, 1) == 1
	} else {
		fields["mcp_fcu_mode_valid"] = false
	}

	// Bits 52-53: Reserved

	// Bit 54: Status of target altitude source
	if ExtractBits(data, 53, 1) == 1 {
		// Bits 55-56: Target altitude source
		source := ExtractBits(data, 54, 2)
		switch source {
		case 0:
			fields["target_altitude_source"] = "Unknown"
		case 1:
			fields["target_altitude_source"] = "Aircraft altitude"
		case 2:
			fields["target_altitude_source"] = "FCU/MCP selected altitude"
		case 3:
			fields["target_altitude_source"] = "FMS selected altitude"
		}
		fields["target_altitude_source_valid"] = true
	} else {
		fields["target_altitude_source_valid"] = false
	}

	return fields, nil
}