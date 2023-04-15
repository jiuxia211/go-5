package router

import (
	"jiuxia/chatDemo/api"
	"jiuxia/chatDemo/middleware"
	"jiuxia/chatDemo/service"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(200, "success")
		})
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		authorized := v1.Group("/")
		authorized.Use(middleware.JwT())
		{
			v1.GET("ws", service.Handler)
			v1.POST("friend/:fid", api.UserAddFriend)
			v1.GET("friend/list", api.UserGetFriendList)
		}

	}
	return r
}
