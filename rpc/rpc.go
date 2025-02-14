package rpc

import "chat_game/utils/common"

const (
	MsgServiceName = "MsgService"
)

type MsgResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MsgServer interface {
	SendMessage(req common.Msg, res *MsgResp) error
}
