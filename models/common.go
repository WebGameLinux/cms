package models

import "time"

const ormUsingError = "<Ormer.Using> unknown db alias name"

// 业务码结构体
type BusinessCode struct {
		Code    int    `json:"code"` // 业务码
		Message string `json:"msg"`  // 业务提示
}

// int64 主键
type IdPrimaryKeyIntN struct {
		Id int64 `orm:"column(id);pk;auto;description(序号)" json:"id"`
}

// int 主键
type IdPrimaryKeyInt struct {
		Id int `orm:"column(id);pk;auto;description(序号)" json:"id"`
}

// 记录删除日期
type SoftDeleteDate struct {
		DeletedAt time.Time `orm:"column(deleted_at);type(datetime);null;description(删除时间)" json:"deleted_at"`
}

// 创建记录日期
type CreateDate struct {
		CreatedAt time.Time `orm:"column(created_at);type(datetime);auto_now_add;description(创建时时间)" json:"created_at"`
		UpdatedAt time.Time `orm:"column(updated_at);type(datetime);auto_now;description(更新时间)" json:"updated_at"`
}

// 记录删除时间戳
type SoftDeleteTimeStamp struct {
		DeletedAt time.Time `orm:"column(deleted_at);type(time);null;description(删除时间)" json:"deleted_at"`
}

// 创建记录时间戳
type CreateTimeStamp struct {
		CreatedAt time.Time `orm:"column(created_at);type(time);auto_now_add;description(创建时时间)" json:"created_at"`
		UpdatedAt time.Time `orm:"column(updated_at);type(time);auto_now;description(更新时间)" json:"updated_at"`
}

// 排序字段
type OrderSortInt struct {
		Sort int `orm:"column(sort);type(int);description(排序字段,越大越前);default(0)" json:"sort"`
}

// 	权重字段
type OrderWeightInt struct {
		Weight int `orm:"column(weight);type(int);description(权重字段);default(0)" json:"weight"`
}

// 唯一id
type UniqueSeqKey struct {
		SeqId string `orm:"column(seq_id);unique;size(128);description(用户分布式唯一ID)" json:"seq_id"`
}

// 唯一id
type StateKey struct {
		State int `orm:"column(state);description(状态字段,0:初始状态);default(0)" json:"state"`
}

// 唯一id
type UidKey struct {
		Uid int64 `orm:"column(uid);description(用户id,对应用户表id);default(0)" json:"uid"`
}
