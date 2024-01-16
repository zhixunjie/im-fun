package msg_body

type MsgType int32

const (
	// MsgTypeNone 基本的消息类型
	MsgTypeNone     MsgType = iota
	MsgTypeText             // 文本消息
	MsgTypeCustom           // 自定义消息
	MsgTypeImage            // 图片消息
	MsgTypeAudio            // 音频消息
	MsgTypeVideo            // 视频消息
	MsgTypeFile             // 文件消息
	MsgTypeFace             // 表情消息
	MsgTypeLocation         // 位置消息
	MsgTypeTips             // 提示消息
)

// MsgBody
// 设计参考：
// - https://cloud.tencent.com/document/product/269/2720
// - https://cloud.tencent.com/document/product/269/2282
type MsgBody[T ContentEntity] struct {
	MsgType    MsgType        `json:"msg_type"`
	MsgContent *MsgContent[T] `json:"msg_content"`
}

type MsgContent[T ContentEntity] struct {
	Content   *T  `json:"content"`              // 不同消息类型对应不同的结构体
	CheckFail int `json:"check_fail,omitempty"` // 让客户端展示感叹号！
}

type ContentEntity interface {
	GetType() MsgType
}

func NewMsgBody[T ContentEntity](t MsgType, Content *T) *MsgBody[T] {
	body := &MsgBody[T]{
		MsgType: t,
		MsgContent: &MsgContent[T]{
			Content: Content,
		},
	}

	return body
}

//type MsgContent struct {
//	// 不同消息类型对应不同的结构体
//	TextContent   *TextContent   `json:"text_content,omitempty"`   // 文本消息
//	TipsContent   *TipsContent   `json:"tips_content,omitempty"`   // 提示消息
//	CustomContent *CustomContent `json:"custom_content,omitempty"` // 自定义消息
//	ImageContent  *ImageContent  `json:"image_content,omitempty"`  // 图片消息
//	AudioContent  *AudioContent  `json:"audio_content,omitempty"`  // 音频消息
//
//	// 其他信息
//	CheckFail int `json:"check_fail,omitempty"` // 让客户端展示感叹号！
//}
