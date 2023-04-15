package service

import (
	"context"
	"fmt"
	"jiuxia/chatDemo/conf"
	"jiuxia/chatDemo/model/ws"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

func InsertMsg(database, id string, content string, read uint, expire int64) error {
	//插入到mongoDB中
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}
func FindMany(database, sendID, id string, pageSize int) (results []ws.Result, err error) {
	var resultsMe []ws.Trainer
	var resultsYou []ws.Trainer

	sendIDCollection := conf.MongoDBClient.Database(database).Collection(sendID)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)
	var findOptions = &options.FindOptions{}
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSort(bson.D{{"startTime", -1}})
	sendIDCurcor, err := sendIDCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		panic(err)
	}
	idTimeCurcor, err := idCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		panic(err)
	}
	err = sendIDCurcor.All(context.TODO(), &resultsYou)
	if err != nil {
		panic(err)
	}
	err = idTimeCurcor.All(context.TODO(), &resultsMe)
	if err != nil {
		panic(err)
	}
	results, _ = AppendAndSort(resultsMe, resultsYou)
	return
}
func FindGroupMany(database, groupId string, pageSize int) (results []ws.Result, err error) {
	var resultsGroup []ws.Trainer
	groupCollection := conf.MongoDBClient.Database(database).Collection(groupId)
	var findOptions = &options.FindOptions{}
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSort(bson.D{{"startTime", -1}})
	groupCurcor, err := groupCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	err = groupCurcor.All(context.TODO(), &resultsGroup)
	if err != nil {
		panic(err)
	}
	for _, r := range resultsGroup {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      groupId,
		}
		results = append(results, result)
	}
	return results, err
}
func FirsFindtMsg(database string, sendId string, id string) (results []ws.Result, err error) {
	// 首次查询(把对方发来的所有未读都取出来)
	var resultsMe []ws.Trainer
	var resultsYou []ws.Trainer
	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	filter := bson.D{{"read", 0}}
	var findOptions = &options.FindOptions{}
	findOptions.SetSort(bson.D{{"startTime", 1}})
	sendIdCursor, err := sendIdCollection.Find(context.TODO(), filter, findOptions)
	if sendIdCursor == nil {
		fmt.Println("过滤器里啥都没有")
		return
	}
	var unReads []ws.Trainer
	err = sendIdCursor.All(context.TODO(), &unReads)
	if err != nil {
		panic(err)
	}

	if len(unReads) > 0 {
		timeFilter := bson.M{
			"startTime": bson.M{
				"$gte": unReads[0].StartTime,
			},
		}
		fmt.Println("unReads里有东西")
		sendIdTimeCursor, _ := sendIdCollection.Find(context.TODO(), timeFilter)
		idTimeCursor, _ := idCollection.Find(context.TODO(), timeFilter)
		err = sendIdTimeCursor.All(context.TODO(), &resultsYou)
		err = idTimeCursor.All(context.TODO(), &resultsMe)
		results, err = AppendAndSort(resultsMe, resultsYou)
	} else {
		fmt.Println("unReads里没有东西")
		results, err = FindMany(database, sendId, id, 10)
	}
	overTimeFilter := bson.D{
		{"$and", bson.A{
			bson.D{{"endTime", bson.M{"&lt": time.Now().Unix()}}},
			bson.D{{"read", bson.M{"$eq": 1}}},
		}},
	}
	_, _ = sendIdCollection.DeleteMany(context.TODO(), overTimeFilter)
	_, _ = idCollection.DeleteMany(context.TODO(), overTimeFilter)
	// 将所有的维度设置为已读
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"read": 1},
	})
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"ebdTime": time.Now().Unix() + int64(3*month)},
	})
	return
}
func AppendAndSort(resultMe, resultYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultMe {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, r := range resultYou {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	return
}
