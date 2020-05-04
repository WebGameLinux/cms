package string

import (
		"fmt"
		"github.com/WebGameLinux/cms/utils/math"
		"golang.org/x/crypto/scrypt"
		"strconv"
		"strings"
		"time"
)

type PasswordHashOptions struct {
		Salt            []byte
		N, R, P, keyLen int
		AutoSaltBool    bool
}

var (
		PassSalt        string
		encodeErr       error
		sortMapArr                           = []string{"p", "N", "keyLen", "r"}
		sortIndex                            = []int{1, 3, 0, 4}
		PasswordOptions *PasswordHashOptions = &PasswordHashOptions{AutoSaltBool: true}
		nArr                                 = []int{512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}
)

func (options *PasswordHashOptions) Init() {
		if options.Booted() {
				return
		}
		options.compute()
}

func (options *PasswordHashOptions) compute() {
		options.P = 1
		// N r p keylen
		// 32768, 8, 1, 32
		// 16384, 8, 1, 32
		options.R = 8
		options.keyLen = 32
		options.N = options.computeN()
		if options.AutoSaltBool && options.IsEmptySalt() {
				options.Salt = []byte(options.AutoSalt())
		}
}

//  N&(N-1) N&(N-1) != 0
func (options *PasswordHashOptions) computeN() int {
		return math.RandIntArrayValue(nArr)
}

// 是否初始化
func (options *PasswordHashOptions) Booted() bool {
		if options.keyLen == 0 {
				return false
		}
		if options.P == 0 || options.R == 0 || options.N == 0 {
				return false
		}
		return true
}

func (options *PasswordHashOptions) Salting() string {
		if !options.Booted() {
				options.Init()
		}
		if options.AutoSaltBool && options.IsEmptySalt() {
				options.Salt = []byte(options.AutoSalt())
		}
		return string(options.Salt)
}

func (options *PasswordHashOptions) GetByte(key string) []byte {
		if !options.Booted() {
				options.Init()
		}
		if key == "salt" {
				return options.Salt
		}
		return nil
}

func (options *PasswordHashOptions) GetInt(key string) int {
		if !options.Booted() {
				options.Init()
		}
		switch key {
		case "p":
				fallthrough
		case "P":
				return options.P
		case "N":
				fallthrough
		case "n":
				return options.N
		case "r":
				fallthrough
		case "R":
				return options.R
		case "keyLen":
				fallthrough
		case "KeyLen":
				fallthrough
		case "keylen":
				return options.keyLen
		}
		return 0
}

func (options *PasswordHashOptions) IsEmptySalt() bool {
		return len(options.Salt) == 0
}

func (options *PasswordHashOptions) AutoSalt() string {
		if PassSalt == "" && 0 == len(options.Salt) {
				now := time.Now()
				salt := fmt.Sprintf("%d%d", now.Minute(), now.Second())
				options.Salt = []byte(salt)
				return salt
		}
		options.Salt = []byte(PassSalt)
		return PassSalt
}

func (options *PasswordHashOptions) Set(key string, value interface{}) {
		switch key {
		case "salt":
				fallthrough
		case "Salt":
				if str, ok := value.(string); ok {
						options.Salt = []byte(str)
						return
				}
				if salt, ok := value.([]byte); ok {
						options.Salt = salt
				}
				return
		case "p":
				fallthrough
		case "P":
				if v, ok := value.(string); ok {
						if num, err := strconv.Atoi(v); err == nil {
								options.P = num
						}
						return
				}
				if num, ok := value.(int); ok {
						options.P = num
						return
				}
				if num, ok := value.(*int); ok {
						options.P = *num
						return
				}
				return
		case "N":
				fallthrough
		case "n":
				if num, ok := value.(int); ok {
						options.N = num
						return
				}
				if num, ok := value.(*int); ok {
						options.N = *num
						return
				}
				if v, ok := value.(string); ok {
						if num, err := strconv.Atoi(v); err == nil {
								options.N = num
						}
				}
				return
		case "r":
				fallthrough
		case "R":
				if num, ok := value.(int); ok {
						options.R = num
						return
				}
				if num, ok := value.(*int); ok {
						options.R = *num
						return
				}
				if v, ok := value.(string); ok {
						if num, err := strconv.Atoi(v); err == nil {
								options.R = num
						}
				}
				return
		case "KeyLen":
				fallthrough
		case "keyLen":
				fallthrough
		case "keylen":
				if num, ok := value.(int); ok {
						options.keyLen = num
						return
				}
				if num, ok := value.(*int); ok {
						options.keyLen = *num
						return
				}
				if v, ok := value.(string); ok {
						if num, err := strconv.Atoi(v); err == nil {
								options.keyLen = num
						}
				}
		}
		return
}

func PasswordSortMapper() []string {
		return sortMapArr
}

func passwordString(pass []byte, option *PasswordHashOptions) string {
		var p = fmt.Sprintf("%x", pass)
		if p == "" || option == nil {
				return ""
		}
		for i, k := range PasswordSortMapper() {
				if i%2 == 0 {
						p = fmt.Sprintf("%d.%s", option.GetInt(k), p)
				} else {
						p = fmt.Sprintf("%s.%d", p, option.GetInt(k))
				}
		}
		p = p + "." + ByteBase64Encode(option.GetByte("salt"))
		return p
}

func PasswordArgs(text string) ([]byte, *PasswordHashOptions) {
		if text == "" || !strings.Contains(text, ".") {
				return nil, nil
		}

		var vars = strings.SplitN(text, ".", -1)
		if len(vars) != 6 {
				return nil, nil
		}
		pass := vars[2]
		opt := new(PasswordHashOptions)
		salt := Base64Decode2Byte(vars[5])
		index := getSortIndexArray()
		for i, key := range PasswordSortMapper() {
				opt.Set(key, vars[index[i]])
		}
		opt.Salt = salt
		return []byte(pass), opt
}

func getSortIndexArray() []int {
		return sortIndex
}

func PasswordHash(password string, options ...*PasswordHashOptions) string {
		if len(options) == 0 {
				options = append(options, PasswordOptions)
		}
		var (
				dk     []byte
				err    error
				opt    = options[0]
				salt   = opt.GetByte("salt")
				N      = opt.GetInt("N")
				r      = opt.GetInt("r")
				p      = opt.GetInt("p")
				keyLen = opt.GetInt("keyLen")
		)
		if dk, err = scrypt.Key([]byte(password), salt, N, r, p, keyLen); err == nil {
				return passwordString(dk, opt)
		}
		encodeErr = err
		return ""
}

// 获取异常
func GetEncodeError() error {
		defer func() { encodeErr = nil }()
		return encodeErr
}

func PasswordVerify(text string, password string, options ...*PasswordHashOptions) bool {
		if len(options) == 0 {
				_, option := PasswordArgs(password)
				if option == nil {
						return false
				}
				options = append(options, option)
		}
		return PasswordHash(text, options...) == password
}
