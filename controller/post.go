package controller

import (
	"blue_bell/logic"
	"blue_bell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1. 获取参数和参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("ShouldBindJSON post failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取用户id
	uid, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = uid

	//2. 创建帖子 ->业务层
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回处理结果
	ResponseSuccess(c, nil)
}

// 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1. 获取帖子id（从URL中）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id查出帖子的具体内容
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回结果
	ResponseSuccess(c, data)
}

// 展示帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页数据 ->用于展示帖子
	page, size := getPageInfo(c)

	// 1.获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//2.返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 展示帖子列表-->按照时间或者分数
// 需要根据前端传来的参数，按照时间或者分数展示帖子
// 1.获取前端的参数
// 2. 查询redis的post数据，获得post列表展示顺序
// 3.根据id查找mysql数据库，将post列表信息补全
// 4.分页展示post列表
func GetPostListHandler2(c *gin.Context) {
	//1. 获取前端传来的参数  get 请求
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("GetPostListHandler2 ShouldBind post failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 1.获取数据
	data, err := logic.GetPostListNew(p)

	if err != nil {
		zap.L().Error("logic.GetPostListNew failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//2.返回响应
	ResponseSuccess(c, data)
}

// 根据社区id 查询post列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	//1. 获取前端传来的参数  get 请求
//	p := &models.ParamCommunityPostList{
//		ParamPostList: &models.ParamPostList{
//			Page:  1,
//			Size:  10,
//			Order: models.OrderTime,
//		},
//		CommunityID: 0,
//	}
//
//	if err := c.ShouldBind(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler ShouldBind post failed", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	// 1.获取数据
//	data, err := logic.GetCommunityPostList2(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostList2 failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	//2.返回响应
//	ResponseSuccess(c, data)
//}
