package data

import (
	k "github.com/zhixunjie/im-fun/pkg/goredis/key"
)

// 待划分：划分出cache操作层（需要wire上Redis对象）

const Prefix = "im:logic:"

// 分布式锁：保证version_id和数据库写入的时序一致性
const (
	TimelineMessageLock k.Key = Prefix + "timeline:message:lock:{session_id}"
	TimelineContactLock k.Key = Prefix + "timeline:contact:lock:{contact_id}"
)

const (
	UserToken k.Key = Prefix + "user:token:{uid}"
)

const (
	KeyExpire = 3600
)

// mark tcp connection
const (
	// TcpUserAllSession Hash：
	// 格式：userId -> [ tcpSessionId : serverId ]
	TcpUserAllSession k.Key = Prefix + "tcp:user:all:session:{uid}"
	// TcpSessionToSrv String
	// 格式：tcpSessionId -> serverId
	TcpSessionToSrv k.Key = Prefix + "tcp:session:to:server:{tcp_session_id}"
	// TcpServerOnline String
	// 格式：serverId -> online
	TcpServerOnline k.Key = Prefix + "tcp:server:online:{server_id}"
)
