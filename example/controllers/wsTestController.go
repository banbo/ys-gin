package controllers

import (
	"time"

	"github.com/banbo/ys-gin/controller"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type WsTestController struct {
	controller.WsController
}

func (w *WsTestController) Router(e *gin.Engine) {
	//初始化melody服务
	w.Melody = melody.New()

	//路由
	e.GET("/ws_test/:uuid", func(ctx *gin.Context) { //每个连接uuid不同
		w.HandleRequest(ctx.Writer, ctx.Request)
	})

	//处理新连接
	w.HandleConnect(func(s *melody.Session) {
		s.Write(w.WsRespOK("初始化数据"))
	})

	//处理客户端消息
	w.HandleMessage(func(s *melody.Session, msg []byte) {
		w.BroadcastFilter(w.WsRespOK("按路径筛选广播，路径："+s.Request.URL.Path), func(q *melody.Session) bool {
			return s.Request.URL.Path == q.Request.URL.Path
		})
	})

	//定时广播
	go func() {
		for {
			w.Broadcast(w.WsRespOK("定时全体广播"))
			time.Sleep(5 * time.Second)
		}
	}()
}
