package yaml

import (
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/johnstarich/env2config"
	"github.com/johnstarich/env2config/formats/internal"
	"gopkg.in/yaml.v3"
)

func init() {
	env2config.RegisterFormat("yaml", &yamlMarshaler{})
}

type yamlMarshaler struct{}

func (*yamlMarshaler) Marshal(w io.Writer, value interface{}) error {
	value = internal.Walk(value, parseValues)
	value = internal.Walk(value, omitNils)
	return yaml.NewEncoder(w).Encode(value)
}

func (*yamlMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return yaml.NewDecoder(r).Decode(dest)
}

func omitNils(v interface{}) interface{} {
	if v == nil {
		return &yaml.Node{Kind: yaml.ScalarNode}
	}
	return v
}

func parseValues(v interface{}) interface{} {
	str, isString := v.(string)
	if !isString {
		return v
	}

	switch {
	case str == "true":
		return true
	case str == "false":
		return false
	case strings.TrimFunc(str, unicode.IsNumber) == "":
		integer, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic(err)
		}
		return integer
	default:
		return v
	}
}
