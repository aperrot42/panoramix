package transform

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
		Latitude  float64 `json:"latitude_deg"`
		Longitude float64 `json:"longitude_deg"`
		Altitude  float64 `json:"altitude_m,omitempty"`
		Source    string  `json:"source"` // Position: "POLAR" or "CARTESIAN", Altitude: "_ALT:3D_RADAR", "_ALT:FLIGHT_LEVEL"
	} `json:"wgs84_position"`

	// Metadata
	Timestamp PrecisionTime `json:"timestamp"`
	RadarSIC  uint8     `json:"radar_sic"`
	RadarSAC  uint8     `json:"radar_sac"`
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

// ExtractFromMessageTyped extracts aircraft observations from a typed ASTERIX message
func (ae *PositionExtractor) ExtractFromMessageTyped(msg *asterix.CAT048Message, timestamp time.Time) (*Position, error) {
	// Get radar position
	radarPos, ok := ae.radarRegistry.GetRadarPosition(msg.SIC, msg.SAC)
	if !ok {
		return nil, fmt.Errorf("radar position not found for SIC=%d SAC=%d - please configure radar position in radar_config.yaml", msg.SIC, msg.SAC)
	}

	// Extract data based on message category
	if msg.Category == 48 {
		obs, err := ae.extractFromCAT048Typed(msg, radarPos, timestamp)
		if err != nil {
			return nil, err
		}
		return &obs, nil
	}

	return nil, nil // No data extracted from other categories
}
