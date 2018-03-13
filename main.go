package main

import (
	"./web/models"
	."./utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
)

var DB = make(map[string]string)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var tmpInt = 1


func setupRouter() *gin.Engine {
	// Disable Console Color
	//gin.DisableConsoleColor()
	rootBasePath := "D:/GoWorkSpare/myWebByGo"

	r := gin.Default()
	//静态文件
	r.Static("/static", rootBasePath+"/web/static")
	//模板文件
	r.LoadHTMLGlob(rootBasePath + "/web/templates/*")

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/", func(c *gin.Context) {
		tmpIndex := models.TemplateIndex{
			Title: "chatRoom",
		}
		c.HTML(http.StatusOK, "chatRoom.html", tmpIndex)
	})

	r.GET("/ws", func(context *gin.Context) {
		wsCon,err := upgrader.Upgrade(context.Writer,context.Request,nil)
		if(err != nil){
			return
		}
		for {
			_, p, err := wsCon.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			var trans TransModel
			json.Unmarshal(p,&trans)
			switch trans.Method {
				case "setName":
					//log.Println("更新客户列表")
					UpdateClient(trans.Uuid,trans.Data.(string),wsCon)
				case "sendMsg":
					BroadcastMsg(trans.Data)
				case "Logout":
					RemoveClient(trans.Uuid)
			}
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8083")
}
