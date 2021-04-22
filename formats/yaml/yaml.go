package yaml

import (
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/johnstarich/env2config"
	"gopkg.in/yaml.v3"
)

func init() {
	env2config.RegisterFormat("yaml", &yamlMarshaler{})
}

type yamlMarshaler struct{}

func (*yamlMarshaler) Marshal(w io.Writer, value interface{}) error {
	value = walk(value, parseValues)
	value = walk(value, omitNils)
	return yaml.NewEncoder(w).Encode(value)
}

func (*yamlMarshaler) Unmarshal(r io.Reader, dest interface{}) error {
	return yaml.NewDecoder(r).Decode(dest)
}

func walk(v interface{}, fn func(v interface{}) interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{}, len(v))
		for key, value := range v {
			newMap[key] = walk(value, fn)
		}
		return fn(newMap)
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for index, value := range v {
			newSlice[index] = walk(value, fn)
		}
		return fn(newSlice)
	default:
		return fn(v)
	}
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
