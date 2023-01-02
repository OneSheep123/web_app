package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"

	"github.com/gin-gonic/gin"
)

func TestCreatePost(t *testing.T) {

	gin.SetMode(gin.TestMode)

	engine := gin.Default()
	url := "/api/v1/post"
	engine.POST(url, CreatePost)

	body := `{
		"community_id": 1,
		"title": "test",
		"content": :"just a test"
    }`
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, request)

	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal(w.Body.Bytes(), res) failed, err:%v\n", err)
	}
	t.Logf("解析出来的数据为:%v\n", *res)
}
