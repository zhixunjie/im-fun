package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// SortNum 两个数字排序，小的数字在前面，大的数字在后面
func SortNum(a, b uint64) (uint64, uint64) {
	if a > b {
		return b, a
	}
	return a, b
}

// LargerNum 返回更大的一个数字
func LargerNum(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func PrettyPrint(str []byte) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, str, "", " "); err != nil {
		return
	}

	fmt.Println(prettyJSON.String())
}

func Max[T int64 | uint64](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T int64 | uint64](a, b T) T {
	if a < b {
		return a
	}

	return b
}
