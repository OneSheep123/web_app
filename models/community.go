package models

import "time"

type Community struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
}

type CommunityDetail struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
	// omitempty: 数据为空的时候, 不返回对应字段数据
	Introduction string `json:"introduction,omitempty" db:"introduction"`
	// 这里使用time.Time的时候, 记得在mysql连接加上parseTime=True
	CreateTime time.Time `json:"create_time" db:"create_time"`
}
