package logic

import (
	"blue_bell/dao/mysql"
	"blue_bell/models"
	"blue_bell/pkg/jwt"
	"blue_bell/pkg/snowflake"
)

// 注册业务执行层，存在业务执行的代码，调用dao层完成相关业务
func SignUp(p *models.ParamSignUp) (err error) {
	// 1.用户是否存在--防止重复注册
	if err = mysql.CheckUserExist(p.Username); err != nil { // mysql查询语句出错
		return err
	}

	//2.生成用户的uid
	userID := snowflake.GetID()
	// 构造一个User实例方便传入数据库中
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.将用户的数据保存数据库中
	return mysql.InsertUser(user)
}

// 登录业务执行层
func Login(p *models.ParamLogin) (user *models.User, err error) {
	//1. 先创建user实例方便数据库操作
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// dao层执行相关业务
	if err = mysql.Login(user); err != nil {
		return user, err
	}
	// 登录没有问题，那么生成token
	user.Token, err = jwt.GenToken(user.UserID, user.Username)
	return
}
