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
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	return enc.Encode(value)
}

func (*jsonMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return json.NewDecoder(r).Decode(dest)
}
