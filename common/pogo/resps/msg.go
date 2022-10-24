package resps

import "tinypro/common/types"

type MyMsg struct {
	SN          int64         `json:"sn,string"`
	MsgType     types.MsgEnum `json:"msg_type"`
	MsgTypeDesc string        `json:"msg_type_desc"`
	MsgContent  string        `json:"msg_content"`
	CreatedAt   string        `json:"created_at"`
}
