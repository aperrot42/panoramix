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

	cat034, ok := msg.Record.(*CAT034Message)
	if !ok {
		t.Fatal("expected Record to be *CAT034Message")
	}

	if cat034.DataSourceIdentifier == nil {
		t.Fatal("DataSourceIdentifier is nil")
	}

	if msg.Record.GetSAC() != 40 {
		t.Errorf("SAC = %d, want 40", msg.Record.GetSAC())
	}
	if msg.Record.GetSIC() != 33 {
		t.Errorf("SIC = %d, want 33", msg.Record.GetSIC())
	}

	t.Logf("Successfully decoded CAT 034 message with SAC=%d, SIC=%d", msg.Record.GetSAC(), msg.Record.GetSIC())
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

	cat034, ok := msg.Record.(*CAT034Message)
	if !ok {
		t.Fatal("expected Record to be *CAT034Message")
	}
	if cat034.DataSourceIdentifier == nil {
		t.Fatal("DataSourceIdentifier is nil")
	}
	if msg.Record.GetSAC() != 40 || msg.Record.GetSIC() != 33 {
		t.Errorf("SAC=%d SIC=%d, want SAC=40 SIC=33", msg.Record.GetSAC(), msg.Record.GetSIC())
	}
	if msg.Category != 34 {
		t.Errorf("Expected category 34, got %d", msg.Category)
	}

	t.Logf("Decoded CAT 034: SAC=%d, SIC=%d, Category=%d", msg.Record.GetSAC(), msg.Record.GetSIC(), msg.Category)
}