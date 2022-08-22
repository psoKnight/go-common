package utils

import (
	"github.com/ybzhanghx/copier"
)

// CopyStructByTag 按标签复制struct
// fromStruct：待复制struct 对象
// toStruct：需复制到的struct 对象
func CopyStructByTag(fromStruct interface{}, toStruct interface{}, tag string) error {
	if err := copier.CopyByTag(toStruct, fromStruct, tag); err != nil {
		return err
	}
	return nil
}
