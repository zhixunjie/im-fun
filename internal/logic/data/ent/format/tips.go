package format

import "encoding/json"

// TipsContent 提示消息
type TipsContent struct {
	Text   string `json:"text,omitempty"`    // 文本内容
	ImgUrl string `json:"img_url,omitempty"` // 附带图片
}

func (c TipsContent) GetType() MsgType {
	return MsgTypeTips
}

func (c TipsContent) Decode(buf []byte) error {
	return json.Unmarshal(buf, &c)
}
