package env2config

import (
	"io"

	"github.com/pkg/errors"
)

type Marshaler interface {
	Marshal(io.Writer, interface{}) error
}

var defaultRegistry = &registry{marshalers: make(map[string]Marshaler)}

type registry struct {
	marshalers map[string]Marshaler
}

func (r *registry) RegisterFormat(format string, marshaler Marshaler) {
	if marshaler == nil {
		panic("Marshaler must not be nil")
	}
	r.marshalers[format] = marshaler
}

func (r *registry) MarshalFormat(format string, w io.Writer, value interface{}) error {
	marshaler := r.marshalers[format]
	if marshaler == nil {
		return errors.Errorf("Unsupported file format: %q", format)
	}
	return marshaler.Marshal(w, value)
}

func RegisterFormat(format string, marshaler Marshaler) {
	defaultRegistry.RegisterFormat(format, marshaler)
}
