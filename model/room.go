package model

import (
	"sync"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/syncmap"
)

type Event struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

type EventType string

const Join = "join"
const Leave = "leave"
const SkipF = "skipf" // Forward
const SkipB = "skipb" // Back
const Unknown = "unknown"
const Pause = "pause"
const Play = "play"
const Seek = "seek"
const Vid = "vid"

type Room struct {
	Mutex    sync.Mutex
	Name     string
	Users    syncmap.Map
	Playlist []string
	Index    int
	InChan   chan Event
}

func (r *Room) Run() {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	sugar := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	sugar.Info("Running room")
	for {
		select {
		case msg := <-r.InChan:
			sugar.Info("message", zap.String("type", msg.Type), zap.String("data", msg.Data))
			switch msg.Type {
			case Pause, Play, Seek:
				break
			case Join:
				var msgs []Event
				r.Users.Range(func(key, value interface{}) bool {
					if un, ok := key.(string); ok {
						msgs = append(msgs, Event{
							Type: Join,
							Data: un,
						})
					}
					return true
				})
				userchan := make(chan Event)
				r.Users.Store(msg.Data, userchan)
				sugar.Info("Stored user")

				msgs = append(msgs, Event{
					Type: "vid",
					Data: r.Playlist[r.Index],
				})
				for _, m := range msgs {
					userchan <- m
				}
				sugar.Info("Sent vid and users to new user")

				msg = Event{
					Type: Join,
					Data: msg.Data,
				}
				sugar.Info("Created event")

				sugar.Info("Sending join to all others")
				break
			case Leave:
				r.Users.Delete(msg.Data)
				break
			case SkipB:
				r.Index--
				msg = Event{
					Type: Vid,
					Data: r.Playlist[r.Index],
				}
				break
			case SkipF:
				r.Index++
				msg = Event{
					Type: Vid,
					Data: r.Playlist[r.Index],
				}
				break
			default:
				sugar.Info("got unknown message")
				msg = Event{Type: Unknown}
			}
			r.Users.Range(func(key, value interface{}) bool {
				c, ok := value.(chan Event)
				if !ok {
					return false
				}

				sugar.Info("Sending message for ", zap.String("user", key.(string)))
				c <- msg

				return true
			})
		default:

		}
	}

}
