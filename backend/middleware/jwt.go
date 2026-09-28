package middleware

import (
	"net/http"
	"strings"

	"go-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

// ContextKey 存储登录用户 ID
const ContextKeyUserID = "userID"
const ContextKeyUsername = "username"

// JWT 鉴权中间件：Authorization: Bearer <token>
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			header = "Bearer " + c.Query("token")
		}
		token := extractToken(header)
		if token == "" {
			utils.Error(c, http.StatusUnauthorized, "缺少 Authorization 头（Bearer token）")
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(secret, token)
		if err != nil {
			msg := "token 无效"
			if err == utils.ErrTokenExpired {
				msg = "token 已过期，请重新登录"
			}
			utils.Error(c, http.StatusUnauthorized, msg)
			c.Abort()
			return
		}
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Next()
	}
}

func extractToken(header string) string {
	header = strings.TrimSpace(header)
	if strings.HasPrefix(header, "Bearer ") {
		header = strings.TrimPrefix(header, "Bearer ")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}
