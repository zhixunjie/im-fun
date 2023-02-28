package apihttp

// SetupRouter 设置-路由
func (s *Server) SetupRouter() {
	router := s.engine
	// 设置-单个路由
	router.GET("/ping", s.pingHandler)

	// 设置-路由组
	g1 := router.Group("/message")
	{
		g1.POST("/send", s.sendHandler)
		g1.GET("/fetch", s.fetchHandler)
	}
}
