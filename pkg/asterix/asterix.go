package asterix

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Message represents a generic ASTERIX message
type RawAsterixMessage struct {
	Category byte
	Length   uint16
	Payload  []byte
}

type AsterixMessage struct {
	Category byte
	Sic      uint8
	Sac      uint8
	Items    map[string]interface{}
	FSPEC    []byte
}

var decoders = map[byte]Decoder{
	34: &CAT034Decoder{},
	48: &CAT048Decoder{},
}

func Decode(r io.Reader) (*AsterixMessage, error) {
	msg, err := ParseMessage(r)
	if err != nil {
		return nil, err
	}
	return Dispatch(msg)
}

// ParseMessage reads a message from the reader
func ParseMessage(r io.Reader) (*RawAsterixMessage, error) {
	header := make([]byte, 3)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(header[1:3])
	if length < 3 {
		return nil, fmt.Errorf("invalid message length: %d", length)
	}
	payload := make([]byte, length-3)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return &RawAsterixMessage{
		Category: header[0],
		Length:   length,
		Payload:  payload,
	}, nil
}

// RegisterDecoder registers or replaces a category decoder.
func RegisterDecoder(category byte, decoder Decoder) {
	decoders[category] = decoder
}

// Dispatch dispatches the message to the appropriate decoder
func Dispatch(msg *RawAsterixMessage) (*AsterixMessage, error) {
	decoder, ok := decoders[msg.Category]
	if !ok {
		return nil, fmt.Errorf("no decoder for category %d", msg.Category)
	}
	result, err := decoder.Decode(msg)
	if err != nil {
		return nil, err
	}
	return result, nil
}
