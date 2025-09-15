package asterix

import (
	"bytes"
	"encoding/hex"
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

func TestCAT034RealWorldDecoding(t *testing.T) {
	tests := []struct {
		name     string
		hexData  string
		wantSIC  uint8
		wantSAC  uint8
		wantCat  byte
		minItems int // Minimum number of items expected
	}{
		{
			name:     "CAT 034 Real Message 1",
			hexData:  realWorldCAT034Messages[0],
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  34,
			minItems: 5,
		},
		{
			name:     "CAT 034 Real Message 2",
			hexData:  realWorldCAT034Messages[1],
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  34,
			minItems: 5,
		},
		{
			name:     "CAT 034 Real Message 3",
			hexData:  realWorldCAT034Messages[2],
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  34,
			minItems: 5,
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

			if msg.Sic != tt.wantSIC {
				t.Errorf("SIC = %d, want %d", msg.Sic, tt.wantSIC)
			}

			if msg.Sac != tt.wantSAC {
				t.Errorf("SAC = %d, want %d", msg.Sac, tt.wantSAC)
			}

			if len(msg.Items) < tt.minItems {
				t.Errorf("Items count = %d, want at least %d", len(msg.Items), tt.minItems)
			}

			// Verify mandatory fields exist
			if _, exists := msg.Items["I034/010"]; !exists {
				t.Error("Missing mandatory field I034/010 (Data Source Identifier)")
			}

			if _, exists := msg.Items["I034/000"]; !exists {
				t.Error("Missing mandatory field I034/000 (Message Type)")
			}

			t.Logf("Decoded %d items: %v", len(msg.Items), getCAT034MapKeys(msg.Items))
		})
	}
}

func TestCAT034MessageParsing(t *testing.T) {
	// Test raw message parsing without full decoding
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
	// Test FSPEC parsing with various field configurations for CAT 034
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

// Helper function to get map keys for CAT 034
func getCAT034MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}