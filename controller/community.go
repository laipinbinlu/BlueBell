package controller

import (
	"blue_bell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 社区相关业务handler

func CommunityHandler(c *gin.Context) {
	//1.查询当前社区的所有消息(community_id,community_name 等等),最后以列表的形式返回  ->logic层业务处理
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易将后端服务详细错误返回给前端
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 根据社区id得到更加详细的社区信息
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取用户传入的id参数
	idstr := c.Param("id")
	// 校验参数
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam) // 请求参数错误
		return
	}

	// 执行业务->根据社区id查询到社区的详细信息
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易将后端服务详细错误返回给前端
		return
	}

	// 执行成功
	ResponseSuccess(c, data)
}
