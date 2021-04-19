package env2config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Name   string
	Opts   Opts
	Values Values

	registry *registry
}

type Opts struct {
	File   string `required:"true"`
	Format string `required:"true"`
}

type Values map[string]string

var _ envconfig.Setter = Values{}

func (v Values) Set(value string) error { return nil }

func New(name string) (Config, error) {
	return newConfig(name, parseEnv(os.Environ()), defaultRegistry)
}

func newConfig(name string, env map[string]string, registry *registry) (Config, error) {
	name = strings.ToLower(name)
	if name == "" {
		return Config{}, errors.New("Config name is required")
	}
	trimmed := name
	trimmed = strings.TrimFunc(trimmed, unicode.IsLetter)
	trimmed = strings.TrimFunc(trimmed, unicode.IsNumber)
	if trimmed != "" {
		return Config{}, errors.Errorf("Config names must only use letters or numbers: %q", name)
	}
	config := Config{
		registry: registry,
	}
	err := envconfig.Process(name, &config)
	config.Name = name
	config.Values = configValues(name, env)
	return config, err
}

func (c Config) Write() error {
	values := c.writableValues()
	f, err := os.OpenFile(c.Opts.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.registry.MarshalFormat(c.Opts.Format, f, values)
}

func (c Config) writableValues() interface{} {
	result := make(map[string]interface{})
	for key, value := range c.Values {
		current := result
		keys := parseKeyPath(key)
		for i := 0; i < len(keys)-1; i++ {
			key := keys[i]
			_, exists := current[key]
			if !exists {
				current[key] = make(map[string]interface{})
			}
			current = current[key].(map[string]interface{})
		}
		lastKey := keys[len(keys)-1]
		current[lastKey] = value
	}

	return mapsToArrays(result)
}

// parseKeyPath splits 'key' by '.' and returns the paths. Skips over escapes like '\.'.
func parseKeyPath(key string) []string {
	cursor := 0
	var paths []string
	for ix, r := range key {
		if r == '.' && (ix == 0 || key[ix-1] != '\\') {
			paths = append(paths, key[cursor:ix])
			cursor = ix + 1
		}
	}
	if cursor < len(key) {
		paths = append(paths, key[cursor:])
	}
	return paths
}

func mapsToArrays(m map[string]interface{}) interface{} {
	isArray := true
	for key, value := range m {
		if mapValue, isMap := value.(map[string]interface{}); isMap {
			m[key] = mapsToArrays(mapValue)
		}
		if strings.TrimFunc(key, unicode.IsNumber) != "" {
			isArray = false
		}
	}
	if !isArray {
		return m
	}
	// Must be an array, or empty set.
	// Empty set shouldn't be possible based on logic of writableValues().
	values := make([]interface{}, len(m))
	for key, value := range m {
		index, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			panic(err) // string is all digits so must parse as int
		}
		values[index] = value
	}
	fmt.Println("flattening", values)
	return values
}
