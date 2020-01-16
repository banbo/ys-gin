package client

import (
	"context"

	"github.com/banbo/ys-gin/errors"

	"github.com/banbo/ys-gin/example2/proto"
)

var HelloClient helloClient

type helloClient struct {
}

func (*helloClient) SayHello(request *proto.HelloRequest) (*proto.HelloResponse, error) {
	conn, err := NewExampleConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	helloClient := proto.NewHelloClient(conn)
	rsp, err := helloClient.SayHello(context.Background(), request)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return rsp, nil
}

func (*helloClient) SayBye(request *proto.ByeRequest) (*proto.ByeResponse, error) {
	conn, err := NewExampleConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	helloClient := proto.NewHelloClient(conn)
	rsp, err := helloClient.SayBye(context.Background(), request)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return rsp, nil
}
