package controller

import (
	"errors"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/service"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(ctx *gin.Context) {
	param := new(models.ParamSignUp)
	// 1. 获取参数和参数校验
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Error("[SignUpHandler]解析参数异常", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2. 业务处理
	if err := service.SignUp(param); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(ctx, CodeUserExist)
			return
		}
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 3. 返回响应结果
	ResponseSuccess(ctx, struct{}{})
}

// Login 登录函数
func Login(ctx *gin.Context) {
	param := new(models.ParamLogin)
	// 1. 获取参数和参数校验
	if err := ctx.ShouldBindJSON(param); err != nil {
		zap.L().Error("[SignUpHandler]解析参数异常", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 进行业务处理
	token, err := service.Login(param)
	if err != nil {
		zap.L().Error("service.Login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
			return
		}
		ResponseError(ctx, CodeInvalidPassword)
		return
	}
	// 3.返回响应
	ResponseSuccess(ctx, token)
}
