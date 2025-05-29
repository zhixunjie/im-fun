package format

import (
	"encoding/json"
	"errors"
	"fmt"
)

func DecodeMsgBody(base *MsgBody) (msgContent MsgContent, err error) {
	if base == nil {
		err = errors.New("msg body is invalid")
		return
	}
	// 创建：newInstance
	msgContent = newInstance(base.MsgType)
	if msgContent == nil {
		err = fmt.Errorf("unknown type: %v", base.MsgType)
		return
	}

	buf, err := json.Marshal(base.MsgContent)
	if err != nil {
		err = fmt.Errorf("marshal msg content: %w", err)
		return
	}

	if err = json.Unmarshal(buf, msgContent); err != nil {
		err = fmt.Errorf("unmarshal json(data): %w", err)
		return
	}

	return
}

// newInstance 根据类型创建新实例
func newInstance(t MsgType) MsgContent {
	if factory, exists := JsonDecoderMap[t]; exists {
		return factory()
	}
	return nil
}
