package redis

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// IsUserLogin 获取用户登录状态
func IsUserLogin(userId int) bool {
	return rdb.Get(strconv.Itoa(userId)).Err() != redis.Nil
}

// SaveUserLoginState 保存用户状态信息
func SaveUserLoginState(userId int, token string) {
	rdb.SetNX(strconv.Itoa(userId), token, time.Hour*24)
}
