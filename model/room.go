package model

type User struct {
	Name string
	Addr string
}

type Room struct {
	Name     string
	Users    []User
	Playlist []string
	Index    int
}
