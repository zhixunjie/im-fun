package format

// 自定义消息

type CustomContent struct {
	EventType EventType `json:"event_type"`
	Data      string    `json:"data,omitempty"` // 自定义消息（一般会采用JSON格式）
}

func (c CustomContent) GetType() MsgType {
	return MsgTypeCustom
}

//......
// 根据不同的Event，把Data解析为不同的结构体
// note：解析一条自定义消息，需要两次Unmarshal

type EventType int

const (
	EventNONE    EventType = 0
	EventLevelUp EventType = 1 // 人物升级（对应结构体: ImJsonLevel）
	EventDrop    EventType = 2 // 物品掉落（对应结构体: ImJsonDrop）
)

type ImJsonLevel struct {
	UID       uint64 // 谁升级了？
	CurrLevel int    // 当前等级是？
}

type ImJsonDrop struct {
	UID     uint64 // 掉落给谁？
	GoodsId int    // 物品id
}
