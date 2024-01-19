package http

// SetupRouter 设置-路由
func (s *Server) SetupRouter() {
	router := s.engine
	// 设置-单个路由
	router.GET("/ping", s.ping)

	// message
	g1 := router.Group("/message")
	{
		g1.POST("/send", s.sendMessage)           // 发送消息（普通消息） TODO 结合缓存机制优化
		g1.POST("/send/system", s.sendMessage)    // TODO 发送消息（系统消息）
		g1.POST("/fetch", s.fetchMessage)         // version_id拉取：消息列表 TODO 结合缓存机制优化
		g1.POST("/clean", s.fetchMessage)         // TODO 清空聊天记录
		g1.POST("/has/read", s.fetchMessage)      // TODO 消息已读
		g1.POST("/update/status", s.fetchMessage) // TODO 修改消息状态：消息删除 & 撤回消息
	}

	// contact
	g2 := router.Group("/contact")
	{
		g2.POST("/fetch", s.fetchContact)     // version_id拉取：会话列表 TODO 结合缓存机制优化
		g2.POST("/delete", s.fetchContact)    // TODO 删除一个会话
		g2.POST("/top/stick", s.fetchContact) // TODO 会话置顶
	}

	// push
	group := s.engine.Group("/im")
	{
		group.POST("/send/user/keys", s.sendToUserKeys) // 发送：给指定的用户key
		group.POST("/send/user/ids", s.sendToUserIds)   // 发送：给指定的用户id
		group.POST("/send/user/room", s.sendToRoom)     // 广播：给房间的所有用户
		group.POST("/send/user/all", s.sendToAll)       // 广播：给所有用户
	}
}
