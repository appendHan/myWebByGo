package ws4Chatroom

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type client struct {
	hub       *hub
	clientCon *websocket.Conn
	send      chan []byte
	*ClientModel
}
type ClientModel struct {
	Uuid             string
	UserName         string
	RegistrationTime time.Time
}

type TransModel struct {
	Method string
	Uuid   string
	Data   interface{}
}

var clientMap = make(map[int]*client)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(s *hub, w http.ResponseWriter, r *http.Request) {
	wsCon, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	clientModel := &ClientModel{"", "无名氏", time.Now()} //数据模型
	client := &client{s, wsCon, make(chan []byte, 256), clientModel}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}

func (c *client) readPump() {
	for {
		_, p, err := c.clientCon.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var trans TransModel
		json.Unmarshal(p, &trans)
		switch trans.Method {
		case "setName":
			c.ClientModel.setName(trans.Uuid, trans.Data.(string))
			clientRefresh(c)
		case "sendMsg":
			handleBroadcast(c, trans.Data)
		case "Logout":
			c.hub.unRegister <- c
			clientRefresh(c)
			return
		}
	}
}

func (c *client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// 通道已关闭
				//c.clientCon.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println("疯狂写入中", message)
			c.clientCon.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (cm *ClientModel) setName(uuid string, name string) {
	if name != "" {
		cm.UserName = name
	}
	cm.Uuid = uuid
}

//刷新客户端列表并广播
func clientRefresh(c *client) {
	//todo 这个方式不好，最好有hub维护一份统一的客户端列表，而不是每次重新计算
	var clientMs []*ClientModel
	for i := range c.hub.clients {
		clientMs = append(clientMs, i.ClientModel)
	}
	trans := TransModel{"ClientUpdate", c.Uuid, clientMs}
	data, _ := json.Marshal(trans)
	c.hub.broadcast <- data
}

type MessageModel struct {
	UserName string
	Message  string
}

func handleBroadcast(c *client, msgJson interface{}) {
	mm, _ := msgJson.(map[string]interface{})
	//todo 这个interface解析方式太繁琐
	msg := MessageModel{c.UserName, mm["Message"].(string)}
	trans := TransModel{"BroadcoastMsg", c.Uuid, msg}
	data, _ := json.Marshal(trans)
	c.hub.broadcast <- data
}
