package response

type PingResp struct {
	Base
	Pong string `json:"pong"`
}
