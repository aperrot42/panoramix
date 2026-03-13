package asterix

import (
	"fmt"
	"time"
)

// DataSourceIdentifier represents I048/010
type DataSourceIdentifier struct {
	SAC uint8 `json:"sac"`
	SIC uint8 `json:"sic"`
}

// TimeOfDay represents I048/140
type TimeOfDay struct {
	Duration time.Duration `json:"duration"` // time since midnight UTC (raw * 1/128 s)
}

// TargetReportDescriptor is defined in cat048.go (already a named type)

// MeasuredPositionPolar represents I048/040
type MeasuredPositionPolar struct {
	Rho   float64 `json:"rho"`   // range in NM (raw * 1/256)
	Theta float64 `json:"theta"` // azimuth in degrees (raw * 360/2^16)
}

// Mode3ACode represents I048/070
type Mode3ACode struct {
	Validated bool   `json:"validated"`
	Garbled   bool   `json:"garbled"`
	Local     bool   `json:"local"`
	Code      uint16 `json:"code"` // 12-bit transponder code
}

func (c Mode3ACode) OctalString() string { return fmt.Sprintf("%04o", c.Code) }

// FlightLevel represents I048/090
type FlightLevel struct {
	FL        float64 `json:"fl"` // flight level (raw * 1/4)
	Validated bool    `json:"validated"`
	Garbled   bool    `json:"garbled"`
}

// RadarPlotCharacteristics is defined in cat048.go (already a named type)

// CalculatedPositionCartesian represents I048/042
type CalculatedPositionCartesian struct {
	X float64 `json:"x"` // NM (raw * 1/128)
	Y float64 `json:"y"` // NM (raw * 1/128)
}

// Mode2Code represents I048/050
type Mode2Code struct {
	Validated bool   `json:"validated"`
	Garbled   bool   `json:"garbled"`
	Local     bool   `json:"local"`
	Code      uint16 `json:"code"` // 12-bit transponder code
}

func (c Mode2Code) OctalString() string { return fmt.Sprintf("%04o", c.Code) }

// Mode1Code represents I048/055
type Mode1Code struct {
	Validated bool  `json:"validated"`
	Garbled   bool  `json:"garbled"`
	Local     bool  `json:"local"`
	Code      uint8 `json:"code"` // 5-bit code
}

func (c Mode1Code) OctalString() string { return fmt.Sprintf("%02o", c.Code) }

// Mode2CodeConfidence represents I048/060
type Mode2CodeConfidence struct {
	QA4 bool `json:"qa4"`
	QA2 bool `json:"qa2"`
	QA1 bool `json:"qa1"`
	QB4 bool `json:"qb4"`
	QB2 bool `json:"qb2"`
	QB1 bool `json:"qb1"`
	QC4 bool `json:"qc4"`
	QC2 bool `json:"qc2"`
	QC1 bool `json:"qc1"`
	QD4 bool `json:"qd4"`
	QD2 bool `json:"qd2"`
	QD1 bool `json:"qd1"`
}

// Mode1CodeConfidence represents I048/065
type Mode1CodeConfidence struct {
	QA4 bool `json:"qa4"`
	QA2 bool `json:"qa2"`
	QA1 bool `json:"qa1"`
	QB2 bool `json:"qb2"`
	QB1 bool `json:"qb1"`
}

// Mode3ACodeConfidence represents I048/080
type Mode3ACodeConfidence struct {
	QA4 bool `json:"qa4"`
	QA2 bool `json:"qa2"`
	QA1 bool `json:"qa1"`
	QB4 bool `json:"qb4"`
	QB2 bool `json:"qb2"`
	QB1 bool `json:"qb1"`
	QC4 bool `json:"qc4"`
	QC2 bool `json:"qc2"`
	QC1 bool `json:"qc1"`
	QD4 bool `json:"qd4"`
	QD2 bool `json:"qd2"`
	QD1 bool `json:"qd1"`
}

// ModeCCodeConfidence represents I048/100
type ModeCCodeConfidence struct {
	Validated bool   `json:"validated"`
	Garbled   bool   `json:"garbled"`
	Code      uint16 `json:"code"`
	QC1       bool   `json:"qc1"`
	QA1       bool   `json:"qa1"`
	QC2       bool   `json:"qc2"`
	QA2       bool   `json:"qa2"`
	QC4       bool   `json:"qc4"`
	QA4       bool   `json:"qa4"`
	QB1       bool   `json:"qb1"`
	QD1       bool   `json:"qd1"`
	QB2       bool   `json:"qb2"`
	QD2       bool   `json:"qd2"`
	QB4       bool   `json:"qb4"`
	QD4       bool   `json:"qd4"`
}

// HeightMeasured3D represents I048/110
type HeightMeasured3D struct {
	Height float64 `json:"height"` // feet (raw * 25)
}

// RadialDopplerSpeed represents I048/120
type RadialDopplerSpeed struct {
	RawData []byte `json:"raw_data"` // compound data item, raw bytes
}

// TrackNumber represents I048/161
type TrackNumber struct {
	Number uint16 `json:"number"` // 12-bit track number
}

// TrackStatus represents I048/170
type TrackStatus struct {
	CNF          bool `json:"cnf"`           // Confirmed vs Tentative Track
	RAD          byte `json:"rad"`           // Type of Sensor(s) maintaining Track
	DOU          bool `json:"dou"`           // Confidence in plot to track association
	MAH          bool `json:"mah"`           // Manoeuvre detection in Horizontal Sense
	CDM          byte `json:"cdm"`           // Climbing/Descending Mode
	HasExtension bool `json:"-"`             // whether second octet was present
	TRE          bool `json:"tre,omitempty"` // Signal for End_of_Track
	GHO          bool `json:"gho,omitempty"` // Ghost vs. true target
	SUP          bool `json:"sup,omitempty"` // Track maintained with neighbouring info
	TCC          bool `json:"tcc,omitempty"` // Type of plot coordinate transformation
}

// CalculatedTrackVelocity represents I048/200
type CalculatedTrackVelocity struct {
	Groundspeed float64 `json:"groundspeed"` // NM/s (raw * 2^-14), ~0.22 kt per unit
	Heading     float64 `json:"heading"`     // degrees from geographic north (raw * 360/2^16)
}

// TrackQuality represents I048/210
type TrackQuality struct {
	SigmaX float64 `json:"sigma_x"` // NM (raw * 1/128)
	SigmaY float64 `json:"sigma_y"` // NM (raw * 1/128)
	SigmaV float64 `json:"sigma_v"` // NM/s (raw * 2^-14)
	SigmaH float64 `json:"sigma_h"` // degrees (raw * 360/2^12)
}

// CommunicationsCapability represents I048/230
type CommunicationsCapability struct {
	COM  byte `json:"com"`  // Communications capability
	STAT byte `json:"stat"` // Flight Status
	SI   bool `json:"si"`   // SI/II Transponder Capability
	MSSC bool `json:"mssc"` // Mode-S Specific Service Capability
	ARC  bool `json:"arc"`  // Altitude reporting capability
	AIC  bool `json:"aic"`  // Aircraft identification capability
	B1A  bool `json:"b1a"`  // BDS 1,0 bit 16
	B1B  byte `json:"b1b"`  // BDS 1,0 bits 37/40
}

// WarningErrorConditions represents I048/030
type WarningErrorConditions struct {
	Codes []uint8 `json:"codes"`
}

// ACASResolutionAdvisory represents I048/260
type ACASResolutionAdvisory struct {
	ACASRA string `json:"acas_ra"` // hex-encoded 7-byte data
}

// BDSRegisterData represents I048/250 Mode S MB Data
type BDSRegisterData struct {
	Repetition uint8                  `json:"repetition"`
	Registers  map[string]BDSRegister `json:"registers"`
}

// BDSRegister represents a single BDS register entry
type BDSRegister struct {
	BDSCode    byte   `json:"bds_code"`         // raw register address byte (e.g. 0x50 for BDS 5,0)
	RawData    []byte `json:"-"`                 // raw 7-byte payload (not serialized, for pipeline filters)
	BDSDataRaw string `json:"bds_data_raw"`     // hex-encoded 7-byte payload
	Decoded    any    `json:"decoded,omitempty"` // populated by BDS enrichment filter
	Error      string `json:"error,omitempty"`   // decoding error message, if any
}
