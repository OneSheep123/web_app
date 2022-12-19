package service

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
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

func Login(param *models.ParamLogin) (token string, err error) {
	user := &models.User{
		UserName: param.Username,
		PassWord: param.Password,
	}
	if err = mysql.Login(user); err != nil {
		return
	}
	return jwt.GenToken(user.UserID, user.UserName)
}
