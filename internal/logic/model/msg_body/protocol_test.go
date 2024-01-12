package msg_body

import (
	"encoding/json"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"testing"
)

// 文本消息
func TestText(t *testing.T) {
	body := MsgBody{
		MsgType: MsgTypeText,
		MsgContent: &MsgContent{
			TextContent: &TextContent{
				Text: "哈哈哈",
			},
		},
	}
	buf, err := json.Marshal(&body)
	if err != nil {
		return
	}
	utils.PrettyPrint(buf)
}

// 图片消息
func TestImage(t *testing.T) {
	body := MsgBody{
		MsgType: MsgTypeImage,
		MsgContent: &MsgContent{
			ImageContent: &ImageContent{
				List: []Image{
					{
						Url:    "https://1.png",
						Width:  11,
						Height: 22,
						Size:   33,
					},
					{
						Url:    "https://2.png",
						Width:  11,
						Height: 22,
						Size:   33,
					},
				},
			},
		},
	}
	buf, err := json.Marshal(&body)
	if err != nil {
		return
	}
	utils.PrettyPrint(buf)
}

// 音频消息
func TestAudio(t *testing.T) {
	body := MsgBody{
		MsgType: MsgTypeAudio,
		MsgContent: &MsgContent{
			AudioContent: &AudioContent{
				Url:      "https://xxxx.mp3",
				Duration: 1,
				Text:     "我是音频",
			},
		},
	}
	buf, err := json.Marshal(&body)
	if err != nil {
		return
	}
	utils.PrettyPrint(buf)
}

// 提示消息
func TestTips(t *testing.T) {
	body := MsgBody{
		MsgType: MsgTypeTips,
		MsgContent: &MsgContent{
			TipsContent: &TipsContent{
				Text:   "提示消息：对方已通过认证",
				ImgUrl: "https://1.png",
			},
		},
	}
	buf, err := json.Marshal(&body)
	if err != nil {
		return
	}
	utils.PrettyPrint(buf)
}

// 发送高亮文字
func TestHighLightText(t *testing.T) {
	text := "尊敬的用户，感谢您的关注，如有疑问请联系在线客服！"
	//highLightArr := []string{"尊敬的用户", "如有疑问"}

	body := MsgBody{
		MsgType: MsgTypeText,
		MsgContent: &MsgContent{
			TextContent: &TextContent{
				Text: text,
				HighLights: []HighLight{
					map[string]any{
						"text":  "尊敬的用户",
						"link":  "https://1111",
						"color": "#0046FF",
						"range": GetRange(text, "尊敬的用户"),
					},
					map[string]any{
						"text":   "如有疑问",
						"link":   "https://222",
						"color":  "#0046FF",
						"offset": GetRange(text, "如有疑问"),
					},
				},
			},
		},
	}
	buf, err := json.Marshal(&body)
	if err != nil {
		return
	}
	utils.PrettyPrint(buf)
}
