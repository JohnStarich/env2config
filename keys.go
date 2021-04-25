package env2config

import (
	"sort"
	"strconv"
	"strings"
)

const (
	keySeparator    = '.'
	keySeparatorStr = string(keySeparator)
)

// parseKeyPath splits 'key' by '.' and returns the paths. Skips over escapes like '\.'.
func parseKeyPath(key string) []string {
	cursor := 0
	var paths []string
	for ix, r := range key {
		if r == keySeparator && (ix == 0 || key[ix-1] != '\\') {
			paths = append(paths, unescapeKey(key[cursor:ix]))
			cursor = ix + 1
		}
	}
	if cursor < len(key) {
		paths = append(paths, unescapeKey(key[cursor:]))
	}
	return paths
}

func unescapeKey(key string) string {
	return strings.Replace(key, `\.`, keySeparatorStr, -1)
}

func escapeKey(key string) string {
	return strings.Replace(key, keySeparatorStr, `\.`, -1)
}

func deleteKeyPath(v interface{}, keyPath []string) (newValue interface{}, deleteMe bool) {
	if len(keyPath) == 0 {
		return nil, true
	}
	key := keyPath[0]
	nextKeyPath := keyPath[1:]
	switch v := v.(type) {
	case map[string]interface{}:
		newVal, shouldDelete := deleteKeyPath(v[key], nextKeyPath)
		v[key] = newVal
		if shouldDelete {
			delete(v, key)
		}
		return v, false
	case []interface{}:
		keyInt, err := strconv.ParseUint(key, 10, 64)
		if err != nil || keyInt >= uint64(len(v)) {
			return v, false
		}
		newVal, shouldDelete := deleteKeyPath(v[keyInt], nextKeyPath)
		v[keyInt] = newVal
		if shouldDelete {
			v = append(v[:keyInt], v[keyInt+1:]...)
		}
		return v, false
	default:
		return v, false
	}
}

// sortTemplateDeleteKeys sorts keys so that they can all be honored correctly.
// Edge cases come into play when deleting array elements, since the indexes change.
func sortTemplateDeleteKeys(deleteKeys []string) {
	keyPaths := make(map[string][]string)
	possibleIndexes := make(map[string]uint64)
	for _, deleteKey := range deleteKeys {
		keys := parseKeyPath(deleteKey)
		keyPaths[deleteKey] = keys
		keyStr := ""
		for _, key := range keys {
			keyStr += escapeKey(key) + keySeparatorStr
			possibleIndex, _ := strconv.ParseUint(key, 10, 64)
			possibleIndexes[keyStr] = possibleIndex
		}
	}

	sort.Slice(deleteKeys, func(a, b int) bool {
		keyA, keyB := deleteKeys[a], deleteKeys[b]
		pathA, pathB := keyPaths[keyA], keyPaths[keyB]
		if len(pathA) != len(pathB) {
			return len(pathA) > len(pathB) // (1) delete the longest paths first
		}

		keyStrA, keyStrB := "", ""
		for ix := range pathA {
			keyStrA += escapeKey(pathA[ix]) + keySeparatorStr
			keyStrB += escapeKey(pathB[ix]) + keySeparatorStr
			possibleIndexA, possibleIndexB := possibleIndexes[keyStrA], possibleIndexes[keyStrB]
			if possibleIndexA != possibleIndexB {
				return possibleIndexA > possibleIndexB // (2) then delete biggest index
			}
		}
		return false // (3) if no key path segments differed, no need to sort further
	})
}
