package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/banbo/ys-gin/controller"
)

var allowOrigins []string

func SetCorsOrigin(origins []string) {
	allowOrigins = origins
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.GetHeader("origin")

		//容许域名
		if len(allowOrigins) == 0 {
			new(controller.Controller).RespErr(ctx, nil, "请调用SetCorsOrigin配置跨域域名")

			ctx.Abort()
			return
		}
		for _, v := range allowOrigins {
			if v == origin {
				ctx.Header("Access-Control-Allow-Origin", v)
				break
			}
		}

		//容许跨域带cookie
		ctx.Header("Access-Control-Allow-Credentials", "true")

		//容许请求header
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, X-Requested-With, Authorization, Token")

		//容许请求方法
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		//容许响应获取的header
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		//放行OPTIONS请求
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}

		ctx.Next()
	}
}
