package mysql

import (
	"blue_bell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "wangbiao.com"

// dao层 将数据库的操作封装成函数
// 交给logic层的业务进行调用

// CheckUserExist 检查用户名是否存在
func CheckUserExist(username string) error {
	// sql查询语句-->判断当前用户是否存在于数据库中
	sqlStr := "select count(user_id) from user where username=?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 将用户插入到数据库中
func InsertUser(user *models.User) (err error) {
	// 对password进行加密
	user.Password = encryptPassword(user.Password)
	// 执行sql语句入库
	sqlStr := " insert into user (user_id,username,password) values(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err
}

// md5加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// 用户登录时查询数据库进行对比
func Login(user *models.User) (err error) {
	oPassword := user.Password
	//1. 先找到该用户
	sqlStr := "select user_id,username,password from user where username=?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows { // 表示用户不存在
		return ErrorUserNotExist
	}
	if err != nil {
		// 数据库查询失败
		return err
	}
	// 进行密码对比，判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	// 密码正确，则登录成功
	//
	return nil

}

// 根据uid查看用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id,username from user where user_id=?"
	err = db.Get(user, sqlStr, uid)
	return
}
