package models

import (
	"errors"
	"time"

	"github.com/goccy/go-json"
)

type Post struct {
	// 记得内存对齐
	Id int64 `json:"id" db:"id"`
	// postId使用雪花算法, 大小可能会超过前端整数表示的大小
	// 因此序列化时进行序列化为字符串, 同时反序列化时, 可以反序列化为整数
	PostID      int64 `json:"post_id,string" db:"post_id"`
	AuthorId    int64 `json:"author_id" db:"author_id"`
	CommunityID int64 `json:"community_id" db:"community_id" binding:"required"`
	// 帖子状态
	Status     int32     `json:"status" db:"status"`
	Title      string    `json:"title" db:"title" binding:"required"`
	Content    string    `json:"content" db:"content" binding:"required"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

func (p *Post) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Title       string `json:"title" db:"title"`
		Content     string `json:"content" db:"content"`
		CommunityID int64  `json:"community_id" db:"community_id"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.Title) == 0 {
		err = errors.New("帖子标题不能为空")
	} else if len(required.Content) == 0 {
		err = errors.New("帖子内容不能为空")
	} else if required.CommunityID == 0 {
		err = errors.New("未指定版块")
	} else {
		p.Title = required.Title
		p.Content = required.Content
		p.CommunityID = required.CommunityID
	}
	return
}

type ApiPostDetail struct {
	*Post
	AuthorName    string `json:"author_name"`
	CommunityName string `json:"community_name"`
}
