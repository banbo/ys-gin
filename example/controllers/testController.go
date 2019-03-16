package controllers

import (
	"github.com/banbo/ys-gin/example/models"

	"github.com/banbo/ys-gin/controller"
	"github.com/banbo/ys-gin/middleware"
	"github.com/banbo/ys-gin/model"
	"github.com/gin-gonic/gin"
)

type TestController struct {
	controller.Controller
}

func (t *TestController) Router(e *gin.Engine) {
	//1.不需要验证参数一致性的接口
	group := e.Group("/test")
	{
		group.GET("/list", t.List)
		group.GET("/get", t.Get)
	}

	//2.需要验证参数一致性的接口
	group2 := *group
	group2.Use(middleware.CheckParamUnanimous())
	{
		group2.POST("/add", t.Add)
		group2.POST("/update", t.Update)
		group2.POST("/delete", t.Delete)
	}
}

//列表
func (t *TestController) List(ctx *gin.Context) {
	//获取参数
	isPage, err := t.GetBool(ctx, "is_page")
	if err != nil {
		t.RespErr(ctx, nil, "参数is_page格式错误")
		return
	}
	pageIndex, err := t.GetInt(ctx, "page_index")
	if isPage && (err != nil || pageIndex <= 0) {
		t.RespErr(ctx, nil, "参数page_index格式错误")
		return
	}
	pageSize, err := t.GetInt(ctx, "page_size")
	if isPage && (err != nil || pageSize <= 0) {
		t.RespErr(ctx, nil, "参数page_size格式错误")
		return
	}

	//筛选参数
	filter := make(map[string]interface{})
	if v, ok := t.GetParam(ctx, "name"); ok && v != "" {
		filter["name"] = v
	}

	//调用model
	var modelList *model.ModelList
	var list []*models.TestModel
	if isPage { //分页
		modelList, list, err = new(models.TestModel).List(pageIndex, pageSize, filter, "")
	} else { //不分页
		modelList, list, err = new(models.TestModel).ListAll(filter, "")
	}
	if err != nil {
		t.RespErr(ctx, nil, err)
		return
	}
	modelList.Items = list

	//返回
	t.RespOK(ctx, modelList)
	return
}

//获取
func (t *TestController) Get(ctx *gin.Context) {
	//获取参数
	uid := t.GetString(ctx, "uid")
	if len(uid) == 0 {
		t.RespErr(ctx, nil, "参数uid格式错误")
		return
	}

	//调用model
	has, testModel, err := new(models.TestModel).Get(uid)
	if err != nil {
		t.RespErr(ctx, nil, err)
		return
	}
	if !has {
		t.RespErr(ctx, nil, "用户不存在")
		return
	}

	//设置返回值
	t.Put(ctx, "user", testModel)

	t.RespOK(ctx, testModel)
	return
}

//新增
func (t *TestController) Add(ctx *gin.Context) {
	//获取参数
	name := t.GetString(ctx, "name")
	if len(name) == 0 {
		t.RespErr(ctx, nil, "参数name格式错误")
		return
	}
	age, err := t.GetInt(ctx, "age")
	if err != nil || age <= 0 {
		t.RespErr(ctx, nil, "参数age格式错误")
		return
	}

	//组织参数
	data := new(models.TestModel)
	data.Name = name
	data.Age = age

	//调用model
	uid, err := new(models.TestModel).Add(data)
	if err != nil {
		t.RespErr(ctx, nil, err)
		return
	}

	t.RespOK(ctx, uid)
	return
}

//更新
func (t *TestController) Update(ctx *gin.Context) {
	//获取参数
	uid := t.GetString(ctx, "uid")
	if len(uid) == 0 {
		t.RespErr(ctx, nil, "参数id格式错误")
		return
	}

	//更新的字段
	params := make(map[string]interface{})
	if v, ok := t.GetParam(ctx, "name"); ok {
		params["name"] = v
	}
	if v, ok := t.GetParam(ctx, "age"); ok {
		params["age"] = v
	}

	//调用model
	err := new(models.TestModel).Update(uid, params)
	if err != nil {
		t.RespErr(ctx, nil, err)
		return
	}

	t.RespOK(ctx, nil)
	return
}

//删除
func (t *TestController) Delete(ctx *gin.Context) {
	//获取参数
	uid := t.GetString(ctx, "uid")
	if len(uid) == 0 {
		t.RespErr(ctx, nil, "参数uid格式错误")
		return
	}

	//调用model
	err := new(models.TestModel).Delete(uid)
	if err != nil {
		t.RespErr(ctx, nil, err)
		return
	}

	t.RespOK(ctx, nil)
	return
}
