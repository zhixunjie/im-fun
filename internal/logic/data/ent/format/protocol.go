package format

// MsgBody 消息体
// 设计参考：
// - https://cloud.tencent.com/document/product/269/2720
// - https://cloud.tencent.com/document/product/269/2282
type MsgBody struct {
	MsgType    MsgType     `json:"msg_type"`
	MsgContent *MsgContent `json:"msg_content"`
}

// MsgContent 不同消息类型，对应不同的结构体（简单好用）
type MsgContent struct {
	CustomContent *CustomContent `json:"custom_content,omitempty"` // 自定义消息
	TextContent   *TextContent   `json:"text_content,omitempty"`   // 文本消息
	TipsContent   *TipsContent   `json:"tips_content,omitempty"`   // 提示消息
	ImageContent  *ImageContent  `json:"image_content,omitempty"`  // 图片消息
	AudioContent  *AudioContent  `json:"audio_content,omitempty"`  // 音频消息
	VideoContent  *VideoContent  `json:"video_content,omitempty"`  // 视频消息
}
type MsgType uint32

// 基本的消息类型
const (
	MsgTypeUnknown  MsgType = 0 // 未知消息
	MsgTypeCustom   MsgType = 1 // 自定义消息 (CustomContent)
	MsgTypeText     MsgType = 2 // 文本消息   (TextContent)
	MsgTypeImage    MsgType = 3 // 图片消息   (ImageContent)
	MsgTypeVideo    MsgType = 4 // 视频消息   (VideoContent)
	MsgTypeAudio    MsgType = 5 // 音频消息   (AudioContent)
	MsgTypeTips     MsgType = 6 // 提示消息   (TipsContent)
	MsgTypeFile     MsgType = 7 // 文件消息
	MsgTypeFace     MsgType = 8 // 表情消息
	MsgTypeLocation MsgType = 9 // 位置消息
)
