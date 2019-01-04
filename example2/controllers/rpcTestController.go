package controllers

import (
	"github.com/banbo/ys-gin/example2/proto"
	"github.com/banbo/ys-gin/example2/rpc/client"

	"github.com/banbo/ys-gin/constant"
	"github.com/banbo/ys-gin/controller"

	"github.com/gin-gonic/gin"
	"github.com/banbo/ys-gin/errors"
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
		r.RespErr(ctx, "参数greeting格式错误")
		return
	}

	// 调用rpc
	rsp, err := client.HelloClient.SayHello(&proto.HelloRequest{Greeting: greeting})
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	//判断返回状态
	if rsp.Code != constant.RESPONSE_CODE_OK {
		r.RespErr(ctx, errors.NewNormal(rsp.Code, rsp.Msg))
		return
	}

	r.Put(ctx, "replay", rsp.Data.Replay)

	r.RespOK(ctx)
	return
}

func (r *RpcTestController) Bye(ctx *gin.Context) {
	// 获取参数
	bye := r.GetString(ctx, "bye")
	if len(bye) == 0 {
		r.RespErr(ctx, "参数bye格式错误")
		return
	}

	// 调用rpc
	rsp, err := client.HelloClient.SayBye(&proto.ByeRequest{Bye: bye})
	if err != nil {
		r.RespErr(ctx, err)
		return
	}

	//判断返回状态
	if rsp.Code != constant.RESPONSE_CODE_OK {
		r.RespErr(ctx, errors.NewNormal(rsp.Code, rsp.Msg))
		return
	}

	r.Put(ctx, "replay", rsp.Data.Replay)

	r.RespOK(ctx)
	return
}
