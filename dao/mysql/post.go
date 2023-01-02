package mysql

import (
	"database/sql"
	"web_app/models"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostID, post.Title, post.Content, post.AuthorId, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// GetPostDetail 获取帖子详情
func GetPostDetail(id int) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select id, post_id, author_id, community_id, status, title, content, create_time from post where id = ?`
	err = db.Get(post, sqlStr, id)
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

func GetPostList(page int, size int) (postList []*models.ApiPostDetail, err error) {
	postList = make([]*models.ApiPostDetail, 0, 2)
	sqlStr := `select id, post_id, author_id, community_id, status, title, content, create_time ` +
		`from post limit ?, ?`
	err = db.Select(&postList, sqlStr, (page-1)*size, size)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	return
}
