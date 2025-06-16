package biz

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/api"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
)

type MessageFilterUseCase struct {
}

func NewMessageFilterUseCase() *MessageFilterUseCase {
	return &MessageFilterUseCase{}
}

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
