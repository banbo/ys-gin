package main

import (
	ysGin "github.com/banbo/ys-gin"

	"github.com/banbo/ys-gin/example2/router"
)

func main() {
	//启动服务
	app := ysGin.NewApp("/Volumes/WorkHD/workspace/go/src/github.com/banbo/ys-gin/example2/example2.conf")

	//设置路由
	router.Init(app.GinEngine)

	app.Run()
}
