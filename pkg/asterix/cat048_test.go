package asterix

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// Real-world CAT 048 test messages
var realWorldCAT048Messages = []string{
	// Short CAT 048 message with 15 fields
	"300046ffff022821545fe8a8602a1648220001906003c14b1a2b15a4f3d8c42003ce599116e01c0060fff9c11f7ff4745093880030aa000040080f19012912046e9aa84020fd",
	// Medium CAT 048 message
	"300070ffff022821545ff2a8627a1a142ddb01886006bd4b1a1a15a4f1cc12e00390b68123217c8450b369cb19e3a4726095800030aa00004004321d68277e03854eac4220fdffdf022821545ff3a82e221a700e00011d6005bd4b51db202cf3e3082008ed0df11261010b2f44402060",
	// Compact CAT 048 message
	"300046ffff022821545ff7a8626e1bdc220004886005bf4932a338a174db90e003c65a392e7ff80060fff8d93bc004da50b8a80030aa01454005bb1f17262708878d084020f5",
	// Single CAT 048 message for basic testing
	"300046ffff022821545fffa88cd51f4426f805f26005bc345691242170dc5520038579f331602c0060ca380030a40000408010eb34e004e15001c630e332af07b40b5c4020f5",
	// Large CAT 048 message with enhanced Mode S fields (I048/250) - 305 bytes
	"300131ffff022821545ffaa870a11d14220005c86005c4501c3b0d43b2e4682003ffd4f7396004e050a609f930a00c0060c8480030a4000040046e24de2a9308094f584020f5ff1e2821545ffc4894f51d64018904d06005c005d9313037ed0859c28440ffff122821545ffba86d791dcc249b05776004c640815c15a678d0c12003fff57f39a004dc50a9aa0930204c0160c4600030aa000040005f248f28be08445858400661000120f5ffff022821545ffca816f71e102af4054f6003cb451e9034257419282003801fd3352004e150ff5a1d30e00c0260c2680030a4000040030807ba087f07dcf6b04020fdffff022821545ffaa816fb1d60205c05286002d13c0ca751527731882003803a2b3b6004df50d12a21303fffff60c07e4270a800004009f1079608a2086ea1a84020f5",
	// Very large CAT 048 message with enhanced Mode S fields - 448 bytes
	"3001c0ffff422821546016a8738a2848260600ce6004bc4b18f51445f2c724a003f01a9d12fd2c5250d2a92f0fa1f43b6097700030a400004003a030411fc4027cb690422220fdffff022821546016a869e528342b9605c76004c644014114a578cc52e00384f5ad3b6054e650abaa0b3260140060c8480030aa00004000da2c2c1d33087259804020f5ffff022821546016a86d9b2898260201026005c04b16134d74b35cc82003803a6f22200c8a50d3c9f91ae3045e6097700030a400004005b52e021dc8048da4f84220f5ffff022821546017a82a6428dc2bd402d76005ca781e130c3078db2820038052192fe004c8508efa6528220c4560aee00030aa00004004e711de0b6706c420204020fdffdf022821546016a81a5f28602e00011e6005cb4b4db6202cb3c7982006590b08073a00c25a244020e0ffff022821546017a84ea028c02e0000c26005bf4b1b132022c7420820038373510a4004295098189707ff88006089c80030a401044002402116153d018a33104420e0ffff022821546017a8751c28a0220002806003c1345692242170d926600399680030aa000040f5dd2726ff0ca550e89a09213e87d1600a8c312e1fc80609dc344420f5",
}

func TestCAT048RealWorldDecoding(t *testing.T) {
	tests := []struct {
		name     string
		hexData  string
		wantSIC  uint8
		wantSAC  uint8
		wantCat  byte
		minItems int
	}{
		{
			name:     "CAT 048 Real Short Message",
			hexData:  realWorldCAT048Messages[0],
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  48,
			minItems: 10,
		},
		{
			name:     "CAT 048 Real Medium Message",
			hexData:  realWorldCAT048Messages[1],
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  48,
			minItems: 10,
		},
		{
			name:     "CAT 048 Real Compact Message",
			hexData:  realWorldCAT048Messages[2],
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  48,
			minItems: 10,
		},
		{
			name:     "CAT 048 Large Message with Enhanced Mode S",
			hexData:  realWorldCAT048Messages[4], // 305 bytes with I048/250
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  48,
			minItems: 10,
		},
		{
			name:     "CAT 048 Very Large Message with Enhanced Mode S",
			hexData:  realWorldCAT048Messages[5], // 448 bytes with I048/250
			wantSIC:  33,
			wantSAC:  40,
			wantCat:  48,
			minItems: 10,
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
			if _, exists := msg.Items["I048/010"]; !exists {
				t.Error("Missing mandatory field I048/010 (Data Source Identifier)")
			}

			if _, exists := msg.Items["I048/140"]; !exists {
				t.Error("Missing mandatory field I048/140 (Time of Day)")
			}

			t.Logf("Decoded %d items: %v", len(msg.Items), getCAT048MapKeys(msg.Items))
		})
	}
}

func TestCAT048SpecificFieldDecoding(t *testing.T) {
	// Test specific field types and values
	data, _ := hex.DecodeString(realWorldCAT048Messages[0])
	msg, err := Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	// Test Data Source Identifier (I048/010)
	if dsi, exists := msg.Items["I048/010"]; exists {
		if dsiMap, ok := dsi.(map[string]interface{}); ok {
			if sac, ok := dsiMap["SAC"].(uint8); ok {
				if sac != 40 {
					t.Errorf("SAC = %d, want 40", sac)
				}
			}
			if sic, ok := dsiMap["SIC"].(uint8); ok {
				if sic != 33 {
					t.Errorf("SIC = %d, want 33", sic)
				}
			}
		}
	}

	// Test Time of Day exists and is reasonable
	if tod, exists := msg.Items["I048/140"]; exists {
		if todVal, ok := tod.(float64); ok {
			// Time should be positive and reasonable (0-86400 seconds in a day)
			if todVal < 0 || todVal > 86400 {
				t.Errorf("Time of Day = %f, should be between 0 and 86400", todVal)
			}
		}
	}
}

func TestCAT048EnhancedModeSFields(t *testing.T) {
	// Test large messages with enhanced Mode S fields
	tests := []struct {
		name                string
		hexData             string
		expectedFields      []string
		expectedMessageSize int
	}{
		{
			name:                "Large CAT 048 with I048/250",
			hexData:             realWorldCAT048Messages[4], // 305 bytes
			expectedFields:      []string{"I048/250"},       // Mode S MB Data
			expectedMessageSize: 305,
		},
		{
			name:                "Very Large CAT 048 with I048/250",
			hexData:             realWorldCAT048Messages[5], // 448 bytes
			expectedFields:      []string{"I048/250"},       // Mode S MB Data
			expectedMessageSize: 448,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hexData)
			if err != nil {
				t.Fatalf("Failed to decode hex data: %v", err)
			}

			// Verify message size
			if len(data) != tt.expectedMessageSize {
				t.Errorf("Message size = %d, want %d", len(data), tt.expectedMessageSize)
			}

			msg, err := Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("Failed to decode ASTERIX message: %v", err)
			}

			// Verify enhanced Mode S fields are present
			for _, field := range tt.expectedFields {
				if _, exists := msg.Items[field]; !exists {
					t.Errorf("Missing expected enhanced Mode S field: %s", field)
				}
			}

			// Test that I048/250 (Mode S MB Data) contains expected structure
			if mbData, exists := msg.Items["I048/250"]; exists {
				if mbDataSlice, ok := mbData.([]interface{}); ok {
					if len(mbDataSlice) == 0 {
						t.Error("I048/250 Mode S MB Data is empty")
					}
					t.Logf("I048/250 contains %d Mode S reports", len(mbDataSlice))
				}
			}

			t.Logf("Successfully decoded %d-byte message with %d fields: %v",
				len(data), len(msg.Items), getCAT048MapKeys(msg.Items))
		})
	}
}

func TestCAT048FSPECParsing(t *testing.T) {
	// Test FSPEC parsing with various field configurations for CAT 048
	tests := []struct {
		name     string
		hexData  string
		category byte
		minFSPEC int
	}{
		{
			name:     "CAT 048 Extended FSPEC",
			hexData:  realWorldCAT048Messages[0],
			category: 48,
			minFSPEC: 3, // Should have multiple FSPEC octets
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

// Helper function to get map keys for CAT 048
func getCAT048MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
