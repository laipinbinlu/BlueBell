package models

// 定义请求参数的结构体

// 注册参数结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// 登录参数结构体
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 帖子投票的结构体
type ParamVoteData struct {
	// uid   直接获取
	PostID    string `json:"post_id" binding:"required"`              // 帖子投票id
	Direction int8   `json:"direction,string" binding:"oneof=1 -1 0"` // 投票结果:赞成（1）反对（-1）取消（0）
}

// 请求帖子列表参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`
	Page        int64  `form:"page" json:"page"`
	Size        int64  `form:"size" json:"size"`
	Order       string `form:"order" json:"order"`
}

var (
	OrderTime  = "time"
	OrderScore = "score"
)
