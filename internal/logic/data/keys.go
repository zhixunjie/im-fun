package data

// 待划分：划分出cache操作层（需要wire上Redis对象）

const Prefix = "im:logic:"

// 分布式锁：保证version_id和数据库写入的时序一致性
const (
	TimelineMessageLock = Prefix + "timeline:message:lock:{%v}"
	TimelineContactLock = Prefix + "timeline:contact:lock:{%v}"
)

const (
	UserToken = Prefix + "user:token:{%v}"
)

const (
	KeyExpire = 3600
)

// mark tcp connection
const (
	// TcpUserAllSession Hash：
	// 格式：uniId -> [ tcpSessionId : serverId ]
	TcpUserAllSession = Prefix + "tcp:user:all:session:{%v}"
	// TcpSessionToSrv String
	// 格式：tcpSessionId -> serverId
	TcpSessionToSrv = Prefix + "tcp:session:to:server:{%v}"
	// TcpServerOnline String
	// 格式：serverId -> online
	TcpServerOnline = Prefix + "tcp:server:online:{%v}"
)
