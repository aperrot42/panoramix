package asterix

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

func TestCAT034TypedDecoding(t *testing.T) {
	// Use hex data from the working test
	hexData := "220012f62821025460022084404600840000"
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex data: %v", err)
	}

	// Parse the full message first to get the correct payload
	msg, err := Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}
	
	// Extract the raw message from the parsed result
	rawMsg := &RawAsterixMessage{
		Category: msg.Category,
		Payload:  data[3:], // Skip header (category + length)
	}

	t.Run("CAT 034 Typed Decoding", func(t *testing.T) {
		decoder := &CAT034Decoder{}

		// Test typed decoding
		typedMsg, err := decoder.DecodeTyped(rawMsg)
		if err != nil {
			t.Fatalf("DecodeTyped failed: %v", err)
		}

		// Verify basic structure
		if typedMsg.Category != 34 {
			t.Errorf("Expected category 34, got %d", typedMsg.Category)
		}

		// Verify that SAC/SIC are populated from DataSourceIdentifier
		if typedMsg.DataSourceIdentifier == nil {
			t.Error("DataSourceIdentifier should not be nil")
		} else {
			if typedMsg.SAC != typedMsg.DataSourceIdentifier.SAC {
				t.Errorf("SAC mismatch: top-level %d vs DataSourceIdentifier %d", 
					typedMsg.SAC, typedMsg.DataSourceIdentifier.SAC)
			}
			if typedMsg.SIC != typedMsg.DataSourceIdentifier.SIC {
				t.Errorf("SIC mismatch: top-level %d vs DataSourceIdentifier %d", 
					typedMsg.SIC, typedMsg.DataSourceIdentifier.SIC)
			}
		}

		// Verify MessageType is present
		if typedMsg.MessageType == nil {
			t.Error("MessageType should not be nil")
		}

		// Verify TimeOfDay is present and is a valid duration
		if typedMsg.TimeOfDay == nil {
			t.Error("TimeOfDay should not be nil")
		} else {
			if *typedMsg.TimeOfDay < 0 || *typedMsg.TimeOfDay > 24*time.Hour {
				t.Errorf("TimeOfDay out of valid range: %v", *typedMsg.TimeOfDay)
			}
		}

		t.Logf("Successfully decoded typed CAT 034 message with SAC=%d, SIC=%d", 
			typedMsg.SAC, typedMsg.SIC)
	})
}

func TestCAT034TypedVsOriginalDecoding(t *testing.T) {
	// Test that DecodeTyped and Decode produce consistent SAC/SIC values
	hexData := "220012f62821025460022084404600840000"
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex data: %v", err)
	}
	
	decoder := &CAT034Decoder{}
	rawMsg := &RawAsterixMessage{
		Category: 34,
		Payload:  data[3:], // Skip category and length
	}

	// Original decoding
	originalMsg, err := decoder.Decode(rawMsg)
	if err != nil {
		t.Fatalf("Original Decode failed: %v", err)
	}

	// Typed decoding
	typedMsg, err := decoder.DecodeTyped(rawMsg)
	if err != nil {
		t.Fatalf("DecodeTyped failed: %v", err)
	}

	// Compare SAC/SIC values
	if originalMsg.Sac != typedMsg.SAC {
		t.Errorf("SAC mismatch: original %d vs typed %d", originalMsg.Sac, typedMsg.SAC)
	}
	if originalMsg.Sic != typedMsg.SIC {
		t.Errorf("SIC mismatch: original %d vs typed %d", originalMsg.Sic, typedMsg.SIC)
	}
	if originalMsg.Category != typedMsg.Category {
		t.Errorf("Category mismatch: original %d vs typed %d", originalMsg.Category, typedMsg.Category)
	}

	t.Logf("Both decoders produce consistent results: SAC=%d, SIC=%d, Category=%d", 
		originalMsg.Sac, originalMsg.Sic, originalMsg.Category)
}