package format

// MsgBody 消息体
// 设计参考：
// - https://cloud.tencent.com/document/product/269/2720
// - https://cloud.tencent.com/document/product/269/2282
type MsgBody struct {
	MsgType    MsgType    `json:"msg_type"`
	MsgContent MsgContent `json:"msg_content"`
}

type MsgType uint32

const (
	// MsgTypeNone 基本的消息类型
	MsgTypeNone     MsgType = 0 // 未知消息
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

type MsgContent interface {
	GetType() MsgType
	Decode([]byte) error
}

// JsonDecoderMap 某个类型对应的json解析器
var JsonDecoderMap = map[MsgType]func() MsgContent{
	MsgTypeCustom: func() MsgContent { return new(CustomContent) },
	MsgTypeText:   func() MsgContent { return new(TextContent) },
	MsgTypeTips:   func() MsgContent { return new(TipsContent) },
	MsgTypeImage:  func() MsgContent { return new(ImageContent) },
	MsgTypeVideo:  func() MsgContent { return new(VideoContent) },
	MsgTypeAudio:  func() MsgContent { return new(AudioContent) },
}
