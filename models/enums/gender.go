package enums

import (
		"fmt"
		"strconv"
)

type Gender int8

const (
		Unknown Gender = iota
		Man
		WoMan
		Other
)

type Dict map[string]string

var (
		genderMap = map[Gender]Dict{
				Unknown: {"zh": "保密未知", "en": "unknown"},
				Man:     {"zh": "男", "en": "man"},
				WoMan:   {"zh": "女", "en": "woman"},
				Other:   {"zh": "其他", "en": "other"},
		}
		local = "zh"
)

func Parse(str string) Gender {
		for g, dict := range genderMap {
				for _, v := range dict {
						if v == str {
								return g
						}
				}
		}
		return Unknown
}

func SetLocal(lang LocalLang) {
		local = string(lang)
}

func GetLocal() string {
		return string(local)
}

func GetGenderMap() map[Gender]Dict {
		return genderMap
}

func ParseInt(gender int) Gender {
		var g = Gender(gender)
		if g < Unknown || g > Other {
				return Unknown
		}
		return g
}

func (g Gender) String() string {
		if v, ok := genderMap[g][local]; ok {
				return v
		}
		return genderMap[Unknown][local]
}

func (g Gender) Json() string {
		return fmt.Sprintf(`{"gender":%d,"dict":%s}`, g, genderMap[g])
}

func (g Gender) Parse(gender string) Gender {
		return Parse(gender)
}

func (g Gender) Map() Dict {
		var m = make(Dict)
		m[strconv.Itoa(g.Int())] = g.String()
		return m
}

func (g Gender) Int() int {
		return int(g)
}
