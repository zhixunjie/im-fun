package http

// SetupRouter 设置-路由
func (s *Server) SetupRouter() {
	router := s.engine
	// 设置-单个路由
	router.GET("/ping", s.ping)

	// message
	message := router.Group("/message")
	{
		message.POST("/send", s.MessageSend)           // ✅发送消息（普通消息） TODO 结合缓存机制优化
		message.POST("/send/system", s.MessageSend)    // TODO 发送消息（系统消息）
		message.POST("/fetch", s.MessageFetch)         // ✅version_id拉取：消息列表 TODO 结合缓存机制优化
		message.POST("/clear", s.MessageClearHistory)  // ✅清空聊天记录
		message.POST("/has/read", s.MessageFetch)      // TODO 消息已读
		message.POST("/update/status", s.MessageFetch) // TODO 修改消息状态：消息删除 & 撤回消息
	}

	// contact
	contact := router.Group("/contact")
	{
		contact.POST("/fetch", s.ContactFetch)     // ✅version_id拉取：会话列表 TODO 结合缓存机制优化
		contact.POST("/delete", s.ContactFetch)    // TODO 删除一个会话
		contact.POST("/top/stick", s.ContactFetch) // TODO 会话置顶
	}

	// push
	im := s.engine.Group("/im")
	{
		im.POST("/send/to/users", s.sendToUsers)             // 发送：给指定的用户key
		im.POST("/send/to/users/by/ids", s.sendToUsersByIds) // 发送：给指定的用户id
		im.POST("/send/to/room", s.sendToRoom)               // 广播：给房间的所有用户
		im.POST("/send/to/all", s.sendToAll)                 // 广播：给所有用户
	}
}
