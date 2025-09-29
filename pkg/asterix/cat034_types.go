package asterix

import "time"

// CAT034Message represents a complete decoded CAT 034 ASTERIX message
type CAT034Message struct {
	Category byte   `json:"category"`
	SAC      uint8  `json:"sac"`
	SIC      uint8  `json:"sic"`
	FSPEC    []byte `json:"fspec,omitempty"`

	// Data Items (following ASTERIX CAT 034 specification)
	DataSourceIdentifier *DataSourceIdentifier034 `json:"data_source_identifier,omitempty"` // I034/010
	MessageType          *uint8                   `json:"message_type,omitempty"`           // I034/000
	TimeOfDay            *time.Duration           `json:"time_of_day,omitempty"`            // I034/030
	SectorNumber         *SectorNumber034         `json:"sector_number,omitempty"`          // I034/020
	AntennaRotationSpeed *AntennaRotationSpeed034 `json:"antenna_rotation_speed,omitempty"` // I034/041
	SystemConfiguration  *SystemConfiguration034  `json:"system_configuration,omitempty"`   // I034/050
	SystemProcessingMode *SystemProcessingMode034 `json:"system_processing_mode,omitempty"` // I034/060
	MessageCountValues   *MessageCountValues034   `json:"message_count_values,omitempty"`   // I034/070
	GenericPolarWindow   *GenericPolarWindow034   `json:"generic_polar_window,omitempty"`   // I034/100
	DataFilter           *DataFilter034           `json:"data_filter,omitempty"`            // I034/110
	Position3D           *Position3D034           `json:"position_3d,omitempty"`            // I034/120
	CollimationError     *CollimationError034     `json:"collimation_error,omitempty"`      // I034/090
	ReservedExpansion    *ReservedField034        `json:"reserved_expansion,omitempty"`     // I034/RE
	SpecialPurpose       *SpecialPurposeField034  `json:"special_purpose,omitempty"`        // I034/SP
}

// I034/010 - Data Source Identifier
type DataSourceIdentifier034 struct {
	SAC uint8 `json:"sac"` // System Area Code
	SIC uint8 `json:"sic"` // System Identification Code
}

// I034/020 - Sector Number
type SectorNumber034 struct {
	Sector     uint8   `json:"sector"`      // Raw sector number (0-255)
	AzimuthDeg float64 `json:"azimuth_deg"` // Azimuth in degrees (0-360)
}

// I034/041 - Antenna Rotation Speed
type AntennaRotationSpeed034 struct {
	RotationPeriodS float64 `json:"rotation_period_s"` // Period in seconds
	RawValue        uint16  `json:"raw_value"`         // Raw value in 1/128 second units
}

// I034/050 - System Configuration and Status
type SystemConfiguration034 struct {
	COM bool `json:"com"` // Common Part present
	PSR bool `json:"psr"` // PSR Sensor present
	SSR bool `json:"ssr"` // SSR Sensor present
	MDS bool `json:"mds"` // Mode S Sensor present
	FX  bool `json:"fx"`  // Field Extension

	// Conditional subfields
	COMData *COMConfigData034 `json:"com_data,omitempty"`
	PSRData *PSRConfigData034 `json:"psr_data,omitempty"`
	SSRData *SSRConfigData034 `json:"ssr_data,omitempty"`
	MDSData *MDSConfigData034 `json:"mds_data,omitempty"`
}

// COM Configuration Data
type COMConfigData034 struct {
	NOGO   bool `json:"nogo"`    // Operational Release Status
	RDPC   bool `json:"rdpc"`    // Radar Data Processor Chain
	RDPR   bool `json:"rdpr"`    // Radar Data Processor Ready
	OVLRDP bool `json:"ovl_rdp"` // Radar Data Processor Overload
	OVLXMT bool `json:"ovl_xmt"` // Transmission Subsystem Overload
	MSC    bool `json:"msc"`     // Monitoring System Connected
	TSV    bool `json:"tsv"`     // Time Source Validity
}

// PSR Configuration Data
type PSRConfigData034 struct {
	ANT  bool  `json:"ant"`   // Selected antenna
	CHAB uint8 `json:"ch_ab"` // Channel A/B selection (2 bits)
	OVL  bool  `json:"ovl"`   // Overload condition
	MSC  bool  `json:"msc"`   // Monitoring System Connected
}

// SSR Configuration Data
type SSRConfigData034 struct {
	ANT  bool  `json:"ant"`   // Selected antenna
	CHAB uint8 `json:"ch_ab"` // Channel A/B selection (2 bits)
	OVL  bool  `json:"ovl"`   // Overload condition
	MSC  bool  `json:"msc"`   // Monitoring System Connected
}

// MDS Configuration Data
type MDSConfigData034 struct {
	ANT    bool  `json:"ant"`     // Selected antenna
	CHAB   uint8 `json:"ch_ab"`   // Channel A/B selection (2 bits)
	OVLSUR bool  `json:"ovl_sur"` // Surveillance Overload
	MSC    bool  `json:"msc"`     // Monitoring System Connected
	SCF    bool  `json:"scf"`     // Channel A/B selection for Surveillance Co-ordination Function
	DLF    bool  `json:"dlf"`     // Channel A/B selection for Data Link Function
	OVLSCF bool  `json:"ovl_scf"` // Surveillance Co-ordination Function Overload
	OVLDLF bool  `json:"ovl_dlf"` // Data Link Function Overload
}

// I034/060 - System Processing Mode
type SystemProcessingMode034 struct {
	COM bool `json:"com"` // Common Part present
	PSR bool `json:"psr"` // PSR Sensor present
	SSR bool `json:"ssr"` // SSR Sensor present
	MDS bool `json:"mds"` // Mode S Sensor present
	FX  bool `json:"fx"`  // Field Extension

	// Conditional subfields
	COMData *COMProcessingData034 `json:"com_data,omitempty"`
	PSRData *PSRProcessingData034 `json:"psr_data,omitempty"`
	SSRData *SSRProcessingData034 `json:"ssr_data,omitempty"`
	MDSData *MDSProcessingData034 `json:"mds_data,omitempty"`
}

// COM Processing Mode Data
type COMProcessingData034 struct {
	REDRDP uint8 `json:"red_rdp"` // Reduction steps in RDP (3 bits)
	REDXMT uint8 `json:"red_xmt"` // Reduction steps in XMT (3 bits)
}

// PSR Processing Mode Data
type PSRProcessingData034 struct {
	POL    bool  `json:"pol"`     // Polarization
	REDRAD uint8 `json:"red_rad"` // Reduction steps (3 bits)
	STC    uint8 `json:"stc"`     // STC Map (2 bits)
}

// SSR Processing Mode Data
type SSRProcessingData034 struct {
	REDRAD uint8 `json:"red_rad"` // Reduction steps (3 bits)
}

// MDS Processing Mode Data
type MDSProcessingData034 struct {
	REDRAD uint8 `json:"red_rad"` // Reduction steps (3 bits)
	CLU    bool  `json:"clu"`     // Cluster state
}

// I034/070 - Message Count Values
type MessageCountValues034 struct {
	Repetition uint8               `json:"repetition"` // Number of counter entries
	Counters   []MessageCounter034 `json:"counters"`   // List of message counters
}

// Individual Message Counter
type MessageCounter034 struct {
	Type    uint16 `json:"type"`    // Message type (5 bits)
	Counter uint16 `json:"counter"` // Counter value (11 bits)
}

// I034/100 - Generic Polar Window
type GenericPolarWindow034 struct {
	RhoStartNM    float64 `json:"rho_start_nm"`    // Start range in nautical miles
	RhoEndNM      float64 `json:"rho_end_nm"`      // End range in nautical miles
	ThetaStartDeg float64 `json:"theta_start_deg"` // Start azimuth in degrees
	ThetaEndDeg   float64 `json:"theta_end_deg"`   // End azimuth in degrees
}

// I034/110 - Data Filter
type DataFilter034 struct {
	Type uint8 `json:"type"` // Filter type
}

// I034/120 - 3D-Position of Data Source
type Position3D034 struct {
	HeightM      float64 `json:"height_m"`      // Height in meters
	LatitudeDeg  float64 `json:"latitude_deg"`  // Latitude in degrees
	LongitudeDeg float64 `json:"longitude_deg"` // Longitude in degrees
}

// I034/090 - Collimation Error
type CollimationError034 struct {
	RangeErrorNM    float64 `json:"range_error_nm"`    // Range error in nautical miles
	AzimuthErrorDeg float64 `json:"azimuth_error_deg"` // Azimuth error in degrees
}

// I034/RE - Reserved Expansion Field
type ReservedField034 struct {
	Note string `json:"note"` // Description
	Data string `json:"data"` // Hex representation of data
}

// I034/SP - Special Purpose Field
type SpecialPurposeField034 struct {
	Note string `json:"note"` // Description
	Data string `json:"data"` // Hex representation of data
}
