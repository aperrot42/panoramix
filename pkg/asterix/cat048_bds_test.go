package asterix

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func decodeCAT048WithBDS(t *testing.T, hexData string) *AsterixMessage {
	t.Helper()

	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Fatalf("Failed to decode hex data: %v", err)
	}

	reader := bytes.NewReader(data)
	for reader.Len() > 0 {
		msg, err := Decode(reader)
		if err != nil {
			t.Fatalf("Failed to decode ASTERIX message: %v", err)
		}
		if cat048, ok := msg.Record.(*Cat048Message); ok && cat048.BDSRegister != nil {
			return msg
		}
	}
	t.Fatal("No message with BDSRegister found")
	return nil
}

func TestDecodeBDSRegisterData_RawOnly(t *testing.T) {
	msg := decodeCAT048WithBDS(t, realWorldCAT048Messages[4])

	cat048, ok := msg.Record.(*Cat048Message)
	if !ok {
		t.Fatal("expected Record to be *Cat048Message")
	}
	raw := cat048.BDSRegister
	if raw == nil {
		t.Fatal("BDSRegister is nil")
	}

	if raw.Repetition == 0 {
		t.Fatal("Repetition should be > 0")
	}

	if len(raw.Registers) == 0 {
		t.Fatal("Registers map should not be empty")
	}

	for key, reg := range raw.Registers {
		if len(key) < 3 || key[:2] != "0x" {
			t.Errorf("Key %q should be hex format like 0x50", key)
		}
		if reg.BDSCode == 0 && key != "0x00" {
			t.Errorf("BDSCode should be set for key %s", key)
		}
		if len(reg.RawData) != 7 {
			t.Errorf("RawData for %s has length %d, want 7", key, len(reg.RawData))
		}
	}

	t.Logf("Raw-only: %d registers, keys: %v", len(raw.Registers), mapKeys(raw.Registers))
}

func TestDecodeBDSRegisterData_BDSKeyIsHexAddress(t *testing.T) {
	msg := decodeCAT048WithBDS(t, realWorldCAT048Messages[4])

	cat048, ok := msg.Record.(*Cat048Message)
	if !ok || cat048.BDSRegister == nil {
		t.Fatal("Cat048 or BDSRegister is nil")
	}

	for key, reg := range cat048.BDSRegister.Registers {
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
