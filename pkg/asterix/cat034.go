package asterix

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type CAT034Decoder struct{}

func (d *CAT034Decoder) Decode(msg *RawAsterixMessage) (*AsterixMessage, error) {
	if msg.Category != 34 {
		return nil, fmt.Errorf("expected category 34, got %d", msg.Category)
	}

	fspec, restPayload, err := extractFSPEC(msg.Payload)
	if err != nil {
		return nil, err
	}

	decoded, _, err := WalkFSPEC(fspec, restPayload, cat034Items)
	if err != nil {
		return nil, err
	}

	// Extract SAC/SIC from decoded data
	var sac, sic uint8
	if dsid, ok := decoded["I034/010"].(*DataSourceIdentifier034); ok && dsid != nil {
		sac = dsid.SAC
		sic = dsid.SIC
	}

	AsterixMessage := &AsterixMessage{
		Category: msg.Category,
		Sic:      sic,
		Sac:      sac,
		Items:    decoded,
		FSPEC:    fspec,
	}
	return AsterixMessage, nil
}


var cat034Items = map[int]struct {
	Name    string
	Decoder ItemDecoder[any]
}{
	// CAT 034 User Application Profile - 14 FRNs
	1:  {"I034/010", func(data []byte) (any, int, error) { return decodeDataSourceIdentifier034(data) }}, // FRN 1: Data Source Identifier
	2:  {"I034/000", func(data []byte) (any, int, error) { return decodeMessageType034(data) }},          // FRN 2: Message Type
	3:  {"I034/030", func(data []byte) (any, int, error) { return decodeTimeOfDay034(data) }},            // FRN 3: Time of Day
	4:  {"I034/020", func(data []byte) (any, int, error) { return decodeSectorNumber034(data) }},         // FRN 4: Sector Number
	5:  {"I034/041", func(data []byte) (any, int, error) { return decodeAntennaRotationSpeed034(data) }}, // FRN 5: Antenna Rotation Speed
	6:  {"I034/050", func(data []byte) (any, int, error) { return decodeSystemConfiguration034(data) }},  // FRN 6: System Configuration and Status
	7:  {"I034/060", func(data []byte) (any, int, error) { return decodeSystemProcessingMode034(data) }}, // FRN 7: System Processing Mode
	8:  {"I034/070", func(data []byte) (any, int, error) { return decodeMessageCountValues034(data) }},   // FRN 8: Message Count Values
	9:  {"I034/100", func(data []byte) (any, int, error) { return decodeGenericPolarWindow034(data) }},   // FRN 9: Generic Polar Window
	10: {"I034/110", func(data []byte) (any, int, error) { return decodeDataFilter034(data) }},           // FRN 10: Data Filter
	11: {"I034/120", func(data []byte) (any, int, error) { return decode3DPositionOfSource034(data) }},   // FRN 11: 3D-Position of Data Source
	12: {"I034/090", func(data []byte) (any, int, error) { return decodeCollimationError034(data) }},     // FRN 12: Collimation Error
	13: {"I034/RE", func(data []byte) (any, int, error) { return decodeReservedExpansion034(data) }},     // FRN 13: Reserved Expansion Field
	14: {"I034/SP", func(data []byte) (any, int, error) { return decodeSpecialPurpose034(data) }},        // FRN 14: Special Purpose Field
}

// I034/010 - Data Source Identifier
func decodeDataSourceIdentifier034(data []byte) (*DataSourceIdentifier034, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I034/010")
	}
	return &DataSourceIdentifier034{
		SAC: data[0],
		SIC: data[1],
	}, 2, nil
}

// I034/000 - Message Type
func decodeMessageType034(data []byte) (*uint8, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I034/000")
	}

	messageType := data[0]
	return &messageType, 1, nil
}

// I034/030 - Time of Day
func decodeTimeOfDay034(data []byte) (*time.Duration, int, error) {
	if len(data) < 3 {
		return nil, 0, fmt.Errorf("too short for I034/030")
	}

	// 24-bit time value
	timeValue := uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	// Unit: 1/128 seconds
	duration := time.Duration(float64(timeValue) * float64(time.Second) / 128.0)
	return &duration, 3, nil
}

// I034/020 - Sector Number
func decodeSectorNumber034(data []byte) (*SectorNumber034, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I034/020")
	}

	// 8-bit sector number, LSB = 360°/2^8 ≈ 1.41°
	sectorNumber := data[0]
	azimuth := float64(sectorNumber) * 360.0 / 256.0

	return &SectorNumber034{
		Sector:     sectorNumber,
		AzimuthDeg: azimuth,
	}, 1, nil
}

// I034/041 - Antenna Rotation Speed
func decodeAntennaRotationSpeed034(data []byte) (*AntennaRotationSpeed034, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I034/041")
	}

	// 16-bit rotation period in 1/128 second units
	rotationPeriod := binary.BigEndian.Uint16(data)
	periodSeconds := float64(rotationPeriod) / 128.0

	return &AntennaRotationSpeed034{
		RotationPeriodS: periodSeconds,
		RawValue:        rotationPeriod,
	}, 2, nil
}

// I034/050 - System Configuration and Status (Compound Data Item)
func decodeSystemConfiguration034(data []byte) (*SystemConfiguration034, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I034/050")
	}

	offset := 0
	primarySubfield := data[offset]
	offset++

	result := &SystemConfiguration034{
		COM: (primarySubfield>>7)&0x01 == 1, // Common Part
		PSR: (primarySubfield>>4)&0x01 == 1, // PSR Sensor
		SSR: (primarySubfield>>3)&0x01 == 1, // SSR Sensor
		MDS: (primarySubfield>>2)&0x01 == 1, // Mode S Sensor
		FX:  primarySubfield&0x01 == 1,      // Field Extension
	}

	// Handle subfields based on presence bits
	if (primarySubfield>>7)&0x01 == 1 { // COM subfield present
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing COM subfield data")
		}
		comData := data[offset]
		offset++
		result.COMData = &COMConfigData034{
			NOGO:   (comData>>7)&0x01 == 1,
			RDPC:   (comData>>6)&0x01 == 1,
			RDPR:   (comData>>5)&0x01 == 1,
			OVLRDP: (comData>>4)&0x01 == 1,
			OVLXMT: (comData>>3)&0x01 == 1,
			MSC:    (comData>>2)&0x01 == 1,
			TSV:    (comData>>1)&0x01 == 1,
		}
	}

	// PSR subfield
	if (primarySubfield>>4)&0x01 == 1 {
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing PSR subfield data")
		}
		psrData := data[offset]
		offset++
		result.PSRData = &PSRConfigData034{
			ANT:  (psrData>>7)&0x01 == 1,
			CHAB: (psrData >> 5) & 0x03,
			OVL:  (psrData>>4)&0x01 == 1,
			MSC:  (psrData>>3)&0x01 == 1,
		}
	}

	// SSR subfield
	if (primarySubfield>>3)&0x01 == 1 {
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing SSR subfield data")
		}
		ssrData := data[offset]
		offset++
		result.SSRData = &SSRConfigData034{
			ANT:  (ssrData>>7)&0x01 == 1,
			CHAB: (ssrData >> 5) & 0x03,
			OVL:  (ssrData>>4)&0x01 == 1,
			MSC:  (ssrData>>3)&0x01 == 1,
		}
	}

	// MDS subfield (2 octets)
	if (primarySubfield>>2)&0x01 == 1 {
		if offset+1 >= len(data) {
			return result, offset, fmt.Errorf("missing MDS subfield data")
		}
		mdsData1 := data[offset]
		mdsData2 := data[offset+1]
		offset += 2
		result.MDSData = &MDSConfigData034{
			ANT:    (mdsData1>>7)&0x01 == 1,
			CHAB:   (mdsData1 >> 5) & 0x03,
			OVLSUR: (mdsData1>>4)&0x01 == 1,
			MSC:    (mdsData1>>3)&0x01 == 1,
			SCF:    (mdsData1>>2)&0x01 == 1,
			DLF:    (mdsData1>>1)&0x01 == 1,
			OVLSCF: (mdsData2>>7)&0x01 == 1,
			OVLDLF: (mdsData2>>6)&0x01 == 1,
		}
	}

	return result, offset, nil
}

// I034/060 - System Processing Mode (Compound Data Item)
func decodeSystemProcessingMode034(data []byte) (*SystemProcessingMode034, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I034/060")
	}

	offset := 0
	primarySubfield := data[offset]
	offset++

	result := &SystemProcessingMode034{
		COM: (primarySubfield>>7)&0x01 == 1, // Common Part
		PSR: (primarySubfield>>4)&0x01 == 1, // PSR Sensor
		SSR: (primarySubfield>>3)&0x01 == 1, // SSR Sensor
		MDS: (primarySubfield>>2)&0x01 == 1, // Mode S Sensor
		FX:  primarySubfield&0x01 == 1,      // Field Extension
	}

	// Handle subfields based on presence bits
	if (primarySubfield>>7)&0x01 == 1 { // COM subfield present
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing COM processing mode subfield")
		}
		comData := data[offset]
		offset++
		result.COMData = &COMProcessingData034{
			REDRDP: (comData >> 4) & 0x07, // Reduction steps RDP
			REDXMT: (comData >> 1) & 0x07, // Reduction steps XMT
		}
	}

	// PSR processing mode subfield
	if (primarySubfield>>4)&0x01 == 1 {
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing PSR processing mode subfield")
		}
		psrData := data[offset]
		offset++
		result.PSRData = &PSRProcessingData034{
			POL:    (psrData>>7)&0x01 == 1, // Polarization
			REDRAD: (psrData >> 4) & 0x07,  // Reduction steps
			STC:    (psrData >> 2) & 0x03,  // STC Map
		}
	}

	// SSR processing mode subfield
	if (primarySubfield>>3)&0x01 == 1 {
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing SSR processing mode subfield")
		}
		ssrData := data[offset]
		offset++
		result.SSRData = &SSRProcessingData034{
			REDRAD: (ssrData >> 5) & 0x07, // Reduction steps
		}
	}

	// MDS processing mode subfield
	if (primarySubfield>>2)&0x01 == 1 {
		if offset >= len(data) {
			return result, offset, fmt.Errorf("missing MDS processing mode subfield")
		}
		mdsData := data[offset]
		offset++
		result.MDSData = &MDSProcessingData034{
			REDRAD: (mdsData >> 5) & 0x07,  // Reduction steps
			CLU:    (mdsData>>4)&0x01 == 1, // Cluster state
		}
	}

	return result, offset, nil
}

// I034/070 - Message Count Values (Repetitive Data Item)
func decodeMessageCountValues034(data []byte) (*MessageCountValues034, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I034/070")
	}

	rep := data[0] // Repetition factor
	offset := 1
	counters := make([]MessageCounter034, 0, rep)

	for i := uint8(0); i < rep && offset+1 < len(data); i++ {
		if offset+2 > len(data) {
			break
		}

		// Extract 2 bytes: TYP (5 bits) + COUNTER (11 bits)
		word := binary.BigEndian.Uint16(data[offset : offset+2])
		typ := (word >> 11) & 0x1F // bits 16-12
		counter := word & 0x7FF    // bits 11-1

		counters = append(counters, MessageCounter034{
			Type:    typ,
			Counter: counter,
		})

		offset += 2
	}

	return &MessageCountValues034{
		Repetition: rep,
		Counters:   counters,
	}, offset, nil
}

// I034/100 - Generic Polar Window
func decodeGenericPolarWindow034(data []byte) (*GenericPolarWindow034, int, error) {
	if len(data) < 8 {
		return nil, 0, fmt.Errorf("too short for I034/100")
	}

	rhoStart := binary.BigEndian.Uint16(data[0:2])   // LSB = 1/256 NM
	rhoEnd := binary.BigEndian.Uint16(data[2:4])     // LSB = 1/256 NM
	thetaStart := binary.BigEndian.Uint16(data[4:6]) // LSB = 360°/2^16
	thetaEnd := binary.BigEndian.Uint16(data[6:8])   // LSB = 360°/2^16

	return &GenericPolarWindow034{
		RhoStartNM:    float64(rhoStart) / 256.0,
		RhoEndNM:      float64(rhoEnd) / 256.0,
		ThetaStartDeg: float64(thetaStart) * 360.0 / 65536.0,
		ThetaEndDeg:   float64(thetaEnd) * 360.0 / 65536.0,
	}, 8, nil
}

// I034/110 - Data Filter
func decodeDataFilter034(data []byte) (*DataFilter034, int, error) {
	if len(data) < 1 {
		return nil, 0, fmt.Errorf("too short for I034/110")
	}

	filterType := data[0]
	return &DataFilter034{
		Type: filterType,
	}, 1, nil
}

// I034/120 - 3D-Position of Data Source
func decode3DPositionOfSource034(data []byte) (*Position3D034, int, error) {
	if len(data) < 8 {
		return nil, 0, fmt.Errorf("too short for I034/120")
	}

	// Height (16 bits, signed, LSB = 1 meter)
	heightRaw := int16(binary.BigEndian.Uint16(data[0:2]))
	height := float64(heightRaw)

	// Latitude (24 bits, signed, LSB = 180/2^23 degrees)
	latRaw := int32(data[2])<<16 | int32(data[3])<<8 | int32(data[4])
	if latRaw&0x800000 != 0 { // Sign extend 24-bit to 32-bit
		latRaw = latRaw - 0x1000000 // Two's complement for 24-bit signed
	}
	latitude := float64(latRaw) * 180.0 / (1 << 23)

	// Longitude (24 bits, signed, LSB = 180/2^23 degrees)
	lonRaw := int32(data[5])<<16 | int32(data[6])<<8 | int32(data[7])
	if lonRaw&0x800000 != 0 { // Sign extend 24-bit to 32-bit
		lonRaw = lonRaw - 0x1000000 // Two's complement for 24-bit signed
	}
	longitude := float64(lonRaw) * 180.0 / (1 << 23)

	return &Position3D034{
		HeightM:      height,
		LatitudeDeg:  latitude,
		LongitudeDeg: longitude,
	}, 8, nil
}

// I034/090 - Collimation Error
func decodeCollimationError034(data []byte) (*CollimationError034, int, error) {
	if len(data) < 2 {
		return nil, 0, fmt.Errorf("too short for I034/090")
	}

	// Range error (8 bits, signed, LSB = 1/128 NM)
	rangeErrorRaw := int8(data[0])
	rangeError := float64(rangeErrorRaw) / 128.0

	// Azimuth error (8 bits, signed, LSB = 360°/2^14)
	azimuthErrorRaw := int8(data[1])
	azimuthError := float64(azimuthErrorRaw) * 360.0 / (1 << 14)

	return &CollimationError034{
		RangeErrorNM:    rangeError,
		AzimuthErrorDeg: azimuthError,
	}, 2, nil
}

// I034/RE - Reserved Expansion Field
func decodeReservedExpansion034(data []byte) (*ReservedField034, int, error) {
	return &ReservedField034{
		Note: "Reserved Expansion Field - implementation specific",
		Data: hex.EncodeToString(data),
	}, len(data), nil
}

// I034/SP - Special Purpose Field
func decodeSpecialPurpose034(data []byte) (*SpecialPurposeField034, int, error) {
	return &SpecialPurposeField034{
		Note: "Special Purpose Field - implementation specific",
		Data: hex.EncodeToString(data),
	}, len(data), nil
}
