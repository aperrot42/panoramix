package bds

// BDS 6,0 - Heading and Speed Report

// BDS60Decoder decodes BDS 6,0 - Heading and Speed Report
type BDS60Decoder struct{}

func (d *BDS60Decoder) BDSCode() (uint8, uint8) {
	return 6, 0
}

func (d *BDS60Decoder) Decode(data []byte) (any, error) {
	decoded := BDS60Decoded{}

	// Bit 1: Status of magnetic heading
	if ExtractBits(data, 0, 1) == 1 {
		// Bits 2-12: Magnetic heading (11 bits)
		headingRaw := ExtractBits(data, 1, 11)
		magHeading := float64(headingRaw) * 90.0 / 512.0 // Resolution: 90/512 degrees
		decoded.MagneticHeadingDeg = magHeading
		decoded.MagneticHeadingValid = true
	} else {
		decoded.MagneticHeadingValid = false
	}

	// Bit 13: Status of indicated airspeed
	if ExtractBits(data, 12, 1) == 1 {
		// Bits 14-23: Indicated airspeed (10 bits)
		iasRaw := ExtractBits(data, 13, 10)
		indicatedAirspeed := float64(iasRaw) // Resolution: 1 knot
		decoded.IndicatedAirspeedKt = indicatedAirspeed
		decoded.IndicatedAirspeedValid = true
	} else {
		decoded.IndicatedAirspeedValid = false
	}

	// Bit 24: Status of Mach number
	if ExtractBits(data, 23, 1) == 1 {
		// Bits 25-34: Mach number (10 bits)
		machRaw := ExtractBits(data, 24, 10)
		machNumber := float64(machRaw) * 2.048 / 512.0 // Resolution: 2.048/512
		decoded.MachNumber = machNumber
		decoded.MachNumberValid = true
	} else {
		decoded.MachNumberValid = false
	}

	// Bit 35: Status of barometric altitude rate
	if ExtractBits(data, 34, 1) == 1 {
		// Bits 36-45: Barometric altitude rate (10 bits, signed)
		rateRaw := ExtractSignedBits(data, 35, 10)
		baroRate := float64(rateRaw) * 32.0 // Resolution: 32 ft/min
		decoded.BaroAltitudeRateFpm = baroRate
		decoded.BaroAltitudeRateValid = true
	} else {
		decoded.BaroAltitudeRateValid = false
	}

	// Bit 46: Status of inertial vertical velocity
	if ExtractBits(data, 45, 1) == 1 {
		// Bits 47-56: Inertial vertical velocity (10 bits, signed)
		ivvRaw := ExtractSignedBits(data, 46, 10)
		inertialVV := float64(ivvRaw) * 32.0 // Resolution: 32 ft/min
		decoded.InertialVerticalVelocityFpm = inertialVV
		decoded.InertialVerticalVelocityValid = true
	} else {
		decoded.InertialVerticalVelocityValid = false
	}

	return decoded, nil
}
