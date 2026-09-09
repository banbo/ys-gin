package models

import (
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/banbo/ys-gin/errors"
	"github.com/banbo/ys-gin/id"
	"github.com/banbo/ys-gin/log"
	"github.com/banbo/ys-gin/model"

	"github.com/banbo/ys-gin/example/constants"
)

type UserModel struct {
	model.Model `gorm:"-"`
	Uid         string `gorm:"column:uid;primary_key" json:"uid"`
	Name        string `gorm:"column:name" json:"name"`
	Age         int    `gorm:"column:age" json:"age"`
}

// 表名
func (UserModel) TableName() string {
	return "users"
}

// 列表，分页
func (u *UserModel) List(pageIndex int, pageSize int, filter map[string]interface{}, orderBy string) (*model.ModelList, []*UserModel, error) {
	db, err := model.Engineer.Get(u)
	if err != nil {
		return nil, nil, err
	}

	session := db.Where("1=1")

	//筛选
	if v, ok := filter["name"]; ok {
		session = session.Where("name = ?", v)
	}

	//排序
	if orderBy != "" {
		session = session.Order(orderBy)
	} else {
		session = session.Order("uid DESC")
	}

	//获取分页
	var total int64
	err = session.Model(&UserModel{}).Count(&total).Error
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	limit, offset, modelList := u.Paging(pageIndex, pageSize, int(total))

	//获取列表
	var list []*UserModel
	err = session.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return modelList, list, nil
}

// 列表，不分页
func (u *UserModel) ListAll(filter map[string]interface{}, orderBy string) (*model.ModelList, []*UserModel, error) {
	db, err := model.Engineer.Get(u)
	if err != nil {
		return nil, nil, err
	}

	session := db.Where("1=1")

	//筛选
	if v, ok := filter["name"]; ok {
		session = session.Where("name = ?", v)
	}

	//排序
	if orderBy != "" {
		session = session.Order(orderBy)
	} else {
		session = session.Order("uid DESC")
	}

	var list []*UserModel
	err = session.Find(&list).Error
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return u.NoPaging(len(list), list), list, nil
}

// 获取
func (u *UserModel) Get(uid string) (bool, *UserModel, error) {
	db, err := model.Engineer.Get(u)
	if err != nil {
		return false, nil, err
	}

	testModel := new(UserModel)

	err = db.Where("uid=?", uid).First(testModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil
		}
		return false, nil, errors.NewSys(err)
	}

	return true, testModel, nil
}

// 新增
func (u *UserModel) Add(testModel *UserModel) (string, error) {
	db, err := model.Engineer.Get(u)
	if err != nil {
		return "", err
	}

	//生成uid
	testModel.Uid = strconv.FormatInt(id.IdWorker.Generate(), 10)

	err = db.Create(testModel).Error
	if err != nil {
		return "", errors.NewSys(err)
	}

	return testModel.Uid, nil
}

// 更新
func (u *UserModel) Update(uid string, params map[string]interface{}) error {
	db, err := model.Engineer.Get(u)
	if err != nil {
		return err
	}

	//判断是否存在
	has, _, err := u.Get(uid)
	if err != nil {
		return err
	}
	if !has {
		log.Logger.Error("用户不存在或已删除，uid：", uid)
		return errors.NewNormal(constants.RESPONSE_CODE_NO_USER, "用户不存在或已删除")
	}

	//设置更新字段
	data := make(map[string]interface{})
	if v, ok := params["name"]; ok {
		data["name"] = v
	}
	if v, ok := params["age"]; ok {
		data["age"] = v
	}

	err = db.Model(&UserModel{}).Where("uid = ?", uid).Updates(data).Error
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

// 删除
func (u *UserModel) Delete(uid string) error {
	db, err := model.Engineer.Get(u)
	if err != nil {
		return err
	}

	err = db.Where("uid = ?", uid).Delete(&UserModel{}).Error
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
