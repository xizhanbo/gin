package services

import (
	"errors"
	"fmt"
	"micro-gin/app/common/request"
	"micro-gin/app/models"
	"micro-gin/global"
	"micro-gin/utils"
	"strconv"
)

type testGenStruct struct {
}

var TestGenStruct = new(testGenStruct)

// Register 注册
func (testGenStruct *testGenStruct) Register(params request.AddTestGenStruct) (err error, user models.TestGenStruct) {
	var result = global.App.DB.Where("mobile = ?", params.Mobile).Select("id").First(&models.User{})
	if result.RowsAffected != 0 {
		err = errors.New("手机号已存在")
		return
	}
	user = models.TestGenStruct{Name: params.Name, Mobile: params.Mobile, Password: utils.BcryptMake([]byte(params.Password))}
	err = global.App.DB.Create(&user).Error
	return
}

// GetUserInfo 获取用户信息
func (testGenStruct *testGenStruct) testGenStruct(id string) (err error, user models.TestGenStruct) {
	fmt.Println("----------")
	intId, err := strconv.Atoi(id)
	err = global.App.DB.First(&user, intId).Error
	if err != nil {
		err = errors.New("数据不存在")
	}
	return
}
