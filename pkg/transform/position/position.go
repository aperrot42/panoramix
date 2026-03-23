package position

import (
	"fmt"
	"time"
)

// PrecisionTime provides millisecond-precision JSON marshaling for timestamps
type PrecisionTime struct {
	time.Time
}

// MarshalJSON implements custom JSON marshaling with millisecond precision
func (pt PrecisionTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pt.Time.UTC().Format("2006-01-02T15:04:05.000Z") + `"`), nil
}

// RawPlot holds the raw measurement data needed for position extraction.
// All position/altitude fields are optional — nil means not present in the source message.
type RawPlot struct {
	SIC   uint8
	SAC   uint8
	Polar *PolarCoord // range + azimuth from radar
	Cart  *CartCoord  // cartesian offset from radar
	Alt3D *float64    // height from 3D radar (ft)
	FL    *float64    // flight level
}

// PolarCoord holds a polar measurement from a radar.
type PolarCoord struct {
	RhoNM    float64 // range in nautical miles
	ThetaDeg float64 // azimuth in degrees from north, clockwise
}

// CartCoord holds a cartesian offset from a radar.
type CartCoord struct {
	XNM float64 // east offset in nautical miles
	YNM float64 // north offset in nautical miles
}

// Position represents a computed aircraft position in WGS84 coordinates.
type Position struct {
	WGS84Position struct {
		Latitude_deg   float64 `json:"latitude_deg"`
		Longitude_deg  float64 `json:"longitude_deg"`
		PostionSource  string  `json:"position_source"`          // "POLAR" or "CARTESIAN"
		AltitudeFt     float64 `json:"altitude_ft,omitempty"`
		AltitudeSource string  `json:"altitude_source,omitempty"` // "3D_RADAR" or "FLIGHT_LEVEL"
	} `json:"wgs84_position"`

	Timestamp PrecisionTime `json:"timestamp"`
	RadarSIC  uint8         `json:"radar_sic"`
	RadarSAC  uint8         `json:"radar_sac"`
}

// PositionExtractor computes WGS84 positions from raw radar measurements.
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

// Extract computes a WGS84 position from raw radar measurements.
func (ae *PositionExtractor) Extract(plot RawPlot, timestamp time.Time) (*Position, error) {
	radarPos, ok := ae.radarRegistry.GetRadarPosition(plot.SIC, plot.SAC)
	if !ok {
		return nil, fmt.Errorf("radar position not found for SIC=%d SAC=%d", plot.SIC, plot.SAC)
	}

	obs := Position{
		Timestamp: PrecisionTime{timestamp},
		RadarSIC:  plot.SIC,
		RadarSAC:  plot.SAC,
	}

	hasPosition := false

	// Try polar coordinates first — more accurate
	if plot.Polar != nil {
		rangeM := NauticalMilesToMeters(plot.Polar.RhoNM)
		lat, lon := PolarToWGS84(radarPos, rangeM, plot.Polar.ThetaDeg)
		obs.WGS84Position.Latitude_deg = lat
		obs.WGS84Position.Longitude_deg = lon
		obs.WGS84Position.PostionSource = "POLAR"
		hasPosition = true
	}

	// Fall back to cartesian
	if !hasPosition && plot.Cart != nil {
		xM := NauticalMilesToMeters(plot.Cart.XNM)
		yM := NauticalMilesToMeters(plot.Cart.YNM)
		lat, lon := CartesianToWGS84(radarPos, xM, yM)
		obs.WGS84Position.Latitude_deg = lat
		obs.WGS84Position.Longitude_deg = lon
		obs.WGS84Position.PostionSource = "CARTESIAN"
		hasPosition = true
	}

	if !hasPosition {
		return nil, fmt.Errorf("no position data in plot")
	}

	// Altitude priority: 3D radar > flight level
	if plot.Alt3D != nil {
		obs.WGS84Position.AltitudeFt = *plot.Alt3D
		obs.WGS84Position.AltitudeSource = "3D_RADAR"
	} else if plot.FL != nil {
		obs.WGS84Position.AltitudeFt = *plot.FL * 100 // FL to feet
		obs.WGS84Position.AltitudeSource = "FLIGHT_LEVEL"
	}

	return &obs, nil
}