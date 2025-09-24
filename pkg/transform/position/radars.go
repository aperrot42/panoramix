package position

import (
	"fmt"
)

// RadarPosition represents the radar station position in WGS84
type RadarPosition struct {
	Latitude  float64 `json:"latitude_deg"`
	Longitude float64 `json:"longitude_deg"`
	Height    float64 `json:"height_m"`
	SIC       uint8   `json:"sic"`
	SAC       uint8   `json:"sac"`
}

// RadarRegistry manages radar station positions
type RadarRegistry struct {
	positions map[string]RadarPosition // key: "SIC-SAC"
}

// NewRadarRegistry creates a new radar registry
func NewRadarRegistry() *RadarRegistry {
	return &RadarRegistry{
		positions: make(map[string]RadarPosition),
	}
}

// AddRadarPosition registers a radar station position
func (rr *RadarRegistry) AddRadarPosition(pos RadarPosition) {
	key := radarKey(pos.SIC, pos.SAC)
	rr.positions[key] = pos
}

// GetRadarPosition retrieves a radar position by SIC/SAC
func (rr *RadarRegistry) GetRadarPosition(sic, sac uint8) (RadarPosition, bool) {
	key := radarKey(sic, sac)
	pos, found := rr.positions[key]
	return pos, found
}

// GetRadarPositionBySacSic is a convenience method to get radar position by SAC/SIC
func (rr *RadarRegistry) GetRadarPositionBySacSic(sac, sic uint8) (RadarPosition, bool) {
	return rr.GetRadarPosition(sic, sac)
}

// radarKey creates a unique key for radar identification
func radarKey(sic, sac uint8) string {
	return fmt.Sprintf("%d-%d", sic, sac)
}
