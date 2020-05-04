package math

import (
		"crypto/rand"
		"math"
		"math/big"
)

func RandIntRangeN(min, max int64) int64 {
		if min > max {
				max, min = min, max
		}
		if min < 0 {
				f64Min := math.Abs(float64(min))
				i64Min := int64(f64Min)
				result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
				return result.Int64() - i64Min
		}
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
}

func RandIntRange(min, max int) int {
		return int(RandIntRangeN(int64(min), int64(max)))
}

// 随机获取数组中的一个值
func RandIntArrayValue(array []int) int {
		var any []interface{}
		for _, v := range array {
				any = append(any, v)
		}
		return RandArrayValue(any).(int)
}

// 随机生成 int array
func RandIntRangeArray(min, max, size int, minStep ...int) []int {
		var arr []int
		if min > max {
				max, min = min, max
		}
		if size == 0 {
				return arr
		}
		if len(minStep) == 0 {
				minStep = append(minStep, (max-min)/size)
		}
		m := max - minStep[0]
		for size > 0 {
				if m < min || m > max {
						m = min
				}
				n := RandIntRange(m, max)
				arr = append(arr, n)
				m = m - minStep[0]
				size--
		}
		return arr
}

func RandArrayValue(array []interface{}) interface{} {
		var max = len(array)
		if max == 0 {
				return nil
		}
		index := RandIntRange(0, max-1)
		return array[index]
}

func ShuffleArray(array []interface{}) []interface{} {
		var size = len(array)
		if size <= 0 {
				return array
		}
		for i := size - 1; i > 0; i-- {
				num := RandInt(i + 1)
				array[i], array[num] = array[num], array[i]
		}
		return array
}

func WrapperArray(any ...interface{}) []interface{} {
		return any
}

func RandInt(max ...int) int {
		if len(max) == 0 {
				max = append(max, math.MaxInt64)
		}
		if n, err := rand.Int(rand.Reader, big.NewInt(int64(max[0]))); err == nil {
				return int(n.Int64())
		}
		return 0
}
