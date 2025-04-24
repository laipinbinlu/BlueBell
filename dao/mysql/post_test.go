package mysql

import (
	"blue_bell/models"
	"blue_bell/settings"
	"testing"
)

func init() {
	dbcfg := settings.MysqlConfig{ // 写的是测试的时候的数据库的地址，并不是实际数据库的地址
		Host:      "127.0.0.1",
		Port:      3306,
		User:      "root",
		Password:  "123456",
		DBName:    "bluebell",
		MaxActive: 10,
		MaxIdle:   1,
	}
	err := Init(&dbcfg)
	if err != nil {
		panic(err)
	}
}

// 需要初始化db
func TestCreatePost(t *testing.T) {
	p := models.Post{
		ID:          12345,
		AuthorID:    123123,
		CommunityID: 1,
		Title:       "test",
		Content:     "就是一个测试",
	}
	if err := CreatePost(&p); err != nil {
		t.Fatalf("create post failed, err:%v\n", err)
	}
	t.Logf("test create post succeed")
}
