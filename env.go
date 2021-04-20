package env2config

import "strings"

func parseEnv(envPairs []string) map[string]string {
	env := make(map[string]string)
	for _, kv := range envPairs {
		kvItems := strings.SplitN(kv, "=", 2)
		if len(kvItems) < 2 {
			panic("Invalid key value pair from os.Environ()")
		}
		key, value := kvItems[0], kvItems[1]
		if _, isSet := env[key]; !isSet {
			env[key] = value
		}
	}
	return env
}

func configEnvValues(name string, env map[string]string) map[string]string {
	env = filterEnvPrefix(name, env)
	removeEnvOpts(env)
	return env
}

func filterEnvPrefix(prefix string, env map[string]string) map[string]string {
	prefix += "_"
	prefixLen := len(prefix)

	prefixedEnv := make(map[string]string)
	for key, value := range env {
		if len(key) > prefixLen && strings.EqualFold(prefix, key[:prefixLen]) {
			key = key[prefixLen:]
			prefixedEnv[key] = value
		}
	}
	return prefixedEnv
}

func removeEnvOpts(m map[string]string) {
	const optsPrefix = "OPTS_"
	for key := range m {
		if len(key) > len(optsPrefix) && strings.EqualFold(key[:len(optsPrefix)], optsPrefix) {
			delete(m, key)
		}
	}
}
