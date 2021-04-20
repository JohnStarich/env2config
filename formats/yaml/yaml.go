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
	value = omitNils(value)
	return yaml.NewEncoder(w).Encode(value)
}

func (*yamlMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return yaml.NewDecoder(r).Decode(dest)
}

func omitNils(v interface{}) interface{} {
	switch v := v.(type) {
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode}
	case map[string]interface{}:
		newMap := make(map[string]interface{}, len(v))
		for key, value := range v {
			newMap[key] = omitNils(value)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for index, value := range v {
			newSlice[index] = omitNils(value)
		}
		return newSlice
	default:
		return v
	}
}
