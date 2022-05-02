package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/manfromth3m0oN/csgo/ch"
	"github.com/manfromth3m0oN/csgo/model"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/syncmap"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func createRoom(c *gin.Context) {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	sugar := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	sugar.Info("called create room")

	var req model.CreateRoomReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		sugar.Error("Failed to get req json ", zap.Error(err))
		c.JSON(400, gin.H{
			"message": "borked req",
		})
		return
	}

	users := syncmap.Map{}

	playlist, err := ch.GetMedia(req.BoardSn, req.ThreadID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "bad request",
		})
	}

	newRoom := model.Room{
		Name:     req.Name,
		Users:    users,
		Playlist: playlist,
		Index:    0,
		InChan:   make(chan model.Event),
	}
	rooms[req.Name] = &newRoom
	sugar.Info("Added room")

	go newRoom.Run()

	c.JSON(http.StatusOK, gin.H{"room": newRoom.Name})
	return
}

func handleRoom(c *gin.Context) {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	sugar := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	sugar.Info("Started room handler")
	roomName := c.Param("room")
	me := c.Query("username")
	if me == "" {
		sugar.Warn("No username param")
		return
	}
	sugar.Named(me)

	room := rooms[roomName]
	sugar.Info("got room ", zap.String("room", room.Name))
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		sugar.Error("failed to upgrade connection")
		panic(err)
	}
	sugar.Info("Upgraded websocket")
	defer ws.Close()
	ws.SetCloseHandler(func(code int, text string) error {
		sugar.Warn("Socket closed", zap.Int("code", code), zap.String("text", text))
		room.InChan <- model.Event{
			Type: model.Leave,
			Data: me,
		}
		return errors.New("Websocket closed")
	})

	var wg sync.WaitGroup

	// Read pump
	wg.Add(1)
	go func() {
		for {
			var msg model.Event
			err := ws.ReadJSON(&msg)
			if err != nil {
				return
			}
			sugar.Info("Read message ", zap.String("type", msg.Type), zap.String("data", msg.Data))

			room.InChan <- msg
		}
	}()

	// Write pump
	wg.Add(1)
	go func() {
		var meChan chan model.Event
		for {
			if meChan == nil {
				meChanIf, ok := room.Users.Load(me)
				if !ok {
					sugar.Warn("Failed to get interface")
					continue
				}
				meChan, ok = meChanIf.(chan model.Event)
				if !ok {
					sugar.Warn("Failed to convert inf to channel")
					break
				}
			}
			select {
			case msg := <-meChan:
				sugar.Info("Writing message", zap.String("type", msg.Type), zap.String("data", msg.Data))
				ws.WriteJSON(msg)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
