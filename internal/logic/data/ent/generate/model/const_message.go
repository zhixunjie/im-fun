package model

const (
	TotalDb           = 10
	TotalTableMessage = 100 // message表：分表个数（一共10个数据库，每个数据库100个表）
	TotalTableContact = 100 // contact表：分表个数（一共10个数据库，每个数据库100个表）
)

// 消息读取状态
const (
	MsgNotRead = 0 // 未读
	MsgRead    = 1 // 已读
)

// 消息状态
const (
	MsgStatusNormal   = 0 // 正常
	MsgStatusDel      = 1 // 删除
	MsgStatusWithdraw = 2 // 后台删除
)
