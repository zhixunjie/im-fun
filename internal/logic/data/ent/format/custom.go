package format

import "encoding/json"

// 自定义消息

type CustomContent struct {
	Data string `json:"data,omitempty"` // 自定义消息（一般会采用JSON格式）
}

func (c CustomContent) GetType() MsgType {
	return MsgTypeCustom
}

func (c CustomContent) Decode(buf []byte) error {
	return json.Unmarshal(buf, &c)
}

type CustomEvent struct {
	Type EventType `json:"event_type"` // 事件类型
	Data string    `json:"data"`       // JSON字符串
}

//......
// 根据不同的Event，把Data解析为不同的结构体
// note：解析一条自定义消息，需要两次Unmarshal

type EventType int

const (
	EventNONE    EventType = 0
	EventLevelUp EventType = 1 // 人物升级 (EventJsonLevelUp)
	EventDropSth EventType = 2 // 物品掉落 (EventJsonDropSth)
)

type EventJsonLevelUp struct {
	UID       uint64 // 谁升级了？
	CurrLevel int    // 当前等级是？
}

type EventJsonDropSth struct {
	UID     uint64 // 掉落给谁？
	GoodsId int    // 物品id
}
