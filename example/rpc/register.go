package rpc

import (
	"google.golang.org/grpc"

	"github.com/banbo/ys-gin/example/proto"
	"github.com/banbo/ys-gin/example/rpc/server"
)

func Register(rpcSrv *grpc.Server) {
	proto.RegisterHelloServer(rpcSrv, &server.HelloServer{})
}
