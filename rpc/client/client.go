package rpc

import (
	"context"
	"net/rpc"

	mrpc "chat_game/rpc"
	"chat_game/utils/common"
)

type MsgServiceClient struct {
	*rpc.Client
}

func NewMsgServiceClient(network, address string) (*MsgServiceClient, error) {
	conn, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &MsgServiceClient{Client: conn}, nil
}

func (m *MsgServiceClient) SendMessage(ctx context.Context, req common.WsHubMsg) (*mrpc.MsgResp, error) {
	res := &mrpc.MsgResp{}
	err := m.Client.Call(mrpc.MsgServiceName+".SendMessage", req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
