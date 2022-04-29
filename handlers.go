package main

import (
	"net/http"

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

	users := make([]string, 0)
	users = append(users, req.UName)
	log.Info("Added user")

	newRoom := model.Room{
		Name:     req.Name,
		Users:    users,
		Playlist: []string{"https://storage.googleapis.com/downloads.webmproject.org/media/video/webmproject.org/big_buck_bunny_trailer_480p_logo.webm"},
		Index:    0,
	}
	rooms[req.Name] = &newRoom
	log.Info("Added room")

	go newRoom.Run()

	c.JSON(http.StatusOK, gin.H{"room": newRoom.Name})
	return
}

func handleRoom(c *gin.Context) {
	roomName := c.Param("room")
	var joinReq model.RoomJoinReq
	err := c.BindJSON(&joinReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	room := rooms[roomName]
	room.Mutex.Lock()
	room.Join(joinReq.Username)
	room.Mutex.Unlock()
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
		outMsg := <-room.OutChan
		ws.WriteJSON(outMsg)

		var inMsg model.Event
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Warn("Read err ", err)
		}
		log.Info(msg)

		room.InChan <- inMsg
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}
