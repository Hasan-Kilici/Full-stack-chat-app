package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	User User
	Room string
}

type User struct {
	ID    string
	Name  string
}