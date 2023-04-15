package conf

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
)

var (
	AppMode    string
	HttpPort   string
	DbUser     string
	DbPassWord string
	DbHost     string
	DbName     string
	DbPort     string
	Path       string

	MongoDBName   string
	MongoDBAddr   string
	MongoDBPwd    string
	MongoDBPort   string
	MongoDBClient *mongo.Client
)

func Init() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误", err)
	}
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbName = file.Section("mysql").Key("DbName").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	Path = strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort,
		")/", DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	//"用户名:密码@tcp(ip:port)/dbName?charset=utf8mb4&parseTime=True&loc=Local"
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
}
func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	fmt.Println("MongoDB连接成功")

}
