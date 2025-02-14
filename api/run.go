package api

import (
	"context"
	"net"
	"net/rpc"

	"chat_game/config"
	"chat_game/log"
	mrpc "chat_game/rpc"
)

func ServerRpc(ctx context.Context) {
	log.Info(ctx, "rpc server start")

	rpc.RegisterName(mrpc.MsgServiceName, hub)

	appConfig := config.GetAppConfig()

	listener, err := net.Listen(appConfig.Rpc.Network, appConfig.Rpc.Addr)
	if err != nil {
		panic(err)
	}

	hub.SetRpcAddr(listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go rpc.ServeConn(conn)
	}
}
