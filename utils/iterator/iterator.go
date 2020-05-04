package iterator

type Decoder interface {
		Unmarshal(v interface{}) error
}

type IIteratorNode interface {
		Key() interface{}              // 键
		Value() interface{}            // 值
		CompareKey(interface{}) bool   // 是否对应键
		CompareValue(interface{}) bool // 是否对应键
		KeyString() string             // 键转字符串
}

type IIterator interface {
		Reset() IIterator            // 重置
		NextAble() bool              // 是否还可以遍历
		Cursor() interface{}         // 当前游标位置
		SetCursor(interface{})       // 设置游标位置
		Next() (IIteratorNode, bool) // 下一个值
		Current() IIteratorNode      // 当前值
}

type IAggregate interface {
		Iterator() IIterator
}
