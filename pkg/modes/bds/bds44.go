package bds

// BDS44Decoder decodes BDS 4,4 - Meteorological Routine Air Report
type BDS44Decoder struct{}

func (d *BDS44Decoder) BDSCode() (uint8, uint8) {
	return 4, 4
}

func (d *BDS44Decoder) Decode(data []byte) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// Bit 1: Figure of Merit/Source
	fom := ExtractBits(data, 0, 1)
	if fom == 0 {
		fields["wind_source"] = "INS"
	} else {
		fields["wind_source"] = "FMS"
	}

	// Bit 2: Status of wind speed and direction
	if ExtractBits(data, 1, 1) == 1 {
		// Bits 3-11: Wind speed (9 bits)
		windSpeed := ExtractBits(data, 2, 9) // Resolution: 1 knot
		fields["wind_speed_kt"] = windSpeed
		fields["wind_speed_ms"] = KnotsToMetersPerSecond(float64(windSpeed))
		
		// Bits 12-20: Wind direction (9 bits)
		windDir := float64(ExtractBits(data, 11, 9)) * 180.0 / 256.0 // Resolution: 180/256 degrees
		fields["wind_direction_deg"] = windDir
		fields["wind_valid"] = true
	} else {
		fields["wind_valid"] = false
	}

	// Bit 21: Status of static air temperature
	if ExtractBits(data, 20, 1) == 1 {
		// Bits 22-32: Static air temperature (11 bits, signed)
		tempRaw := ExtractSignedBits(data, 21, 11)
		temperature := float64(tempRaw) * 0.25 // Resolution: 0.25 degrees Celsius
		fields["static_air_temperature_c"] = temperature
		fields["static_air_temperature_valid"] = true
	} else {
		fields["static_air_temperature_valid"] = false
	}

	// Bit 33: Status of average static pressure
	if ExtractBits(data, 32, 1) == 1 {
		// Bits 34-44: Average static pressure (11 bits)
		pressRaw := ExtractBits(data, 33, 11)
		pressure := float64(pressRaw) // Resolution: 1 hPa
		fields["average_static_pressure_hpa"] = pressure
		fields["average_static_pressure_valid"] = true
	} else {
		fields["average_static_pressure_valid"] = false
	}

	// Bit 45: Status of turbulence
	if ExtractBits(data, 44, 1) == 1 {
		// Bits 46-47: Turbulence (2 bits)
		turb := ExtractBits(data, 45, 2)
		switch turb {
		case 0:
			fields["turbulence"] = "None"
		case 1:
			fields["turbulence"] = "Light"
		case 2:
			fields["turbulence"] = "Moderate"
		case 3:
			fields["turbulence"] = "Severe"
		}
		fields["turbulence_valid"] = true
	} else {
		fields["turbulence_valid"] = false
	}

	// Bit 48: Status of humidity
	if ExtractBits(data, 47, 1) == 1 {
		// Bits 49-54: Humidity (6 bits)
		humidity := float64(ExtractBits(data, 48, 6)) * 100.0 / 64.0 // Resolution: 100/64 percent
		fields["humidity_percent"] = humidity
		fields["humidity_valid"] = true
	} else {
		fields["humidity_valid"] = false
	}

	// Bits 55-56: Reserved
	
	return fields, nil
}