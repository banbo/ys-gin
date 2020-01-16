package controllers

import (
	"github.com/banbo/ys-gin/controller"
	"github.com/banbo/ys-gin/middleware"
	"github.com/banbo/ys-gin/model"
	"github.com/gin-gonic/gin"

	"github.com/banbo/ys-gin/example/models"
)

type UserController struct {
	controller.Controller
}

func (u *UserController) Router(e *gin.Engine) {
	//1.不需要验证参数一致性的接口
	group := e.Group("/user")
	{
		group.GET("/list", u.List)
		group.GET("/get", u.Get)
	}

	//2.需要验证参数一致性的接口
	group2 := *group
	group2.Use(middleware.CheckParamUnanimous())
	{
		group2.POST("/add", u.Add)
		group2.POST("/update", u.Update)
		group2.POST("/delete", u.Delete)
	}
}

//列表
func (u *UserController) List(ctx *gin.Context) {
	//获取参数
	isPage, err := u.GetBool(ctx, "is_page")
	if err != nil {
		u.RespErr(ctx, nil, "参数is_page格式错误")
		return
	}
	pageIndex, err := u.GetInt(ctx, "page_index")
	if isPage && (err != nil || pageIndex <= 0) {
		u.RespErr(ctx, nil, "参数page_index格式错误")
		return
	}
	pageSize, err := u.GetInt(ctx, "page_size")
	if isPage && (err != nil || pageSize <= 0) {
		u.RespErr(ctx, nil, "参数page_size格式错误")
		return
	}

	//筛选参数
	filter := make(map[string]interface{})
	if v, ok := u.GetParam(ctx, "name"); ok && v != "" {
		filter["name"] = v
	}

	//调用model
	var modelList *model.ModelList
	var list []*models.UserModel
	if isPage { //分页
		modelList, list, err = new(models.UserModel).List(pageIndex, pageSize, filter, "")
	} else { //不分页
		modelList, list, err = new(models.UserModel).ListAll(filter, "")
	}
	if err != nil {
		u.RespErr(ctx, nil, err)
		return
	}
	modelList.Items = list

	//返回
	u.RespOK(ctx, modelList)
	return
}

//获取
func (u *UserController) Get(ctx *gin.Context) {
	//获取参数
	uid := u.GetString(ctx, "uid")
	if len(uid) == 0 {
		u.RespErr(ctx, nil, "参数uid格式错误")
		return
	}

	//调用model
	has, testModel, err := new(models.UserModel).Get(uid)
	if err != nil {
		u.RespErr(ctx, nil, err)
		return
	}
	if !has {
		u.RespErr(ctx, nil, "用户不存在")
		return
	}

	//设置返回值
	u.Put(ctx, "user", testModel)

	u.RespOK(ctx, testModel)
	return
}

//新增
func (u *UserController) Add(ctx *gin.Context) {
	//获取参数
	name := u.GetString(ctx, "name")
	if len(name) == 0 {
		u.RespErr(ctx, nil, "参数name格式错误")
		return
	}
	age, err := u.GetInt(ctx, "age")
	if err != nil || age <= 0 {
		u.RespErr(ctx, nil, "参数age格式错误")
		return
	}

	//组织参数
	data := new(models.UserModel)
	data.Name = name
	data.Age = age

	//调用model
	uid, err := new(models.UserModel).Add(data)
	if err != nil {
		u.RespErr(ctx, nil, err)
		return
	}

	u.RespOK(ctx, uid)
	return
}

//更新
func (u *UserController) Update(ctx *gin.Context) {
	//获取参数
	uid := u.GetString(ctx, "uid")
	if len(uid) == 0 {
		u.RespErr(ctx, nil, "参数id格式错误")
		return
	}

	//更新的字段
	params := make(map[string]interface{})
	if v, ok := u.GetParam(ctx, "name"); ok {
		params["name"] = v
	}
	if v, ok := u.GetParam(ctx, "age"); ok {
		params["age"] = v
	}

	//调用model
	err := new(models.UserModel).Update(uid, params)
	if err != nil {
		u.RespErr(ctx, nil, err)
		return
	}

	u.RespOK(ctx, nil)
	return
}

//删除
func (u *UserController) Delete(ctx *gin.Context) {
	//获取参数
	uid := u.GetString(ctx, "uid")
	if len(uid) == 0 {
		u.RespErr(ctx, nil, "参数uid格式错误")
		return
	}

	//调用model
	err := new(models.UserModel).Delete(uid)
	if err != nil {
		u.RespErr(ctx, nil, err)
		return
	}

	u.RespOK(ctx, nil)
	return
}
