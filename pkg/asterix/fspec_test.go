package asterix

import (
	"fmt"
	"testing"
)

func TestParseFSPEC(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		rest     []byte
		wantErr  bool
	}{
		{
			name:     "Single FSPEC byte without extension",
			input:    []byte{0b10100000, 0x12, 0x34},
			expected: []byte{0b10100000},
			rest:     []byte{0x12, 0x34},
			wantErr:  false,
		},
		{
			name:     "Two FSPEC bytes with extension",
			input:    []byte{0b10100001, 0b10000000, 0x56},
			expected: []byte{0b10100001, 0b10000000},
			rest:     []byte{0x56},
			wantErr:  false,
		},
		{
			name:     "FSPEC terminated unexpectedly",
			input:    []byte{0b10000001}, // FX == 1 but no next byte
			expected: nil,
			rest:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fspec, rest, err := extractFSPEC(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFSPEC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equalBytes(fspec, tt.expected) {
				t.Errorf("parseFSPEC() fspec = %v, expected %v", fspec, tt.expected)
			}
			if !equalBytes(rest, tt.rest) {
				t.Errorf("parseFSPEC() rest = %v, expected %v", rest, tt.rest)
			}
		})
	}
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestWalkFSPEC(t *testing.T) {
	fieldTable := map[int]DataItem{
		1: NewDataItem("I048/010", func(data []byte) (interface{}, int, error) {
			if len(data) < 2 {
				return nil, 0, fmt.Errorf("too short")
			}
			return []byte{data[0], data[1]}, 2, nil
		}),
		3: NewDataItem("I048/140", func(data []byte) (interface{}, int, error) {
			if len(data) < 3 {
				return nil, 0, fmt.Errorf("too short")
			}
			return int(data[0])<<16 | int(data[1])<<8 | int(data[2]), 3, nil
		}),
	}

	// FSPEC: bit 1 and bit 3 set (bit positions 7 and 5), FX=0
	// Field 1: 2 bytes → 0x12 0x34
	// Field 3: 3 bytes → 0x00 0x01 0x02
	data := []byte{0b10100000, 0x12, 0x34, 0x00, 0x01, 0x02}

	fspec, rest, err := extractFSPEC(data)
	if err != nil {
		t.Fatalf("unexpected fspec error: %v", err)
	}

	fields, consumed, err := WalkFSPEC(fspec, rest, fieldTable)
	if err != nil {
		t.Fatalf("unexpected walkerror: %v", err)
	}
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
	if consumed != 5 {
		t.Errorf("expected 5 bytes consumed, got %d", consumed)
	}
}
