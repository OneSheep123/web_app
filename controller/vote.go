package controller

import (
	"fmt"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

/*
投票算法：http://www.ruanyifeng.com/blog/2012/03/ranking_algorithm_reddit.html
*/

// CreateVoteHandler 创建投票
func CreateVoteHandler(ctx *gin.Context) {
	vote := new(models.ParamVoteData)
	if err := ctx.ShouldBindJSON(vote); err != nil {
		zap.L().Error("[CreateVoteHandler]解析参数异常", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, service.RemoveTopStruct(errs.Translate(service.Trans)))
		return
	}

	userID, err := getCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNoLogin)
		return
	}
	if err := redis.PostVote(strconv.Itoa(int(vote.PostId)), fmt.Sprint(userID), float64(vote.Direction)); err != nil {
		zap.L().Error("发起投票失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, nil)
}
