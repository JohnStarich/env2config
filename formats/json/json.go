package json

import (
	"encoding/json"
	"io"

	"github.com/johnstarich/env2config"
)

func init() {
	env2config.RegisterFormat("json", &jsonMarshaler{})
}

type jsonMarshaler struct{}

func (*jsonMarshaler) Marshal(w io.Writer, value interface{}) error {
	return json.NewEncoder(w).Encode(value)
}
