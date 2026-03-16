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
	rec := msg.Record
	if rec == nil {
		return nil, nil
	}

	radarPos, ok := ae.radarRegistry.GetRadarPosition(rec.GetSIC(), rec.GetSAC())
	if !ok {
		return nil, fmt.Errorf("radar position not found for SIC=%d SAC=%d - please configure radar position in radar_config.yaml", rec.GetSIC(), rec.GetSAC())
	}

	switch r := rec.(type) {
	case *asterix.Cat048Message:
		obs, err := ae.extractFromCAT048(r, radarPos, timestamp)
		if err != nil {
			return nil, err
		}
		return &obs, nil
	default:
		return nil, nil
	}
}

// extractFromCAT048 extracts a single aircraft observation from CAT 048 messages
func (ae *PositionExtractor) extractFromCAT048(rec *asterix.Cat048Message, radarPos RadarPosition, timestamp time.Time) (Position, error) {
	obs := Position{
		Timestamp: PrecisionTime{timestamp},
		RadarSIC:  rec.GetSIC(),
		RadarSAC:  rec.GetSAC(),
	}

	hasPosition := false

	// Try polar coordinates first (I048/040) - more accurate
	if rec.MeasuredPosition != nil {
		rangeM := NauticalMilesToMeters(rec.MeasuredPosition.Rho)
		lat, lon := PolarToWGS84(radarPos, rangeM, rec.MeasuredPosition.Theta)
		obs.WGS84Position.Latitude_deg = lat
		obs.WGS84Position.Longitude_deg = lon
		obs.WGS84Position.PostionSource = "POLAR"
		hasPosition = true
	}

	// Try Cartesian coordinates (I048/042) if no polar
	if !hasPosition && rec.CalculatedPosition != nil {
		xM := NauticalMilesToMeters(rec.CalculatedPosition.X)
		yM := NauticalMilesToMeters(rec.CalculatedPosition.Y)
		lat, lon := CartesianToWGS84(radarPos, xM, yM)
		obs.WGS84Position.Latitude_deg = lat
		obs.WGS84Position.Longitude_deg = lon
		obs.WGS84Position.PostionSource = "CARTESIAN"
		hasPosition = true
	}

	if !hasPosition {
		return obs, fmt.Errorf("no position data found in CAT 048 message")
	}

	// Extract altitude with priority: 3D Radar > Flight Level

	// Priority 1: I048/110 - Height Measured by 3D Radar (most accurate)
	if rec.Height3D != nil {
		obs.WGS84Position.AltitudeFt = rec.Height3D.Height
		obs.WGS84Position.AltitudeSource = "3D_RADAR"
	}

	// Priority 2: I048/090 - Flight Level (current barometric altitude)
	if obs.WGS84Position.AltitudeSource == "" && rec.FlightLevel != nil {
		obs.WGS84Position.AltitudeFt = rec.FlightLevel.FL * 100 // FL to feet
		obs.WGS84Position.AltitudeSource = "FLIGHT_LEVEL"
	}

	return obs, nil
}