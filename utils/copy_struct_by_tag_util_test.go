package utils

import (
	"testing"
	"time"
)

func TestCopyStructByTag(t *testing.T) {
	obj := &structC{Desc: ""}
	obj1 := &structA{
		Items:          []string{"233"},
		UserId:         "",
		PubTime:        time.Now().UnixNano(),
		Ext:            obj,
		StructBIsEmpty: "",
	}
	obj2 := &structB{}

	t.Logf("obj2 before copy: %+v.", obj2)

	err := CopyStructByTag(obj1, obj2, "mson")
	if err != nil {
		t.Errorf("Copy struct by tag err: %v.", err)
		return
	}

	t.Logf("obj2 after copy: %+v.", obj2)
}

type structA struct {
	Items   []string `mson:"Item_[]string"`
	UserId  string   `mson:"UserId_string"`
	PubTime int64    `mson:"PubTime_int64"`
	Ext     *structC

	StructBIsEmpty string `mson:"StructBIsEmpty_string"`
}

type structB struct {
	Item    []string  `mson:"Item_[]string"`
	UserId  int64     `mson:"UserId_int64"` // 和structA 类型不同，不会copy 生效
	PubTime time.Time `mson:"PubTime_Time"` // 和structA tag 类型不同，不会copy 生效
	Ext     *structC

	StructAIsEmpty string `mson:"StructAIsEmpty_string"`
}

type structC struct {
	Desc string `mson:"Desc_string"`
}
