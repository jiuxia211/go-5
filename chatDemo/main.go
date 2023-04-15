package main

import (
	"jiuxia/chatDemo/conf"
	"jiuxia/chatDemo/model"
	"jiuxia/chatDemo/router"
	"jiuxia/chatDemo/service"
)

func main() {
	conf.Init()      //加载配置文件
	conf.MongoDB()   //连接MongoDB
	model.Database() //连接mysql
	go service.Manager.Start()
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
