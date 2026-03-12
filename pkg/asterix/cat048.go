package asterix

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type CAT048Decoder struct {
	BDSDecoder BDSDecoder // optional; nil means raw-only BDS output
}

func (d *CAT048Decoder) Decode(msg *RawAsterixMessage) (*AsterixMessage, error) {
	if msg.Category != 48 {
		return nil, fmt.Errorf("expected category 48, got %d", msg.Category)
	}

	fspec, restPayload, err := extractFSPEC(msg.Payload)
	if err != nil {
		return nil, err
	}

	items := cat048Items
	if d.BDSDecoder != nil {
		items = make(map[int]DataItem, len(cat048Items))
		for k, v := range cat048Items {
			items[k] = v
		}
		items[10] = NewDataItem("I048/250", func(data []byte) (interface{}, int, error) {
			return decodeBDSRegisterDataWith(d.BDSDecoder, data)
		})
	}

	decoded, _, err := WalkFSPEC(fspec, restPayload, items)
	if err != nil {
		return nil, err
	}

	AsterixMessage := &AsterixMessage{
		Category: msg.Category,
		Sic:      decoded["I048/010"].(map[string]uint8)["SIC"],
		Sac:      decoded["I048/010"].(map[string]uint8)["SAC"],
		Items:    decoded,
		FSPEC:    fspec,
	}
	return AsterixMessage, nil
}

var cat048Items = map[int]DataItem{
	// First FSPEC octet (FRN 1-7 + FX)
	1: NewDataItem("I048/010", decodeDataSourceIdentifier),               // FRN 1: Data Source Identifier
	2: NewDataItem("I048/140", decodeTimeOfDay),                          // FRN 2: Time-of-Day
	3: NewDataItem("I048/020", decodeTargetReportDescriptor),             // FRN 3: Type and Properties of the Target Report
	4: NewDataItem("I048/040", decodeMeasuredPositionInPolarCoordinates), // FRN 4: Measured Position in Slant Polar Coordinates
	5: NewDataItem("I048/070", decodeMode3ACode),                         // FRN 5: Mode-3/A Code in Octal Representation
	6: NewDataItem("I048/090", decodeFlightLevel),                        // FRN 6: Flight Level in Binary Representation
	7: NewDataItem("I048/130", decodeRadarPlotCharacteristics),           // FRN 7: Radar Plot Characteristics
	//FX = Field Extension Indicator

	// Second FSPEC octet (FRN 8-14 + FX)
	8:  NewDataItem("I048/220", decodeAircraftAddress),             // FRN 8: Aircraft Address
	9:  NewDataItem("I048/240", decodeAircraftIdentification),      // FRN 9: Aircraft Identification
	10: NewDataItem("I048/250", decodeBDSRegisterData),             // FRN 10: Mode S MB Data
	11: NewDataItem("I048/161", decodeTrackNumber),                 // FRN 11: Track Number
	12: NewDataItem("I048/042", decodeCalculatedPositionCartesian), // FRN 12: Calculated Position in Cartesian Coordinates
	13: NewDataItem("I048/200", decodeCalculatedTrackVelocity),     // FRN 13: Calculated Track Velocity in Polar Representation
	14: NewDataItem("I048/170", decodeTrackStatus),                 // FRN 14: Track Status
	//FX = Field Extension Indicator

	// Third FSPEC octet (FRN 15-21 + FX)
	15: NewDataItem("I048/210", decodeTrackQuality),             // FRN 15: Track Quality
	16: NewDataItem("I048/030", decodeWarningErrorConditions),   // FRN 16: Warning/Error Conditions/Target Classification
	17: NewDataItem("I048/080", decodeMode3ACodeConfidence),     // FRN 17: Mode-3/A Code Confidence Indicator
	18: NewDataItem("I048/100", decodeModeCodeConfidence),       // FRN 18: Mode-C Code and Confidence Indicator
	19: NewDataItem("I048/110", decodeHeightMeasured3D),         // FRN 19: Height Measured by 3D Radar
	20: NewDataItem("I048/120", decodeRadialDopplerSpeed),       // FRN 20: Radial Doppler Speed
	21: NewDataItem("I048/230", decodeCommunicationsCapability), // FRN 21: Communications / ACAS Capability and Flight Status
	//FX = Field Extension Indicator

	// Fourth FSPEC octet (FRN 22-28 + FX)
	22: NewDataItem("I048/260", decodeACASResolutionAdvisory), // FRN 22: ACAS Resolution Advisory Report
	23: NewDataItem("I048/055", decodeMode1Code),              // FRN 23: Mode-1 Code in Octal Representation
	24: NewDataItem("I048/050", decodeMode2Code),              // FRN 24: Mode-2 Code in Octal Representation
	25: NewDataItem("I048/065", decodeMode1CodeConfidence),    // FRN 25: Mode-1 Code Confidence Indicator
	26: NewDataItem("I048/060", decodeMode2CodeConfidence),    // FRN 26: Mode-2 Code Confidence Indicator
	27: NewDataItem("I048/SP", decodeSpecialPurposeField),     // FRN 27: Special Purpose Field
	28: NewDataItem("I048/RE", decodeReservedExpansionField),  // FRN 28: Reserved Expansion Field
	//FX = Field Extension Indicator
}

func decodeDataSourceIdentifier(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/010")
	}
	return map[string]uint8{
		"SAC": data[0],
		"SIC": data[1],
	}, 2, nil
}

type TargetReportDescriptor struct {
	// Octet 1
	TYP byte // bits 8–7-6
	SIM bool
	RDP bool
	SPI bool
	RAB bool

	// Octet 2 (optional)
	TST    bool
	ERR    bool
	XPP    bool
	ME     bool
	MI     bool
	FOEFRI bool
}

func decodeTargetReportDescriptor(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("I048/020: not enough data")
	}
	offset := 0
	b1 := data[offset]
	offset++

	// typ
	//= 000 No detection
	//= 001 Single PSR detection
	//= 010 Single SSR detection
	//= 011 SSR + PSR detection
	//= 100 Single ModeS All-Call
	//= 101 Single ModeS Roll-Call
	//= 110 ModeS All-Call + PSR
	//= 111 ModeS Roll-Call +PSR

	trd := TargetReportDescriptor{
		TYP: (b1 >> 5) & 0x07, // bits 8–7-6
		SIM: readBit(b1, 5),
		RDP: readBit(b1, 4),
		SPI: readBit(b1, 3),
		RAB: readBit(b1, 2),
	}

	if readBit(b1, 1) {
		if len(data) < 2 {
			return nil, 0, fmt.Errorf("I048/020: FX=1 but only one byte present")
		}
		b2 := data[offset]
		offset++

		trd.TST = readBit(b2, 8)
		trd.ERR = readBit(b2, 7)
		trd.XPP = readBit(b2, 6)
		trd.ME = readBit(b2, 5)
		trd.MI = readBit(b2, 4)
		trd.FOEFRI = readBit(b2, 3)
		// bit 2 = spare
		// bit 1 = FX — not used (no third byte defined)
	}

	return trd, offset, nil
}

func decodeTimeOfDay(data []byte) (interface{}, int, error) {
	if len(data) < 3 {
		return nil, 0, fmt.Errorf("too short for I048/140")
	}
	v := uint24(data[:3])
	// Unit: 1/128 seconds so convert to nanoseconds
	duration := time.Duration((time.Duration(v) * time.Second) / 128)

	return map[string]interface{}{
		"raw_value":   v,        // Raw value in 1/128 second units
		"duration_ns": duration, // Converted as time.Duration
	}, 3, nil
}

// readBit reads the nth (8-1) bit of the byte 8 is the leftmost bit.
func readBit(data byte, nth uint8) bool {
	return (data & (1 << (nth - 1))) != 0
}

// I048/130 — Radar Plot Characteristics
type RadarPlotCharacteristics struct {
	SRL *uint8 // Subfield #1: SSR Plot Runlength
	SRR *uint8 // Subfield #2: Number of replies (MSSR)
	SAM *uint8 // Subfield #3: Amplitude of replies (MSSR)
	PRL *uint8 // Subfield #4: PSR Plot Runlength
	PAM *uint8 // Subfield #5: Amplitude of PSR plot
	RPD *int8  // Subfield #6: Difference in range (PSR - SSR)
}

// decodeRadarPlotCharacteristics parses I048/130 starting at data[0]
func decodeRadarPlotCharacteristics(data []byte) (interface{}, int, error) {
	var rpc RadarPlotCharacteristics
	offset := 0

	// read sub fspec
	var fspec []byte
	for {
		if offset >= len(data) {
			return rpc, offset, fmt.Errorf("unexpected end of data reading FSPEC of I048/130")
		}
		f := data[offset]
		fspec = append(fspec, f)
		offset++
		if f&0x01 == 0 {
			break // FX == 0 → fin
		}
	}

	// sub FSPEC bits 8-2
	subfieldIndex := 0
	for _, f := range fspec {
		for bit := uint8(8); bit >= 2; bit-- {
			subfieldIndex++
			if f&(1<<(bit-1)) != 0 {
				if offset >= len(data) {
					return rpc, offset, fmt.Errorf("missing data for subfield #%d of I048/130", subfieldIndex)
				}
				value := data[offset]
				switch subfieldIndex {
				case 1: // SRL
					rpc.SRL = &value
				case 2: // SRR
					rpc.SRR = &value
				case 3: // SAM
					rpc.SAM = &value
				case 4: // PRL
					rpc.PRL = &value
				case 5: // PAM
					rpc.PAM = &value
				case 6: // RPD
					signed := int8(value)
					rpc.RPD = &signed
				default:
					return rpc, offset, fmt.Errorf("subfield #%d not defined in I048/130", subfieldIndex)
				}
				offset++
			}
		}
	}

	return rpc, offset, nil
}

func decodeFlightLevel(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("I048/090: not enough data")
	}

	b1 := data[0]
	b2 := data[1]

	v := readBit(b1, 8) // bit 16
	g := readBit(b1, 7) // bit 15

	// Bits 14–1: 14-bit signed integer
	raw := (uint16(b1&0x3F) << 8) | uint16(b2) // clear top 2 bits (V & G)

	// Sign-extend
	signedVal := SignExtend14To16(raw)

	return struct {
		FlightLevel float64
		RawValue    int16
		Validated   bool
		Garbled     bool
	}{
		FlightLevel: float64(signedVal) * 0.25,
		RawValue:    signedVal,
		Validated:   !v,
		Garbled:     g,
	}, 2, nil
}

func decodeAircraftAddress(data []byte) (interface{}, int, error) {
	if len(data) < 3 {
		return nil, 0, fmt.Errorf("too short for I048/220")
	}
	addr := uint24(data[0:3])
	return fmt.Sprintf("%06x", addr), 3, nil
}

func decodeAircraftIdentification(data []byte) (interface{}, int, error) {
	if len(data) < 6 {
		return "", 0, fmt.Errorf("not enough data for Aircraft Identification")
	}

	raw := binary.BigEndian.Uint64(append([]byte{0, 0}, data[:6]...)) // 48 bits
	var result []rune
	for i := 0; i < 8; i++ {
		shift := uint(42 - i*6)
		char := (raw >> shift) & 0x3F
		result = append(result, ia5Char(byte(char)))
	}
	return strings.TrimSpace(string(result)), 6, nil
}

func ia5Char(b byte) rune {
	switch {
	case b >= 1 && b <= 26:
		return rune('A' + b - 1)
	case b >= 48 && b <= 57:
		return rune('0' + b - 48)
	case b == 32:
		return ' '
	case b == 45:
		return '-'
	default:
		return ' ' // undefined or padding
	}
}

func decodeMeasuredPositionInPolarCoordinates(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/040")
	}
	rhoRaw := binary.BigEndian.Uint16(data[:2])
	thetaRaw := binary.BigEndian.Uint16(data[2:4])
	rho := float64(rhoRaw) / 256.0               // LSB = 1/256 NM
	theta := float64(thetaRaw) * 360.0 / 65536.0 // LSB = 360°/2^16

	return map[string]interface{}{
		"rho_nm":    rho,
		"rho_raw":   rhoRaw, // Raw value in 1/256 NM units
		"theta_deg": theta,
		"theta_raw": thetaRaw, // Raw value in 360°/2^16 units
	}, 4, nil
}

func decodeMode3ACode(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/070")
	}
	b1 := data[0]

	return struct {
		Validated bool
		Garbled   bool
		Local     bool
		Code      string
	}{
		Validated: !readBit(b1, 8), // V bit - 0=validated
		Garbled:   readBit(b1, 7),  // G bit
		Local:     readBit(b1, 6),  // L bit
		Code:      FormatTransponderCodeFromBytes(data, 0),
	}, 2, nil
}

// Legacy functions - now using binary.go utilities
func int24(b []byte) int32 {
	return ReadInt24BE(b, 0)
}

func uint24(b []byte) uint32 {
	return ReadUint24BE(b, 0)
}

func twosComplement24(x int32) int32 {
	return TwosComplement24(x)
}

// Additional decoder functions for missing data items

// I048/042 - Calculated Position in Cartesian Co-ordinates
func decodeCalculatedPositionCartesian(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/042")
	}
	// X-Component (signed, two's complement)
	xComp := int16(binary.BigEndian.Uint16(data[0:2]))
	// Y-Component (signed, two's complement)
	yComp := int16(binary.BigEndian.Uint16(data[2:4]))

	return map[string]interface{}{
		"x_nm":  float64(xComp) / 128.0, // LSB = 1/128 NM
		"x_raw": xComp,                  // Raw value in 1/128 NM units
		"y_nm":  float64(yComp) / 128.0, // LSB = 1/128 NM
		"y_raw": yComp,                  // Raw value in 1/128 NM units
	}, 4, nil
}

// I048/050 - Mode-2 Code in Octal Representation
func decodeMode2Code(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/050")
	}
	b1 := data[0]

	return struct {
		Validated bool
		Garbled   bool
		Local     bool
		Code      string
	}{
		Validated: !readBit(b1, 8), // V bit - 0=validated
		Garbled:   readBit(b1, 7),  // G bit
		Local:     readBit(b1, 6),  // L bit
		Code:      fmt.Sprintf("%04o", binary.BigEndian.Uint16(data)&0x0FFF),
	}, 2, nil
}

// I048/055 - Mode-1 Code in Octal Representation
func decodeMode1Code(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/055")
	}
	b := data[0]

	return struct {
		Validated bool
		Garbled   bool
		Local     bool
		Code      uint8
	}{
		Validated: !readBit(b, 8), // V bit - 0=validated
		Garbled:   readBit(b, 7),  // G bit
		Local:     readBit(b, 6),  // L bit
		Code:      b & 0x1F,       // bits 5-1
	}, 1, nil
}

// I048/060 - Mode-2 Code Confidence Indicator
func decodeMode2CodeConfidence(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/060")
	}
	return map[string]interface{}{
		"qa4": readBit(data[0], 4),
		"qa2": readBit(data[0], 3),
		"qa1": readBit(data[0], 2),
		"qb4": readBit(data[0], 1),
		"qb2": readBit(data[1], 8),
		"qb1": readBit(data[1], 7),
		"qc4": readBit(data[1], 6),
		"qc2": readBit(data[1], 5),
		"qc1": readBit(data[1], 4),
		"qd4": readBit(data[1], 3),
		"qd2": readBit(data[1], 2),
		"qd1": readBit(data[1], 1),
	}, 2, nil
}

// I048/065 - Mode-1 Code Confidence Indicator
func decodeMode1CodeConfidence(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/065")
	}
	b := data[0]
	return map[string]interface{}{
		"qa4": readBit(b, 5),
		"qa2": readBit(b, 4),
		"qa1": readBit(b, 3),
		"qb2": readBit(b, 2),
		"qb1": readBit(b, 1),
	}, 1, nil
}

// I048/080 - Mode-3/A Code Confidence Indicator
func decodeMode3ACodeConfidence(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/080")
	}
	return map[string]interface{}{
		"qa4": readBit(data[0], 4),
		"qa2": readBit(data[0], 3),
		"qa1": readBit(data[0], 2),
		"qb4": readBit(data[0], 1),
		"qb2": readBit(data[1], 8),
		"qb1": readBit(data[1], 7),
		"qc4": readBit(data[1], 6),
		"qc2": readBit(data[1], 5),
		"qc1": readBit(data[1], 4),
		"qd4": readBit(data[1], 3),
		"qd2": readBit(data[1], 2),
		"qd1": readBit(data[1], 1),
	}, 2, nil
}

// I048/100 - Mode-C Code and Code Confidence Indicator
func decodeModeCodeConfidence(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/100")
	}

	b1 := data[0]

	return struct {
		Validated bool
		Garbled   bool
		Code      uint16
		Quality   map[string]bool
	}{
		Validated: !readBit(b1, 8),
		Garbled:   readBit(b1, 7),
		Code:      binary.BigEndian.Uint16(data[0:2]) & 0x0FFF,
		Quality: map[string]bool{
			"qc1": readBit(data[2], 4),
			"qa1": readBit(data[2], 3),
			"qc2": readBit(data[2], 2),
			"qa2": readBit(data[2], 1),
			"qc4": readBit(data[3], 8),
			"qa4": readBit(data[3], 7),
			"qb1": readBit(data[3], 6),
			"qd1": readBit(data[3], 5),
			"qb2": readBit(data[3], 4),
			"qd2": readBit(data[3], 3),
			"qb4": readBit(data[3], 2),
			"qd4": readBit(data[3], 1),
		},
	}, 4, nil
}

// I048/110 - Height Measured by a 3D Radar
func decodeHeightMeasured3D(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/110")
	}

	// 14-bit signed value in two's complement
	raw := binary.BigEndian.Uint16(data) & 0x3FFF
	var height int16
	if raw&0x2000 != 0 { // Check sign bit (bit 14)
		height = int16(raw | 0xC000) // Sign extend
	} else {
		height = int16(raw)
	}

	return map[string]interface{}{
		"height_ft":  float64(height) * 25.0, // LSB = 25 ft
		"height_raw": height,                 // Raw value in 25 ft units
	}, 2, nil
}

// I048/120 - Radial Doppler Speed (simplified implementation)
func decodeRadialDopplerSpeed(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/120")
	}

	// This is a compound data item - simplified implementation
	// Full implementation would need to handle multiple subfields
	return map[string]interface{}{
		"note": "Compound data item - simplified implementation",
	}, len(data), nil
}

// I048/161 - Track Number
func decodeTrackNumber(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/161")
	}

	trackNum := binary.BigEndian.Uint16(data) & 0x0FFF // 12 bits
	return trackNum, 2, nil
}

// I048/170 - Track Status
func decodeTrackStatus(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/170")
	}

	b := data[0]
	offset := 1

	status := map[string]interface{}{
		"cnf": readBit(b, 8),   // Confirmed vs Tentative Track
		"rad": (b >> 5) & 0x03, // Type of Sensor(s) maintaining Track
		"dou": readBit(b, 5),   // Confidence in plot to track association
		"mah": readBit(b, 4),   // Manoeuvre detection in Horizontal Sense
		"cdm": (b >> 1) & 0x03, // Climbing/Descending Mode
	}

	// Check for extensions
	if readBit(b, 1) && len(data) > offset {
		b2 := data[offset]
		offset++
		status["tre"] = readBit(b2, 8) // Signal for End_of_Track
		status["gho"] = readBit(b2, 7) // Ghost vs. true target
		status["sup"] = readBit(b2, 6) // Track maintained with neighbouring info
		status["tcc"] = readBit(b2, 5) // Type of plot coordinate transformation
	}

	return status, offset, nil
}

// I048/200 - Calculated Track Velocity in Polar Co-ordinates
func decodeCalculatedTrackVelocity(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/200")
	}

	groundSpeed := binary.BigEndian.Uint16(data[0:2])
	heading := binary.BigEndian.Uint16(data[2:4])

	return map[string]interface{}{
		"groundspeed_kt":  float64(groundSpeed) * 0.22,        // LSB = (2^-14) NM/s ≈ 0.22 kt
		"groundspeed_raw": groundSpeed,                        // Raw value in (2^-14) NM/s units
		"heading_deg":     float64(heading) * 360.0 / 65536.0, // LSB = 360°/2^16
		"heading_raw":     heading,                            // Raw value in 360°/2^16 units
	}, 4, nil
}

// I048/210 - Track Quality
func decodeTrackQuality(data []byte) (interface{}, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/210")
	}

	return map[string]interface{}{
		"sigma_x_nm":  float64(data[0]) / 128.0,   // Standard deviation on X axis
		"sigma_x_raw": data[0],                    // Raw value in 1/128 NM units
		"sigma_y_nm":  float64(data[1]) / 128.0,   // Standard deviation on Y axis
		"sigma_y_raw": data[1],                    // Raw value in 1/128 NM units
		"sigma_v_kt":  float64(data[2]) * 0.22,    // Standard deviation on groundspeed
		"sigma_v_raw": data[2],                    // Raw value in (2^-14) NM/s units
		"sigma_h_deg": float64(data[3]) * 0.08789, // Standard deviation on heading
		"sigma_h_raw": data[3],                    // Raw value in 360°/2^12 units
	}, 4, nil
}

// I048/230 - Communications/ACAS Capability and Flight Status
func decodeCommunicationsCapability(data []byte) (interface{}, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/230")
	}

	b1 := data[0]
	b2 := data[1]

	return map[string]interface{}{
		"com":  (b1 >> 5) & 0x07, // Communications capability
		"stat": (b1 >> 2) & 0x07, // Flight Status
		"si":   readBit(b1, 2),   // SI/II Transponder Capability
		"mssc": readBit(b2, 8),   // Mode-S Specific Service Capability
		"arc":  readBit(b2, 7),   // Altitude reporting capability
		"aic":  readBit(b2, 6),   // Aircraft identification capability
		"b1a":  readBit(b2, 5),   // BDS 1,0 bit 16
		"b1b":  b2 & 0x0F,        // BDS 1,0 bits 37/40
	}, 2, nil
}

// decodeBDSRegisterData decodes I048/250 Mode S MB Data (raw-only, no BDS interpretation).
func decodeBDSRegisterData(data []byte) (interface{}, int, error) {
	return decodeBDSRegisterDataWith(nil, data)
}

// decodeBDSRegisterDataWith decodes I048/250 Mode S MB Data.
// If bdsDecoder is non-nil, each register is decoded into a typed struct.
func decodeBDSRegisterDataWith(bdsDecoder BDSDecoder, data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/250")
	}

	rep := data[0]
	offset := 1
	registers := make(map[string]BDSRegister)

	for i := uint8(0); i < rep; i++ {
		if offset+8 > len(data) {
			break
		}

		bdsData := data[offset : offset+7]
		bdsAddr := data[offset+7]
		bdsKey := fmt.Sprintf("0x%02x", bdsAddr)

		reg := BDSRegister{
			BDSCode:    bdsAddr,
			RawData:    bdsData,
			BDSDataRaw: hex.EncodeToString(bdsData),
		}

		if bdsDecoder != nil {
			decoded, err := bdsDecoder.DecodeBDS(bdsAddr, bdsData)
			if err != nil {
				reg.Error = err.Error()
			} else {
				reg.Decoded = decoded
			}
		}

		registers[bdsKey] = reg
		offset += 8
	}

	return BDSRegisterData{
		Repetition: rep,
		Registers:  registers,
	}, offset, nil
}

// I048/260 - ACAS Resolution Advisory Report
func decodeACASResolutionAdvisory(data []byte) (interface{}, int, error) {
	if len(data) < 7 {
		return nil, 0, fmt.Errorf("too short for I048/260")
	}

	// Extract 56-bit ACAS RA data
	return map[string]interface{}{
		"acas_ra": hex.EncodeToString(data[:7]),
	}, 7, nil
}

// I048/030 - Warning/Error Conditions and Target Classification
func decodeWarningErrorConditions(data []byte) (interface{}, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/030")
	}

	offset := 0
	codes := make([]uint8, 0)

	for offset < len(data) {
		b := data[offset]
		code := (b >> 1) & 0x7F // bits 8-2
		codes = append(codes, code)
		offset++

		// Check FX bit (bit 1) - if 0, end of data item
		if (b & 0x01) == 0 {
			break
		}
	}

	return map[string]interface{}{
		"codes": codes,
	}, offset, nil
}

// Placeholder implementations for Special Purpose and Reserved Expansion Fields
func decodeSpecialPurposeField(data []byte) (interface{}, int, error) {
	return map[string]interface{}{
		"note": "Special Purpose Field - implementation specific",
		"data": hex.EncodeToString(data),
	}, len(data), nil
}

func decodeReservedExpansionField(data []byte) (interface{}, int, error) {
	return map[string]interface{}{
		"note": "Reserved Expansion Field - implementation specific",
		"data": hex.EncodeToString(data),
	}, len(data), nil
}
