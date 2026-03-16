package asterix

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"testing"
)

func TestDispatch(t *testing.T) {
	// Create a test message

	msg := &RawAsterixMessage{
		Category: 48,

		Payload: []byte{
			0b10100000, // FSPEC: bit 1 + bit 3 set
			0x12, 0x34, // I048/010 - Data Source Identifier
			0x01, 0x00, 0x00, // I048/1
			0x00, 0x00, 0x00, // I048/2
		},
		Length: 10,
	}

	// Dispatch the message
	res, err := Dispatch(msg)
	if err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}
	if res == nil {
		t.Fatal("Dispatch returned nil result")
	}
	if res.Category != 48 {
		t.Errorf("expected category 48, got %d", res.Category)
	}
	cat048, ok := res.Record.(*Cat048Message)
	if !ok {
		t.Fatal("expected Record to be *Cat048Message")
	}
	// Check that DataSource was decoded (FSPEC bit 1)
	if cat048.DataSource == nil {
		t.Error("expected Cat048.DataSource to be non-nil")
	}

}

func TestDispatch_InvalidCategory(t *testing.T) {
	msg := &RawAsterixMessage{
		Category: 47,
		Payload: []byte{
			0b10100000, // FSPEC: bit 1 + bit 3 set
			0x12, 0x34, // I048/010 - Data Source Identifier
			0x01, 0x00, 0x00, // I048/1
			0x00, 0x00, 0x00, // I048/2
		},
		Length: 10,
	}

	decoded, err := Dispatch(msg)
	if err == nil {
		t.Fatal("expected error for invalid category, got nil")
	}
	if decoded != nil {
		t.Fatal("expected nil result for invalid category")
	}
}

func TestParseMessage_Valid(t *testing.T) {
	var buf bytes.Buffer
	// Header: Category = 48, Length = 10 (3-byte header + 7-byte payload)
	buf.WriteByte(48)
	binary.Write(&buf, binary.BigEndian, uint16(10))
	buf.Write([]byte{0xA0, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06})

	msg, err := ParseMessage(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Category != 48 {
		t.Errorf("expected category 48, got %d", msg.Category)
	}
	if msg.Length != 10 {
		t.Errorf("expected length 10, got %d", msg.Length)
	}
	if len(msg.Payload) != 7 {
		t.Errorf("expected payload length 7, got %d", len(msg.Payload))
	}
}

func TestParseMessage_InvalidLength(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteByte(48)
	binary.Write(&buf, binary.BigEndian, uint16(2)) // Length < 3

	_, err := ParseMessage(&buf)
	if err == nil {
		t.Fatal("expected error for invalid length, got nil")
	}
}

func TestParseMessage_EOF(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteByte(48)
	binary.Write(&buf, binary.BigEndian, uint16(10)) // Length = 10 but missing payload

	_, err := ParseMessage(&buf)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if err != io.ErrUnexpectedEOF && err != io.EOF {
		t.Fatalf("expected EOF-related error, got: %v", err)
	}
}

func TestInvalidMessages(t *testing.T) {
	tests := []struct {
		name    string
		hexData string
		wantErr bool
	}{
		{
			name:    "Empty message",
			hexData: "",
			wantErr: true,
		},
		{
			name:    "Too short message",
			hexData: "22",
			wantErr: true,
		},
		{
			name:    "Invalid category",
			hexData: "ff0012f62821025460022084404600840000",
			wantErr: true,
		},
		{
			name:    "Invalid length",
			hexData: "220000f62821025460022084404600840000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			if tt.hexData != "" {
				data, err = hex.DecodeString(tt.hexData)
				if err != nil {
					t.Fatalf("Failed to decode hex data: %v", err)
				}
			}

			_, err = Decode(bytes.NewReader(data))
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}