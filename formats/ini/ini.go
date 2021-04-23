package ini

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/johnstarich/env2config"
	"gopkg.in/ini.v1"
)

func init() {
	env2config.RegisterFormat("ini", &iniMarshaler{})
}

type iniMarshaler struct{}

func (*iniMarshaler) Marshal(w io.Writer, value interface{}) error {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf) // Encode to toml first, since ini library can't handle map types.
	err := enc.Encode(value)
	if err != nil {
		return err
	}
	f, err := ini.Load(&buf)
	if err != nil {
		return err
	}
	_, err = f.WriteToIndent(w, "  ")
	return err
}

func (*iniMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return ini.MapTo(dest, r)
}
