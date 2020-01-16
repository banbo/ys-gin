package controllers

import (
	"github.com/banbo/ys-gin/constant"
	"github.com/banbo/ys-gin/controller"
	"github.com/banbo/ys-gin/errors"
	"github.com/gin-gonic/gin"

	"github.com/banbo/ys-gin/example2/proto"
	"github.com/banbo/ys-gin/example2/rpc/client"
)

type RpcTestController struct {
	controller.Controller
}

func (t *RpcTestController) Router(e *gin.Engine) {
	group := e.Group("/rpc_test")
	{
		group.GET("/hello", t.Hello)
		group.GET("/bye", t.Bye)
	}
}

func (r *RpcTestController) Hello(ctx *gin.Context) {
	// 获取参数
	greeting := r.GetString(ctx, "greeting")
	if len(greeting) == 0 {
		r.RespErr(ctx, nil, "参数greeting格式错误")
		return
	}

	// 调用rpc
	rsp, err := client.HelloClient.SayHello(&proto.HelloRequest{Greeting: greeting})
	if err != nil {
		r.RespErr(ctx, nil, err)
		return
	}

	//判断返回状态
	if rsp.Code != constant.RESPONSE_CODE_OK {
		r.RespErr(ctx, nil, errors.NewNormal(rsp.Code, rsp.Msg))
		return
	}

	r.RespOK(ctx, rsp.Data.Replay)
	return
}

func (r *RpcTestController) Bye(ctx *gin.Context) {
	// 获取参数
	bye := r.GetString(ctx, "bye")
	if len(bye) == 0 {
		r.RespErr(ctx, nil, "参数bye格式错误")
		return
	}

	// 调用rpc
	rsp, err := client.HelloClient.SayBye(&proto.ByeRequest{Bye: bye})
	if err != nil {
		r.RespErr(ctx, nil, err)
		return
	}

	//判断返回状态
	if rsp.Code != constant.RESPONSE_CODE_OK {
		r.RespErr(ctx, nil, errors.NewNormal(rsp.Code, rsp.Msg))
		return
	}

	r.RespOK(ctx, rsp.Data.Replay)
	return
}
