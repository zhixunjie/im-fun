package http

// SetupRouter 设置-路由
func (s *Server) SetupRouter() {
	router := s.engine
	// 设置-单个路由
	router.GET("/ping", s.ping)

	// message
	g1 := router.Group("/message")
	{
		g1.POST("/send", s.send)
		g1.GET("/fetch", s.fetch)
	}

	// push
	router.POST("/push/user/keys", s.pushUserKeys) // 发送给指定的用户key
	router.POST("/push/user/ids", s.pushUserIds)   // 发送给指定的用户id
	router.POST("/push/user/room", s.pushUserRoom) // 广播给房间的所有用户
	router.POST("/push/user/all", s.pushUserAll)   // 广播给所有用户
}
