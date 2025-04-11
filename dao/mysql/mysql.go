package mysql

import (
	"blue_bell/settings"
	"fmt"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(sqlconfig *settings.MysqlConfig) (err error) {
	// dsn 由配置文件提供   按照格式组合形成dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		sqlconfig.User, sqlconfig.Password, sqlconfig.Host, sqlconfig.Port, sqlconfig.DBName)

	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(sqlconfig.MaxActive)
	db.SetMaxIdleConns(sqlconfig.MaxIdle)

	return
}

func Close() {
	_ = db.Close()
}
