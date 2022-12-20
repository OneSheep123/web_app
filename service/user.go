package service

import (
	"errors"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"
)

var (
	ErrorUserLogin = errors.New("用户已经登录")
)

func SignUp(up *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(up.Username); err != nil {
		return
	}
	// 2. 生成UID
	userId := snowflake.GenID()
	u := &models.User{
		UserID:   userId,
		UserName: up.Username,
		PassWord: up.Password,
	}
	return mysql.CreateUser(u)
}

func Login(param *models.ParamLogin) (userId int, err error) {
	user := &models.User{
		UserName: param.Username,
		PassWord: param.Password,
	}
	if err != nil {
		return -1, err
	}
	userId, err = mysql.Login(user)
	if err != nil {
		return -1, err
	}
	return userId, nil
}
