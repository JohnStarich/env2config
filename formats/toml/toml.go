package toml

import (
	"io"

	"github.com/BurntSushi/toml"
	"github.com/johnstarich/env2config"
)

func init() {
	toml := &tomlMarshaler{}
	env2config.RegisterFormat("ini", toml)
	env2config.RegisterFormat("toml", toml)
}

type tomlMarshaler struct{}

func (*tomlMarshaler) Marshal(w io.Writer, value interface{}) error {
	return toml.NewEncoder(w).Encode(value)
}

func (*tomlMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	_, err := toml.DecodeReader(r, dest)
	return err
}
