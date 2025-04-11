package logic

import (
	"blue_bell/dao/mysql"
	"blue_bell/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	// 调用dao层，查询数据库，并且将社区数据返回
	return mysql.GetCommunityList()

}

// 业务:根据id查询社区详细信息的业务
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// 直接调用dao层
	return mysql.GetCommunityDetailByID(id)
}
