package models

import (
		"github.com/WebGameLinux/cms/models/enums"
		"time"
)

// 用户数据 模型
type User struct {
		Id           int64        `orm:"column(id);pk;auto;description(序号)" json:"id"`
		Email        string       `orm:"size(128);description(邮箱地址);null" json:"email"`
		UserName     string       `orm:"column(username);unique;size(100);description(用户名)" json:"username"`
		Mobile       string       `orm:"column(mobile);index;size(20);description(手机号);null" json:"mobile"`
		PasswordHash string       `orm:"size(128);column(password_hash);description(密文密码)" json:"passwordHash"`
		Version      int          `orm:"column(version);default(1);description(记录版本号)" json:"version"`
		Gender       enums.Gender `orm:"column(gender);default(0);description(用户性别,0:未知,1:男,2:女,3:其他)" json:"gender"`
		UniqueSeqKey
		CreateDate
		SoftDeleteDate
}

// 日志模型
type Logs struct {
		IdPrimaryKeyIntN
		Payloads string    `orm:"column(payloads);type(text);size(1024);null;description(日志记录)" json:"payloads"`
		HappenAt time.Time `orm:"column(happen_at);type(datetime);null;description(产生时间)" json:"happen_at"`
		StateKey
		UniqueSeqKey
		CreateDate
		SoftDeleteDate
}

// 菜单项数据模型
type Menu struct {
		Id   int    `orm:"column(id);pk;auto;description(序号)" json:"id"`
		Name string `orm:"column(name);unique;description(菜单名称)" json:"name"`
		Info string `orm:"column(info);type(text);size(1024);null;description(菜单信息)" json:"info"`
		Pid  int    `orm:"column(pid);index;description(父级id);default(0)" json:"pid"`
		OrderSortInt
}

// 菜单树
type MenuTree struct {
		Menu
		SubMenus []MenuTree `json:"subMenus"`
}

// 菜单数据库模型
type MenuModel struct {
		Menu
		CreateDate
		SoftDeleteDate
}

// 权限模型
type RBac struct {

}
