package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"web_app/models"
)

const secret = "miqimiaomiaowu"

var (
	ErrorUserExist = errors.New("用户已存在")
)

// CheckUserExist 根据用户名，判读用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(*) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// CreateUser 创建用户数据
func CreateUser(user *models.User) (err error) {
	// 注意：这里使用:xxx得和结构体中得db一致！！！
	sqlStr := "insert into user(user_id, username, password) values (:user_id,:username,:password)"
	user.PassWord = encryptPassword(user.PassWord)
	_, err = db.NamedExec(sqlStr, user)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
