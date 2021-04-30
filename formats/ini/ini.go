package ini

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/johnstarich/env2config"
	"github.com/johnstarich/env2config/formats/internal"
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
	f = internal.Walk(f, omitEmptyStrings).(*ini.File)
	_, err = f.WriteToIndent(w, "  ")
	return err
}

func (*iniMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return ini.MapTo(dest, r)
}

func omitEmptyStrings(v interface{}) interface{} {
	switch v := v.(type) {
	case *ini.File:
		for _, section := range v.Sections() {
			_ = internal.Walk(section, omitEmptyStrings)
		}
		return v
	case *ini.Section:
		for _, key := range v.Keys() {
			if key.Value() == "" {
				// convert empty strings to boolean empty keys
				_, err := v.NewBooleanKey(key.Name())
				if err != nil {
					panic(err)
				}
			}
			_ = internal.Walk(key, omitEmptyStrings)
		}
		return v
	default:
		return v
	}
}
