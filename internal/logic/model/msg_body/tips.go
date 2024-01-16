package msg_body

type TipsContent struct {
	Text   string `json:"text,omitempty"`    // 文本内容
	ImgUrl string `json:"img_url,omitempty"` // 附带图片
}
