package service

import (
	"fmt"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	id := snowflake.GenID()
	post.PostID = id
	// 创建帖子
	if err = mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	community, err := mysql.GetCommunityNameById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
		return err
	}
	if err := redis.CreatePost(
		fmt.Sprint(post.PostID),
		fmt.Sprint(post.AuthorId),
		post.Title,
		TruncateByWords(post.Content, 120),
		community.CommunityName); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return
}

// GetPostDetail 获取帖子详情
func GetPostDetail(id int) (res *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostDetail(id)
	res = &models.ApiPostDetail{
		AuthorName:    "",
		CommunityName: "",
	}
	res.Post = post
	// 查询用户名
	user, err := mysql.GetUserNameByAuthorId(post.AuthorId)
	if err != nil {
		zap.L().Error("获取用户名失败, 失败原因为", zap.Error(err))
		return nil, err
	}
	res.AuthorName = user.UserName
	community, err := mysql.GetCommunityNameById(post.CommunityID)
	if err != nil {
		zap.L().Error("获取社区失败, 失败原因为", zap.Error(err))
		return nil, err
	}
	res.CommunityName = community.CommunityName
	return
}

func GetPostList2(page, size int) (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		fmt.Println(err)
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
			continue
		}
		post.AuthorName = user.UserName
		community, err := mysql.GetCommunityNameById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityNameById() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		post.CommunityName = community.CommunityName
		data = append(data, post)
	}
	return
}
