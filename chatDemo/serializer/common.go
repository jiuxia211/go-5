package serializer

import "jiuxia/chatDemo/model"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
type User struct {
	UserName string `json:"username"`
	ID       string `json:"id"`
}
type FriendList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}
type Friend struct {
	Uid string `json:"uid"`
}

func BuildFriendList(items []model.Friend) (friendList []Friend) {
	for _, item := range items {
		var friend Friend
		friend.Uid = item.Fid
		friendList = append(friendList, friend)
	}
	return friendList
}
