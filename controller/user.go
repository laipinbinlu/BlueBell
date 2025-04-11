package controller

import (
	"blue_bell/dao/mysql"
	"blue_bell/logic"
	"blue_bell/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1. 获取用户参数（gin框架获取），参数校验（是否满足基本要求）
	var p = new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("ShouldBindJSON with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok { // 不是validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 请求参数有误，直接返回响应给前端，翻译器
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))

		return
	}

	//2. 相关业务（用户注册）处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3. 返回处理的结果（response）
	ResponseSuccess(c, nil)
}

// 用户登录
func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("login shouldbindjson with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok { // 不是validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	//2. 业务逻辑处理，执行对应的登录业务
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login with failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3. 返回执行结果
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // "变为string类型防止json失真"
		"user_name": user.Username,
		"token":     user.Token,
	})
}
