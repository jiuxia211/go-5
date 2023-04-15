package service

import (
	"encoding/json"
	"fmt"
	"jiuxia/chatDemo/cache"
	"jiuxia/chatDemo/conf"
	"jiuxia/chatDemo/pkg/e"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const month = 60 * 60 * 24 * 30

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}
type Client struct {
	ID      string
	SendID  string
	GroupID string
	Socket  *websocket.Conn
	Send    chan []byte
}
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}
func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	gid := c.Query("gid")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建用户实例
	client := &Client{
		ID:      CreateID(uid, toUid),
		SendID:  CreateID(toUid, uid),
		GroupID: "group" + gid,
		Socket:  conn,
		Send:    make(chan []byte),
	}
	//用户注册到Manager上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		//c.Socket.ReadMessage  接收字符串
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			fmt.Println("数据格式不正确")
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		if sendMsg.Type == 1 {
			//发送消息
			r1, _ := cache.RedisClient.Get(c.ID).Result()
			r2, _ := cache.RedisClient.Get(c.SendID).Result()
			if r1 > "3" && r2 == "" {
				//1给2 发了三条 2没有回复，就停止1发送
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "发送消息达到限制",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(c.ID)
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result()
			}
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
				Type:    sendMsg.Type,
			}
		} else if sendMsg.Type == 2 {
			//获取历史消息
			results, err := FindMany(conf.MongoDBName, c.SendID, c.ID, 10)
			if err != nil {
				panic(err)
			}
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}

		} else if sendMsg.Type == 3 { //查看历史记录
			results, err := FirsFindtMsg(conf.MongoDBName, c.SendID, c.ID)
			if err != nil {
				panic(err)
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		} else if sendMsg.Type == 4 { //群聊
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
				Type:    sendMsg.Type,
			}
		} else if sendMsg.Type == 5 { //获取群聊历史信息
			results, err := FindGroupMany(conf.MongoDBName, c.GroupID, 10)
			if err != nil {
				panic(err)
			}
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}

	}
}
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replymsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: string(message),
			}
			msg, _ := json.Marshal(replymsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)

		}
	}
}
