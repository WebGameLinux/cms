package reflects

import "time"

var parseError error

func ParseInt(v interface{}) int {
		return 0
}

func ParseStr(v interface{}) string {
		return Any2Str(v)
}

func ParseBool(v interface{}) bool {

		return false
}

func ParseIntN(v interface{}) int64 {

		return 0
}

func ParseInt32(v interface{}) int32 {

		return 0
}

func ParseFloat(v interface{}) float32 {

		return 0
}

func ParseFloatN(v interface{}) float64 {

		return 0
}

func ParseTime(v interface{}) *time.Time {

		return nil
}

func ParseDuration(v interface{}) time.Duration {

		return 0
}

func ParseMapper(v interface{}) *map[string]interface{} {

		return nil
}

func ParseIntArray(v interface{}) []int {

		return nil
}

func ParseAnyArray(v interface{}) []interface{} {

		return nil
}

func ParseBoolArray(v interface{}) []bool {

		return nil
}

func ParseStringArray(v interface{}) []string {
		return nil
}

func GetParseError() error {
		var err = parseError
		if err != nil {
				parseError = nil
		}
		return err
}
