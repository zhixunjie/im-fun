package biz

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/api"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
)

type MessageFilterUseCase struct {
}

func NewMessageFilterUseCase() *MessageFilterUseCase {
	return &MessageFilterUseCase{}
}

// FilterMessageUser 过滤发送者
func (b *MessageFilterUseCase) FilterMessageUser(id1, id2 *gmodel.ComponentId) (err error) {
	if id1 == nil || id2 == nil {
		err = fmt.Errorf("%w,id not exists", api.ErrSenderOrReceiverNotAllow)
		return
	}
	if id1.GetId() == 0 || id2.GetId() == 0 {
		err = fmt.Errorf("%w,ID is zero", api.ErrSenderOrReceiverNotAllow)
		return
	}
	if id1.Equal(id2) {
		err = fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
		return
	}
	// 此接口不适合群聊场景
	if id1.IsGroup() || id2.IsGroup() {
		err = fmt.Errorf("group not allowed %w", api.ErrSenderOrReceiverNotAllow)
		return
	}
	return
}

// FilterGroupMessageUser 过滤发送者（群聊场景）
func (b *MessageFilterUseCase) FilterGroupMessageUser(id1, id2 *gmodel.ComponentId) (err error) {
	if id1 == nil || id2 == nil {
		err = fmt.Errorf("%w,id not exists", api.ErrSenderOrReceiverNotAllow)
		return
	}
	if id1.GetId() == 0 || id2.GetId() == 0 {
		err = fmt.Errorf("%w,ID is zero", api.ErrSenderOrReceiverNotAllow)
		return
	}
	if id1.Equal(id2) {
		err = fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
		return
	}
	// 此接口只适合群聊场景
	if !id1.IsGroup() && !id2.IsGroup() {
		err = fmt.Errorf("group not allowed %w", api.ErrSenderOrReceiverNotAllow)
		return
	}
	return
}

// FilterMsgContent 限制消息内容
func (b *MessageFilterUseCase) FilterMsgContent(msgBody *format.MsgBody) (err error) {
	if msgBody == nil {
		err = fmt.Errorf("msgBody nil")
		return
	}
	content := msgBody.MsgContent

	// check: message type
	typeLimit := []format.MsgType{
		format.MsgTypeCustom,
		format.MsgTypeText,
		format.MsgTypeImage,
		format.MsgTypeVideo,
		format.MsgTypeTips,
	}
	if !lo.Contains(typeLimit, msgBody.MsgType) {
		return api.ErrMessageTypeNotAllowed
	}

	// check: message content
	switch msgBody.MsgType {
	case format.MsgTypeCustom:
		if content.CustomContent == nil {
			err = fmt.Errorf("%w,CustomContent nil", api.ErrMessageContentNotAllowed)
			return
		}
		if content.CustomContent.Data == "" {
			err = fmt.Errorf("%w,data is empty", api.ErrMessageContentNotAllowed)
			return
		}
	case format.MsgTypeText:
		if content.TextContent == nil {
			err = fmt.Errorf("%w,TextContent nil", api.ErrMessageContentNotAllowed)
			return
		}
		if content.TextContent.Text == "" {
			err = fmt.Errorf("%w,text empty", api.ErrMessageContentNotAllowed)
			return
		}
	case format.MsgTypeImage:
		if content.ImageContent == nil {
			err = fmt.Errorf("%w,ImageContent nil", api.ErrMessageContentNotAllowed)
			return
		}
		if len(content.ImageContent.ImageInfos) == 0 {
			err = fmt.Errorf("%w,image array empty", api.ErrMessageContentNotAllowed)
			return
		}
	case format.MsgTypeVideo:
		if content.VideoContent == nil {
			err = fmt.Errorf("%w,VideoContent nil", api.ErrMessageContentNotAllowed)
			return
		}
		if content.VideoContent.VideoUrl == "" {
			return fmt.Errorf("%w,video url empty", api.ErrMessageContentNotAllowed)
		}
		if content.VideoContent.VideoSecond == 0 {
			return fmt.Errorf("%w,video second is zero", api.ErrMessageContentNotAllowed)
		}
	case format.MsgTypeAudio:
		if content.AudioContent == nil {
			err = fmt.Errorf("%w,AudioContent nil", api.ErrMessageContentNotAllowed)
			return
		}
		if content.AudioContent.Url == "" {
			return fmt.Errorf("%w,audio url empty", api.ErrMessageContentNotAllowed)
		}
		if content.AudioContent.Second == 0 {
			return fmt.Errorf("%w,audio second is zero", api.ErrMessageContentNotAllowed)
		}
	case format.MsgTypeTips:
		if content.TipsContent == nil {
			err = fmt.Errorf("%w,TextContent nil", api.ErrMessageContentNotAllowed)
			return
		}
		if content.TipsContent.Text == "" {
			return fmt.Errorf("%w,tip's text empty", api.ErrMessageContentNotAllowed)
		}
	}
	return
}

// FilterSendUserType 限制发送方和接收方的类型
func (b *MessageFilterUseCase) FilterSendUserType(sender, receiver *gmodel.ComponentId) (err error) {
	allowSenderType := []gmodel.ContactIdType{
		gmodel.TypeUser,
		gmodel.TypeRobot,
		gmodel.TypeSystem,
	}

	allowReceiverType := []gmodel.ContactIdType{
		gmodel.TypeUser,
		gmodel.TypeRobot,
	}

	// check: sender type
	if !lo.Contains(allowSenderType, sender.GetType()) {
		err = api.ErrSenderTypeNotAllow
		return
	}
	// check: receiver type
	if !lo.Contains(allowReceiverType, receiver.GetType()) {
		err = api.ErrReceiverTypeNotAllow
		return
	}
	return
}

// FilterFetchUserType 限制拉取消息的用户类型
func (b *MessageFilterUseCase) FilterFetchUserType(owner, peer *gmodel.ComponentId) (err error) {
	allowOwnerType := []gmodel.ContactIdType{
		gmodel.TypeUser,
	}

	allowPeerType := []gmodel.ContactIdType{
		gmodel.TypeUser,
		gmodel.TypeRobot,
		gmodel.TypeGroup,
	}

	// check: owner type
	if !lo.Contains(allowOwnerType, owner.GetType()) {
		err = api.ErrSenderTypeNotAllow
		return
	}
	// check: peer type
	if !lo.Contains(allowPeerType, peer.GetType()) {
		err = api.ErrReceiverTypeNotAllow
		return
	}
	return
}
