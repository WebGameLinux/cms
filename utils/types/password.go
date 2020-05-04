package types

import string2 "github.com/WebGameLinux/cms/utils/string"

type PasswordHashed string
type PasswordHashByte []byte

func StrToPassword(str string) PasswordHashed {
		return PasswordHashed(str)
}

func ByteToPassword(pass []byte) PasswordHashByte {
		return PasswordHashByte(pass)
}

func (p PasswordHashed) String() string {
		return string(p)
}

func (p PasswordHashed) Verify(str string, opts ...*string2.PasswordHashOptions) bool {
		return string2.PasswordVerify(str, p.String(), opts...)
}

func (p PasswordHashed) Args() *string2.PasswordHashOptions {
		if v, opt := string2.PasswordArgs(p.String()); len(v) != 0 && opt != nil {
				return opt
		}
		return nil
}

func (p PasswordHashed) Hash(options ...*string2.PasswordHashOptions) string {
		return string2.PasswordHash(p.String(), options...)
}

func (p PasswordHashed) Export() string {
		if v, opt := string2.PasswordArgs(p.String()); len(v) != 0 && opt != nil && opt.N > 0 {
				return string(v)
		}
		return ""
}

func (p PasswordHashed) IsHashed() bool {
		if v, opt := string2.PasswordArgs(p.String()); len(v) != 0 && opt != nil && opt.N > 0 {
				return true
		}
		return false
}

func (p PasswordHashByte) String() string {
		return string([]byte(p))
}

func (p PasswordHashByte) Verify(str string, opts ...*string2.PasswordHashOptions) bool {
		return string2.PasswordVerify(str, p.String(), opts...)
}

func (p PasswordHashByte) Args() *string2.PasswordHashOptions {
		if v, opt := string2.PasswordArgs(p.String()); len(v) != 0 && opt != nil {
				return opt
		}
		return nil
}

func (p PasswordHashByte) Hash(options ...*string2.PasswordHashOptions) string {
		return string2.PasswordHash(p.String(), options...)
}

func (p PasswordHashByte) Export() string {
		if v, opt := string2.PasswordArgs(p.String()); len(v) != 0 && opt != nil && opt.N > 0 {
				return string(v)
		}
		return ""
}

func (p PasswordHashByte) IsHashed() bool {
		if v, opt := string2.PasswordArgs(p.String()); len(v) != 0 && opt != nil && opt.N > 0 {
				return true
		}
		return false
}
