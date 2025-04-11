package mysql

import (
	"blue_bell/models"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	// 执行对应相关业务的sql语句
	sqlStr := "select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community  in db")
			err = nil
		}
	}
	return communityList, err
}

// 根据id查看社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlstr := "select community_id, community_name, introduction, create_time from community where community_id = ?"
	err = db.Get(community, sqlstr, id)
	if err != nil {
		if err == sql.ErrNoRows { // id不存在
			err = ErrorInvalidID
		}
	}
	return community, err
}
