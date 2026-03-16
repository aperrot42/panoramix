package pipeline

import "github.com/aperrot42/panoramix/pkg/internal_format"

// IFSource adapts an internal_format.Reader into a Source[internal_format.Record].
// Records reference internal buffers and are valid until the next Next() call.
type IFSource struct {
	reader *internal_format.Reader
}

func NewIFSource(reader *internal_format.Reader) *IFSource {
	return &IFSource{reader: reader}
}

func (s *IFSource) Next() (internal_format.Record, bool, error) {
	rec, ok, err := s.reader.NextRecord()
	if err != nil {
		return internal_format.Record{}, false, err
	}
	if !ok {
		return internal_format.Record{}, false, nil
	}
	return rec, true, nil
}
