package models

import (
	"cloud.google.com/go/firestore"
	"github.com/gorilla/websocket"
)

//MsgPayload is to format the request sent from client
type MsgPayload struct {
	Conn      *websocket.Conn
	Client    *firestore.Client
	Type      string
	ID        string
	PostID    string
	Username  string
	Email     string
	ImageURL  string
	Text      string
	Bearer    string
	Date      string
	AvatarURL string
}
