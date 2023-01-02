package models

// 定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票请求参数
type ParamVoteData struct {
	PostId    int64 `json:"post_id,string" binding:"required"` // 帖子id
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"`
}

// ParamOffset 请求参数
type ParamOffset struct {
	Page int `json:"page" form:"page" example:"1"`  // 页码
	Size int `json:"size" form:"size" example:"12"` // 页长
}
