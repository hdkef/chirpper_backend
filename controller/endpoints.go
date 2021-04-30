package controller

import (
	"chirpper_backend/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type EndPoints struct {
}

//const for ping pong time
const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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

//this map contains ID as key and websocket conn as value, use to know who online
var onlineMap map[string]*websocket.Conn = make(map[string]*websocket.Conn)

//define all channels needed so concurrent can happen
//channel for posting new chirp
var postFromClientChan chan MsgPayload = make(chan MsgPayload)
var initFromClientChan chan MsgPayload = make(chan MsgPayload)

//to upgrade protocol
var upgrader websocket.Upgrader = websocket.Upgrader{
	// CheckOrigin: verifyToken,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//CheckToken the first thing after logging in, to check token
func (x *EndPoints) CheckToken(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		fmt.Println("CheckToken")

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
		fmt.Println("establishWS()")

		ws, err := upgrader.Upgrade(res, req, res.Header())
		if err != nil {
			log.Println(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		fmt.Println("beginning ", runtime.NumGoroutine())

		go readMsg(ws, cancel, client)
		go routeMsg(ws, ctx)
	}
}

//readMsg is to read incoming msg from client, each ws.conn be assigned one readMsg goroutine
func readMsg(ws *websocket.Conn, cancel context.CancelFunc, client *firestore.Client) {

	fmt.Println("readMsg")

	var payload MsgPayload = MsgPayload{
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

	fmt.Println("routeMsg")

	defer func() {
		fmt.Println("after ctx.Done() ", runtime.NumGoroutine())
	}()

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
