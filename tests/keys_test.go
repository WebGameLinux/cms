package test

import (
		"fmt"
		. "github.com/WebGameLinux/cms/utils/reflects"
		"testing"
		"time"
)

type SoftDeleteDate struct {
		DeletedAt time.Time `orm:"column(deleted_at);type(date);null;description(删除时间)" json:"deleted_at"`
}

type CreateDate struct {
		CreatedAt time.Time `orm:"column(created_at);type(date);auto_now_add;description(创建时时间)" json:"created_at"`
		UpdatedAt time.Time `orm:"column(updated_at);type(date);auto_now;description(更新时间)" json:"updated_at"`
}

type SoftDeleteTimeStamp struct {
		DeletedAt time.Time `orm:"column(deleted_at);type(time);null;description(删除时间)" json:"deleted_at"`
}

type User struct {
		Id           int64  `orm:"column(id);pk;auto;description(序号)" json:"id"`
		Email        string `orm:"size(128);description(邮箱地址)" json:"email"`
		UserName     string `orm:"column(username);unique;size(100);description(用户名)" json:"username"`
		PasswordHash string `orm:"size(128);column(password_hash);description(密文密码)" json:"passwordHash"`
		SeqId        string `orm:"column(seq_id);size(128);description(用户分布式唯一ID)" json:"seq_id"`
		Version      int    `orm:"column(version);default(1);description(记录版本号)" json:"version"`
		CreateDate
		SoftDeleteDate
}

func TestGetItemsTypes(t *testing.T) {
		var user = new(User)
		user.PasswordHash = "user"
		user.UserName = "hello"
		fmt.Println(GetItemsTypes(user))
}

func TestGetItemsValues(t *testing.T) {
		var user = new(User)
		user.PasswordHash = "user"
		user.UserName = "hello"
		user.Version = 1
		fmt.Println(GetItemsValues(user))
}

func TestGetAllItemsTypes(t *testing.T) {
		var user = new(User)
		user.PasswordHash = "user"
		user.UserName = "hello"
		user.Version = 1
		fmt.Println(GetItemsAllTypes(user))
}

func TestGetAllItemsValues(t *testing.T) {
		var user = new(User)
		user.PasswordHash = "user"
		user.UserName = "hello"
		user.Version = 1
		fmt.Println(GetItemsAllValues(user))
}
