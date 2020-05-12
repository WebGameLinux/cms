package models

import (
		"time"
)

const (
		SeqKey = "seq_id"
)

// 用户数据 模型
type User struct {
		Id           int64  `orm:"column(id);pk;auto;description(序号)" json:"id"`
		Email        string `orm:"size(128);description(邮箱地址);null" json:"email"`
		UserName     string `orm:"column(username);unique;size(100);description(用户名)" json:"username"`
		Mobile       string `orm:"column(mobile);index;size(20);description(手机号);null" json:"mobile"`
		PasswordHash string `orm:"size(128);column(password_hash);description(密文密码)" json:"passwordHash"`
		Version      int    `orm:"column(version);default(1);description(记录版本号)" json:"version"`
		Gender       int    `orm:"column(gender);default(0);description(用户性别,0:未知,1:男,2:女,3:其他)" json:"gender"`
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

// 附件信息
type Attachments struct {
		Id          int64  `orm:"column(id);pk;auto;description(序号)" json:"id"`
		Hash        string `orm:"column(hash);size(32);description(文件hash值)" json:"hash"`
		AccessUrl   string `orm:"column(access_url);size(255);description(对外访问路径)" json:"access_url"`
		SavePath    string `orm:"column(save_path);size(255);description(文件存储路径)" json:"save_path"`
		RelateTable string `orm:"column(relate_table);size(100);description(关联的表名);null" json:"relate_table"`
		RelateId    int    `orm:"column(relate_id);description(关联表对应ID);default(0)" json:"relate_id"`
		UploaderId  int64  `orm:"column(updater_id);description(上传文件用户id,0:默认为系统);default(0)" json:"uploader_id"`
		DeleterId   int64  `orm:"column(deleter_Id);description(删除文件用户id,0:默认为系统);default(0)" json:"deleter_id"`
		FileInfo
		UniqueSeqKey
		CreateDate
		SoftDeleteDate
}

// 文件信息
type FileInfo struct {
		FileName  string `orm:"column(filename);size(255);description(文件名)" json:"filename"`
		Size      string `orm:"column(size);size(50);description(文件大小,自动计算单位,b,kb,mb,gb,tb)" json:"size"`
		FileType  string `orm:"column(file_type);size(10);description(文件类型,doc,image,video,audio,html,config,code,password,binary,data)" json:"file_type"`
		Extension string `orm:"column(extension);size(10);description(文件扩展名,png,jpg,jpeg,mp3,mp4,ini,pdf,doc,docs...)" json:"extension"`
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
