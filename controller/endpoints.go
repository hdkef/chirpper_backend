package controller

import (
	"chirpper_backend/models"
	"chirpper_backend/utils"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/websocket"
)

type EndPoints struct {
}

//const for ping pong time
const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

//this map contains ID as key and websocket conn as value, use to know who online
var onlineMap map[string]*websocket.Conn = make(map[string]*websocket.Conn)

//define all channels needed so concurrent can happen
//channel for posting new chirp
var postFromClientChan chan models.MsgPayload = make(chan models.MsgPayload)
var initFromClientChan chan models.MsgPayload = make(chan models.MsgPayload)

//to upgrade protocol
var upgrader websocket.Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//CheckToken the first thing after logging in, to check token
func (x *EndPoints) CheckToken(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		valid := verifyToken(req)
		if valid != true {
			utils.ResClearSite(&res)
			utils.ResError(res, http.StatusUnauthorized, errors.New("INVALID TOKEN"))
			return
		}

		utils.ResOK(res, "TOKEN VALID")
	}
}

//establishWS is after checked token to establish 2 way comm
func (x *EndPoints) EstablishWS(client *firestore.Client) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		ws, err := upgrader.Upgrade(res, req, res.Header())
		if err != nil {
			log.Println(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		go readMsg(ws, cancel, client)
		go routeMsg(ws, ctx)
	}
}

//readMsg is to read incoming msg from client, each ws.conn be assigned one readMsg goroutine
func readMsg(ws *websocket.Conn, cancel context.CancelFunc, client *firestore.Client) {

	var payload models.MsgPayload = models.MsgPayload{
		Conn:   ws,
		Client: client,
	}

	defer cancel()

	for {
		err := ws.ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			break
		}
		switch payload.Type {
		case "postFromClient":
			postFromClientChan <- payload
		case "initFromClient":
			initFromClientChan <- payload
		}
	}
}

//routeMsg is executed once to route every msg come to channel to corresponding function
func routeMsg(ws *websocket.Conn, ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-postFromClientChan:
			postFromClient(msg)
		case msg := <-initFromClientChan:
			initFromClient(msg)
		}
	}

}
