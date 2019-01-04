package router

import (
	"github.com/banbo/ys-gin/example2/controllers"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) {
	new(controllers.RpcTestController).Router(engine)
}
