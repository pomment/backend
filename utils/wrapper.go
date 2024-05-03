package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func AjaxSuccess(c *gin.Context, content interface{}) {
	c.JSON(200, Response{
		Code: 200,
		Data: content,
	})
}

func AjaxSuccessWithReturns(c *gin.Context, content interface{}) ([]byte, error) {
	c.JSON(200, Response{
		Code: 200,
		Data: content,
	})
	res, err := json.Marshal(Response{
		Code: 200,
		Data: content,
	})
	return res, err
}

func AjaxError(c *gin.Context, code int, err error) {
	fmt.Printf("Error: %s\n", err.Error())
	var httpCode = code
	if httpCode < 200 || httpCode >= 600 {
		httpCode = 500
	}
	c.JSON(httpCode, Response{
		Code: code,
		Data: nil,
	})
}

func AjaxWrapJSONString(c *gin.Context, content string) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, fmt.Sprintf("{\"code\":200,\"data\":%s}", content))
}

func AjaxJSONString(c *gin.Context, content string) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, content)
}
