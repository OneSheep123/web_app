package controller

import (
	"strconv"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePost 创建帖子
func CreatePost(ctx *gin.Context) {
	post := new(models.Post)
	if err := ctx.ShouldBindJSON(post); err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	// 获取用户id
	id, err := getCurrentUserID(ctx)
	if err != nil {
		zap.L().Error("获取用户id错误", zap.Error(err))
		ResponseError(ctx, CodeUserNotExist)
		return
	}
	post.AuthorId = id
	err = service.CreatePost(post)
	if err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, nil)
}

// GetPostDetail 获取帖子详情
func GetPostDetail(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		zap.L().Error("参数有误", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	detail, err := service.GetPostDetail(id)
	if err != nil {
		zap.L().Error("service.GetPostDetail() failed", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeSuccess, err.Error())
		return
	}
	ResponseSuccess(ctx, detail)
}

// GetPostList 获取帖子列表
func GetPostList(ctx *gin.Context) {
	sizeParam, isSize := ctx.GetQuery("size")
	pageParam, isPage := ctx.GetQuery("page")
	if sizeParam == "" || !isSize || pageParam == "" || !isPage {
		zap.L().Error("获取帖子列表, 请求参数有误",
			zap.String("size", sizeParam),
			zap.String("page", pageParam))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	page, _ := strconv.Atoi(pageParam)
	size, _ := strconv.Atoi(sizeParam)
	list, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("获取帖子列表失败",
			zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, list)
}

func PostList2Handler(ctx *gin.Context) {
	offset := new(models.ParamOffset)
	if err := ctx.ShouldBindQuery(offset); err != nil {
		zap.L().Error("获取帖子列表, 请求参数有误", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	data, err := service.GetPostList2(offset.Page, offset.Size)
	if err != nil {
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, data)
}
