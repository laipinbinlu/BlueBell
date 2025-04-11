package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 单元测试

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler) // 测试案例

	body := `{
		"community_id" : 1,
		"title": "test",
		"content": "just a test"
	}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body))) // 模拟请求

	w := httptest.NewRecorder() // 得到回复结构体

	r.ServeHTTP(w, req) // 服务

	assert.Equal(t, 200, w.Code)

	assert.Contains(t, w.Body.String(), "需要登录")
}
