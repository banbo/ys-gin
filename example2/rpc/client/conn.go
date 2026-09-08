package client

import (
	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewExampleConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(conf.Configer.BeeConfiger.String("rpc_client::example_svr"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return conn, nil
}
