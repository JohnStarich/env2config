package internal

func Walk(v interface{}, fn func(v interface{}) interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{}, len(v))
		for key, value := range v {
			newMap[key] = Walk(value, fn)
		}
		return fn(newMap)
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for index, value := range v {
			newSlice[index] = Walk(value, fn)
		}
		return fn(newSlice)
	default:
		return fn(v)
	}
}
