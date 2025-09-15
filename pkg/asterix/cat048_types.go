package asterix

import "time"

// CAT048Message represents a complete decoded CAT 048 ASTERIX message
type CAT048Message struct {
	Category byte   `json:"category"`
	SAC      uint8  `json:"sac"`
	SIC      uint8  `json:"sic"`
	FSPEC    []byte `json:"fspec,omitempty"`
	
	// Data Items (following ASTERIX CAT 048 specification)
	DataSourceIdentifier       *DataSourceIdentifier048       `json:"data_source_identifier,omitempty"`       // I048/010
	TimeOfDay                  *time.Duration                 `json:"time_of_day,omitempty"`                  // I048/140
	TargetReportDescriptor     *TargetReportDescriptor048     `json:"target_report_descriptor,omitempty"`     // I048/020
	MeasuredPosition           *MeasuredPosition048           `json:"measured_position,omitempty"`            // I048/040
	AircraftIdentification     *string                        `json:"aircraft_identification,omitempty"`     // I048/240
	ModeACode                  *ModeACode048                  `json:"mode_a_code,omitempty"`                  // I048/070
	FlightLevel                *FlightLevel048                `json:"flight_level,omitempty"`                 // I048/090
	RadarPlotCharacteristics   *RadarPlotCharacteristics048   `json:"radar_plot_characteristics,omitempty"`   // I048/130
	AircraftDerivedData        *AircraftDerivedData048        `json:"aircraft_derived_data,omitempty"`        // I048/220
	TargetIdentification       *TargetIdentification048       `json:"target_identification,omitempty"`       // I048/240
	ModeSMBData                *ModeSMBData048                `json:"mode_s_mb_data,omitempty"`               // I048/250
	CalculatedPosition         *CalculatedPosition048         `json:"calculated_position,omitempty"`          // I048/042
	CalculatedPositionCartesian *CalculatedPositionCartesian048 `json:"calculated_position_cartesian,omitempty"` // I048/041
	Mode3ACodeConfidence       *Mode3ACodeConfidence048       `json:"mode_3a_code_confidence,omitempty"`      // I048/080
	ModeC                      *ModeC048                      `json:"mode_c,omitempty"`                       // I048/100
	ModeCConfidence            *ModeCConfidence048            `json:"mode_c_confidence,omitempty"`            // I048/110
	HeightMeasured3D           *HeightMeasured3D048           `json:"height_measured_3d,omitempty"`           // I048/120
	RadialDopplerSpeed         *RadialDopplerSpeed048         `json:"radial_doppler_speed,omitempty"`         // I048/200
	CommunicationsCapability   *CommunicationsCapability048   `json:"communications_capability,omitempty"`    // I048/230
	SystemStatus               *SystemStatus048               `json:"system_status,omitempty"`                // I048/260
	WarningErrorConditions     *WarningErrorConditions048     `json:"warning_error_conditions,omitempty"`     // I048/030
	ReservedExpansion          *ReservedField048              `json:"reserved_expansion,omitempty"`           // I048/RE
	SpecialPurpose             *SpecialPurposeField048        `json:"special_purpose,omitempty"`              // I048/SP
}

// I048/010 - Data Source Identifier
type DataSourceIdentifier048 struct {
	SAC uint8 `json:"sac"` // System Area Code
	SIC uint8 `json:"sic"` // System Identification Code
}

// I048/040 - Measured Position in Slant Polar Coordinates
type MeasuredPosition048 struct {
	Rho      float64 `json:"rho_nm"`      // Range in nautical miles
	RhoRaw   uint16  `json:"rho_raw"`     // Raw range value (LSB = 1/256 NM)
	Theta    float64 `json:"theta_deg"`   // Azimuth in degrees
	ThetaRaw uint16  `json:"theta_raw"`   // Raw azimuth value (LSB = 360°/2^16)
}

// I048/020 - Target Report Descriptor
type TargetReportDescriptor048 struct {
	TYP  uint8 `json:"typ"`  // Type of detection (3 bits)
	SIM  bool  `json:"sim"`  // Simulated target report
	RDP  bool  `json:"rdp"`  // RDP Chain 1 or 2
	SPI  bool  `json:"spi"`  // Special Position Identification
	RAB  bool  `json:"rab"`  // Report from aircraft or field monitor
	TST  bool  `json:"tst"`  // Test target
	ERR  bool  `json:"err"`  // Extended Range present
	XPP  bool  `json:"xpp"`  // X-Pulse present
	ME   bool  `json:"me"`   // Military Emergency
	MI   bool  `json:"mi"`   // Military Identification
	FOEFRI uint8 `json:"foe_fri"` // Foe/Friend Identification (2 bits)
	FX   bool  `json:"fx"`   // Field extension
}

// I048/070 - Mode-3/A Code
type ModeACode048 struct {
	V    bool   `json:"v"`    // Validated
	G    bool   `json:"g"`    // Garbled
	L    bool   `json:"l"`    // Smoothed
	Mode3A uint16 `json:"mode_3a"` // Mode-3/A reply (12 bits)
}

// I048/090 - Flight Level
type FlightLevel048 struct {
	V    bool    `json:"v"`    // Validated
	G    bool    `json:"g"`    // Garbled
	Code uint16  `json:"code"` // Flight Level in binary representation
	FlightLevelFt float64 `json:"flight_level_ft"` // Flight level in feet
}

// I048/130 - Radar Plot Characteristics
type RadarPlotCharacteristics048 struct {
	SRL *RadarPlotSRL048 `json:"srl,omitempty"` // SSR plot runlength
	SRR *uint8           `json:"srr,omitempty"` // Number of received replies for M(SSR)
	SAM *int8            `json:"sam,omitempty"` // Amplitude of M(SSR) reply
	PRL *RadarPlotPRL048 `json:"prl,omitempty"` // Primary plot runlength
	PAM *int8            `json:"pam,omitempty"` // Amplitude of Primary plot
	RPD *RadarPlotRPD048 `json:"rpd,omitempty"` // Difference in range
	APD *RadarPlotAPD048 `json:"apd,omitempty"` // Difference in azimuth
}

// SSR Plot Runlength subfield
type RadarPlotSRL048 struct {
	SRL uint8 `json:"srl"` // SSR plot runlength
}

// Primary Plot Runlength subfield
type RadarPlotPRL048 struct {
	PRL uint8 `json:"prl"` // Primary plot runlength
}

// Range Difference subfield
type RadarPlotRPD048 struct {
	RPD int8 `json:"rpd"` // Difference in range
}

// Azimuth Difference subfield
type RadarPlotAPD048 struct {
	APD int8 `json:"apd"` // Difference in azimuth
}

// I048/220 - Aircraft Derived Data
type AircraftDerivedData048 struct {
	ADR *AircraftAddress048         `json:"adr,omitempty"` // Aircraft address
	ID  *AircraftIdentification048  `json:"id,omitempty"`  // Aircraft identification
	MHG *MagneticHeading048         `json:"mhg,omitempty"` // Magnetic heading
	IAS *IndicatedAirspeed048       `json:"ias,omitempty"` // Indicated airspeed
	TAS *TrueAirspeed048            `json:"tas,omitempty"` // True airspeed
	SAL *SelectedAltitude048        `json:"sal,omitempty"` // Selected altitude
	FSS *FinalStateSelectedAlt048   `json:"fss,omitempty"` // Final state selected altitude
	TIS *TrajectoryIntentStatus048  `json:"tis,omitempty"` // Trajectory intent status
	TID *TrajectoryIntentData048    `json:"tid,omitempty"` // Trajectory intent data
	COM *CommunicationStatus048     `json:"com,omitempty"` // Communications/ACAS capability and flight status
	SAB *StatusReportedByADS048     `json:"sab,omitempty"` // Status reported by ADS-B
	ACS *ACASResolutionAdvisory048  `json:"acs,omitempty"` // ACAS resolution advisory report
	BVR *BarometricVerticalRate048  `json:"bvr,omitempty"` // Barometric vertical rate
	GVR *GeometricVerticalRate048   `json:"gvr,omitempty"` // Geometric vertical rate
	RAN *RollAngle048               `json:"ran,omitempty"` // Roll angle
	TAR *TrackAngleRate048          `json:"tar,omitempty"` // Track angle rate
	TAN *TrackAngle048              `json:"tan,omitempty"` // Track angle
	GSP *GroundSpeed048             `json:"gsp,omitempty"` // Ground speed
	VUN *VelocityUncertainty048     `json:"vun,omitempty"` // Velocity uncertainty
	MET *MeteorologicalData048      `json:"met,omitempty"` // Meteorological data
	EMC *EmitterCategory048         `json:"emc,omitempty"` // Emitter category
	POS *Position048                `json:"pos,omitempty"` // Position
	GAL *GeometricAltitude048       `json:"gal,omitempty"` // Geometric altitude
	PUN *PositionUncertainty048     `json:"pun,omitempty"` // Position uncertainty
	MB  *ModeSMBData048             `json:"mb,omitempty"`  // Mode S MB data
	IAR *IndicatedAirspeedRate048   `json:"iar,omitempty"` // Indicated airspeed rate
	MAC *MachNumber048              `json:"mac,omitempty"` // Mach number
	BPS *BarometricPressure048      `json:"bps,omitempty"` // Barometric pressure setting
}

// Aircraft Address subfield
type AircraftAddress048 struct {
	Address uint32 `json:"address"` // 24-bit aircraft address
}

// Aircraft Identification subfield  
type AircraftIdentification048 struct {
	Callsign string `json:"callsign"` // 8-character callsign
}

// Magnetic Heading subfield
type MagneticHeading048 struct {
	HeadingDeg float64 `json:"heading_deg"` // Magnetic heading in degrees
}

// Indicated Airspeed subfield
type IndicatedAirspeed048 struct {
	IAS uint16 `json:"ias"` // Indicated airspeed
}

// True Airspeed subfield
type TrueAirspeed048 struct {
	TAS uint16 `json:"tas"` // True airspeed
}

// Selected Altitude subfield
type SelectedAltitude048 struct {
	SAS   bool    `json:"sas"`   // Source availability
	Source uint8   `json:"source"` // Source (2 bits)
	Alt   uint16  `json:"alt"`   // Altitude
	AltFt float64 `json:"alt_ft"` // Altitude in feet
}

// Final State Selected Altitude subfield
type FinalStateSelectedAlt048 struct {
	MV    bool    `json:"mv"`    // Manage vertical mode
	AH    bool    `json:"ah"`    // Altitude hold mode  
	AM    bool    `json:"am"`    // Approach mode
	Alt   uint16  `json:"alt"`   // Altitude
	AltFt float64 `json:"alt_ft"` // Altitude in feet
}

// Trajectory Intent Status subfield
type TrajectoryIntentStatus048 struct {
	NAV bool `json:"nav"` // LNAV mode
	NVB bool `json:"nvb"` // LNAV/VNAV mode
}

// Trajectory Intent Data subfield
type TrajectoryIntentData048 struct {
	TCA bool   `json:"tca"` // TCP number availability
	NC  bool   `json:"nc"`  // TCP compliance
	TCP uint8  `json:"tcp"` // TCP number (6 bits)
	Alt uint16 `json:"alt"` // Altitude
	Lat float64 `json:"lat"` // Latitude
	Lon float64 `json:"lon"` // Longitude
	PT  uint8  `json:"pt"`  // Point type (4 bits)
	TD  uint8  `json:"td"`  // TD (2 bits)
	TRA bool   `json:"tra"` // Turn radius availability
	TOA bool   `json:"toa"` // TOA availability
	TOV uint16 `json:"tov"` // Time over point
	TTR uint16 `json:"ttr"` // TCP turn radius
}

// Communication Status subfield
type CommunicationStatus048 struct {
	COM  uint8 `json:"com"`  // Communications capability (3 bits)
	STAT uint8 `json:"stat"` // Flight status (3 bits)
	SI   bool  `json:"si"`   // SI/II transponder capability
	MSSC bool  `json:"mssc"` // Mode S specific service capability
	ARC  bool  `json:"arc"`  // Altitude reporting capability
	AIC  bool  `json:"aic"`  // Aircraft identification capability
	B1A  bool  `json:"b1a"`  // BDS 1,0 bit 16
	B1B  uint8 `json:"b1b"`  // BDS 1,0 bits 37/40 (4 bits)
}

// Status Reported by ADS-B subfield
type StatusReportedByADS048 struct {
	AC  uint8 `json:"ac"`  // Aircraft operational status (2 bits)
	MN  uint8 `json:"mn"`  // Maximum airspeed (2 bits)
	DC  uint8 `json:"dc"`  // Data correspondence (2 bits)
	GBS bool  `json:"gbs"` // Ground bit set
	CFS bool  `json:"cfs"` // Common flag set
	STAT uint8 `json:"stat"` // Flight status (3 bits)
}

// ACAS Resolution Advisory Report subfield
type ACASResolutionAdvisory048 struct {
	COM  uint8  `json:"com"`  // Communications capability (3 bits)
	STAT uint8  `json:"stat"` // Flight status (3 bits)
	SI   bool   `json:"si"`   // SI/II transponder capability
	MSSC bool   `json:"mssc"` // Mode S specific service capability
	ARC  bool   `json:"arc"`  // Altitude reporting capability
	AIC  bool   `json:"aic"`  // Aircraft identification capability
	B1A  bool   `json:"b1a"`  // BDS 1,0 bit 16
	B1B  uint8  `json:"b1b"`  // BDS 1,0 bits 37/40 (4 bits)
	AC   uint8  `json:"ac"`   // Aircraft operational status (2 bits)
	MN   uint8  `json:"mn"`   // Maximum airspeed (2 bits)
	DC   uint8  `json:"dc"`   // Data correspondence (2 bits)
	GBS  bool   `json:"gbs"`  // Ground bit set
	CFS  bool   `json:"cfs"`  // Common flag set
	STAT2 uint8 `json:"stat2"` // Flight status (3 bits)
	ARA  []byte `json:"ara"`  // Active resolution advisories
	RAC  uint8  `json:"rac"`  // RA complement (4 bits)
	RAT  bool   `json:"rat"`  // RA terminated
	MTE  bool   `json:"mte"`  // Multiple threat encounter
	TTI  uint8  `json:"tti"`  // Threat type indicator (2 bits)
	TID  uint32 `json:"tid"`  // Threat identity data (26 bits)
}

// Barometric Vertical Rate subfield
type BarometricVerticalRate048 struct {
	RE    bool    `json:"re"`    // Range exceeded indicator
	Rate  int16   `json:"rate"`  // Barometric vertical rate
	RateFtMin float64 `json:"rate_ft_min"` // Rate in feet per minute
}

// Geometric Vertical Rate subfield
type GeometricVerticalRate048 struct {
	RE    bool    `json:"re"`    // Range exceeded indicator  
	Rate  int16   `json:"rate"`  // Geometric vertical rate
	RateFtMin float64 `json:"rate_ft_min"` // Rate in feet per minute
}

// Roll Angle subfield
type RollAngle048 struct {
	AngleDeg float64 `json:"angle_deg"` // Roll angle in degrees
}

// Track Angle Rate subfield
type TrackAngleRate048 struct {
	RateDegS float64 `json:"rate_deg_s"` // Track angle rate in degrees per second
}

// Track Angle subfield
type TrackAngle048 struct {
	AngleDeg float64 `json:"angle_deg"` // Track angle in degrees
}

// Ground Speed subfield
type GroundSpeed048 struct {
	RE       bool    `json:"re"`       // Range exceeded indicator
	Speed    uint16  `json:"speed"`    // Ground speed
	SpeedKt  float64 `json:"speed_kt"` // Speed in knots
}

// Velocity Uncertainty subfield
type VelocityUncertainty048 struct {
	Uncertainty uint8 `json:"uncertainty"` // Velocity uncertainty
}

// Meteorological Data subfield
type MeteorologicalData048 struct {
	WS  *WindSpeed048     `json:"ws,omitempty"`  // Wind speed
	WD  *WindDirection048 `json:"wd,omitempty"`  // Wind direction  
	TMP *Temperature048   `json:"tmp,omitempty"` // Temperature
	TRB *Turbulence048    `json:"trb,omitempty"` // Turbulence
}

// Wind Speed subfield
type WindSpeed048 struct {
	SpeedKt float64 `json:"speed_kt"` // Wind speed in knots
}

// Wind Direction subfield  
type WindDirection048 struct {
	DirectionDeg float64 `json:"direction_deg"` // Wind direction in degrees
}

// Temperature subfield
type Temperature048 struct {
	TempC float64 `json:"temp_c"` // Temperature in Celsius
}

// Turbulence subfield
type Turbulence048 struct {
	Turbulence uint8 `json:"turbulence"` // Turbulence level
}

// Emitter Category subfield
type EmitterCategory048 struct {
	Category uint8 `json:"category"` // Emitter category
}

// Position subfield
type Position048 struct {
	LatitudeDeg  float64 `json:"latitude_deg"`  // Latitude in degrees
	LongitudeDeg float64 `json:"longitude_deg"` // Longitude in degrees
}

// Geometric Altitude subfield
type GeometricAltitude048 struct {
	AltitudeFt float64 `json:"altitude_ft"` // Geometric altitude in feet
}

// Position Uncertainty subfield  
type PositionUncertainty048 struct {
	Uncertainty uint8 `json:"uncertainty"` // Position uncertainty
}

// Indicated Airspeed Rate subfield
type IndicatedAirspeedRate048 struct {
	Rate int16 `json:"rate"` // Indicated airspeed rate
}

// Mach Number subfield
type MachNumber048 struct {
	Mach float64 `json:"mach"` // Mach number
}

// Barometric Pressure Setting subfield
type BarometricPressure048 struct {
	PressureHPa float64 `json:"pressure_hpa"` // Pressure in hPa
}

// I048/240 - Target Identification
type TargetIdentification048 struct {
	Identification string `json:"identification"` // 8-character identification
}

// I048/250 - Mode S MB Data
type ModeSMBData048 struct {
	RepetitionFactor uint8            `json:"repetition_factor"` // Number of MB data
	MBData           []ModeSMBItem048 `json:"mb_data"`           // List of MB data items
}

// Individual Mode S MB Data Item
type ModeSMBItem048 struct {
	Data   []byte `json:"data"`   // 8-byte MB data
	BDS1   uint8  `json:"bds1"`   // BDS register 1 (4 bits)
	BDS2   uint8  `json:"bds2"`   // BDS register 2 (4 bits)
	DataHex string `json:"data_hex"` // Hex representation
}

// I048/042 - Calculated Position in WGS-84 Co-ordinates
type CalculatedPosition048 struct {
	LatitudeDeg  float64 `json:"latitude_deg"`  // Latitude in degrees
	LongitudeDeg float64 `json:"longitude_deg"` // Longitude in degrees
}

// I048/041 - Calculated Position in Cartesian Co-ordinates  
type CalculatedPositionCartesian048 struct {
	X float64 `json:"x"` // X coordinate in meters
	Y float64 `json:"y"` // Y coordinate in meters
}

// I048/080 - Mode-3/A Code Confidence Indicator
type Mode3ACodeConfidence048 struct {
	QA4 bool `json:"qa4"` // Confidence in bit QA4 of Mode-3/A reply
	QA2 bool `json:"qa2"` // Confidence in bit QA2 of Mode-3/A reply  
	QA1 bool `json:"qa1"` // Confidence in bit QA1 of Mode-3/A reply
	QB4 bool `json:"qb4"` // Confidence in bit QB4 of Mode-3/A reply
	QB2 bool `json:"qb2"` // Confidence in bit QB2 of Mode-3/A reply
	QB1 bool `json:"qb1"` // Confidence in bit QB1 of Mode-3/A reply
	QC4 bool `json:"qc4"` // Confidence in bit QC4 of Mode-3/A reply
	QC2 bool `json:"qc2"` // Confidence in bit QC2 of Mode-3/A reply
	QC1 bool `json:"qc1"` // Confidence in bit QC1 of Mode-3/A reply
	QD4 bool `json:"qd4"` // Confidence in bit QD4 of Mode-3/A reply
	QD2 bool `json:"qd2"` // Confidence in bit QD2 of Mode-3/A reply
	QD1 bool `json:"qd1"` // Confidence in bit QD1 of Mode-3/A reply
}

// I048/100 - Mode-C Code and Code Confidence Indicator
type ModeC048 struct {
	V          bool    `json:"v"`           // Validated
	G          bool    `json:"g"`           // Garbled  
	Code       uint16  `json:"code"`        // Mode-C reply in Gray notation
	QC1        bool    `json:"qc1"`         // Confidence in bit QC1
	QA1        bool    `json:"qa1"`         // Confidence in bit QA1
	QC2        bool    `json:"qc2"`         // Confidence in bit QC2
	QA2        bool    `json:"qa2"`         // Confidence in bit QA2
	QC4        bool    `json:"qc4"`         // Confidence in bit QC4
	QA4        bool    `json:"qa4"`         // Confidence in bit QA4
	QB1        bool    `json:"qb1"`         // Confidence in bit QB1
	QD1        bool    `json:"qd1"`         // Confidence in bit QD1
	QB2        bool    `json:"qb2"`         // Confidence in bit QB2
	QD2        bool    `json:"qd2"`         // Confidence in bit QD2
	QB4        bool    `json:"qb4"`         // Confidence in bit QB4
	QD4        bool    `json:"qd4"`         // Confidence in bit QD4
	FlightLevelFt float64 `json:"flight_level_ft"` // Flight level in feet
}

// I048/110 - Mode-C Code Confidence Indicator  
type ModeCConfidence048 struct {
	QC1 bool `json:"qc1"` // Confidence in bit QC1
	QA1 bool `json:"qa1"` // Confidence in bit QA1
	QC2 bool `json:"qc2"` // Confidence in bit QC2
	QA2 bool `json:"qa2"` // Confidence in bit QA2
	QC4 bool `json:"qc4"` // Confidence in bit QC4
	QA4 bool `json:"qa4"` // Confidence in bit QA4
	QB1 bool `json:"qb1"` // Confidence in bit QB1
	QD1 bool `json:"qd1"` // Confidence in bit QD1
	QB2 bool `json:"qb2"` // Confidence in bit QB2
	QD2 bool `json:"qd2"` // Confidence in bit QD2
	QB4 bool `json:"qb4"` // Confidence in bit QB4
	QD4 bool `json:"qd4"` // Confidence in bit QD4
}

// I048/120 - Measured Height
type HeightMeasured3D048 struct {
	HeightFt float64 `json:"height_ft"` // 3D height in feet
}

// I048/200 - Calculated Track Velocity in Polar Co-ordinates
type RadialDopplerSpeed048 struct {
	SpeedKt float64 `json:"speed_kt"` // Calculated groundspeed in knots
	HeadingDeg float64 `json:"heading_deg"` // Calculated heading in degrees
}

// I048/230 - Communications/ACAS Capability and Flight Status
type CommunicationsCapability048 struct {
	COM  uint8 `json:"com"`  // Communications capability (3 bits)
	STAT uint8 `json:"stat"` // Flight status (3 bits)  
	SI   bool  `json:"si"`   // SI/II transponder capability
	MSSC bool  `json:"mssc"` // Mode S specific service capability
	ARC  bool  `json:"arc"`  // Altitude reporting capability
	AIC  bool  `json:"aic"`  // Aircraft identification capability
	B1A  bool  `json:"b1a"`  // BDS 1,0 bit 16
	B1B  uint8 `json:"b1b"`  // BDS 1,0 bits 37/40 (4 bits)
}

// I048/260 - ACAS Resolution Advisory Report
type SystemStatus048 struct {
	NOGO   bool `json:"nogo"`    // Operational release status of the system
	OVL    bool `json:"ovl"`     // Overload indicator
	TSV    bool `json:"tsv"`     // Time source validity
	DIV    bool `json:"div"`     // Diversity
	TTF    bool `json:"ttf"`     // Transponder test failure
}

// I048/030 - Warning/Error Conditions and Target Classification
type WarningErrorConditions048 struct {
	WarningErrors uint8 `json:"warning_errors"` // Warning/Error condition value
}

// I048/RE - Reserved Expansion Field
type ReservedField048 struct {
	Note string `json:"note"` // Description
	Data string `json:"data"` // Hex representation of data
}

// I048/SP - Special Purpose Field
type SpecialPurposeField048 struct {
	Note string `json:"note"` // Description  
	Data string `json:"data"` // Hex representation of data
}