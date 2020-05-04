package mapper

import (
		"fmt"
		"testing"
)

func TestCompare(t *testing.T) {
		var val1 = "123"
		var val2 *string
		var val3 = "122"
		val2 = &val1
		fmt.Println(Compare(val1, val2))
		fmt.Println(Compare(val1, val1))
		fmt.Println(Compare(val1, val3))
		fmt.Println(Compare(val2, val2))
		fmt.Println(Compare(val2, val1))
		fmt.Println(Compare(val2, val3))
		fmt.Println(Compare(map[string]string{"1":"123","2":"12"}, map[string]string{"1":"123","2":"12"}))


}
