package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	mrand "math/rand"
)

var (
	chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// RandString 随机生成指定长度的string
func RandString(l int) string {
	bs := []byte{}
	for i := 0; i < l; i++ {
		bs = append(bs, chars[mrand.Intn(len(chars))])
	}
	return string(bs)
}

// RandByte 随机生成指定长度的byte
func RandByte(l int) []byte {
	bs := []byte{}
	for i := 0; i < l; i++ {
		bs = append(bs, chars[mrand.Intn(len(chars))])
	}
	return bs
}

// RandCountForDiff 生成指定范围内的指定个数(不同的数字)
func RandCountForDiff(min, max int64, count int) ([]int64, error) {

	if min > max {
		return []int64{}, errors.New("Please check min and max.")
	} else if max-min < int64(count) {
		return []int64{}, errors.New("Please check relationship between area and count.")
	}

	var (
		allCount map[int64]int64
		result   []int64
	)
	allCount = make(map[int64]int64)
	maxBigInt := big.NewInt(max)
	for {
		// rand
		i, _ := rand.Int(rand.Reader, maxBigInt)
		number := i.Int64()
		// 是否大于下标
		if i.Int64() >= min {
			// 是否已经存在
			_, ok := allCount[number]
			if !ok {
				result = append(result, number)
				// 添加到map
				allCount[number] = number
			}
		}
		if len(result) >= count {
			return result, nil
		}
	}
}

// RandByArea 随机生成指定范围内的数
func RandByArea(min, max int64) int32 {
	maxBigInt := big.NewInt(max)
	for {
		i, _ := rand.Int(rand.Reader, maxBigInt)
		if i.Int64() >= min {
			return int32(i.Int64())
		}
	}
}
