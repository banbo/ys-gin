package router

import (
	"github.com/banbo/ys-gin/middleware"
	"github.com/gin-gonic/gin"

	"github.com/banbo/ys-gin/example2/controllers"
)

func Init(engine *gin.Engine) {
	//设置跨域
	middleware.SetCorsOrigin([]string{"http://127.0.0.1"})
	engine.Use(middleware.Cors())

	new(controllers.RpcTestController).Router(engine)
}
