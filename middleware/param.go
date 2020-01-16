package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/banbo/ys-gin/conf"
	"github.com/banbo/ys-gin/controller"
)

// 验证参数一致
func CheckParamUnanimous() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取所有参数
		ctx.Request.ParseForm()
		ctx.Request.ParseMultipartForm(32 << 20) // 32M
		params := ctx.Request.Form

		// 获取用户提交的签名
		formSign := ctx.GetHeader("sign")
		if formSign == "" {
			formSign = params.Get("sign")
		}
		if formSign == "" {
			new(controller.Controller).RespErr(ctx, nil, "缺少参数sign")

			ctx.Abort()
			return
		}
		// 签名不参与加密
		params.Del("sign")

		// 参数排序
		paramKeys := make([]string, len(params))

		i := 0
		for k := range params {
			paramKeys[i] = k
			i++
		}

		sort.Strings(paramKeys)

		// 拼成：key1=value2&key2=value2&secret=xxx
		var paramStr string
		for _, v := range paramKeys {
			paramStr += fmt.Sprintf("%s=%s&", v, params.Get(v))
		}
		// 最后拼上secret
		paramStr += fmt.Sprintf("secret=%s", conf.Configer.ApiConf.ParamSecret)

		// 把拼接的参数签名后进行对比
		s := sha1.New()
		s.Write([]byte(paramStr))
		sign := hex.EncodeToString(s.Sum(nil))

		// 打印出签名，方便调试
		if gin.Mode() != gin.ReleaseMode {
			fmt.Println("服务器生成的sign：", sign)
		}

		// 对比服务器生成的签名和用户提交的签名
		if sign != formSign {
			new(controller.Controller).RespErr(ctx, nil, "参数不一致")

			ctx.Abort()
			return
		}

		//如果传了timestamp，检查sign有效期
		timestamp := params.Get("timestamp")
		if timestamp != "" {
			timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				new(controller.Controller).RespErr(ctx, nil, "参数timestamp格式错误")

				ctx.Abort()
				return
			}

			if time.Now().Unix()-timestampInt > 30 { //有效期30分钟
				new(controller.Controller).RespErr(ctx, nil, "sign已失效")

				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
