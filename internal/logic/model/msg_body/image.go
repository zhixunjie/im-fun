package msg_body

type ImageContent struct {
	List []Image `json:"list,omitempty"` // 图片列表

}
type Image struct {
	Url    string `json:"url,omitempty"`    // 图片链接
	Width  int32  `json:"width,omitempty"`  // 图片宽
	Height int32  `json:"height,omitempty"` // 图片高
	Size   int32  `json:"size"`             // 图片大小
}