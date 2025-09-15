package bds

// BDS 6,0 - Heading and Speed Report

// BDS60Decoder decodes BDS 6,0 - Heading and Speed Report
type BDS60Decoder struct{}

func (d *BDS60Decoder) BDSCode() (uint8, uint8) {
	return 6, 0
}

func (d *BDS60Decoder) Decode(data []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// Bit 1: Status of magnetic heading
	if ExtractBits(data, 0, 1) == 1 {
		// Bits 2-12: Magnetic heading (11 bits)
		headingRaw := ExtractBits(data, 1, 11)
		magHeading := float64(headingRaw) * 90.0 / 512.0 // Resolution: 90/512 degrees
		fields["magnetic_heading_deg"] = magHeading
		fields["magnetic_heading_valid"] = true
	} else {
		fields["magnetic_heading_valid"] = false
	}

	// Bit 13: Status of indicated airspeed
	if ExtractBits(data, 12, 1) == 1 {
		// Bits 14-23: Indicated airspeed (10 bits)
		iasRaw := ExtractBits(data, 13, 10)
		indicatedAirspeed := float64(iasRaw) // Resolution: 1 knot
		fields["indicated_airspeed_kt"] = indicatedAirspeed
		fields["indicated_airspeed_ms"] = KnotsToMetersPerSecond(indicatedAirspeed)
		fields["indicated_airspeed_valid"] = true
	} else {
		fields["indicated_airspeed_valid"] = false
	}

	// Bit 24: Status of Mach number
	if ExtractBits(data, 23, 1) == 1 {
		// Bits 25-34: Mach number (10 bits)
		machRaw := ExtractBits(data, 24, 10)
		machNumber := float64(machRaw) * 2.048 / 512.0 // Resolution: 2.048/512
		fields["mach_number"] = machNumber
		fields["mach_number_valid"] = true
	} else {
		fields["mach_number_valid"] = false
	}

	// Bit 35: Status of barometric altitude rate
	if ExtractBits(data, 34, 1) == 1 {
		// Bits 36-45: Barometric altitude rate (10 bits, signed)
		rateRaw := ExtractSignedBits(data, 35, 10)
		baroRate := float64(rateRaw) * 32.0 // Resolution: 32 ft/min
		fields["baro_altitude_rate_fpm"] = baroRate
		fields["baro_altitude_rate_ms"] = baroRate * 0.00508 // Convert to m/s
		fields["baro_altitude_rate_valid"] = true
	} else {
		fields["baro_altitude_rate_valid"] = false
	}

	// Bit 46: Status of inertial vertical velocity
	if ExtractBits(data, 45, 1) == 1 {
		// Bits 47-56: Inertial vertical velocity (10 bits, signed)
		ivvRaw := ExtractSignedBits(data, 46, 10)
		inertialVV := float64(ivvRaw) * 32.0 // Resolution: 32 ft/min
		fields["inertial_vertical_velocity_fpm"] = inertialVV
		fields["inertial_vertical_velocity_ms"] = inertialVV * 0.00508 // Convert to m/s
		fields["inertial_vertical_velocity_valid"] = true
	} else {
		fields["inertial_vertical_velocity_valid"] = false
	}

	return fields, nil
}