package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aperrot42/panoramix/pkg/asterix"
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
	MessageNumber int                   `json:"message_number"`
	Category      byte                  `json:"category"`
	SIC           uint8                 `json:"sic"`
	SAC           uint8                 `json:"sac"`
	Record        any                   `json:"record"`
	FSPEC         string                `json:"fspec,omitempty"`
	BDS           *bds.DecodedRegisters `json:"bds,omitempty"`
	Position      *position.Position    `json:"position,omitempty"`
	FSPECFields   []string              `json:"fspec_fields,omitempty"`
}

func outputJSON(data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

func outputText(v OutputMessage) {
	fmt.Printf("Msg %3d: CAT=%d SIC=%d SAC=%d\n",
		v.MessageNumber, v.Category, v.SIC, v.SAC)
}

func main() {
	filename := flag.String("filename", "", "Input file containing raw ASTERIX data")
	limit := flag.Int("limit", 0, "Maximum number of messages to decode (0 = unlimited)")
	jsonOutput := flag.Bool("json", false, "Output as JSON instead of text")
	positionFilter := flag.Bool("position", false, "Add computed WGS84 position to CAT 048 messages")
	fspecFields := flag.Bool("fspec-fields", false, "Add FSPEC available fields list to messages")
	radarConfig := flag.String("radar-config", "radar_config.yaml", "Radar configuration file path")
	flag.Parse()

	if *filename == "" {
		fmt.Fprintln(os.Stderr, "Usage: panoramix --filename <asterix-file> [--json] [--limit N] [--position] [--fspec-fields]")
		os.Exit(1)
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var positionExtractor *position.PositionExtractor
	if *positionFilter {
		positionExtractor, err = position.NewPositionExtractorFromConfig(*radarConfig)
		if err != nil {
			log.Fatalf("Failed to initialize position extractor: %v", err)
		}
	}

	count := 0
	for {
		if *limit > 0 && count >= *limit {
			break
		}

		asterixMsg, err := asterix.Decode(file)
		if err != nil {
			if count == 0 {
				log.Fatalf("Failed to decode first message: %v", err)
			}
			break // EOF or unrecoverable error
		}

		count++

		outputMsg := OutputMessage{
			MessageNumber: count,
			Category:      asterixMsg.Category,
			SIC:           asterixMsg.Record.GetSIC(),
			SAC:           asterixMsg.Record.GetSAC(),
			Record:        asterixMsg.Record,
			FSPEC:         hex.EncodeToString(asterixMsg.FSPEC),
		}

		// Decode BDS registers for CAT 048 messages
		if cat048, ok := asterixMsg.Record.(*asterix.Cat048Message); ok && cat048.BDSRegister != nil {
			raw := make([]bds.RawRegister, 0, len(cat048.BDSRegister.Registers))
			for _, reg := range cat048.BDSRegister.Registers {
				raw = append(raw, bds.RawRegister{Code: reg.BDSCode, Data: reg.RawData})
			}
			var decoded bds.DecodedRegisters
			decoded.DecodeAll(raw)
			outputMsg.BDS = &decoded
		}

		// Apply FSPEC transform if requested
		if *fspecFields {
			outputMsg.FSPECFields = fspec.ComputeAvailableFields(asterixMsg)
		}

		if *positionFilter && positionExtractor != nil {
			if cat048, ok := asterixMsg.Record.(*asterix.Cat048Message); ok {
				plot := position.RawPlot{
					SIC: cat048.GetSIC(),
					SAC: cat048.GetSAC(),
				}
				if cat048.MeasuredPosition != nil {
					plot.Polar = &position.PolarCoord{
						RhoNM:    cat048.MeasuredPosition.Rho,
						ThetaDeg: cat048.MeasuredPosition.Theta,
					}
				}
				if cat048.CalculatedPosition != nil {
					plot.Cart = &position.CartCoord{
						XNM: cat048.CalculatedPosition.X,
						YNM: cat048.CalculatedPosition.Y,
					}
				}
				if cat048.Height3D != nil {
					h := cat048.Height3D.Height
					plot.Alt3D = &h
				}
				if cat048.FlightLevel != nil {
					fl := cat048.FlightLevel.FL
					plot.FL = &fl
				}
				pos, err := positionExtractor.Extract(plot, time.Now())
				if err != nil {
					log.Printf("Position extraction failed: %v", err)
				}
				outputMsg.Position = pos
			}
		}

		if *jsonOutput {
			outputJSON(outputMsg)
		} else {
			outputText(outputMsg)
		}
	}
}
