package controller

import (
	"net/http"
	"strconv"

	"github.com/banbo/ys-gin/constant"
	"github.com/banbo/ys-gin/errors"

	"github.com/gin-gonic/gin"
)

const (
	SAVE_DATA_KEY = "save_api_data_key"
)

// 返回的结构
type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// controller
type Controller struct {
}

// 设置返回的数据，key-value
// 使用gin context的Keys保存
// gin context每个请求都会先reset
func (c *Controller) Put(ctx *gin.Context, key string, value interface{}) {
	// lazy init
	if ctx.Keys == nil {
		ctx.Keys = make(map[string]interface{})
	}
	if ctx.Keys[SAVE_DATA_KEY] == nil {
		ctx.Keys[SAVE_DATA_KEY] = make(map[string]interface{})
	}

	ctx.Keys[SAVE_DATA_KEY].(map[string]interface{})[key] = value
}

// 正确的响应
func (c *Controller) RespOK(ctx *gin.Context, data interface{}) {
	resp := &Response{
		Status: "success",
		Code:   constant.RESPONSE_CODE_OK,
		Msg:    "成功",
		Data:   data,
	}

	ctx.JSON(http.StatusOK, resp)
}

// 错误的响应
func (c *Controller) RespErr(ctx *gin.Context, data interface{}, options ...interface{}) {
	resp := &Response{
		Status: "fail",
		Code:   constant.RESPONSE_CODE_ERROR, // 默认是常规错误
		Msg:    "",
		Data:   data,
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

	ctx.JSON(http.StatusOK, resp)
}

// 获取get、post提交的参数
// 参数存在时第二个参数返回true，即使参数的值为空字符串
// 参数不存在时第二个参数返回false
func (c *Controller) GetParam(ctx *gin.Context, key string) (string, bool) {
	var param string
	var ok bool

	switch ctx.Request.Method {
	case "POST":
		fallthrough
	case "PUT":
		param, ok = ctx.GetPostForm(key)
	default:
		param, ok = ctx.GetQuery(key)
	}

	return param, ok
}

// 获取get、post提交的string类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetString(ctx *gin.Context, key string, def ...string) string {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0]
	}

	return param
}

// 获取get、post提交的int类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetInt(ctx *gin.Context, key string, def ...int) (int, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.Atoi(param)
}

// 获取get、post提交的int64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetInt64(ctx *gin.Context, key string, def ...int64) (int64, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseInt(param, 10, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetFloat64(ctx *gin.Context, key string, def ...float64) (float64, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseFloat(param, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetBool(ctx *gin.Context, key string, def ...bool) (bool, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseBool(param)
}
