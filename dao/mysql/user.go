package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"web_app/models"
)

const secret = "miqimiaomiaowu"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
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

// GetUserId 通过user信息获取userId
func GetUserId(user *models.User) (int, error) {
	sqlStr := "select user_id from user where username=?"
	rows, err := db.Queryx(sqlStr, user.UserName)
	if err != nil {
		return -1, err
	}
	userId := 0
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&userId)
	}
	return userId, nil
}

// Login 登录用户
func Login(user *models.User) (userId int, err error) {
	oPassword := user.PassWord // 用户登录的密码
	sqlStr := "select user_id,username,password from user where username=?"
	err = db.Get(user, sqlStr, user.UserName)
	if err == sql.ErrNoRows {
		return -1, ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return -1, err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.PassWord {
		return -1, ErrorInvalidPassword
	}
	return int(user.UserID), nil
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
