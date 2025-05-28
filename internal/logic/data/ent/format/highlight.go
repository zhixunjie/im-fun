package format

import (
	"strings"
	"unicode/utf8"
)

// HighLight 高亮文本
// 注意：高亮文本的位置是基于源字符串的偏移量，而不是基于字节的偏移量
type HighLight struct {
	Text   string `json:"text,omitempty"`   // 高亮：文本
	Link   string `json:"link,omitempty"`   // 高亮：跳转链接
	Color  string `json:"color,omitempty"`  // 高亮：颜色
	Offset [2]int `json:"offset,omitempty"` // 文本偏移量（在源字符串中的偏移）
}

// GetOffset 获取高亮的位置
func GetOffset(src, highLight string) [2]int {
	index := strings.Index(src, highLight)
	start := utf8.RuneCountInString(src[:index])
	lens := utf8.RuneCountInString(highLight)
	end := start + lens - 1

	return [2]int{start, end}
}
