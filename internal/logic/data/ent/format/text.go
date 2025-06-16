package format

// TextContent 文本消息
type TextContent struct {
	Text       string      `json:"text,omitempty"` // 文本内容
	HighLights []HighLight `json:"highLights"`     // 高亮文本（支持多段高亮）
}

func (c *TextContent) GetType() MsgType {
	return MsgTypeText
}
