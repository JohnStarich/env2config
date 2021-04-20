package yaml

import (
	"io"

	"github.com/johnstarich/env2config"
	"gopkg.in/yaml.v3"
)

func init() {
	env2config.RegisterFormat("yaml", &yamlMarshaler{})
}

type yamlMarshaler struct{}

func (*yamlMarshaler) Marshal(w io.Writer, value interface{}) error {
	return yaml.NewEncoder(w).Encode(value)
}

func (*yamlMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return yaml.NewDecoder(r).Decode(dest)
}
