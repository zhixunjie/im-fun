package format

type AudioContent struct {
	Url    string `json:"url,omitempty"`    // 音频：链接
	Second int32  `json:"second,omitempty"` // 音频：时长（秒）
	Uuid   string `json:"uuid,omitempty"`   // 资源标识
	Text   string `json:"text,omitempty"`   // 音频：附带文本
}

func (c AudioContent) GetType() MsgType {
	return MsgTypeAudio
}
