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

// decode calls a field decoder, type-asserts the result, and stores a pointer in dst.
func decode[T any](data []byte, fn func([]byte) (any, int, error), dst **T) (int, error) {
	v, n, err := fn(data)
	if err != nil {
		return 0, err
	}
	r := v.(T)
	*dst = &r
	return n, nil
}

func (d *CAT048Decoder) Decode(msg *RawAsterixMessage) (*AsterixMessage, error) {
	if msg.Category != 48 {
		return nil, fmt.Errorf("expected category 48, got %d", msg.Category)
	}

	fspec, payload, err := extractFSPEC(msg.Payload)
	if err != nil {
		return nil, err
	}

	cat048, err := d.decodeFields(fspec, payload)
	if err != nil {
		return nil, err
	}

	return &AsterixMessage{
		Category: msg.Category,
		FSPEC:    fspec,
		Record:   cat048,
	}, nil
}

// decodeFields walks the FSPEC and populates a Cat048Message directly.
func (d *CAT048Decoder) decodeFields(fspec []byte, payload []byte) (*Cat048Message, error) {
	m := &Cat048Message{}
	cursor := payload
	frn := 1

	for _, fspecByte := range fspec {
		for bit := 7; bit >= 1; bit-- {
			if fspecByte&(1<<uint(bit)) != 0 {
				n, err := d.decodeField(frn, cursor, m)
				if err != nil {
					return nil, fmt.Errorf("FRN %d: %w", frn, err)
				}
				cursor = cursor[n:]
			}
			frn++
		}
	}
	return m, nil
}

func (d *CAT048Decoder) decodeField(frn int, data []byte, m *Cat048Message) (int, error) {
	switch frn {
	case 1: // I048/010
		return decode(data, decodeDataSourceIdentifier, &m.DataSource)
	case 2: // I048/140
		return decode(data, decodeTimeOfDay, &m.TimeOfDay)
	case 3: // I048/020
		return decode(data, decodeTargetReportDescriptor, &m.TargetReport)
	case 4: // I048/040
		return decode(data, decodeMeasuredPositionInPolarCoordinates, &m.MeasuredPosition)
	case 5: // I048/070
		return decode(data, decodeMode3ACode, &m.Mode3A)
	case 6: // I048/090
		return decode(data, decodeFlightLevel, &m.FlightLevel)
	case 7: // I048/130
		return decode(data, decodeRadarPlotCharacteristics, &m.RadarPlot)
	case 8: // I048/220
		return decode(data, decodeAircraftAddress, &m.AircraftAddress)
	case 9: // I048/240
		return decode(data, decodeAircraftIdentification, &m.AircraftIdentification)
	case 10: // I048/250
		if d.BDSDecoder != nil {
			bds := d.BDSDecoder
			return decode(data, func(d []byte) (any, int, error) {
				return decodeBDSRegisterDataWith(bds, d)
			}, &m.BDSRegister)
		}
		return decode(data, decodeBDSRegisterData, &m.BDSRegister)
	case 11: // I048/161
		return decode(data, decodeTrackNumber, &m.TrackNumber)
	case 12: // I048/042
		return decode(data, decodeCalculatedPositionCartesian, &m.CalculatedPosition)
	case 13: // I048/200
		return decode(data, decodeCalculatedTrackVelocity, &m.TrackVelocity)
	case 14: // I048/170
		return decode(data, decodeTrackStatus, &m.TrackStatus)
	case 15: // I048/210
		return decode(data, decodeTrackQuality, &m.TrackQuality)
	case 16: // I048/030
		return decode(data, decodeWarningErrorConditions, &m.WarningError)
	case 17: // I048/080
		return decode(data, decodeMode3ACodeConfidence, &m.Mode3AConfidence)
	case 18: // I048/100
		return decode(data, decodeModeCodeConfidence, &m.ModeCConfidence)
	case 19: // I048/110
		return decode(data, decodeHeightMeasured3D, &m.Height3D)
	case 20: // I048/120
		return decode(data, decodeRadialDopplerSpeed, &m.RadialDoppler)
	case 21: // I048/230
		return decode(data, decodeCommunicationsCapability, &m.CommCapability)
	case 22: // I048/260
		return decode(data, decodeACASResolutionAdvisory, &m.ACASAdvisory)
	case 23: // I048/055
		return decode(data, decodeMode1Code, &m.Mode1)
	case 24: // I048/050
		return decode(data, decodeMode2Code, &m.Mode2)
	case 25: // I048/065
		return decode(data, decodeMode1CodeConfidence, &m.Mode1Confidence)
	case 26: // I048/060
		return decode(data, decodeMode2CodeConfidence, &m.Mode2Confidence)
	case 27: // I048/SP
		_, n, err := decodeSpecialPurposeField(data)
		return n, err
	case 28: // I048/RE
		_, n, err := decodeReservedExpansionField(data)
		return n, err
	}
	return 0, nil
}

func decodeDataSourceIdentifier(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/010")
	}
	return DataSourceIdentifier{
		SAC: data[0],
		SIC: data[1],
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

func decodeTargetReportDescriptor(data []byte) (any, int, error) {
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

func decodeTimeOfDay(data []byte) (any, int, error) {
	if len(data) < 3 {
		return nil, 0, fmt.Errorf("too short for I048/140")
	}
	v := uint24(data[:3])
	// raw * (1/128) seconds → time.Duration
	duration := time.Duration(v) * time.Second / 128

	return TimeOfDay{
		Duration: duration,
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
func decodeRadarPlotCharacteristics(data []byte) (any, int, error) {
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

func decodeFlightLevel(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("I048/090: not enough data")
	}

	b1 := data[0]
	b2 := data[1]

	// Bits 14–1: 14-bit signed integer
	raw := (uint16(b1&0x3F) << 8) | uint16(b2) // clear top 2 bits (V & G)
	signedVal := SignExtend14To16(raw)

	return FlightLevel{
		FL:        float64(signedVal) * 0.25, // raw * (1/4) FL
		Validated: !readBit(b1, 8),           // bit 16: V=0 means validated
		Garbled:   readBit(b1, 7),            // bit 15: G
	}, 2, nil
}

func decodeAircraftAddress(data []byte) (any, int, error) {
	if len(data) < 3 {
		return nil, 0, fmt.Errorf("too short for I048/220")
	}
	addr := uint24(data[0:3])
	return fmt.Sprintf("%06x", addr), 3, nil
}

func decodeAircraftIdentification(data []byte) (any, int, error) {
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

func decodeMeasuredPositionInPolarCoordinates(data []byte) (any, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/040")
	}
	rhoRaw := binary.BigEndian.Uint16(data[:2])
	thetaRaw := binary.BigEndian.Uint16(data[2:4])

	return MeasuredPositionPolar{
		Rho:   float64(rhoRaw) / 256.0,              // raw * (1/256) NM
		Theta: float64(thetaRaw) * 360.0 / 65536.0,  // raw * (360/2^16) degrees
	}, 4, nil
}

func decodeMode3ACode(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/070")
	}
	b1 := data[0]

	return Mode3ACode{
		Validated: !readBit(b1, 8),                        // bit 16: V=0 means validated
		Garbled:   readBit(b1, 7),                         // bit 15: G
		Local:     readBit(b1, 6),                         // bit 14: L
		Code:      binary.BigEndian.Uint16(data) & 0x0FFF, // bits 12-1: 12-bit transponder code
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
func decodeCalculatedPositionCartesian(data []byte) (any, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/042")
	}
	xComp := int16(binary.BigEndian.Uint16(data[0:2]))
	yComp := int16(binary.BigEndian.Uint16(data[2:4]))

	return CalculatedPositionCartesian{
		X: float64(xComp) / 128.0, // raw * (1/128) NM
		Y: float64(yComp) / 128.0, // raw * (1/128) NM
	}, 4, nil
}

// I048/050 - Mode-2 Code in Octal Representation
func decodeMode2Code(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/050")
	}
	b1 := data[0]

	return Mode2Code{
		Validated: !readBit(b1, 8),                        // bit 16: V=0 means validated
		Garbled:   readBit(b1, 7),                         // bit 15: G
		Local:     readBit(b1, 6),                         // bit 14: L
		Code:      binary.BigEndian.Uint16(data) & 0x0FFF, // bits 12-1: 12-bit code
	}, 2, nil
}

// I048/055 - Mode-1 Code in Octal Representation
func decodeMode1Code(data []byte) (any, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/055")
	}
	b := data[0]

	return Mode1Code{
		Validated: !readBit(b, 8), // bit 8: V=0 means validated
		Garbled:   readBit(b, 7),  // bit 7: G
		Local:     readBit(b, 6),  // bit 6: L
		Code:      b & 0x1F,       // bits 5-1
	}, 1, nil
}

// I048/060 - Mode-2 Code Confidence Indicator
func decodeMode2CodeConfidence(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/060")
	}
	return Mode2CodeConfidence{
		QA4: readBit(data[0], 4), QA2: readBit(data[0], 3), QA1: readBit(data[0], 2),
		QB4: readBit(data[0], 1), QB2: readBit(data[1], 8), QB1: readBit(data[1], 7),
		QC4: readBit(data[1], 6), QC2: readBit(data[1], 5), QC1: readBit(data[1], 4),
		QD4: readBit(data[1], 3), QD2: readBit(data[1], 2), QD1: readBit(data[1], 1),
	}, 2, nil
}

// I048/065 - Mode-1 Code Confidence Indicator
func decodeMode1CodeConfidence(data []byte) (any, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/065")
	}
	b := data[0]
	return Mode1CodeConfidence{
		QA4: readBit(b, 5), QA2: readBit(b, 4), QA1: readBit(b, 3),
		QB2: readBit(b, 2), QB1: readBit(b, 1),
	}, 1, nil
}

// I048/080 - Mode-3/A Code Confidence Indicator
func decodeMode3ACodeConfidence(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/080")
	}
	return Mode3ACodeConfidence{
		QA4: readBit(data[0], 4), QA2: readBit(data[0], 3), QA1: readBit(data[0], 2),
		QB4: readBit(data[0], 1), QB2: readBit(data[1], 8), QB1: readBit(data[1], 7),
		QC4: readBit(data[1], 6), QC2: readBit(data[1], 5), QC1: readBit(data[1], 4),
		QD4: readBit(data[1], 3), QD2: readBit(data[1], 2), QD1: readBit(data[1], 1),
	}, 2, nil
}

// I048/100 - Mode-C Code and Code Confidence Indicator
func decodeModeCodeConfidence(data []byte) (any, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/100")
	}

	return ModeCCodeConfidence{
		Validated: !readBit(data[0], 8),
		Garbled:   readBit(data[0], 7),
		Code:      binary.BigEndian.Uint16(data[0:2]) & 0x0FFF,
		QC1:       readBit(data[2], 4), QA1: readBit(data[2], 3),
		QC2:       readBit(data[2], 2), QA2: readBit(data[2], 1),
		QC4:       readBit(data[3], 8), QA4: readBit(data[3], 7),
		QB1:       readBit(data[3], 6), QD1: readBit(data[3], 5),
		QB2:       readBit(data[3], 4), QD2: readBit(data[3], 3),
		QB4:       readBit(data[3], 2), QD4: readBit(data[3], 1),
	}, 4, nil
}

// I048/110 - Height Measured by a 3D Radar
func decodeHeightMeasured3D(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/110")
	}

	// 14-bit signed value in two's complement
	raw := binary.BigEndian.Uint16(data) & 0x3FFF
	var height int16
	if raw&0x2000 != 0 {
		height = int16(raw | 0xC000) // sign extend
	} else {
		height = int16(raw)
	}

	return HeightMeasured3D{
		Height: float64(height) * 25.0, // raw * 25 ft
	}, 2, nil
}

// I048/120 - Radial Doppler Speed (simplified implementation)
func decodeRadialDopplerSpeed(data []byte) (any, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/120")
	}

	// Compound data item — store raw bytes for now
	raw := make([]byte, len(data))
	copy(raw, data)
	return RadialDopplerSpeed{
		RawData: raw,
	}, len(data), nil
}

// I048/161 - Track Number
func decodeTrackNumber(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/161")
	}

	return TrackNumber{
		Number: binary.BigEndian.Uint16(data) & 0x0FFF, // 12-bit track number
	}, 2, nil
}

// I048/170 - Track Status
func decodeTrackStatus(data []byte) (any, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/170")
	}

	b := data[0]
	offset := 1

	ts := TrackStatus{
		CNF: readBit(b, 8),   // Confirmed vs Tentative Track
		RAD: (b >> 5) & 0x03, // Type of Sensor(s) maintaining Track
		DOU: readBit(b, 5),   // Confidence in plot to track association
		MAH: readBit(b, 4),   // Manoeuvre detection in Horizontal Sense
		CDM: (b >> 1) & 0x03, // Climbing/Descending Mode
	}

	if readBit(b, 1) && len(data) > offset {
		b2 := data[offset]
		offset++
		ts.HasExtension = true
		ts.TRE = readBit(b2, 8) // Signal for End_of_Track
		ts.GHO = readBit(b2, 7) // Ghost vs. true target
		ts.SUP = readBit(b2, 6) // Track maintained with neighbouring info
		ts.TCC = readBit(b2, 5) // Type of plot coordinate transformation
	}

	return ts, offset, nil
}

// I048/200 - Calculated Track Velocity in Polar Co-ordinates
func decodeCalculatedTrackVelocity(data []byte) (any, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/200")
	}

	gsRaw := binary.BigEndian.Uint16(data[0:2])
	hdgRaw := binary.BigEndian.Uint16(data[2:4])

	return CalculatedTrackVelocity{
		Groundspeed: float64(gsRaw) / 16384.0,            // raw * (2^-14) NM/s
		Heading:     float64(hdgRaw) * 360.0 / 65536.0,   // raw * (360/2^16) degrees from geographic north
	}, 4, nil
}

// I048/210 - Track Quality
func decodeTrackQuality(data []byte) (any, int, error) {
	if len(data) < 4 {
		return nil, 0, fmt.Errorf("too short for I048/210")
	}

	return TrackQuality{
		SigmaX: float64(data[0]) / 128.0,          // raw * (1/128) NM
		SigmaY: float64(data[1]) / 128.0,          // raw * (1/128) NM
		SigmaV: float64(data[2]) / 16384.0,        // raw * (2^-14) NM/s
		SigmaH: float64(data[3]) * 360.0 / 4096.0, // raw * (360/2^12) degrees
	}, 4, nil
}

// I048/230 - Communications/ACAS Capability and Flight Status
func decodeCommunicationsCapability(data []byte) (any, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I048/230")
	}

	b1 := data[0]
	b2 := data[1]

	return CommunicationsCapability{
		COM:  (b1 >> 5) & 0x07, // Communications capability
		STAT: (b1 >> 2) & 0x07, // Flight Status
		SI:   readBit(b1, 2),   // SI/II Transponder Capability
		MSSC: readBit(b2, 8),   // Mode-S Specific Service Capability
		ARC:  readBit(b2, 7),   // Altitude reporting capability
		AIC:  readBit(b2, 6),   // Aircraft identification capability
		B1A:  readBit(b2, 5),   // BDS 1,0 bit 16
		B1B:  b2 & 0x0F,        // BDS 1,0 bits 37/40
	}, 2, nil
}

// decodeBDSRegisterData decodes I048/250 Mode S MB Data (raw-only, no BDS interpretation).
func decodeBDSRegisterData(data []byte) (any, int, error) {
	return decodeBDSRegisterDataWith(nil, data)
}

// decodeBDSRegisterDataWith decodes I048/250 Mode S MB Data.
// If bdsDecoder is non-nil, each register is decoded into a typed struct.
func decodeBDSRegisterDataWith(bdsDecoder BDSDecoder, data []byte) (any, int, error) {
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
func decodeACASResolutionAdvisory(data []byte) (any, int, error) {
	if len(data) < 7 {
		return nil, 0, fmt.Errorf("too short for I048/260")
	}

	return ACASResolutionAdvisory{
		ACASRA: hex.EncodeToString(data[:7]),
	}, 7, nil
}

// I048/030 - Warning/Error Conditions and Target Classification
func decodeWarningErrorConditions(data []byte) (any, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I048/030")
	}

	offset := 0
	var codes []uint8

	for offset < len(data) {
		b := data[offset]
		codes = append(codes, (b>>1)&0x7F) // bits 8-2
		offset++
		if b&0x01 == 0 { // FX bit - 0 means end
			break
		}
	}

	return WarningErrorConditions{Codes: codes}, offset, nil
}

// Placeholder implementations for Special Purpose and Reserved Expansion Fields
func decodeSpecialPurposeField(data []byte) (any, int, error) {
	return map[string]any{
		"note": "Special Purpose Field - implementation specific",
		"data": hex.EncodeToString(data),
	}, len(data), nil
}

func decodeReservedExpansionField(data []byte) (any, int, error) {
	return map[string]any{
		"note": "Reserved Expansion Field - implementation specific",
		"data": hex.EncodeToString(data),
	}, len(data), nil
}
