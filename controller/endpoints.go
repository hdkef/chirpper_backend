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
	CheckOrigin: verifyToken,
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

		go readMsg(ws, cancel)
		go routeMsg(ws, ctx)
	}
}

//postFromClient is a function to create a new post
func postFromClient(payload MsgPayload) {
	fmt.Println("postFromClient")
	go broadcastNewFeed(payload, []string{"ID2"})
	//post chirp and ref it to user's feed
	//post ref to user's followers
	//send update to user's followers that online
}

//initFromClient handle if ws conn has been established
func initFromClient(payload MsgPayload) {
	fmt.Println("initFromClient")
	onlineMap[payload.ID] = payload.Conn
	go pingPonger(payload)
	//implement send feed
}

//brodcastNewFeed broadcasts new post to online followers
func broadcastNewFeed(payload MsgPayload, followers []string) {

	fmt.Println("broadcastNewFeed")

	msg := struct {
		ImageURL string
		Text     string
	}{
		ImageURL: payload.ImageURL,
		Text:     payload.Text,
	}

	for i := 0; i < len(followers); i++ {
		ws, err := onlineMap[followers[i]]
		if err == true && ws != payload.Conn {
			ws.WriteJSON(msg)
		}
	}
}

//readMsg is to read incoming msg from client, each ws.conn be assigned one readMsg goroutine
func readMsg(ws *websocket.Conn, cancel context.CancelFunc) {

	fmt.Println("readMsg")

	var payload MsgPayload = MsgPayload{
		Conn: ws,
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

//ping the client, if offline then delete id in onlineMap
func pingPonger(payload MsgPayload) {

	fmt.Println("pingPonger")

	payload.Conn.SetPongHandler(func(string) error {
		payload.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	timer := time.NewTicker(pingPeriod)
	defer func() {
		timer.Stop()
		delete(onlineMap, payload.ID)
		fmt.Println("online Map ", onlineMap)
	}()
	for {
		select {
		case <-timer.C:
			fmt.Println("timer tick")
			if err := payload.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
