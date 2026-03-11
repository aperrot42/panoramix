package asterix

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

// mockBDSDecoder is a test double that records calls and returns a fixed result.
type mockBDSDecoder struct {
	calls []byte // recorded bdsRegisterAddress values
}

func (m *mockBDSDecoder) DecodeBDS(bdsRegisterAddress byte, data []byte) (any, error) {
	m.calls = append(m.calls, bdsRegisterAddress)
	bds1 := (bdsRegisterAddress >> 4) & 0x0F
	bds2 := bdsRegisterAddress & 0x0F
	return map[string]any{
		"mock":    true,
		"bds_key": fmt.Sprintf("%d_%d", bds1, bds2),
	}, nil
}

// decodeCAT048WithDecoder is a helper to decode a multi-record hex message
// using a CAT048Decoder with the given BDSDecoder, returning the first message
// that contains I048/250.
func decodeCAT048WithDecoder(t *testing.T, hexData string, bdsDecoder BDSDecoder) *AsterixMessage {
	t.Helper()

	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex data: %v", err)
	}

	decoder := &CAT048Decoder{BDSDecoder: bdsDecoder}
	RegisterDecoder(48, decoder)
	defer RegisterDecoder(48, &CAT048Decoder{}) // restore default

	reader := bytes.NewReader(data)
	for reader.Len() > 0 {
		msg, err := Decode(reader)
		if err != nil {
			t.Fatalf("Failed to decode ASTERIX message: %v", err)
		}
		if _, ok := msg.Items["I048/250"]; ok {
			return msg
		}
	}
	t.Fatal("No message with I048/250 found")
	return nil
}

func TestDecodeBDSRegisterData_NoBDSDecoder(t *testing.T) {
	msg := decodeCAT048WithDecoder(t, realWorldCAT048Messages[4], nil)

	raw, ok := msg.Items["I048/250"].(BDSRegisterData)
	if !ok {
		t.Fatalf("I048/250 is %T, want BDSRegisterData", msg.Items["I048/250"])
	}

	if raw.Repetition == 0 {
		t.Fatal("Repetition should be > 0")
	}

	if len(raw.Registers) == 0 {
		t.Fatal("Registers map should not be empty")
	}

	for key, reg := range raw.Registers {
		// Key should be hex format
		if len(key) < 3 || key[:2] != "0x" {
			t.Errorf("Key %q should be hex format like 0x50", key)
		}
		// BDSCode should be set
		if reg.BDSCode == 0 && key != "0x00" {
			t.Errorf("BDSCode should be set for key %s", key)
		}
		// Raw data should be hex-encoded 7 bytes = 14 hex chars
		if len(reg.BDSDataRaw) != 14 {
			t.Errorf("BDSDataRaw for %s has length %d, want 14", key, len(reg.BDSDataRaw))
		}
		// Without BDSDecoder, Decoded and Error should be empty
		if reg.Decoded != nil {
			t.Errorf("Decoded should be nil without BDSDecoder, got %v for %s", reg.Decoded, key)
		}
		if reg.Error != "" {
			t.Errorf("Error should be empty without BDSDecoder, got %q for %s", reg.Error, key)
		}
	}

	t.Logf("Raw-only: %d registers, keys: %v", len(raw.Registers), mapKeys(raw.Registers))
}

func TestDecodeBDSRegisterData_WithBDSDecoder(t *testing.T) {
	mock := &mockBDSDecoder{}
	msg := decodeCAT048WithDecoder(t, realWorldCAT048Messages[4], mock)

	data, ok := msg.Items["I048/250"].(BDSRegisterData)
	if !ok {
		t.Fatalf("I048/250 is %T, want BDSRegisterData", msg.Items["I048/250"])
	}

	if len(mock.calls) == 0 {
		t.Fatal("BDSDecoder was not called")
	}

	if int(data.Repetition) != len(mock.calls) {
		t.Errorf("BDSDecoder called %d times, want %d (repetition)", len(mock.calls), data.Repetition)
	}

	for key, reg := range data.Registers {
		// Raw fields should still be populated
		if len(reg.BDSDataRaw) != 14 {
			t.Errorf("BDSDataRaw for %s has length %d, want 14", key, len(reg.BDSDataRaw))
		}
		// Decoded should be populated by mock
		decoded, ok := reg.Decoded.(map[string]any)
		if !ok {
			t.Errorf("Decoded for %s should be map[string]any, got %T", key, reg.Decoded)
			continue
		}
		if decoded["mock"] != true {
			t.Errorf("Decoded for %s missing mock=true", key)
		}
		// Error should be empty (mock doesn't return errors)
		if reg.Error != "" {
			t.Errorf("Error should be empty, got %q for %s", reg.Error, key)
		}
	}

	t.Logf("With decoder: %d registers decoded, decoder called %d times", len(data.Registers), len(mock.calls))
}

// TestDecodeBDSRegisterData_BDSKeyIsHexAddress verifies the map key matches
// the hex-encoded BDSCode.
func TestDecodeBDSRegisterData_BDSKeyIsHexAddress(t *testing.T) {
	msg := decodeCAT048WithDecoder(t, realWorldCAT048Messages[4], nil)

	data := msg.Items["I048/250"].(BDSRegisterData)
	for key, reg := range data.Registers {
		expectedKey := fmt.Sprintf("0x%02x", reg.BDSCode)
		if key != expectedKey {
			t.Errorf("Map key %q does not match BDSCode 0x%02x (expected key %q)", key, reg.BDSCode, expectedKey)
		}
	}
}

func mapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
