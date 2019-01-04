package models

import (
	"strconv"

	"github.com/banbo/ys-gin/example/constant"

	"github.com/banbo/ys-gin/errors"
	"github.com/banbo/ys-gin/id"
	"github.com/banbo/ys-gin/log"
	"github.com/banbo/ys-gin/model"
)

type TestModel struct {
	model.Model `xorm:"-"`
	Uid         string `xorm:"uid pk" json:"uid"`
	Name        string `xorm:"name" json:"name"`
	Age         int    `xorm:"age" json:"age"`
}

//表名
func (TestModel) TableName() string {
	return "test"
}

//列表，分页
func (t *TestModel) List(pageIndex int, pageSize int, filter map[string]string) (*model.ModelList, []*TestModel, error) {
	session := model.Orm.Where("1=1")

	//筛选
	if v, ok := filter["name"]; ok {
		session.And("name = ?", v)
	}

	//获取分页
	sessionCp := session.Clone()
	total, err := sessionCp.Count(new(TestModel))
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	limit, offset, modelList := t.Paging(pageIndex, pageSize, int(total))

	//获取列表
	var list []*TestModel
	err = session.Limit(limit, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return modelList, list, nil
}

//列表，不分页
func (t *TestModel) ListAll(filter map[string]string) (*model.ModelList, []*TestModel, error) {
	session := model.Orm.Where("1=1")

	//筛选
	if v, ok := filter["name"]; ok {
		session.Where("name = ?", v)
	}

	var list []*TestModel
	err := session.Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return t.NoPaging(len(list), list), list, nil
}

//获取
func (t *TestModel) Get(uid string) (bool, *TestModel, error) {
	testModel := new(TestModel)

	has, err := model.Orm.Where("uid=?", uid).Get(testModel)
	if err != nil {
		return false, nil, errors.NewSys(err)
	}

	return has, testModel, nil
}

//新增
func (t *TestModel) Add(testModel *TestModel) (string, error) {
	//生成uid
	testModel.Uid = strconv.FormatInt(id.IdWorker.Generate(), 10)

	_, err := model.Orm.Insert(testModel)
	if err != nil {
		return "", errors.NewSys(err)
	}

	return testModel.Uid, nil
}

//更新
func (t *TestModel) Update(uid string, params map[string]interface{}) error {
	//判断是否存在
	has, _, err := t.Get(uid)
	if err != nil {
		return err
	}
	if !has {
		log.Logger.Error("用户不存在或已删除，uid：", uid)
		return errors.NewNormal(constant.RESPONSE_CODE_NO_USER, "用户不存在或已删除")
	}

	//设置更新字段
	data := make(map[string]interface{})
	if v, ok := params["name"]; ok {
		data["name"] = v
	}
	if v, ok := params["age"]; ok {
		data["age"] = v
	}

	_, err = model.Orm.Table(t).Update(data)
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

//删除
func (t *TestModel) Delete(uid string) error {
	_, err := model.Orm.Id(uid).Delete(t)
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
