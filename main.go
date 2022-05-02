package main

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/manfromth3m0oN/csgo/model"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var rooms map[string]*model.Room
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	rooms = make(map[string]*model.Room)
	r := gin.New()
	conf := zap.NewDevelopmentEncoderConfig()
	conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(CORSMiddleware())
	r.GET("/ping", ping)
	r.GET("/room/:room", handleRoom)
	r.POST("/create-room", createRoom)
	r.Run(":3000")
}
