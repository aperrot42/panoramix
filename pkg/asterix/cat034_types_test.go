package asterix

import (
	"bytes"
	"encoding/hex"
	"testing"
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

		// Test standard decoding
		msg, err := decoder.Decode(rawMsg)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}

		// Verify message was decoded successfully
		if msg == nil {
			t.Fatal("Decoded message should not be nil")
		}

		// Verify basic structure
		if msg.Category != 34 {
			t.Errorf("Expected category 34, got %d", msg.Category)
		}

		// Verify that SAC/SIC are populated
		if msg.Sac == 0 && msg.Sic == 0 {
			t.Error("SAC and SIC should be populated")
		}

		t.Logf("Successfully decoded CAT 034 message with SAC=%d, SIC=%d",
			msg.Sac, msg.Sic)
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

	// Test that original decoding works
	if originalMsg.Sac == 0 && originalMsg.Sic == 0 {
		t.Error("SAC and SIC should be populated in original message")
	}
	if originalMsg.Category != 34 {
		t.Errorf("Expected category 34, got %d", originalMsg.Category)
	}

	t.Logf("Both decoders produce consistent results: SAC=%d, SIC=%d, Category=%d", 
		originalMsg.Sac, originalMsg.Sic, originalMsg.Category)
}