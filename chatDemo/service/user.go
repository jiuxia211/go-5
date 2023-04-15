package service

import (
	"jiuxia/chatDemo/model"
	"jiuxia/chatDemo/pkg/utils"
	"jiuxia/chatDemo/serializer"
	"strconv"

	"github.com/jinzhu/gorm"
)

type UserService struct {
	UserName string `form:"username" json:"user_name" `
	Password string `form:"password" json:"password" `
}
type FriendService struct {
}

func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).
		First(&user).Count(&count)
	if count == 1 {
		return serializer.Response{
			Status: 400,
			Msg:    "该用户已注册！",
		}
	}
	user.UserName = service.UserName
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    err.Error(),
		}
	}
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "数据库创建用户失败",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "创建成功",
	}
}
func (service *UserService) Login() serializer.Response {
	var user model.User
	if err := model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).
		First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: 400,
				Msg:    "用户不存在",
			}
		} else {
			return serializer.Response{
				Status: 500,
				Msg:    "数据库错误",
			}
		}

	}
	if !user.CheckPassword(service.Password) {
		return serializer.Response{
			Status: 400,
			Msg:    "密码错误",
		}
	}
	token, err := utils.GenerateToken(service.UserName, user.ID, service.Password)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "token签发错误",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "登录成功",
		Data: serializer.TokenData{
			User: serializer.User{
				UserName: service.UserName,
				ID:       strconv.Itoa(int(user.ID)),
			},
			Token: token,
		},
	}
}
func (service *FriendService) AddFriend(uid, fid string) serializer.Response {
	var friend model.Friend
	friend.Uid = uid
	friend.Fid = fid
	if err := model.DB.Create(&friend).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "数据库创建好友关系失败",
		}
	}
	friend.Uid = fid
	friend.Fid = uid
	return serializer.Response{
		Status: 200,
		Msg:    "添加好友成功",
	}
}
func (service *FriendService) GetFriendList(uid string) serializer.Response {
	var friends []model.Friend
	count := 0
	model.DB.Model(&model.Friend{}).Where("uid=?", uid).Count((&count)).Find(&friends)
	return serializer.Response{
		Status: 200,
		Msg:    "好友的uid如下",
		Data: serializer.FriendList{
			Item:  serializer.BuildFriendList(friends),
			Total: uint(count),
		},
	}

}
