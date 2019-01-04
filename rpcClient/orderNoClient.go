package rpcClient

import (
	"context"

	"github.com/banbo/ys-gin/errors"
	"github.com/banbo/ys-gin/proto_datasvr"
)

var OrderNoClient orderNoClient

type orderNoClient struct {
}

func (*orderNoClient) Gen(request *proto_datasvr.GenOrderNoRequest) (*proto_datasvr.GenOrderNoResponse, error) {
	conn, err := NewOrderNoConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	orderNoClient := proto_datasvr.NewOrderNoClient(conn)
	rsp, err := orderNoClient.Gen(context.Background(), request)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return rsp, nil
}
