package test

import (
		"fmt"
		_ "github.com/WebGameLinux/cms/bootstarp"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/mapper"
		"testing"
)

func TestUserQuery(t *testing.T) {
		model := models.GetUser()
		user := model.GetByKey("id", 1)

		if user2, ok := model.NewModel().(*models.User); ok {
				fmt.Println(user2)
		}
		users := model.Lists(mapper.Mapper{})
		if model.GetError() != nil {
				t.Error("lists查询失败")
		} else {
				fmt.Println(users)
		}
		if user == nil {
				t.Error("查询失败")
		}
}
