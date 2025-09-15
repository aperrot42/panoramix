package asterix

// ItemDecoder decodes a data field and returns its strongly-typed value, how many bytes it consumed, and error if any.
type ItemDecoder[T any] func(data []byte) (T, int, error)

// Decoder is the interface for category decoders.
type Decoder interface {
	Decode(msg *RawAsterixMessage) (*AsterixMessage, error)
}
