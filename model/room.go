package model

import (
	"sync"

	"github.com/manfromth3m0oN/csgo/util"
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
	Users    []string
	Playlist []string
	Index    int
	InChan   chan Event
	OutChan  chan Event
}

func (r *Room) Run() {
	for {
		msg := <-r.InChan
		switch msg.Type {
		case Pause, Play, Seek:
			break
		case Join:
			r.Users = append(r.Users, msg.Data)
			break
		case Leave:
			r.Users = util.RemoveFromSlice(r.Users, msg.Data)
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
			msg = Event{Type: Unknown}
		}
		r.OutChan <- msg
	}
}

func (r *Room) Join(username string) {
	r.Users = append(r.Users, username)
	msg := Event{
		Type: "join",
		Data: username,
	}
	r.OutChan <- msg
}
