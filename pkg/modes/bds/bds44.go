package bds

// BDS44Decoder decodes BDS 4,4 - Meteorological Routine Air Report
type BDS44Decoder struct{}

func (d *BDS44Decoder) BDSCode() (uint8, uint8) {
	return 4, 4
}

func (d *BDS44Decoder) Decode(data []byte) (any, error) {
	decoded := BDS44Decoded{}

	// Bit 1: Figure of Merit/Source
	fom := ExtractBits(data, 0, 1)
	if fom == 0 {
		decoded.WindSource = "INS"
	} else {
		decoded.WindSource = "FMS"
	}

	// Bit 2: Status of wind speed and direction
	if ExtractBits(data, 1, 1) == 1 {
		// Bits 3-11: Wind speed (9 bits)
		windSpeed := ExtractBits(data, 2, 9) // Resolution: 1 knot
		decoded.WindSpeedKt = windSpeed
		// Bits 12-20: Wind direction (9 bits)
		windDir := float64(ExtractBits(data, 11, 9)) * 180.0 / 256.0 // Resolution: 180/256 degrees
		decoded.WindDirectionDeg = windDir
		decoded.WindValid = true
	} else {
		decoded.WindValid = false
	}

	// Bit 21: Status of static air temperature
	if ExtractBits(data, 20, 1) == 1 {
		// Bits 22-32: Static air temperature (11 bits, signed)
		tempRaw := ExtractSignedBits(data, 21, 11)
		temperature := float64(tempRaw) * 0.25 // Resolution: 0.25 degrees Celsius
		decoded.StaticAirTemperatureC = temperature
		decoded.StaticAirTemperatureValid = true
	} else {
		decoded.StaticAirTemperatureValid = false
	}

	// Bit 33: Status of average static pressure
	if ExtractBits(data, 32, 1) == 1 {
		// Bits 34-44: Average static pressure (11 bits)
		pressRaw := ExtractBits(data, 33, 11)
		pressure := float64(pressRaw) // Resolution: 1 hPa
		decoded.AverageStaticPressureHpa = pressure
		decoded.AverageStaticPressureValid = true
	} else {
		decoded.AverageStaticPressureValid = false
	}

	// Bit 45: Status of turbulence
	if ExtractBits(data, 44, 1) == 1 {
		// Bits 46-47: Turbulence (2 bits)
		turb := ExtractBits(data, 45, 2)
		decoded.TurbulenceRaw = turb
		switch turb {
		case 0:
			decoded.Turbulence = "None"
		case 1:
			decoded.Turbulence = "Light"
		case 2:
			decoded.Turbulence = "Moderate"
		case 3:
			decoded.Turbulence = "Severe"
		}
		decoded.TurbulenceValid = true
	} else {
		decoded.TurbulenceValid = false
	}

	// Bit 48: Status of humidity
	if ExtractBits(data, 47, 1) == 1 {
		// Bits 49-54: Humidity (6 bits)
		humidity := float64(ExtractBits(data, 48, 6)) * 100.0 / 64.0 // Resolution: 100/64 percent
		decoded.HumidityPercent = humidity
		decoded.HumidityValid = true
	} else {
		decoded.HumidityValid = false
	}

	// Bits 55-56: Reserved

	return decoded, nil
}
