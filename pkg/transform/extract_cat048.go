package transform

import (
	"fmt"
	"time"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

// extractFromCAT048 extracts a single aircraft observation from CAT 048 messages
func (ae *PositionExtractor) extractFromCAT048(msg *asterix.AsterixMessage, radarPos RadarPosition, timestamp time.Time) (Position, error) {

	// Create base observation with raw ASTERIX data
	obs := Position{
		Timestamp: PrecisionTime{timestamp},
		RadarSIC:  msg.Sic,
		RadarSAC:  msg.Sac,
	}

	// Extract altitude with priority: 3D Radar > BDS registers > Flight Level
	altitudeSource := ""

	// Priority 1: I048/110 - Height Measured by 3D Radar (most accurate)
	if heightData, ok := msg.Items["I048/110"].(map[string]interface{}); ok {
		if height, heightOk := heightData["height_ft"].(float64); heightOk {
			obs.WGS84Position.Altitude = height * 0.3048 // Convert feet to meters
			altitudeSource = "3D_RADAR"
		}
	}

	// Priority 2: Currently no other current altitude sources available in BDS registers
	// Note: BDS 4.0 selected_altitude_ft and fms_altitude_ft are TARGET altitudes, not current altitude

	// Priority 2: I048/090 - Flight Level (current barometric altitude)
	if altitudeSource == "" {
		if flightLevel, ok := msg.Items["I048/090"].(struct {
			FlightLevel float64
			RawValue    int16
			Validated   bool
			Garbled     bool
		}); ok {
			obs.WGS84Position.Altitude = FlightLevelToMeters(flightLevel.FlightLevel)
			altitudeSource = "FLIGHT_LEVEL"
		}
	}

	// Store altitude source for debugging
	if altitudeSource != "" {
		obs.WGS84Position.Source = obs.WGS84Position.Source + "_ALT:" + altitudeSource
	}

	// Compute WGS84 position from coordinates
	hasPosition := false

	// Try polar coordinates first (I048/040) - more accurate
	if polarPos, ok := msg.Items["I048/040"].(map[string]interface{}); ok {
		if rng, rngOk := polarPos["rho_nm"].(float64); rngOk {
			if az, azOk := polarPos["theta_deg"].(float64); azOk {
				rangeM := NauticalMilesToMeters(rng)
				lat, lon := PolarToWGS84(radarPos, rangeM, az)
				obs.WGS84Position.Latitude = lat
				obs.WGS84Position.Longitude = lon
				obs.WGS84Position.Source = "POLAR"
				hasPosition = true
			}
		}
	}

	// Try Cartesian coordinates (I048/042) if no polar
	if !hasPosition {
		if cartPos, ok := msg.Items["I048/042"].(map[string]interface{}); ok {
			if x, xOk := cartPos["x_nm"].(float64); xOk {
				if y, yOk := cartPos["y_nm"].(float64); yOk {
					xM := NauticalMilesToMeters(x)
					yM := NauticalMilesToMeters(y)
					lat, lon := CartesianToWGS84(radarPos, xM, yM)
					obs.WGS84Position.Latitude = lat
					obs.WGS84Position.Longitude = lon
					obs.WGS84Position.Source = "CARTESIAN"
					hasPosition = true
				}
			}
		}
	}

	if !hasPosition {
		return obs, fmt.Errorf("no position data found in CAT 048 message")
	}

	return obs, nil
}

// extractFromCAT048Typed extracts a single aircraft observation from typed CAT 048 messages
func (ae *PositionExtractor) extractFromCAT048Typed(msg *asterix.CAT048Message, radarPos RadarPosition, timestamp time.Time) (Position, error) {

	// Create base observation with typed ASTERIX data
	obs := Position{
		Timestamp: PrecisionTime{timestamp},
		RadarSIC:  msg.SIC,
		RadarSAC:  msg.SAC,
	}

	// Extract flight level and compute altitude
	if msg.FlightLevel != nil {
		obs.WGS84Position.Altitude = FeetToMeters(msg.FlightLevel.FlightLevelFt)
	}

	// Compute WGS84 position from coordinates
	hasPosition := false

	// Try polar coordinates first (I048/040) - more accurate
	if msg.MeasuredPosition != nil {
		rangeM := NauticalMilesToMeters(msg.MeasuredPosition.Rho)
		lat, lon := PolarToWGS84(radarPos, rangeM, msg.MeasuredPosition.Theta)
		obs.WGS84Position.Latitude = lat
		obs.WGS84Position.Longitude = lon
		obs.WGS84Position.Source = "POLAR"
		hasPosition = true
	}

	// Try Cartesian coordinates (I048/041) if no polar  
	if !hasPosition && msg.CalculatedPositionCartesian != nil {
		// Note: CalculatedPositionCartesian is already in meters according to spec
		xM := msg.CalculatedPositionCartesian.X
		yM := msg.CalculatedPositionCartesian.Y
		lat, lon := CartesianToWGS84(radarPos, xM, yM)
		obs.WGS84Position.Latitude = lat
		obs.WGS84Position.Longitude = lon
		obs.WGS84Position.Source = "CARTESIAN"
		hasPosition = true
	}

	if !hasPosition {
		return obs, fmt.Errorf("no position data found in CAT 048 message")
	}

	return obs, nil
}
