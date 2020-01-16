package ys_gin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/banbo/ys-gin/cache"
	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/id"
	"github.com/banbo/ys-gin/log"
	"github.com/banbo/ys-gin/model"
)

type App struct {
	GinEngine *gin.Engine
	apiSvr    *http.Server
	RpcSvr    *grpc.Server
	rpcAddr   string
}

func NewApp(configFile string) *App {
	//初始化组件
	conf.NewConfiger(configFile)
	log.NewLogger()

	if conf.Configer.BeeConfiger.String("system::worker_id") != "" {
		id.NewIdWorker(conf.Configer.ApiConf.WorkerID)
	}
	if conf.Configer.BeeConfiger.String("redis::host") != "" {
		cache.NewRedisClient()
	}

	//设置gin运行模式
	gin.SetMode(conf.Configer.ApiConf.RunMode)

	//初始化数据库连接池
	_, err := model.NewEngine()
	if err != nil {
		panic(fmt.Sprintf("初始化连接池出错，%v", err))
	}

	/*switch conf.Configer.DbConf.DriverName {
	case "mysql":
		model.NewMysqlOrm()
	case "postgres":
		model.NewPostgresOrm()
	}*/

	app := &App{}

	//初始化api服务
	if len(conf.Configer.ApiConf.HttpPort) > 0 {
		app.GinEngine = gin.Default()
		app.apiSvr = &http.Server{
			Addr:    ":" + conf.Configer.ApiConf.HttpPort,
			Handler: app.GinEngine,
		}

	}

	//初始化rpc服务
	if len(conf.Configer.ApiConf.RpcPort) > 0 {
		app.rpcAddr = ":" + conf.Configer.ApiConf.RpcPort
		app.RpcSvr = grpc.NewServer()
	}

	return app
}

func (app *App) Run() {
	//启动api服务
	if app.apiSvr != nil {
		go func() {
			err := app.apiSvr.ListenAndServe()
			if err != nil {
				panic(fmt.Sprintf("启动http服务失败，%v", err))
			}
		}()
		fmt.Println("Api Svr Listen and serve on", app.apiSvr.Addr)
	}

	//启动rpc服务
	if app.RpcSvr != nil {
		listen, err := net.Listen("tcp", app.rpcAddr)
		if err != nil {
			panic("监听rpc端口失败，err: " + err.Error())
		}

		go func() {
			err := app.RpcSvr.Serve(listen)
			if err != nil {
				panic(fmt.Sprintf("启动rpc服务失败，%v", err))
			}
		}()
		fmt.Println("Rpc Svr Listen and serve on", app.rpcAddr)
	}

	//监听退出
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-quitChan

	//优雅退出
	if app.apiSvr != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		app.apiSvr.Shutdown(ctx)
	}
	if app.RpcSvr != nil {
		app.RpcSvr.GracefulStop()
	}
}
