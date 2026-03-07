package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

// pushPayload 推送给客户端的消息体，与 /message/fetch 返回的 MsgEntity 格式对齐
type pushPayload struct {
	MsgID   uint64          `json:"msg_id,string"`
	FromUID uint64          `json:"from_user_id"`
	MsgBody *format.MsgBody `json:"msg_body"`
}

// MessageSendWithPush 发送消息并推送长链接
// 1. 从 Authorization token 中解析发送方身份
// 2. 调用 BzMessage.Send 将消息存入 DB
// 3. 调用 bz.SendToUsersByIds 通过 Kafka->Job->Comet 推送给在线接收方
func (s *Server) MessageSendWithPush(ctx *gin.Context) {
	logHead := "MessageSendWithPush|"

	// 1. 鉴权：从 token 解析 fromUserId
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.JsonError(ctx, errors.New("missing Authorization header"))
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.BzUser.CheckToken(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid token: " + err.Error()})
		return
	}
	fromUserId := claims.Uid
	logging.Infof(logHead+"fromUserId=%v", fromUserId)

	// 2. 解析请求体
	var req request.MessageSendWithPushReq
	if err = ctx.ShouldBindJSON(&req); err != nil {
		response.JsonError(ctx, err)
		return
	}
	if req.ReceiverUniId == "" {
		response.JsonError(ctx, errors.New("receiver_uni_id is required"))
		return
	}
	if req.MsgBody == nil {
		response.JsonError(ctx, errors.New("msg_body is required"))
		return
	}

	// 3. 存入 DB
	sender := gmodel.NewUserComponentId(fromUserId)
	receiver := gmodel.NewUserComponentId(cast.ToUint64(req.ReceiverUniId))
	sendRsp, err := s.BzMessage.Send(ctx, &request.MessageSendReq{
		Sender:   sender,
		Receiver: receiver,
		MsgBody:  req.MsgBody,
	})
	if err != nil {
		response.JsonError(ctx, err)
		return
	}

	// 4. 构造结构化推送消息体，推送给在线接收方
	payload, _ := json.Marshal(&pushPayload{
		MsgID:   sendRsp.Data.MsgID,
		FromUID: fromUserId,
		MsgBody: req.MsgBody,
	})
	pushErr := s.bz.SendToUsersByIds(ctx, &request.SendToUsersByIdsReq{
		UniIds:  []string{req.ReceiverUniId},
		Message: string(payload),
	})
	if pushErr != nil {
		// 推送失败不影响主流程（消息已入库），记录日志即可
		logging.Errorf(logHead+"SendToUsersByIds err=%v", pushErr)
	}

	ctx.JSON(http.StatusOK, sendRsp)
}
