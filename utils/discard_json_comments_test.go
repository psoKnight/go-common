package utils

import (
	"testing"
)

func TestDiscardComments(t *testing.T) {

	// 去除注释（默认方式）
	result, err := Discard(`
		{//test comment1
			"name": "测试",
			/**
			test comment2
			1
			2
			3
			end
			*/
			"age":26 //test comment3
			/*****/
		}
	`)
	if err != nil {
		t.Errorf("Discard err: %v.", err)
		return
	}
	t.Logf(result)
}

func TestCustomDiscard(t *testing.T) {

	// 自定义去除注释
	Maches = []Map{
		Map{"start": "$$", "end": "@@"},
	}
	result, err := Discard(`
		{//test comment1
			"name": "测试",
			/**$$
			test comment2
			1
			2
			3@@
			end
			*/
			"age":     26 //test comment3
			/*****/
		}
	`)
	if err != nil {
		t.Errorf("Discard err: %v.", err)
		return
	}
	t.Logf(result)
}
