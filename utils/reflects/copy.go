package reflects

import (
		"bytes"
		"encoding/gob"
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

func Copy(src, dst interface{}) error {
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(src); err != nil {
				return err
		}
		return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
