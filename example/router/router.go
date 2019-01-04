package router

import (
	"github.com/banbo/ys-gin/example/controllers"

	"github.com/gin-gonic/gin"
)

func Init(engine *gin.Engine) {
	new(controllers.TestController).Router(engine)
}
