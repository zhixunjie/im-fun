package msg_body

import (
	"strings"
	"unicode/utf8"
)

type HighLight map[string]any

type Params struct {
	Text   string `json:"text,omitempty"`   // 高亮的文本
	Link   string `json:"link,omitempty"`   // 跳转链接
	Color  string `json:"color,omitempty"`  // 高亮的颜色
	Offset [2]int `json:"offset,omitempty"` // 文本偏移量
}

// GetRange 获取高亮的位置
func GetRange(src, highLight string) [2]int {
	index := strings.Index(src, highLight)
	start := utf8.RuneCountInString(src[:index])
	lens := utf8.RuneCountInString(highLight)
	end := start + lens - 1

	return [2]int{start, end}
}
