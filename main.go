package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/aperrot42/panoramix/pkg/asterix"
	"github.com/aperrot42/panoramix/pkg/internal_format"
	"github.com/aperrot42/panoramix/pkg/modes/bds"
	"github.com/aperrot42/panoramix/pkg/transform/fspec"
	"github.com/aperrot42/panoramix/pkg/transform/position"
)

// PrecisionTime provides millisecond-precision JSON marshaling for timestamps
type PrecisionTime struct {
	time.Time
}

// MarshalJSON implements custom JSON marshaling with millisecond precision
func (pt PrecisionTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pt.Time.UTC().Format("2006-01-02T15:04:05.000Z") + `"`), nil
}

type OutputMessage struct {
	MessageNumber int                            `json:"message_number"`
	Timestamp     PrecisionTime                  `json:"timestamp"`
	Category      byte                           `json:"category"`
	Port          int                            `json:"port"`
	SIC           uint8                          `json:"sic"`
	SAC           uint8                          `json:"sac"`
	Record        any                            `json:"record"`
	FSPEC         string                         `json:"fspec,omitempty"`
	BDS           map[string]bds.DecodedRegister `json:"bds,omitempty"`
	Position      *position.Position             `json:"position,omitempty"`
	FSPECFields   []string                       `json:"fspec_fields,omitempty"`
}

func outputJSON(data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

func outputText(data any) {
	switch v := data.(type) {
	case OutputMessage:
		fmt.Printf("Msg %3d: CAT=%d Port=%d SIC=%d SAC=%d Time=%s\n",
			v.MessageNumber, v.Category, v.Port, v.SIC, v.SAC,
			v.Timestamp.Format("15:04:05.000"))
	case position.Position:
		fmt.Printf("Aircraft: SIC=%d SAC=%d Time=%s WGS84=(%.6f,%.6f) Alt=%.0fm\n",
			v.RadarSIC, v.RadarSAC, v.Timestamp.Format("15:04:05.000"),
			v.WGS84Position.Latitude_deg, v.WGS84Position.Longitude_deg,
			v.WGS84Position.AltitudeFt)
	default:
		fmt.Printf("Unknown data type: %T\n", data)
	}
}

func main() {
	filename := flag.String("filename", "recording.ast", "Input .if radar recording file")
	limit := flag.Int("limit", 0, "Maximum number of messages to parse (0 = unlimited)")
	jsonOutput := flag.Bool("json", false, "Output as JSON instead of text")
	positionFilter := flag.Bool("position", false, "Add computed position information to raw ASTERIX messages")
	fspecFields := flag.Bool("fspec-fields", false, "Add computed FSPEC available fields list to messages")
	radarConfig := flag.String("radar-config", "radar_config.yaml", "Radar configuration file path")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Extract date from filename for timestamp base
	baseDate, err := extractDateFromFilename(*filename)
	if err != nil {
		log.Printf("Warning: Could not extract date from filename, using current date: %v", err)
		baseDate = time.Now().UTC()
	}

	reader := internal_format.NewReaderWithBaseDate(file, baseDate)

	var positionExtractor *position.PositionExtractor
	if *positionFilter {
		var err error
		positionExtractor, err = position.NewPositionExtractorFromConfig(*radarConfig)
		if err != nil {
			log.Fatalf("Failed to initialize position extractor with config: %v", err)
		}
	}

	for {
		if *limit > 0 && reader.Count() >= *limit {
			break
		}

		// Read next .if record
		record, err := reader.ReadRecord()
		if err != nil {
			break // EOF or read error
		}
		if record == nil {
			continue // Skip invalid records
		}

		// Decode ASTERIX message from payload
		asterixMsg, err := asterix.Decode(bytes.NewReader(record.Payload))
		if err != nil {
			continue // Skip decode errors
		}

		outputMsg := OutputMessage{
			MessageNumber: record.MessageNumber,
			Timestamp:     PrecisionTime{record.Timestamp},
			Category:      asterixMsg.Category,
			Port:          record.Port,
			SIC:           asterixMsg.Record.GetSIC(),
			SAC:           asterixMsg.Record.GetSAC(),
			Record:        asterixMsg.Record,
			FSPEC:         hex.EncodeToString(asterixMsg.FSPEC),
		}

		// Decode BDS registers for CAT 048 messages
		if cat048, ok := asterixMsg.Record.(*asterix.Cat048Message); ok && cat048.BDSRegister != nil {
			outputMsg.BDS = bds.DecodeRegisters(cat048.BDSRegister)
		}

		// Apply FSPEC transform if requested
		if *fspecFields {
			outputMsg.FSPECFields = fspec.ComputeAvailableFields(asterixMsg)
		}

		if *positionFilter && positionExtractor != nil {
			pos, err := positionExtractor.ExtractFromMessage(asterixMsg, record.Timestamp)
			if err != nil {
				log.Printf("Position extraction failed: %v", err)
			}
			outputMsg.Position = pos
		}
		if *jsonOutput {
			outputJSON(outputMsg)
		} else {
			outputText(outputMsg)
		}
	}
}

// extractDateFromFilename parses date from .if filename format
// Expected format: ...YYYY-MM-DD_HH[h]MM[m]SS[s]_...
// Example: ANY-OPS_DOLS_A_IP1_G_2025-05-15_12h00m00s_2025-05-15_18h00m59s.if
func extractDateFromFilename(filename string) (time.Time, error) {
	// Look for YYYY-MM-DD pattern in filename
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return time.Time{}, fmt.Errorf("no date found in filename: %s", filename)
	}

	dateStr := matches[1]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date %s: %v", dateStr, err)
	}

	return date.UTC(), nil
}
