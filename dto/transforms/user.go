package transforms

import "github.com/WebGameLinux/cms/models/types"

func TransformGender(v interface{}) types.Gender {
		if n, ok := v.(types.Gender); ok {
				return n
		}
		if n, ok := v.(int); ok {
				return types.ParseInt(n)
		}
		if str, ok := v.(string); ok {
				return types.Parse(str)
		}
		return types.Unknown
}

func TransformMapGender(key string, v interface{}) (string, interface{}) {
		return key, TransformGender(v)
}
