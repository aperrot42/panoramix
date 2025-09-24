package position

import (
	"fmt"
	"time"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

// PrecisionTime provides millisecond-precision JSON marshaling for timestamps
type PrecisionTime struct {
	time.Time
}

// MarshalJSON implements custom JSON marshaling with millisecond precision
func (pt PrecisionTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pt.Time.UTC().Format("2006-01-02T15:04:05.000Z") + `"`), nil
}

// Position represents the position of an aircraft from raw Asterix cat048 using cat034 messages
type Position struct {
	// Computed values (derived from raw data)
	WGS84Position struct {
		Latitude_deg   float64 `json:"latitude_deg"`
		Longitude_deg  float64 `json:"longitude_deg"`
		PostionSource  string  `json:"position_source"` // PostionSource: "POLAR" or "CARTESIAN",
		AltitudeFt     float64 `json:"altitude_ft,omitempty"`
		AltitudeSource string  `json:"altitude_source"` // AltitudeSource: "3D_RADAR", "FLIGHT_LEVEL"
	} `json:"wgs84_position"`

	// Metadata
	Timestamp PrecisionTime `json:"timestamp"`
	RadarSIC  uint8         `json:"radar_sic"`
	RadarSAC  uint8         `json:"radar_sac"`
}

// PositionExtractor extracts position information from ASTERIX cat 034 messages for and cat048 messages
type PositionExtractor struct {
	radarRegistry *RadarRegistry
}

func NewPositionExtractor() *PositionExtractor {
	return &PositionExtractor{
		radarRegistry: NewRadarRegistry(),
	}
}

// NewPositionExtractorFromConfig creates a PositionExtractor with radar positions loaded from config file
func NewPositionExtractorFromConfig(configPath string) (*PositionExtractor, error) {
	registry, err := LoadRadarConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load radar config: %w", err)
	}

	return &PositionExtractor{
		radarRegistry: registry,
	}, nil
}

// AddRadarPosition registers a radar station position
func (ae *PositionExtractor) AddRadarPosition(pos RadarPosition) {
	ae.radarRegistry.AddRadarPosition(pos)
}

// ExtractFromMessage extracts aircraft observations from an ASTERIX message
func (ae *PositionExtractor) ExtractFromMessage(msg *asterix.AsterixMessage, timestamp time.Time) (*Position, error) {
	// Get radar position
	radarPos, ok := ae.radarRegistry.GetRadarPosition(msg.Sic, msg.Sac)
	if !ok {
		return nil, fmt.Errorf("radar position not found for SIC=%d SAC=%d - please configure radar position in radar_config.yaml", msg.Sic, msg.Sac)
	}

	// Extract data based on message category
	switch msg.Category {
	case 48:
		obs, err := ae.extractFromCAT048(msg, radarPos, timestamp)
		if err != nil {
			return nil, err
		}
		return &obs, nil
	default:
		return nil, nil // No data extracted from other categories
	}
}

// extractFromCAT048 extracts a single aircraft observation from CAT 048 messages
func (ae *PositionExtractor) extractFromCAT048(msg *asterix.AsterixMessage, radarPos RadarPosition, timestamp time.Time) (Position, error) {

	// Create base observation with raw ASTERIX data
	obs := Position{
		Timestamp: PrecisionTime{timestamp},
		RadarSIC:  msg.Sic,
		RadarSAC:  msg.Sac,
	}

	// Compute WGS84 position from coordinates
	hasPosition := false

	// Try polar coordinates first (I048/040) - more accurate
	if polarPos, ok := msg.Items["I048/040"].(map[string]interface{}); ok {
		if rng, rngOk := polarPos["rho_nm"].(float64); rngOk {
			if az, azOk := polarPos["theta_deg"].(float64); azOk {
				rangeM := NauticalMilesToMeters(rng)
				lat, lon := PolarToWGS84(radarPos, rangeM, az)
				obs.WGS84Position.Latitude_deg = lat
				obs.WGS84Position.Longitude_deg = lon
				obs.WGS84Position.PostionSource = "POLAR"
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
					obs.WGS84Position.Latitude_deg = lat
					obs.WGS84Position.Longitude_deg = lon
					obs.WGS84Position.PostionSource = "CARTESIAN"
					hasPosition = true
				}
			}
		}
	}

	if !hasPosition {
		return obs, fmt.Errorf("no position data found in CAT 048 message")
	}

	// Extract altitude with priority: 3D Radar > Flight Level
	altitudeSource := ""

	// Priority 1: I048/110 - Height Measured by 3D Radar (most accurate)
	if heightData, ok := msg.Items["I048/110"].(map[string]interface{}); ok {
		if height, heightOk := heightData["height_ft"].(float64); heightOk {
			// TODO fix conversion to get feets
			obs.WGS84Position.AltitudeFt = height
			obs.WGS84Position.AltitudeSource = "3D_RADAR"
		}
	}

	// Priority 2: I048/090 - Flight Level (current barometric altitude)
	if altitudeSource == "" {
		if flightLevel, ok := msg.Items["I048/090"].(struct {
			FlightLevel float64
			RawValue    int16 // 1/4th of FL
			Validated   bool
			Garbled     bool
		}); ok {
			obs.WGS84Position.AltitudeFt = float64(int32(flightLevel.RawValue) * 25) // 0.25FL * 100 => ft
			obs.WGS84Position.AltitudeSource = "FLIGHT_LEVEL"
		}
	}

	return obs, nil
}
