package router

import (
	"github.com/banbo/ys-gin/example/controllers"
	"github.com/banbo/ys-gin/middleware"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) {
	//全局跨域
	engine.Use(middleware.Cors())

	new(controllers.TestController).Router(engine)
}
