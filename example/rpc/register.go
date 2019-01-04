package rpc

import (
	"github.com/banbo/ys-gin/example/proto"
	"github.com/banbo/ys-gin/example/rpc/server"

	"google.golang.org/grpc"
)

func Register(rpcSrv *grpc.Server) {
	proto.RegisterHelloServer(rpcSrv, &server.HelloServer{})
}
