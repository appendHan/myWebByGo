package main

import (
	"./utils/ws4Chatroom"
	"./web/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tmpInt = 1

func setupRouter() *gin.Engine {
	// Disable Console Color
	//gin.DisableConsoleColor()

	r := gin.Default()
	//静态文件
	r.Static("/static", "web/static")
	//模板文件
	r.LoadHTMLGlob("web/templates/*")

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	//通过 channel 建立 单消费者-多生产者模式
	wsServer := ws4Chatroom.NewAndRun()
	chatRoom := r.Group("/chatRoom")
	{
		chatRoom.GET("/", func(c *gin.Context) {
			tmpIndex := models.TemplateIndex{Title: "chatRoom"}
			c.HTML(http.StatusOK, "chatRoom.html", tmpIndex)
		})
		chatRoom.GET("/ws", func(c *gin.Context) {
			ws4Chatroom.ServeWs(wsServer,c.Writer, c.Request)
		})
	}

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8083")
}
