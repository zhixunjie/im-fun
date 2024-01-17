package msg_body

// 自定义消息

type CustomContent struct {
	Data string `json:"data,omitempty"` // 自定义消息（一般会采用JSON格式）
}

func (c CustomContent) GetType() MsgType {
	return MsgTypeCustom
}

//......
// 补充各类自定义的结构体
