package server

import (
	"context"

	"github.com/banbo/ys-gin/rpc"

	"github.com/banbo/ys-gin/example/proto"
)

type HelloServer struct {
	rpc.Server
}

func (h *HelloServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	rsp := &proto.HelloResponse{}

	// 验证参数
	if len(req.Greeting) == 0 {
		h.SetErr(rsp, "参数greeting格式错误")
		return rsp, nil
	}

	// 业务逻辑
	rsp.Data = &proto.HelloResponse_Data{
		Replay: "Hello, " + req.Greeting,
	}

	return rsp, nil
}

func (h *HelloServer) SayBye(ctx context.Context, req *proto.ByeRequest) (*proto.ByeResponse, error) {
	rsp := &proto.ByeResponse{}

	// 验证参数
	if len(req.Bye) == 0 {
		h.SetErr(rsp, "参数bye格式错误")
		return rsp, nil
	}

	// 业务逻辑
	rsp.Data = &proto.ByeResponse_Data{
		Replay: "Bye, " + req.Bye,
	}

	return rsp, nil
}
