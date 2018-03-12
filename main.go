package main

import (
	"./web/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

var DB = make(map[string]string)

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, nil)
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	//gin.DisableConsoleColor()
	rootBasePath := "/Users/zhangsihang/Documents/GitHub/myWebByGo"

	r := gin.Default()
	//静态文件
	r.Static("/static", rootBasePath+"/web/static")
	//模板文件
	r.LoadHTMLGlob(rootBasePath + "/web/templates/*")

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.POST("/read", func(context *gin.Context) {

		})
	}

	r.GET("/", func(c *gin.Context) {
		tmpIndex := models.TemplateIndex{
			Title: "Home",
		}
		c.HTML(http.StatusOK, "index.html", tmpIndex)
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8082")
}
