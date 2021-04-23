package toml

import (
	"io"

	"github.com/BurntSushi/toml"
	"github.com/johnstarich/env2config"
)

func init() {
	env2config.RegisterFormat("toml", &tomlMarshaler{})
}

type tomlMarshaler struct{}

func (t *tomlMarshaler) Marshal(w io.Writer, value interface{}) error {
	return toml.NewEncoder(w).Encode(value)
}

func (*tomlMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	_, err := toml.DecodeReader(r, dest)
	return err
}
