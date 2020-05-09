package mapper

func WrapperNewMap(root string, data map[string]interface{}) map[string]map[string]interface{} {
		var wrapper = make(map[string]map[string]interface{})
		wrapper[root] = data
		return wrapper
}

func WrapperNewMapper(root string, data map[string]interface{}) map[string]interface{} {
		var wrapper = make(map[string]interface{})
		wrapper[root] = data
		return wrapper
}

func NewMap() map[string]interface{} {
		return make(map[string]interface{})
}

func NewKvMap(key string, data interface{}) map[string]interface{} {
		var m = NewMap()
		if key != "" {
				m[key] = data
		}
		return m
}
