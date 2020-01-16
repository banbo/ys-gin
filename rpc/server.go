package rpc

import (
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/banbo/ys-gin/constant"
	"github.com/banbo/ys-gin/errors"
)

type Server struct {
}

// 设置code、msg字段
func (*Server) SetErr(rsp interface{}, options ...interface{}) {
	// 获取code、msg
	code := constant.RESPONSE_CODE_ERROR // 默认是常规错误
	msg := ""

	// 继续确定code、msg
	for _, v := range options {
		switch opt := v.(type) {
		case int:
			code = opt // 当前指定code
		case string:
			msg = opt
		case errors.SysErrorInterface: // 系统错误
			code = opt.Status()

			if gin.Mode() == gin.ReleaseMode { // 生产环境不显示错误细节
				msg = opt.Error()
			} else { // 开发环境显示错误细节
				msg = opt.String()
			}
		case errors.NormalErrorInterface: // 常规错误
			if opt.Status() != 0 { // 常规错误指定了code并且不为0
				code = opt.Status()
			}
			msg = opt.Error()
		case error: // go错误
			msg = opt.Error()
		}
	}

	// 优先使用系统指定msg
	sysMsg := constant.GetResponseMsg(code)
	if len(sysMsg) > 0 {
		msg = sysMsg
	}

	// 获取rsp结构
	v := reflect.ValueOf(rsp)
	if !(v.Kind() == reflect.Ptr && v.Elem().CanSet()) {
		panic("rsp需要传递地址")
	}
	v = v.Elem()

	rspCode := v.FieldByName("Code")
	if !rspCode.IsValid() {
		panic("rsp没有code字段")
	}
	if rspCode.Kind() != reflect.Int32 {
		panic("code字段类型必须为int2")
	}

	rspMsg := v.FieldByName("Msg")
	if !rspMsg.IsValid() {
		panic("rsp没有msg字段")
	}
	if rspMsg.Kind() != reflect.String {
		panic("msg字段类型必须为string")
	}

	// 设置
	rspCode.SetInt(int64(code))
	rspMsg.SetString(msg)
}
