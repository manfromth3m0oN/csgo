package main

import (
	"github.com/gin-gonic/gin"
	"github.com/manfromth3m0oN/csgo/model"
	log "github.com/sirupsen/logrus"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func createRoom(c *gin.Context) {
	log.Info("called create room")

	var req model.CreateRoomReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Error("Failed to get req json ", err.Error())
		c.JSON(400, gin.H{
			"message": "borked req",
		})
		return
	}

	users := make([]model.User, 0)
	users = append(users, model.User{Name: req.UName, Addr: c.Request.RemoteAddr})
	log.Info("Added user")

	newRoom := model.Room{
		Name:     req.Name,
		Users:    users,
		Playlist: []string{"https://storage.googleapis.com/downloads.webmproject.org/media/video/webmproject.org/big_buck_bunny_trailer_480p_logo.webm"},
		Index:    0,
	}
	rooms[req.Name] = newRoom
	log.Info("Added room")

	c.JSON(200, gin.H{"room": newRoom.Name})
	return
}

func handleRoom(c *gin.Context) {
	roomName := c.Param("room")
	room := rooms[roomName]
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	msg := model.InitialMsg{
		Video: room.Playlist[room.Index],
	}

	ws.WriteJSON(msg)
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Info(mt, string(message))
	}
}
