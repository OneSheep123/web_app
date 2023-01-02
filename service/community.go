package service

import (
	"web_app/dao/mysql"
	"web_app/models"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityByID(id int) (community *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetail(id)
}
