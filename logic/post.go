package logic

import (
	"blue_bell/dao/mysql"
	"blue_bell/dao/redis"
	"blue_bell/models"
	"blue_bell/pkg/snowflake"
	"fmt"

	"go.uber.org/zap"
)

// 创建post业务
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post id
	p.ID = snowflake.GetID()
	// 2. 插入到数据库中   mysql 和redis数据库
	//3. 返回结果   --->创建帖子 先创建mysql再在创建redis帖子消息。
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	// redis创建时间表
	err = redis.CreatePost(p.ID, p.CommunityID)
	return err
}

// 根据id查到贴子的具体数据并且返回 ->需要对数据进行修饰，返回原始数据不合理
func GetPostById(pid int64) (p *models.ApiPostDetail, err error) {
	// 查询并且并且满足返回数据的
	//  查询帖子详情
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Error(err))
		return
	}
	// 查询用户信息根据用户id
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorID)), zap.Error(err))
		return
	}
	// 社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
		return
	}
	// 开始拼接数据
	p = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	return
}

// 获取帖子信息->列表形式
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorID)), zap.Error(err))
			continue
		}
		// 社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		// 开始拼接数据
		p := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, p)
	}
	return
}

// 获取帖子信息2
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 2. 查询redis的post数据，获得post列表展示顺序
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() len(ids) == 0")
		return
	}

	// 3.根据post_id查找mysql数据库，将post列表信息补全
	// 按照用户给定的顺序获得帖子信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 4.查询每个帖子的赞成票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子详细信息返回给前端
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorID)), zap.Error(err))
			continue
		}
		// 社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		// 开始拼接数据
		p := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, p)
	}
	return
}

// 社区id查询到post信息
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 2. 查询redis的post数据，获得post列表展示顺序
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder() len(ids) == 0")
		return
	}

	// 3.根据post_id查找mysql数据库，将post列表信息补全
	// 按照用户给定的顺序获得帖子信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 4.查询每个帖子的赞成票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子详细信息返回给前端
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorID)), zap.Error(err))
			continue
		}
		// 社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		// 开始拼接数据
		p := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, p)
	}
	return
}

// 将两个查询的帖子合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据参数判断不同的业务逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}
	return
}
