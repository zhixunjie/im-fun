package cache

import (
	"github.com/zhixunjie/im-fun/internal/logic/data"
)

type ChatMessageCache struct {
	messageRepo *data.MessageRepo
}

func NewChatMessageCache(messageRepo *data.MessageRepo) *ChatMessageCache {
	return &ChatMessageCache{messageRepo: messageRepo}
}

//func (b *ChatMessageCache) SetBitContact(ctx context.Context, owner, peer *gmodel.ComponentId) (err error) {
//	mem := b.messageRepo.RedisClient
//
//	key := data.SessionContact.Format(k.M{
//		"owner": fmt.Sprintf("%v:%v", owner.GetId(), owner.GetType()),
//		"peer":  fmt.Sprintf("%v:%v", peer.GetId(), peer.GetType()),
//	})
//
//	_, err = mem.SetBit(ctx, key, , 1).Result()
//	if err != nil {
//		return err
//	}
//
//	return
//}
