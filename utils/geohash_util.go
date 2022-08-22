package utils

import (
	"errors"
	"fmt"
	"strings"
)

/**
具体的计算方法:
Latitude 的范围是:-90到+90 Longitude 的范围:-180到+180
地球参考球体的周长:40075016.68米
 geohash 长度 	Lat 位数 	Lng 位数 	Lat 误差 	Lng 误差 	km 误差
 1              2           3           ±23         ±23         ±2500
 2              5           5           ±2.8        ±5.6        ±630
 3              7           8           ±0.70       ±0.7 		±78
 4 				10 			10 			±0.087 		±0.18 		±20
 5 				12 			13 			±0.022 		±0.022 		±2.4
 6 				15 			15 			±0.0027 	±0.0055 	±0.61
 7 				17			18 			±0.00068 	±0.00068 	±0.076
 8 				20 			20 			+0.000086 	+0.000172 	+0.01911
 9 				22 			23 			±0.000021 	±0.000021 	±0.00478
 10 			25 			25 			±000000268 	+0.00000536 +0.0005971
 11 			27 			28 			±0.00000067 +000000067 	+0.0001492
 12 			30 			30 			+0.00000008 ±0.00000017 +0.0000186
*/

const (
	MAX_LON = 180
	MIN_LON = -180
	MAX_LAT = 90
	MIN_LAT = -90
)

var geoBase32Chars = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"b", "c", "d", "e", "f", "g", "h", "j", "k", "m",
	"n", "p", "q", "r", "s", "t", "u", "v", "w", "x",
	"y", "z",
}

var geoCodesByChar = map[byte]uint8{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'b': 10,
	'c': 11,
	'd': 12,
	'e': 13,
	'f': 14,
	'g': 15,
	'h': 16,
	'j': 17,
	'k': 18,
	'm': 19,
	'n': 20,
	'p': 21,
	'q': 22,
	'r': 23,
	's': 24,
	't': 25,
	'u': 26,
	'v': 27,
	'w': 28,
	'x': 29,
	'y': 30,
	'z': 31,
}

// 对单个坐标值进行geoHash 编码, 位数：digits
// 注意是使用低digits 位编码
func coordinateHash(coordinate, min, max float64, digits int) uint32 {
	hashValue := uint32(0)
	for i := 0; i < digits; i++ {
		mid := (min + max) / 2
		if coordinate >= mid {
			hashValue |= 1 << (digits - i - 1)
			min = mid
		} else {
			max = mid
		}
	}
	return hashValue
}

// 对单个坐标做反向hash
// 位数: digits 为uint32高digits 位
func coordinateReverseHash(hashValue uint32, min, max float64, digits int) float64 {
	mid := (min + max) / 2
	for i := 1; i <= digits; i++ {
		bitValue, _ := getUint32Bit(hashValue, uint8(i))
		if bitValue == 1 {
			min = mid
		} else {
			max = mid
		}
		mid = (min + max) / 2
	}
	return mid
}

// setUint32Bit 设置uint32指定位的值
// index 是从高位到低位, 1 - 32
func setUint32Bit(num *uint32, index, value uint8) (err error) {
	if index > 32 || index == 0 {
		err = errors.New("Index range is 1 to 32.")
		return
	}
	err = nil
	switch value {
	case 0:
		*num &= ^(1 << (32 - index))
	case 1:
		*num |= 1 << (32 - index)
	default:
		err = errors.New("Value can only be 0 or 1.")
	}
	return
}

// getUint32Bit 获取uint32指定位的值, index 从高位到低位: 1 - 32
func getUint32Bit(num uint32, index uint8) (value uint8, err error) {
	if index > 32 || index == 0 {
		err = errors.New("Index range is 1 to 32.")
		return
	}
	err = nil
	movBit := 32 - index
	value = uint8((num & (1 << movBit)) >> movBit)
	return
}

// setUint64Bit 设置uint64指定位的值
func setUint64Bit(num *uint64, index, value uint8) (err error) {
	if index > 64 || index == 0 {
		err = errors.New("Index range is 1 to 64.")
		return
	}
	err = nil
	switch value {
	case 0:
		*num &= ^(1 << (64 - index))
	case 1:
		*num |= 1 << (64 - index)
	default:
		err = errors.New("Value can only be 0 or 1.")
	}
	return
}

// getUint64Bit 获取uint64指定位的值
func getUint64Bit(num uint64, index uint8) (value uint8, err error) {
	if index > 64 || index == 0 {
		err = errors.New("Index range is 1 to 64.")
		return
	}
	err = nil
	movBit := 64 - index
	value = uint8((num & (1 << movBit)) >> movBit)
	return
}

func base32Encode(src []byte) (string, error) {
	dst := make([]string, len(src))
	for i, c := range src {
		if c > 31 {
			err := errors.New("The source Numbers range from 0 to 31.")
			return "", err
		}
		dst[i] = geoBase32Chars[c]
	}
	return strings.Join(dst, ""), nil
}

func base32Decode(encoded string) ([]byte, error) {
	encodedBytes := []byte(encoded)
	decodedBytes := make([]byte, len(encodedBytes))
	for i, enc := range encodedBytes {
		if dec, ok := geoCodesByChar[enc]; ok {
			decodedBytes[i] = dec
		} else {
			err := errors.New(fmt.Sprintf("Not other characters！Char: %c.", enc))
			return decodedBytes, err
		}
	}
	return decodedBytes, nil
}

// GeoHashEncode geoHash编码
// longitude: 经度
// latitude: 纬度
// precision: 精度, 范围: 1-12
// returns: geoHash 编码后的字符串, 错误信息
func GeoHashEncode(longitude, latitude float64, precision int) (string, error) {
	if precision < 1 || precision > 12 {
		err := errors.New("Precision range from 1 to 12.")
		return "", err
	}
	// 总位数
	totalDigits := precision * 5
	// 单个坐标值hash位数
	lonDigits := totalDigits / 2
	latDigits := lonDigits
	if totalDigits%2 != 0 {
		lonDigits++
	}
	lonHash := coordinateHash(longitude, MIN_LON, MAX_LON, lonDigits)
	latHash := coordinateHash(latitude, MIN_LAT, MAX_LAT, latDigits)
	lonHash = lonHash << (32 - lonDigits)
	latHash = latHash << (32 - latDigits)
	mergeHash := uint64(0)
	var i int
	for i = 1; i <= lonDigits; i++ {
		lonBit, _ := getUint32Bit(lonHash, uint8(i))
		latBit, _ := getUint32Bit(latHash, uint8(i))
		setUint64Bit(&mergeHash, uint8(i*2-1), lonBit)
		setUint64Bit(&mergeHash, uint8(i*2), latBit)
	}
	if totalDigits%2 != 0 {
		lonBit, _ := getUint32Bit(lonHash, uint8(i))
		setUint64Bit(&mergeHash, uint8(i*2-1), lonBit)
	}

	// 分组转换字节
	hashBytes := make([]byte, precision)
	for i = 1; i <= precision; i++ {
		hashBytes[i-1] = uint8((mergeHash >> (64 - i*5)) & 31)
	}

	geoDst, _ := base32Encode(hashBytes)

	return geoDst, nil
}

// 分割hash 为两个坐标的子hash
func splitHash(value uint64, digits int) (uint32, uint32) {
	lonHash := uint32(0)
	latHash := uint32(0)
	var i = 1
	for ; i <= digits-1; i += 2 {
		lonBit, _ := getUint64Bit(value, uint8(i))
		latBit, _ := getUint64Bit(value, uint8(i+1))
		setUint32Bit(&lonHash, uint8(i/2+1), lonBit)
		setUint32Bit(&latHash, uint8(i/2+1), latBit)
	}
	if digits&0x1 == 1 {
		lonBit, _ := getUint64Bit(value, uint8(i))
		setUint32Bit(&lonHash, uint8(i/2+1), lonBit)
	}
	return lonHash, latHash
}

// GeoHashDecode geoHash解码
// encoded: hash编码后的字符串, 长度: 1-12
// returns: 经度, 纬度, 错误信息
func GeoHashDecode(encoded string) (longitude, latitude float64, err error) {
	precision := len(encoded)
	if precision > 12 || precision == 0 {
		err = errors.New("Encoded length can only be 1 to 12.")
		return
	}
	decoded, err := base32Decode(encoded)
	if err != nil {
		return
	}
	decodedHash := uint64(0)
	for i, dec := range decoded {
		movBit := 64 - (i+1)*5
		decodedHash |= uint64(dec) << movBit
	}

	lonHash, latHash := splitHash(decodedHash, precision*5)
	lonDigits := precision * 5 / 2
	latDigits := lonDigits
	if precision&0x1 == 1 {
		lonDigits++
	}

	longitude = coordinateReverseHash(lonHash, MIN_LON, MAX_LON, lonDigits)
	latitude = coordinateReverseHash(latHash, MIN_LAT, MAX_LAT, latDigits)

	return
}
