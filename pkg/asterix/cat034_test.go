package asterix

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

// Real-world CAT 034 test messages
var realWorldCAT034Messages = []string{
	"220012f62821025460022084404600840000", // Simple CAT 034 with basic fields
	"220012f62821025460172884404600840000", // CAT 034 with different time
	"220012f628210254602a3084404600840000", // CAT 034 variant
	"220012f628210254603e3884404600840000", // Another CAT 034
	"220012f62821025460524084404600840000", // CAT 034 with different sequence
}

// countCAT034Fields counts the number of non-nil pointer fields in a CAT034Message.
func countCAT034Fields(c *CAT034Message) int {
	if c == nil {
		return 0
	}
	v := reflect.ValueOf(c).Elem()
	count := 0
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Ptr && !f.IsNil() {
			count++
		}
	}
	return count
}

func TestCAT034RealWorldDecoding(t *testing.T) {
	tests := []struct {
		name      string
		hexData   string
		wantSIC   uint8
		wantSAC   uint8
		wantCat   byte
		minFields int
	}{
		{
			name:      "CAT 034 Real Message 1",
			hexData:   realWorldCAT034Messages[0],
			wantSIC:   33,
			wantSAC:   40,
			wantCat:   34,
			minFields: 5,
		},
		{
			name:      "CAT 034 Real Message 2",
			hexData:   realWorldCAT034Messages[1],
			wantSIC:   33,
			wantSAC:   40,
			wantCat:   34,
			minFields: 5,
		},
		{
			name:      "CAT 034 Real Message 3",
			hexData:   realWorldCAT034Messages[2],
			wantSIC:   33,
			wantSAC:   40,
			wantCat:   34,
			minFields: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hexData)
			if err != nil {
				t.Fatalf("Failed to decode hex data: %v", err)
			}

			msg, err := Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("Failed to decode ASTERIX message: %v", err)
			}

			if msg.Category != tt.wantCat {
				t.Errorf("Category = %d, want %d", msg.Category, tt.wantCat)
			}

			if msg.Record.GetSIC() != tt.wantSIC {
				t.Errorf("SIC = %d, want %d", msg.Record.GetSIC(), tt.wantSIC)
			}

			if msg.Record.GetSAC() != tt.wantSAC {
				t.Errorf("SAC = %d, want %d", msg.Record.GetSAC(), tt.wantSAC)
			}

			cat034, ok := msg.Record.(*CAT034Message)
			if !ok {
				t.Fatal("expected Record to be *CAT034Message")
			}

			fieldCount := countCAT034Fields(cat034)
			if fieldCount < tt.minFields {
				t.Errorf("Cat034 field count = %d, want at least %d", fieldCount, tt.minFields)
			}

			// Verify mandatory fields exist
			if cat034.DataSourceIdentifier == nil {
				t.Error("Missing mandatory field DataSourceIdentifier (I034/010)")
			}

			if cat034.MessageType == nil {
				t.Error("Missing mandatory field MessageType (I034/000)")
			}

			t.Logf("Decoded %d fields in Cat034", fieldCount)
		})
	}
}

func TestCAT034MessageParsing(t *testing.T) {
	data, _ := hex.DecodeString(realWorldCAT034Messages[0])

	rawMsg, err := ParseMessage(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to parse message: %v", err)
	}

	if rawMsg.Category != 34 {
		t.Errorf("Category = %d, want 34", rawMsg.Category)
	}

	if rawMsg.Length != uint16(len(data)) {
		t.Errorf("Length = %d, want %d", rawMsg.Length, len(data))
	}

	if len(rawMsg.Payload) != len(data)-3 {
		t.Errorf("Payload length = %d, want %d", len(rawMsg.Payload), len(data)-3)
	}
}

func TestCAT034FSPECParsing(t *testing.T) {
	tests := []struct {
		name     string
		hexData  string
		category byte
		minFSPEC int
	}{
		{
			name:     "CAT 034 Basic FSPEC",
			hexData:  realWorldCAT034Messages[0],
			category: 34,
			minFSPEC: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := hex.DecodeString(tt.hexData)
			msg, err := Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("Failed to decode message: %v", err)
			}

			if len(msg.FSPEC) < tt.minFSPEC {
				t.Errorf("FSPEC length = %d, want at least %d", len(msg.FSPEC), tt.minFSPEC)
			}

			t.Logf("FSPEC: %x", msg.FSPEC)
		})
	}
}
