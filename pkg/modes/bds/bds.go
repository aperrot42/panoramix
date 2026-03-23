// Package bds provides decoding for Mode S BDS (Comm-B Data Selector) registers
// as specified in ED-73F MOPS for Secondary Surveillance Radar Mode S Transponders
package bds

import (
	"fmt"
)

type BDS10Decoded struct {
	AcasReserved       uint32 `json:"acas_reserved"`
	BDS10Cf            bool   `json:"bds_10_cf"`
	BDS17Cap           bool   `json:"bds_17_cap"`
	CommBBroadcast1Cap bool   `json:"comm_b_broadcast_1_cap"`
	BDS20Cap           bool   `json:"bds_20_cap"`
	BDS21Cap           bool   `json:"bds_21_cap"`
	BDS40Cap           bool   `json:"bds_40_cap"`
	BDS41Cap           bool   `json:"bds_41_cap"`
	BDS42Cap           bool   `json:"bds_42_cap"`
	BDS43Cap           bool   `json:"bds_43_cap"`
	BDS44Cap           bool   `json:"bds_44_cap"`
	BDS45Cap           bool   `json:"bds_45_cap"`
	BDS48Cap           bool   `json:"bds_48_cap"`
	BDS50Cap           bool   `json:"bds_50_cap"`
	BDS51Cap           bool   `json:"bds_51_cap"`
	BDS52Cap           bool   `json:"bds_52_cap"`
	BDS53Cap           bool   `json:"bds_53_cap"`
	BDS54Cap           bool   `json:"bds_54_cap"`
	BDS55Cap           bool   `json:"bds_55_cap"`
	BDS56Cap           bool   `json:"bds_56_cap"`
	BDS5FCap           bool   `json:"bds_5F_cap"`
	BDS60Cap           bool   `json:"bds_60_cap"`
}

type BDS17Decoded struct {
	BDS05Cap bool `json:"bds_05_cap"`
	BDS06Cap bool `json:"bds_06_cap"`
	BDS07Cap bool `json:"bds_07_cap"`
	BDS08Cap bool `json:"bds_08_cap"`
	BDS09Cap bool `json:"bds_09_cap"`
	BDS0ACap bool `json:"bds_0A_cap"`
	BDS20Cap bool `json:"bds_20_cap"`
	BDS21Cap bool `json:"bds_21_cap"`
	BDS40Cap bool `json:"bds_40_cap"`
	BDS41Cap bool `json:"bds_41_cap"`
	BDS42Cap bool `json:"bds_42_cap"`
	BDS43Cap bool `json:"bds_43_cap"`
	BDS44Cap bool `json:"bds_44_cap"`
	BDS45Cap bool `json:"bds_45_cap"`
	BDS48Cap bool `json:"bds_48_cap"`
	BDS50Cap bool `json:"bds_50_cap"`
	BDS51Cap bool `json:"bds_51_cap"`
	BDS52Cap bool `json:"bds_52_cap"`
	BDS53Cap bool `json:"bds_53_cap"`
	BDS54Cap bool `json:"bds_54_cap"`
	BDS55Cap bool `json:"bds_55_cap"`
	BDS56Cap bool `json:"bds_56_cap"`
	BDS5FCap bool `json:"bds_5F_cap"`
	BDS60Cap bool `json:"bds_60_cap"`
}

type BDS20Decoded struct {
	Callsign               string `json:"callsign"`
	AircraftIdentification string `json:"aircraft_identification"`
}

type BDS30Decoded struct {
	ThreatTypeIndicator          bool   `json:"threat_type_indicator"`
	ActiveResolutionAdvisory     uint32 `json:"active_resolution_advisory"`
	AraCorrectiveRa              bool   `json:"ara_corrective_ra"`
	AraDownwardSense             bool   `json:"ara_downward_sense"`
	AraIncreasedRate             bool   `json:"ara_increased_rate"`
	AraSenseReversal             bool   `json:"ara_sense_reversal"`
	AraAltitudeCrossing          bool   `json:"ara_altitude_crossing"`
	AraPositiveRa                bool   `json:"ara_positive_ra"`
	AraVerticalSpeedLimit        bool   `json:"ara_vertical_speed_limit"`
	ResolutionAdvisoryComplement uint32 `json:"resolution_advisory_complement"`
	RaTerminated                 bool   `json:"ra_terminated"`
	MultipleThreatEncounter      bool   `json:"multiple_threat_encounter"`
	ThreatTypeRaw                uint32 `json:"threat_type_raw"`
	ThreatType                   string `json:"threat_type"`
	ThreatIdentityData           uint32 `json:"threat_identity_data"`
	ThreatIdentityDataMid        uint32 `json:"threat_identity_data_mid"`
	ThreatIdentityDataLow        uint32 `json:"threat_identity_data_low"`
	ThreatModeSAddress           uint32 `json:"threat_mode_s_address,omitempty"`
}

type BDS40Decoded struct {
	SelectedAltitudeFt        float64 `json:"selected_altitude_ft,omitempty"`
	SelectedAltitudeValid     bool    `json:"selected_altitude_valid"`
	FmsAltitudeFt             float64 `json:"fms_altitude_ft,omitempty"`
	FmsAltitudeValid          bool    `json:"fms_altitude_valid"`
	BaroPressureMb            float64 `json:"baro_pressure_mb,omitempty"`
	BaroPressureValid         bool    `json:"baro_pressure_valid"`
	McpFcuModeValid           bool    `json:"mcp_fcu_mode_valid"`
	VnavMode                  bool    `json:"vnav_mode,omitempty"`
	AltHoldMode               bool    `json:"alt_hold_mode,omitempty"`
	ApproachMode              bool    `json:"approach_mode,omitempty"`
	TargetAltitudeSourceRaw   uint32  `json:"target_altitude_source_raw,omitempty"`
	TargetAltitudeSource      string  `json:"target_altitude_source,omitempty"`
	TargetAltitudeSourceValid bool    `json:"target_altitude_source_valid"`
}

type BDS44Decoded struct {
	WindSource                 string  `json:"wind_source"`
	WindSpeedKt                uint32  `json:"wind_speed_kt,omitempty"`
	WindDirectionDeg           float64 `json:"wind_direction_deg,omitempty"`
	WindValid                  bool    `json:"wind_valid"`
	StaticAirTemperatureC      float64 `json:"static_air_temperature_c,omitempty"`
	StaticAirTemperatureValid  bool    `json:"static_air_temperature_valid"`
	AverageStaticPressureHpa   float64 `json:"average_static_pressure_hpa,omitempty"`
	AverageStaticPressureValid bool    `json:"average_static_pressure_valid"`
	TurbulenceRaw              uint32  `json:"turbulence_raw,omitempty"`
	Turbulence                 string  `json:"turbulence,omitempty"`
	TurbulenceValid            bool    `json:"turbulence_valid"`
	HumidityPercent            float64 `json:"humidity_percent,omitempty"`
	HumidityValid              bool    `json:"humidity_valid"`
}

type BDS50Decoded struct {
	RollAngleDeg        float64 `json:"roll_angle_deg,omitempty"`
	RollAngleValid      bool    `json:"roll_angle_valid"`
	TrueTrackAngleDeg   float64 `json:"true_track_angle_deg,omitempty"`
	TrueTrackAngleValid bool    `json:"true_track_angle_valid"`
	GroundSpeedKt       float64 `json:"ground_speed_kt,omitempty"`
	GroundSpeedValid    bool    `json:"ground_speed_valid"`
	TrackAngleRateDegS  float64 `json:"track_angle_rate_deg_s,omitempty"`
	TrackAngleRateValid bool    `json:"track_angle_rate_valid"`
	TrueAirspeedKt      float64 `json:"true_airspeed_kt,omitempty"`
	TrueAirspeedValid   bool    `json:"true_airspeed_valid"`
}

type BDS60Decoded struct {
	MagneticHeadingDeg            float64 `json:"magnetic_heading_deg,omitempty"`
	MagneticHeadingValid          bool    `json:"magnetic_heading_valid"`
	IndicatedAirspeedKt           float64 `json:"indicated_airspeed_kt,omitempty"`
	IndicatedAirspeedValid        bool    `json:"indicated_airspeed_valid"`
	MachNumber                    float64 `json:"mach_number,omitempty"`
	MachNumberValid               bool    `json:"mach_number_valid"`
	BaroAltitudeRateFpm           float64 `json:"baro_altitude_rate_fpm,omitempty"`
	BaroAltitudeRateValid         bool    `json:"baro_altitude_rate_valid"`
	InertialVerticalVelocityFpm   float64 `json:"inertial_vertical_velocity_fpm,omitempty"`
	InertialVerticalVelocityValid bool    `json:"inertial_vertical_velocity_valid"`
}

// Decoder is the interface for BDS register decoders that return typed structs.
type Decoder interface {
	BDSCode() (uint8, uint8)
	Decode(data []byte) (any, error)
}

// Registry of BDS decoders
var decoders = map[uint8]Decoder{
	0x10: &BDS10Decoder{}, // BDS 1,0 - Data Link Capability Report
	0x17: &BDS17Decoder{}, // BDS 1,7 - Common Usage GICB Capability Report
	0x20: &BDS20Decoder{}, // BDS 2,0 - Aircraft Identification
	0x30: &BDS30Decoder{}, // BDS 3,0 - ACAS Active Resolution Advisory
	0x40: &BDS40Decoder{}, // BDS 4,0 - Selected Vertical Intention
	0x44: &BDS44Decoder{}, // BDS 4,4 - Meteorological Routine Air Report
	0x50: &BDS50Decoder{}, // BDS 5,0 - Track and Turn Report
	0x60: &BDS60Decoder{}, // BDS 6,0 - Heading and Speed Report
}

// Decode decodes a BDS register and returns a typed struct.
func Decode(bds1, bds2 uint8, data []byte) (any, error) {
	if len(data) < 7 {
		return nil, fmt.Errorf("BDS data too short: %d bytes, need 7", len(data))
	}

	bdsCode := (bds1 << 4) | bds2
	decoder, ok := decoders[bdsCode]
	if !ok {
		return nil, fmt.Errorf("unknown BDS code %d,%d", bds1, bds2)
	}

	return decoder.Decode(data)
}

// Helper functions for bit extraction

// ExtractBits extracts bits from a byte array
func ExtractBits(data []byte, startBit, numBits int) uint32 {
	if numBits > 32 {
		return 0
	}

	result := uint32(0)
	for i := 0; i < numBits; i++ {
		bitPos := startBit + i
		bytePos := bitPos / 8
		bitInByte := 7 - (bitPos % 8)

		if bytePos < len(data) {
			if (data[bytePos] & (1 << bitInByte)) != 0 {
				result |= 1 << (numBits - 1 - i)
			}
		}
	}
	return result
}

// ExtractSignedBits extracts signed bits (two's complement)
func ExtractSignedBits(data []byte, startBit, numBits int) int32 {
	value := ExtractBits(data, startBit, numBits)

	// Check if sign bit is set
	if value&(1<<(numBits-1)) != 0 {
		// Extend sign bit
		mask := ^((1 << numBits) - 1)
		return int32(value | uint32(mask))
	}
	return int32(value)
}
