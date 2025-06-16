package format

import "encoding/json"

// ImageContent 图片消息
type ImageContent struct {
	ImageInfos []*Image `json:"image_infos,omitempty"` // 图片列表
}

type Image struct {
	Type   int32  `json:"type,omitempty"`   // 图片：类型，1-原图、2-大图、3-缩略图
	Url    string `json:"url,omitempty"`    // 图片：链接
	Width  int32  `json:"width,omitempty"`  // 图片：宽
	Height int32  `json:"height,omitempty"` // 图片：高
	Size   int32  `json:"size,omitempty"`   // 图片：大小
	Uuid   string `json:"uuid,omitempty"`   // 资源标识
}

func (c ImageContent) GetType() MsgType {
	return MsgTypeImage
}

func (c ImageContent) Decode(buf []byte) error {
	return json.Unmarshal(buf, &c)
}
