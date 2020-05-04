package getter


type Getter interface {
		GetValue(string) (interface{}, error)
}

func GetAnyGetter(v interface{}) Getter {
		switch v.(type) {
		case map[string]interface{}:
		case map[int]interface{}:
		case []interface{}:

		}
		panic("implGetAnyGetter")
		return nil
}
