package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/manfromth3m0oN/csgo/model"
)

var rooms map[string]*model.Room
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	rooms = make(map[string]*model.Room)
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/ping", ping)
	r.POST("/create-room", createRoom)
	r.GET("/room/:room", handleRoom)
	r.Run(":3000")
}
