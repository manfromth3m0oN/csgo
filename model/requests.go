package model

type CreateRoomReq struct {
	Name     string `json:"room_name"`
	ThreadID string `json:"thread_id"`
	BoardSn  string `json:"board_sn"`
}

type RoomJoinReq struct {
	Username string `json:"username"`
}
