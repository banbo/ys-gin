package rpcClient

import (
	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/errors"

	"google.golang.org/grpc"
)

func NewOrderNoConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(conf.Configer.BeeConfiger.String("rpc_client::order_no_svr"), grpc.WithInsecure())
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return conn, nil
}
