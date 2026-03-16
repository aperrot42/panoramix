package asterix

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestCAT034TypedDecoding(t *testing.T) {
	hexData := "220012f62821025460022084404600840000"
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex data: %v", err)
	}

	msg, err := Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	if msg.Category != 34 {
		t.Errorf("Expected category 34, got %d", msg.Category)
	}

	if msg.Cat034 == nil {
		t.Fatal("Cat034 is nil")
	}

	if msg.Cat034.DataSourceIdentifier == nil {
		t.Fatal("DataSourceIdentifier is nil")
	}

	if msg.Sac != 40 {
		t.Errorf("SAC = %d, want 40", msg.Sac)
	}
	if msg.Sic != 33 {
		t.Errorf("SIC = %d, want 33", msg.Sic)
	}

	t.Logf("Successfully decoded CAT 034 message with SAC=%d, SIC=%d", msg.Sac, msg.Sic)
}

func TestCAT034TypedVsOriginalDecoding(t *testing.T) {
	hexData := "220012f62821025460022084404600840000"
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex data: %v", err)
	}

	decoder := &CAT034Decoder{}
	rawMsg := &RawAsterixMessage{
		Category: 34,
		Payload:  data[3:],
	}

	msg, err := decoder.Decode(rawMsg)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if msg.Cat034 == nil {
		t.Fatal("Cat034 is nil")
	}
	if msg.Cat034.DataSourceIdentifier == nil {
		t.Fatal("DataSourceIdentifier is nil")
	}
	if msg.Sac != 40 || msg.Sic != 33 {
		t.Errorf("SAC=%d SIC=%d, want SAC=40 SIC=33", msg.Sac, msg.Sic)
	}
	if msg.Category != 34 {
		t.Errorf("Expected category 34, got %d", msg.Category)
	}

	t.Logf("Decoded CAT 034: SAC=%d, SIC=%d, Category=%d", msg.Sac, msg.Sic, msg.Category)
}