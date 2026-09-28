package utils

import "github.com/gin-gonic/gin"

// 统一响应结构 {code,msg,data}
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 分页结果
type PageResult struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	List     interface{} `json:"list"`
}

// Success 200
func Success(c *gin.Context, data any) {
	c.JSON(200, Response{Code: 200, Msg: "success", Data: data})
}

// Error 业务错误
func Error(c *gin.Context, code int, msg string) {
	c.JSON(200, Response{Code: code, Msg: msg, Data: nil})
}
