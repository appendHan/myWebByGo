package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	clientCon        *websocket.Conn
	UserName         string
	RegistrationTime time.Time
}

type TransModel struct {
	Method string
	Uuid   int
	Data   interface{}
}

var clientMap = make(map[int]*client)

func UpdateClient(uuid int, userName string, clientCon *websocket.Conn) {
	var name = ""
	if userName == "" {
		name = "无名氏"
	} else {
		name = userName
	}
	if oldClient, ok := clientMap[uuid]; ok {
		oldClient.UserName = name
	} else {
		clientMap[uuid] = &client{clientCon, name, time.Now()}
	}
	clientList := []*client{}
	for _, c := range clientMap {
		clientList = append(clientList, &client{c.clientCon,c.UserName,c.RegistrationTime})
	}
	for _, client := range clientMap {
		trans := TransModel{"ClientUpdate", uuid, clientList}
		data, _ := json.Marshal(trans)
		client.clientCon.WriteMessage(websocket.TextMessage, []byte(data))
	}
}

type MessageModel struct {
	UserName string
	Message  string
}

func BroadcastMsg(msgJson interface{}) {
	mm, _ := msgJson.(map[string]interface{})
	var name = ""
	if mm["UserName"].(string) == "" {
		name = "无名氏"
	} else {
		name = mm["UserName"].(string)
	}
	msg := MessageModel{name, mm["Message"].(string)}
	trans := TransModel{"BroadcoastMsg", 0, msg}
	data, _ := json.Marshal(trans)
	for _, client := range clientMap {
		client.clientCon.WriteMessage(websocket.TextMessage, []byte(data))
	}
}

func RemoveClient(uuid int)  {
	if _,ok := clientMap[uuid]; ok{
		delete(clientMap, uuid)
		clientList := []*client{}
		for _, c := range clientMap {
			clientList = append(clientList, &client{c.clientCon,c.UserName,c.RegistrationTime})
		}
		for _, client := range clientMap {
			trans := TransModel{"ClientUpdate", uuid, clientList}
			data, _ := json.Marshal(trans)
			client.clientCon.WriteMessage(websocket.TextMessage, []byte(data))
		}
	}
}