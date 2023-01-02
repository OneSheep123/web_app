package controller

import (
	"strconv"
	"web_app/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CommunityHandler 社区列表
func CommunityHandler(ctx *gin.Context) {
	list, err := service.GetCommunityList()
	if err != nil {
		zap.L().Error("获取社区列表失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, list)
}

// CommunityDetailHandler 社区详情
func CommunityDetailHandler(c *gin.Context) {
	param := c.Param("id")
	communityID, err := strconv.Atoi(param)
	if err != nil {
		zap.L().Error("参数有误", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	communityList, err := service.GetCommunityByID(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeSuccess, err.Error())
		return
	}
	ResponseSuccess(c, communityList)
}
