package middleware

import (
	"fmt"
	"jiuxia/chatDemo/pkg/utils"
	"jiuxia/chatDemo/serializer"

	"time"

	"github.com/gin-gonic/gin"
)

func JwT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "token获取失败",
			})
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {

			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "token解析失败",
			})
			c.Abort()
			return
		} else if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(400, serializer.Response{
				Status: 400,
				Msg:    "token过期",
			})
			c.Abort()
			return
		}
		fmt.Println(time.Now().Unix())
		fmt.Println(claims.ExpiresAt)
		c.Next()

	}
}
