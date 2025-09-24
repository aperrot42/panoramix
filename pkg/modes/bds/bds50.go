package bds

// BDS 5,0 - Track and Turn Report

// BDS50Decoder decodes BDS 5,0 - Track and Turn Report
type BDS50Decoder struct{}

func (d *BDS50Decoder) BDSCode() (uint8, uint8) {
	return 5, 0
}

func (d *BDS50Decoder) Decode(data []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// Bit 1: Status of roll angle
	if ExtractBits(data, 0, 1) == 1 {
		// Bits 2-11: Roll angle (10 bits, signed)
		rollRaw := ExtractSignedBits(data, 1, 10)
		rollAngle := float64(rollRaw) * 45.0 / 256.0 // Resolution: 45/256 degrees
		fields["roll_angle_deg"] = rollAngle
		fields["roll_angle_valid"] = true
	} else {
		fields["roll_angle_valid"] = false
	}

	// Bit 12: Status of true track angle
	if ExtractBits(data, 11, 1) == 1 {
		// Bits 13-23: True track angle (11 bits)
		trackRaw := ExtractBits(data, 12, 11)
		trackAngle := float64(trackRaw) * 90.0 / 512.0 // Resolution: 90/512 degrees
		fields["true_track_angle_deg"] = trackAngle
		fields["true_track_angle_valid"] = true
	} else {
		fields["true_track_angle_valid"] = false
	}

	// Bit 24: Status of ground speed
	if ExtractBits(data, 23, 1) == 1 {
		// Bits 25-34: Ground speed (10 bits)
		speedRaw := ExtractBits(data, 24, 10)
		groundSpeed := float64(speedRaw) * 2.0 // Resolution: 2 knots
		fields["ground_speed_kt"] = groundSpeed
		fields["ground_speed_valid"] = true
	} else {
		fields["ground_speed_valid"] = false
	}

	// Bit 35: Status of track angle rate
	if ExtractBits(data, 34, 1) == 1 {
		// Bits 36-45: Track angle rate (10 bits, signed)
		rateRaw := ExtractSignedBits(data, 35, 10)
		trackRate := float64(rateRaw) * 8.0 / 256.0 // Resolution: 8/256 degrees/second
		fields["track_angle_rate_deg_s"] = trackRate
		fields["track_angle_rate_valid"] = true
	} else {
		fields["track_angle_rate_valid"] = false
	}

	// Bit 46: Status of true airspeed
	if ExtractBits(data, 45, 1) == 1 {
		// Bits 47-56: True airspeed (10 bits)
		tasRaw := ExtractBits(data, 46, 10)
		trueAirspeed := float64(tasRaw) * 2.0 // Resolution: 2 knots
		fields["true_airspeed_kt"] = trueAirspeed
		fields["true_airspeed_valid"] = true
	} else {
		fields["true_airspeed_valid"] = false
	}

	return fields, nil
}