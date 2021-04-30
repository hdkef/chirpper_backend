package models

import (
	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

//MsgPayload is to format the request sent from client
type MsgPayload struct {
	Conn      *websocket.Conn
	Client    *firestore.Client
	Mapclaims *jwt.MapClaims
	Type      string
	ID        string
	Username  string
	Email     string
	ImageURL  string
	Text      string
}
