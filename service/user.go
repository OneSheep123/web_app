package service

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/snowflake"
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
