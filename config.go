package env2config

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Name   string
	Opts   Opts // NAME_OPTS_*
	Values Values

	registry *registry
}

type Opts struct {
	File         string `required:"true"`
	Format       string `required:"true"`
	TemplateFile string `split_words:"true"`

	Inputs Values // NAME_OPTS_IN_*
}

type Values map[string]string

var _ envconfig.Setter = Values{}

func (v Values) Set(value string) error { return nil }

func New(name string) (Config, error) {
	c, err := newConfig(name, parseEnv(os.Environ()), defaultRegistry)
	if name != "" {
		err = errors.Wrap(err, strings.ToLower(name))
	}
	return c, err
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
	if err != nil {
		return Config{}, err
	}
	config.Name = name
	config.Opts.Inputs = filterEnvPrefix(name+"_opts_in", env)
	config.Values = configEnvValues(name, env)

	var missingInputs []string
	for dest, src := range config.Opts.Inputs {
		value, isSet := env[src]
		if !isSet {
			missingInputs = append(missingInputs, src)
		}
		config.Values[dest] = value
	}
	if len(missingInputs) > 0 {
		sort.Strings(missingInputs)
		return Config{}, errors.Errorf("Missing required environment variables: %s", strings.Join(missingInputs, ", "))
	}
	return config, nil
}

func (c Config) Write() error {
	var template map[string]interface{}
	if c.Opts.TemplateFile != "" {
		f, err := os.Open(c.Opts.TemplateFile)
		if err != nil {
			return err
		}
		defer f.Close()
		err = c.registry.UnmarshalFormat(c.Opts.Format, f, &template)
		if err != nil {
			return err
		}
	}
	values := c.writableValues(template)
	f, err := os.OpenFile(c.Opts.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.registry.MarshalFormat(c.Opts.Format, f, values)
}

func (c Config) writableValues(template map[string]interface{}) interface{} {
	result := template
	if result == nil {
		result = make(map[string]interface{})
	}
	for key, value := range c.Values {
		current := result
		keys := parseKeyPath(key)
		for i := 0; i < len(keys)-1; i++ {
			key := keys[i]
			_, exists := current[key]
			if !exists {
				current[key] = make(map[string]interface{})
			}
			switch next := current[key].(type) {
			case map[string]interface{}:
				current = next
			case []interface{}:
				nextMap := arrayToMap(next)
				current[key] = nextMap
				current = nextMap
			default:
				// unrecognized type, just do simple override
				nextMap := make(map[string]interface{})
				current[key] = nextMap
				current = nextMap
			}
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
	return values
}

func arrayToMap(a []interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(a))
	for index, value := range a {
		m[strconv.FormatInt(int64(index), 10)] = value
	}
	return m
}
