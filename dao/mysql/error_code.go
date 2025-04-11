package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已经存在")
	ErrorUserNotExist    = errors.New("用户名不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorGenIDFailed     = errors.New("创建用户ID失败")
	ErrorInvalidID       = errors.New("无效的ID")
)
