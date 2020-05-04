package string

import (
		"encoding/base64"
)

var errBase64 error

func Base64Encode(str string) string {
		return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Encode2Byte(str string) []byte {
		return []byte(Base64Encode(str))
}

func Base64Decode(str string) string {
		var (
				v   []byte
				err error
		)
		if v, err = base64.StdEncoding.DecodeString(str); err == nil {
				return string(v)
		}
		errBase64 = err
		return ""
}

func Base64Error() error {
		defer func() { errBase64 = nil }()
		return errBase64
}

func Base64Decode2Byte(str string) []byte {
		return []byte(Base64Decode(str))
}

func ByteBase64Encode(b []byte) string {
		return base64.StdEncoding.EncodeToString(b)
}

func ByteBase64Decode(b []byte) string {
		var (
				v   []byte
				err error
		)
		if v, err = base64.StdEncoding.DecodeString(string(b)); err == nil {
				return string(v)
		}
		errBase64 = err
		return ""
}
