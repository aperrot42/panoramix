package asterix

import "time"

// ItemDecoder decodes a data field and returns its strongly-typed value, how many bytes it consumed, and error if any.
type ItemDecoder[T any] func(data []byte) (T, int, error)

// Record is the interface implemented by all category-specific message types.
type Record interface {
	GetSAC() uint8
	GetSIC() uint8
	GetTimeOfDay() time.Duration
}

// Decoder is the interface for category decoders.
type Decoder interface {
	Decode(msg *RawAsterixMessage) (*AsterixMessage, error)
}

// DataItem represents a single ASTERIX data item with its name and decoder
type DataItem struct {
	Name    string
	Decoder ItemDecoder[any]
}

// NewDataItem creates a DataItem with a decoder function that matches the common pattern
func NewDataItem(name string, decoderFunc func([]byte) (any, int, error)) DataItem {
	return DataItem{
		Name:    name,
		Decoder: func(data []byte) (any, int, error) { return decoderFunc(data) },
	}
}

