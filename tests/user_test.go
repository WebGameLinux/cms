package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/array"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/WebGameLinux/cms/utils/types"
		"reflect"
		"testing"
)

func TestCreateUserPaginator(t *testing.T) {
		var p = models.CreateUserPaginator()
		fmt.Printf("%+v\n", p)
		fmt.Println(reflects.ClassName(types.StrToPassword("12323")))
		var m = make(mapper.Mapper)
		var m2 = make(map[string]interface{})
		fmt.Println(reflects.ClassName(m))
		fmt.Println(reflects.ClassName(reflect.Array))
		fmt.Println(reflects.ClassName(reflect.TypeOf(m2)))
		fmt.Println(reflects.ClassName([]interface{}{1, 2, 3}))
		fmt.Println(reflects.ClassName([]int{1, 2, 3}))
		fmt.Println(reflects.ClassName([]string{"12", "", "0", "123"}))
		fmt.Println(reflects.ClassName([]rune("12312312")))
		fmt.Println(reflects.ClassName([]byte("12312312")))
		fmt.Println(reflects.ClassName(map[string]string{"123": "1232", "2": "'"}))
}

func TestIsArray(t *testing.T) {
		fmt.Println(array.IsArray([]interface{}{1, 2, 3}))
		fmt.Println(array.IsArray([]int{1, 2, 3}))
		fmt.Println(array.IsArray([]string{"12", "", "0", "123"}))
		fmt.Println(array.IsArray([]rune("12312312")))
		fmt.Println(array.IsArray([]byte("12312312")))
		fmt.Println(array.IsArray(map[string]string{"123": "1232", "2": "'"}))
}

func TestJoinArray(t *testing.T) {
		var arr []interface{}
		arr = make([]interface{}, 0, 100)
		fmt.Println(len(arr), cap(arr))
		arr = array.JoinArrays(arr, []interface{}{1, 2, 3})
		fmt.Println(arr)
		arr = array.JoinArrays(arr, []int{1, 2, 3})
		fmt.Println(arr)
		arr = array.JoinArrays(arr, []string{"12", "", "0", "123"})
		fmt.Println(arr)
		arr = array.JoinArrays(arr, []rune("12312312"))
		fmt.Println(arr)
		arr = array.JoinArrays(arr, []byte("12312312"))
		fmt.Println(arr)

}
