package server

import (
	"context"
	"testing"
	"github.com/banbo/ys-gin/example/proto"
)

func TestHelloServer_SayHello(t *testing.T) {
	helloServer := &HelloServer{}
	rsp, err := helloServer.SayHello(context.Background(), &proto.HelloRequest{Greeting: "test_man"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp)
}

func TestHelloServer_SayBye(t *testing.T) {
	helloServer := &HelloServer{}
	rsp, err := helloServer.SayBye(context.Background(), &proto.ByeRequest{Bye: "test_man"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp)
}
