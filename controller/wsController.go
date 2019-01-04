package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"github.com/banbo/ys-gin/constant"
	"github.com/banbo/ys-gin/errors"
)

// websocket controller
type WsController struct {
	Controller
	*melody.Melody
}

// 正确的响应
func (c *WsController) WsRespOK(data interface{}) []byte {
	resp := &Response{
		Code: constant.RESPONSE_CODE_OK,
		Msg:  "成功",
		Data: data,
	}

	respByte, _ := json.Marshal(resp)

	return respByte
}

// 错误的响应
func (c *WsController) WsRespErr(data interface{}, options ...interface{}) []byte {
	resp := &Response{
		Code: constant.RESPONSE_CODE_ERROR, // 默认是常规错误
		Msg:  "",
		Data: data,
	}

	// 继续确定code、msg
	for _, v := range options {
		switch opt := v.(type) {
		case int:
			resp.Code = opt // 当前指定code
		case string:
			resp.Msg = opt
		case errors.SysErrorInterface: // 系统错误
			resp.Code = opt.Status()

			if gin.Mode() == gin.ReleaseMode { // 生产环境不显示错误细节
				resp.Msg = opt.Error()
			} else { // 开发环境显示错误细节
				resp.Msg = opt.String()
			}
		case errors.NormalErrorInterface: // 常规错误
			if opt.Status() != 0 { // 常规错误指定了code并且不为0
				resp.Code = opt.Status()
			}
			resp.Msg = opt.Error()
		case error: // go错误
			resp.Msg = opt.Error()
		}
	}

	// 优先使用系统指定msg
	sysMsg := constant.GetResponseMsg(resp.Code)
	if len(sysMsg) > 0 {
		resp.Msg = sysMsg
	}

	respByte, _ := json.Marshal(resp)

	return respByte
}
