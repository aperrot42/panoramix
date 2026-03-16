package fspec

import (
	"testing"

	"github.com/aperrot42/panoramix/pkg/asterix"
)

func TestComputeAvailableFields(t *testing.T) {
	ds := asterix.DataSourceIdentifier{SAC: 2, SIC: 1}
	fl := asterix.FlightLevel{FL: 350, Validated: true}
	rp := asterix.RadarPlotCharacteristics{}
	msg := &asterix.AsterixMessage{
		Category: 48,
		Record: &asterix.Cat048Message{
			DataSource:  &ds,
			FlightLevel: &fl,
			RadarPlot:   &rp,
		},
	}

	fields := ComputeAvailableFields(msg)

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d: %v", len(fields), fields)
	}

	fieldMap := make(map[string]bool)
	for _, f := range fields {
		fieldMap[f] = true
	}

	for _, expected := range []string{"I048/010", "I048/090", "I048/130"} {
		if !fieldMap[expected] {
			t.Errorf("expected field %s not found in %v", expected, fields)
		}
	}
}

func TestComputeAvailableFields_NilRecord(t *testing.T) {
	msg := &asterix.AsterixMessage{
		Category: 48,
	}

	fields := ComputeAvailableFields(msg)
	if fields != nil {
		t.Errorf("expected nil, got %v", fields)
	}
}

func TestComputeAvailableFields_EmptyRecord(t *testing.T) {
	msg := &asterix.AsterixMessage{
		Category: 48,
		Record:   &asterix.Cat048Message{},
	}

	fields := ComputeAvailableFields(msg)
	if len(fields) != 0 {
		t.Errorf("expected empty fields, got %v", fields)
	}
}