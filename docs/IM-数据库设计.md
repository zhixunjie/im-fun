# 数据库设计

## 会话表

~~~sql

~~~

## 消息表

~~~sql

~~~

- msg_type：[私信类型](#私信类型)。

# 私信类型

> seeklove.common.api/user/message.go **存放各种私信发送格式。**

不同的消息类型，对应content字段里面的不同JSON结构。

~~~go
MSG_TYPE_TEXT          = 1  // 文字
MSG_TYPE_IMAGE         = 2  // 图片
MSG_TYPE_BUSINESS_CARD = 3  // 名片
MSG_TYPE_SOUND         = 4  // 语音
MSG_TYPE_STICKER       = 5  // 表情
MSG_TYPE_GIFT          = 6  // 礼物
MSG_TYPE_CARD          = 7  // 卡片
MSG_TYPE_ANNOUNCE      = 11 // 公告 (一般由系统发)
MSG_TYPE_CALL          = 12 // 通话
MSG_TYPE_RED_PACKET    = 13 // 红包
MSG_TYPE_TIPS          = 14 // 提示消息
MSG_TYPE_NEWGIFT       = 15 // 新版礼物
MSG_TYPE_VIDEOCALL     = 16 // 视频通话
MSG_TYPE_MGIFT         = 17 // 搭讪礼物
MSG_TYPE_CG_REDPACK    = 18 // 群聊红包
MSG_TYPE_CGREDPACK_TIP = 19 // 群聊红包提示
MSG_TYPE_LUCKGIFT      = 20 // 幸运礼物
MSG_TYPE_PUSHJUMP      = 21 // 推送跳转
MSG_TYPE_PUSHJUMP_PIC  = 22 // 推送跳转(带图片)
MSG_TYPE_LIGHT_GREET   = 23 // 明搭讪
MSG_TYPE_SOOGIF        = 24 // 表情包
MSG_TYPE_PERSON_CARD   = 26 // 个人资料卡
MSG_TYPE_UNIFIED_JUMP  = 27 // 高亮统跳消息
MSG_TYPE_GUARD_CARD    = 28 // 活动卡片
MSG_TYPE_MALE_GREET    = 29 // 男性搭讪消息
~~~

