package api

import (
	"jiuxia/chatDemo/pkg/utils"
	"jiuxia/chatDemo/service"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var useRegister service.UserService
	if err := c.ShouldBind(&useRegister); err == nil {
		res := useRegister.Register()
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}

func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login()
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
func UserAddFriend(c *gin.Context) {
	var friendService service.FriendService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	fid := c.Param("fid")
	if err := c.ShouldBind(&friendService); err == nil {
		res := friendService.AddFriend(claim.Id, fid)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
func UserGetFriendList(c *gin.Context) {
	var friendService service.FriendService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&friendService); err == nil {
		res := friendService.GetFriendList(claim.Id)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
