package main

import (
	ysGin "github.com/banbo/ys-gin"

	"github.com/banbo/ys-gin/example/router"
)

func main() {
	//启动服务
	app := ysGin.NewApp("./example.yaml")

	//设置路由
	router.Init(app.GinEngine)

	app.Run()
}
