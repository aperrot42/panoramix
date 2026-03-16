package pipeline

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aperrot42/panoramix/pkg/asterix"
	"github.com/aperrot42/panoramix/pkg/internal_format"
	"github.com/aperrot42/panoramix/pkg/modes/bds"
)

const testFile = "../../../../../data/plots_20250701_line_A_recdatach20.if"

func loadTestData(tb testing.TB) []byte {
	tb.Helper()
	data, err := os.ReadFile(testFile)
	if err != nil {
		tb.Skipf("test data not available: %v", err)
	}
	return data
}

func BenchmarkOldPath_IFRead(b *testing.B) {
	data := loadTestData(b)
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
		for {
			rec, err := reader.ReadRecord()
			if err != nil || rec == nil {
				break
			}
		}
	}
}

func BenchmarkNewPath_IFRead(b *testing.B) {
	data := loadTestData(b)
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
		src := NewIFSource(reader)
		Drain(src, func(_ internal_format.Record) error {
			return nil
		})
	}
}

const benchMaxRecords = 5000

func BenchmarkOldPath_FullDecode(b *testing.B) {
	data := loadTestData(b)
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
		n := 0
		for n < benchMaxRecords {
			rec, err := reader.ReadRecord()
			if err != nil || rec == nil {
				break
			}
			asterix.Decode(bytes.NewReader(rec.Payload))
			n++
		}
	}
}

func BenchmarkNewPath_FullPipeline(b *testing.B) {
	data := loadTestData(b)
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
		src := NewIFSource(reader)
		decoded := NewMapFilter(src, AsterixDecodeFilter)
		enriched := NewMapFilter(decoded, NewBDSEnrichFilter(&bds.Adapter{}))
		n := 0
		Drain(enriched, func(_ AsterixResult) error {
			n++
			if n >= benchMaxRecords {
				return io.EOF
			}
			return nil
		})
	}
}

func BenchmarkNewPath_FullPipelineJSON(b *testing.B) {
	data := loadTestData(b)
	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		reader := internal_format.NewReaderWithBaseDate(bytes.NewReader(data), time.Now())
		src := NewIFSource(reader)
		decoded := NewMapFilter(src, AsterixDecodeFilter)
		enriched := NewMapFilter(decoded, NewBDSEnrichFilter(&bds.Adapter{}))
		enc := json.NewEncoder(io.Discard)
		n := 0
		Drain(enriched, func(ar AsterixResult) error {
			n++
			if n >= benchMaxRecords {
				return io.EOF
			}
			return enc.Encode(ar)
		})
	}
}
