package controller

import (
	"net/http"
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
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errors.Translate(trans)),
		})
		return
	}

	// 2. 业务处理
	if err := service.SignUp(param); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "注册失败",
			"data":    struct{}{},
		})
		return
	}
	// 3. 返回响应结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    struct{}{},
	})
}
