package utils

import (
	"bytes"
)

type Map map[string]string
type rMap map[rune]interface{}

// Maches 去除注释的起始位置"start"到终止位置"end"
var Maches = []Map{
	Map{"start": "//", "end": "\n"},
	Map{"start": "/*", "end": "*/"},
}

// Discard 去除注释
// content：带注释的json 类型字符串
// return 参数: 字符串json 格式
func Discard(content string) (string, error) {
	var (
		buffer    bytes.Buffer
		flag      int
		v         rune
		protected bool
	)
	runes := []rune(content)
	flag = -1
	for i := 0; i < len(runes); {
		v = runes[i]
		if flag == -1 {
			// match start
			for f, v := range Maches {
				l := match(&runes, i, v["start"])
				if l != 0 {
					flag = f
					i += l
					break
				}
			}
			if flag == -1 {
				if protected {
					buffer.WriteRune(v)
					if v == '"' {
						protected = true
					}
				} else {
					r := filter(v)
					if r != 0 {
						buffer.WriteRune(v)
					}
				}
			} else {
				continue
			}
		} else {
			// match end
			l := match(&runes, i, Maches[flag]["end"])
			if l != 0 {
				flag = -1
				i += l
				continue
			}
		}
		i++
	}
	return buffer.String(), nil
}

func filter(v rune) rune {
	switch v {
	case ' ':
	case '\n':
	case '\t':
	default:
		return v
	}
	return 0
}

func match(runes *[]rune, i int, dst string) int {
	dstLen := len([]rune(dst))
	//fmt.Println("dstLen:", dstLen, ", index:", i, ",runesLen:", len(*runes))
	//fmt.Println(string((*runes)[i : i+dstLen]))
	if len(*runes)-i >= dstLen && string((*runes)[i:i+dstLen]) == dst {
		return dstLen
	}
	return 0
}

/*
// Stack TODO
type Stack []rune

// Push 进栈
func (s Stack) Push(r rune) {
	s = append(s, r)
}

// Pop 出栈
func (s Stack) Pop() (rune, error) {
	if len(s) == 0 {
		return 0, errors.New("stack is empty")
	}
	v := s[len(s)-1]
	s = s[:len(s)-1]
	return v, nil
}
*/
