package format

type VideoContent struct {
	VideoUrl    string `json:"video_url,omitempty"`    // 视频：链接
	VideoSecond int32  `json:"video_second,omitempty"` // 视频：时长（秒）
	VideoUuid   string `json:"video_uuid,omitempty"`   // 资源标识
	VideoSize   int32  `json:"video_size"`             // 视频：大小
	ThumbUrl    string `json:"thumb_url,omitempty"`    // 视频封面：链接
	ThumbWidth  int32  `json:"thumb_width,omitempty"`  // 视频封面：宽
	ThumbHeight int32  `json:"thumb_height,omitempty"` // 视频封面：高
	ThumbUuid   string `json:"thumb_uuid,omitempty"`   // 资源标识
	ThumbSize   int32  `json:"thumb_size"`             // 视频封面：大小
}

func (c VideoContent) GetType() MsgType {
	return MsgTypeVideo
}
