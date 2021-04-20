package env2config

import (
	"io"

	"github.com/pkg/errors"
)

// Marshaler serializes the given value and writes it to the writer
type Marshaler interface {
	Marshal(io.Writer, interface{}) error
}

// Unmarshaler deserializes contents of the reader into the provided value pointer.
// For a value to be mergeable, it must be of type map[string]interface{} or []interface{}.
type Unmarshaler interface {
	Unmarshal(io.Reader, interface{}) error
}

var defaultRegistry = newRegistry()

type registry struct {
	marshalers   map[string]Marshaler
	unmarshalers map[string]Unmarshaler
}

func newRegistry() *registry {
	return &registry{
		marshalers:   make(map[string]Marshaler),
		unmarshalers: make(map[string]Unmarshaler),
	}
}

func (r *registry) RegisterFormat(format string, marshaler Marshaler) {
	if marshaler == nil {
		panic("Marshaler must not be nil")
	}
	r.marshalers[format] = marshaler
	if unmarshaler, ok := marshaler.(Unmarshaler); ok {
		r.unmarshalers[format] = unmarshaler
	}
}

func (r *registry) MarshalFormat(format string, w io.Writer, value interface{}) error {
	marshaler := r.marshalers[format]
	if marshaler == nil {
		return errors.Errorf("Unsupported file format: %q", format)
	}
	return marshaler.Marshal(w, value)
}

func (r *registry) UnmarshalFormat(format string, reader io.Reader, dest interface{}) error {
	unmarshaler := r.unmarshalers[format]
	if unmarshaler == nil {
		return errors.Errorf("Unsupported template file format: %q", format)
	}
	return unmarshaler.Unmarshal(reader, dest)
}

func RegisterFormat(format string, marshaler Marshaler) {
	defaultRegistry.RegisterFormat(format, marshaler)
}
