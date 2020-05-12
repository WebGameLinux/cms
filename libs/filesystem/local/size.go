package local

import (
		"fmt"
		"math"
		"regexp"
		"sort"
		"strconv"
		"strings"
)

type FileSize string

const (
		Bit               = 1
		Byte              = 8
		Unit              = 1024
		KB                = Unit * Unit
		MB                = KB * Unit
		GB                = MB * Unit
		TB                = GB * Unit
		PB                = TB * Unit
		EB                = PB * Unit
		ZB                = EB * Unit
		YB                = ZB * Unit
		BB                = YB * Unit
		NB                = BB * Unit
		DB                = NB * Unit
		BitUtil           = "b"
		BytUtil           = "Byte"
		KBUtil            = "KB"
		MBUtil            = "MB"
		GBUtil            = "GB"
		TBUtil            = "TB"
		PBUtil            = "PB"
		EBUtil            = "EB"
		ZBUtil            = "ZB"
		YBUtil            = "YB"
		BBUtil            = "BB"
		NBUtil            = "NB"
		DBUtil            = "DB"
		DefaultNumFormat  = "%.2f %s"
		DefaultSizeFormat = "%s %s"
		// NB,DB
		EmptySize     FileSize    = "0 byte"
		EmptySizeNum  FileSizeNum = 0
		MaxNumLen                 = 4
		RegExpPattern             = `([0-9]{1,}.?[0-9]{0,}) ?([bdemnkgtpzBDKEMNGTPZ]?([Bb]))$`
)

var sizeArr []float64
var sizeMatch = regexp.MustCompile(RegExpPattern)

var FileSizeUtils = map[float64]string{
		Bit:  BitUtil,
		Byte: BytUtil,
		KB:   KBUtil,
		MB:   MBUtil,
		GB:   GBUtil,
		TB:   TBUtil,
		PB:   PBUtil,
		EB:   EBUtil,
		ZB:   ZBUtil,
		YB:   YBUtil,
		BB:   BBUtil,
		NB:   NBUtil,
		DB:   DBUtil,
}

func GetFileSizeUtilArr() []float64 {
		if len(sizeArr) != 0 {
				return sizeArr
		}
		for n := range FileSizeUtils {
				sizeArr = append(sizeArr, n)
		}
		sort.Slice(sizeArr, func(i, j int) bool {
				return sizeArr[i] < sizeArr[j]
		})
		return sizeArr
}

func (this FileSize) ParseStrict(size string) []string {
		var ret = this.Parse(size)
		if len(ret) < 2 {
				return ret
		}
		num := ret[0]
		n := len(num)
		hasDot := strings.Contains(num, ".")
		if n > MaxNumLen && !hasDot {
				return []string{}
		}
		if n > 2*MaxNumLen+1 && hasDot {
				return []string{}
		}
		return ret
}

func (this FileSize) Next(level ...int) FileSize {
		size := this.Size()
		num, util := this.Compute(size)
		n, u := GetUtilNextN(util, level...)
		if u == util {
				return FileSize(fmt.Sprintf(DefaultNumFormat, num, util))
		}
		return FileSize(fmt.Sprintf(DefaultNumFormat, size/n, u))
}

func (this FileSize) Prev(level ...int) FileSize {
		size := this.Size()
		num, util := this.Compute(size)
		n, u := GetUtilPrevN(util, level...)
		if u == util {
				return FileSize(fmt.Sprintf(DefaultNumFormat, num, util))
		}
		return FileSize(fmt.Sprintf(DefaultNumFormat, size/n, u))
}

func (this FileSize) Parse(size string) []string {
		var ret []string
		if size == "" || size == "0" {
				return ret
		}
		if !sizeMatch.MatchString(size) {
				return ret
		}
		arr := sizeMatch.FindAllSubmatch([]byte(size), 2)
		if len(arr) > 0 && len(arr[0]) >= 3 {
				num := string(arr[0][1])
				util := string(arr[0][2])
				ret = append(ret, strings.Trim(num, " "), strings.Trim(util, " "))
		}
		return ret
}

// 整数部分
func (this FileSize) SizeInt() int64 {
		ret := this.ParseStrict(string(this))
		if len(ret) > 0 {
				if n, err := strconv.ParseInt(ret[0], 10, 64); err == nil {
						return n
				}
		}
		return 0
}

// 还原 float64 size
func (this FileSize) Size() float64 {
		ret := this.ParseStrict(string(this))
		if len(ret) <= 0 {
				return 0
		}
		size := ret[0]
		unit := ret[1]
		n, err := strconv.ParseFloat(size, 64)
		if err != nil {
				return 0
		}

		b := GetUtilMax(unit)

		return n * b
}

func (this FileSize) ParseInt(size int64) FileSize {
		return this.ParseFloat(float64(size))
}

func (this FileSize) ParseFloat(size float64) FileSize {
		if size <= 0 {
				return EmptySize
		}
		var util string
		size, util = this.Compute(size)
		formatStr := fmt.Sprintf(DefaultNumFormat, size, util)
		return FileSize(formatStr)
}

func (this FileSize) Compute(size float64) (num float64, util string) {
		var (
				index  float64 = -1
				arr            = GetFileSizeUtilArr()
				length         = len(arr)
		)
		for i := 0; i < length; i++ {
				if size > arr[i] {
						continue
				}
				if arr[i] != size {
						i = i - 1
				}
				index = arr[i]
				break
		}
		if index == -1 {
				index = arr[len(arr)-1]
		}
		return size / index, GetMaxSizeUtil(index)
}

func (this FileSize) String() string {
		if this == "" {
				return string(EmptySize)
		}
		return string(this)
}

func (this FileSize) Int() int64 {
		if this == EmptySize || this == "" {
				return 0
		}
		return this.SizeInt()
}

// @todo 更多格式支持
func (this FileSize) Format(str ...string) string {
		if len(str) == 0 {
				str = append(str, DefaultSizeFormat)
		}
		temp := str[0]
		arr := this.ParseStrict(this.String())
		num := len(arr)
		if num <= 0 {
				return this.String()
		}
		if strings.Contains(temp, "%") {
				if num >= 2 {
						return fmt.Sprintf(temp, arr[0], arr[1])
				}
		}
		return this.String()
}

func (this FileSize) Load(size string) FileSize {
		return FileSize(size)
}

func (this FileSize) Check() bool {
		return sizeMatch.MatchString(this.String())
}

type FileSizeNum float64

func (this FileSizeNum) Parse(size string) float64 {
		return FileSize(size).Size()
}

func (this FileSizeNum) Parse2Num(size string) FileSizeNum {
		return FileSizeNum(FileSize(size).Size())
}

func (this FileSizeNum) Next(level ...int) FileSize {
		num, util := this.Compute()
		n, u := GetUtilNextN(util, level...)
		if u == util {
				return FileSize(fmt.Sprintf(DefaultNumFormat, num, util))
		}
		return FileSize(fmt.Sprintf(DefaultNumFormat, float64(this)/n, u))
}

func (this FileSizeNum) Prev(level ...int) FileSize {
		num, util := this.Compute()
		n, u := GetUtilPrevN(util, level...)
		if u == util {
				return FileSize(fmt.Sprintf(DefaultNumFormat, num, util))
		}
		return FileSize(fmt.Sprintf(DefaultNumFormat, float64(this)/n, u))
}

func (this FileSizeNum) PasseInt(size int64) FileSizeNum {
		return FileSizeNum(float64(size))
}

func (this FileSizeNum) ParseFloatN(size float64) FileSizeNum {
		return FileSizeNum(size)
}

func (this FileSizeNum) String() string {
		return EmptySize.ParseFloat(float64(this)).String()
}

func (this FileSizeNum) Format(str ...string) string {
		size, util := this.Compute()
		if len(str) == 0 || !strings.Contains(str[0], "%") {
				return fmt.Sprintf(DefaultNumFormat, size, util)
		}
		return fmt.Sprintf(str[0], size, util)
}

func (this FileSizeNum) Compute() (num float64, util string) {
		var (
				index  float64 = -1
				size           = float64(this)
				arr            = GetFileSizeUtilArr()
				length         = len(arr)
		)

		for i := 0; i < length; i++ {
				if size > arr[i] {
						continue
				}
				if arr[i] != size {
						i = i - 1
				}
				index = arr[i]
				break
		}
		if index == -1 {
				index = arr[len(arr)-1]
		}
		return size / index, GetMaxSizeUtil(index)
}

// 获取度量单临界值
func GetUtilMax(unit string) float64 {
		for max, name := range FileSizeUtils {
				if strings.EqualFold(unit, name) {
						return max
				}
		}
		return 0
}

// 通过大小获取度量单位
func GetMaxSizeUtil(max float64) string {
		var last string
		if max <= 0 {
				return FileSizeUtils[Bit]
		}
		if v, ok := FileSizeUtils[max]; ok {
				return v
		}
		for _, num := range GetFileSizeUtilArr() {
				if num == max {
						return FileSizeUtils[num]
				}
				if num < max {
						last = FileSizeUtils[num]
				}
				if num > max {
						return last
				}
		}
		return FileSizeUtils[BB]
}

// 获取上一个度量单位
func GetUtilPrev(util string) (max float64, prev string) {
		return GetUtilPrevN(util, 1)
}

// 获取下一个度量单位
func GetUtilNext(util string) (max float64, next string) {
		return GetUtilNextN(util, 1)
}

func GetUtilNextN(util string, level ...int) (max float64, next string) {
		if len(level) == 0 {
				level = append(level, 1)
		}
		n := GetUtilMax(util)
		if n >= BB {
				return 0, util
		}
		i := 1
		levelNum := level[0]
		for _, cur := range GetFileSizeUtilArr() {
				if cur <= n {
						continue
				}
				max = cur
				if levelNum > 0 && i >= levelNum {
						break
				}
				i++
		}
		if max == 0 {
				return Bit, util
		}
		return max, GetMaxSizeUtil(max)
}

func GetUtilPrevN(util string, level ...int) (max float64, prev string) {
		var prevArr []float64
		if len(level) == 0 {
				level = append(level, 1)
		}
		levelNum := level[0]
		if levelNum <= 0 {
				return Bit, GetMaxSizeUtil(Bit)
		}
		prevArr = GetUtilPrevArr(util)
		size := len(prevArr)
		if size <= 0 {
				return Bit, util
		}
		index := size - levelNum
		if index < 0 {
				index = size - 1
		}
		max = prevArr[index]
		if max == 0 {
				return 0, util
		}
		return max, GetMaxSizeUtil(max)
}

func GetUtilPrevArr(util string) []float64 {
		n := GetUtilMax(util)
		var prevArr = make([]float64, 0)
		if n <= Bit {
				return prevArr
		}
		for _, cur := range GetFileSizeUtilArr() {
				if cur < n {
						prevArr = append(prevArr, cur)
				}
				if cur >= n {
						break
				}
		}
		return prevArr
}

func GetUtilPrevLevelMax(util string) int {
		return len(GetUtilPrevArr(util))
}

func GetUtilNextLevelMax(util string) int {
		n := GetUtilMax(util)
		arr := GetFileSizeUtilArr()
		if n >= BB {
				return 0
		}
		i := 0
		for _, cur := range arr {
				if cur <= n {
						continue
				}
				i++
		}
		return i
}

func GetTwoUtilMultiple(src string, dst string) float64 {
		if src == dst {
				return 1
		}
		n1 := GetUtilMax(src)
		n2 := GetUtilMax(dst)
		if n1 == n2 {
				return 1
		}
		n1, n2 = math.Max(n1, n2), math.Min(n1, n2)
		return n1 / n2
}
