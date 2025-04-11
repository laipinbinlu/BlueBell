package controller

import (
	"blue_bell/logic"
	"blue_bell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func PostVoteController(c *gin.Context) {
	// 1. 参数的校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // 翻译为中文
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 先获取请求用户的id
	uid, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 2. 执行投票业务
	if err := logic.VoteForPost(uid, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(c, nil)
}
