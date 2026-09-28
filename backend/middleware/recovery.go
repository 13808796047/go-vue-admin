package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Recovery 全局错误捕获（panic 恢复 + 统一 {code,msg,data} 500 响应）
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		_ = err
		c.Abort()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器内部错误", "data": nil})
	})
}

// ContainsAny 字符串包含判断
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
