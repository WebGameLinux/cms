package reflects

import (
		"encoding/json"
		"errors"
)

func CopyMap2Struct(data map[string]interface{}, v interface{}) error {
		if len(data) == 0 {
				return errors.New("empty map")
		}
		if v == nil {
				return errors.New("nil ptr")
		}
		var (
				d   []byte
				err error
		)
		if d, err = json.Marshal(data); err == nil {
				return json.Unmarshal(d, v)
		}
		return err
}
