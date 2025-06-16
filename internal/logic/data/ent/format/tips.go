package format

// TipsContent 提示消息
type TipsContent struct {
	Text   string `json:"text,omitempty"`    // 文本内容
	ImgUrl string `json:"img_url,omitempty"` // 附带图片
}

func (c *TipsContent) GetType() MsgType {
	return MsgTypeTips
}
