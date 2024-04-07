package format

type VideoContent struct {
	Url      string `json:"url,omitempty"`      // 视频链接
	Duration int32  `json:"duration,omitempty"` // 视频的持续时间（秒）
}

func (c VideoContent) GetType() MsgType {
	return MsgTypeVideo
}
