package msg_body

type AudioContent struct {
	Url      string `json:"url,omitempty"`      // 音频链接
	Duration int32  `json:"duration,omitempty"` // 音频的持续时间（秒）
	Text     string `json:"text,omitempty"`     // 音频的附带文本
}

func (c AudioContent) GetType() MsgType {
	return MsgTypeAudio
}
