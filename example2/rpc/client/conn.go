package client

import (
	"github.com/banbo/ys-gin/errors"

	"google.golang.org/grpc"
	"github.com/banbo/ys-gin/conf"
)

func NewExampleConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(conf.Configer.BeeConfiger.String("rpc_client::example_svr"), grpc.WithInsecure())
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return conn, nil
}
