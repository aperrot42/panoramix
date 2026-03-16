package pipeline

import (
	"github.com/aperrot42/panoramix/pkg/asterix"
	"github.com/aperrot42/panoramix/pkg/internal_format"
)

// AsterixDecodeFilter decodes an IF record's payload into an AsterixMessage.
// Records that fail to decode are silently skipped.
func AsterixDecodeFilter(rec internal_format.Record) (AsterixResult, bool, error) {
	msg, err := asterix.DecodeFromBytes(rec.Payload)
	if err != nil {
		return AsterixResult{}, false, nil // skip decode errors
	}
	return AsterixResult{
		Record:  rec,
		Message: *msg,
	}, true, nil
}
