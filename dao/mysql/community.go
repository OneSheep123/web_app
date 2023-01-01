package mysql

import (
	"database/sql"
	"web_app/models"

	"go.uber.org/zap"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	return
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(id string) (community *models.CommunityDetail, err error) {
	// 获取单个的时候会，这里需要进行new一下
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time from community where community_id = ?`
	err = db.Get(community, sqlStr, id)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}
