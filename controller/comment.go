package controller

import (
	"chirpper_backend/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/websocket"
)

//commentFromClientChan is channel for comment payload
var commentFromClientChan chan models.MsgPayload = make(chan models.MsgPayload)

//commentFromClientChan is channel for init comment payload
var initCommentFromClientChan chan models.MsgPayload = make(chan models.MsgPayload)

//postIDMap is map contains postID subscriber
var postIDMap map[string]map[string]*websocket.Conn = make(map[string]map[string]*websocket.Conn)

//establishComment will establish websocket inside comment
func (x *EndPoints) EstablishComment(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		fmt.Println("establishComment()")

		ws, err := upgrader.Upgrade(res, req, res.Header())
		if err != nil {
			log.Println(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		fmt.Println("beginning ", runtime.NumGoroutine())

		go readComment(ws, cancel, client)
		go routeComment(ws, ctx)

	}
}

//initCommentFromClient handle if user went to comment
func initCommentFromClient(payload models.MsgPayload) {

	fmt.Println("initCommentFromClient()")

	// valid := verifyTokenString(payload.Bearer)

	// if valid == false {
	// 	fmt.Println("invalid bearer ", payload.Bearer)
	// 	payload.Conn.Close()
	// 	return
	// }

	postIDMap[payload.PostID] = map[string]*websocket.Conn{
		payload.ID: payload.Conn,
	}

	go commentPingPonger(payload)
	go initCommentFromServer(payload)

}

//commentFromClient is to handle initation comment from client
func initCommentFromServer(payload models.MsgPayload) {

	fmt.Println("initCommentFromServer")

	db := NewDBRepo(payload.Client)
	result, err := db.FindAllSubColByID("chirps", payload.PostID, "comment")
	if err != nil {
		fmt.Println(err)
		return
	}
	payloadToBeSent := struct {
		Type string
		Data []map[string]interface{}
	}{
		Type: "initCommentFromServer",
		Data: result,
	}
	payload.Conn.WriteJSON(payloadToBeSent)
}

//commentFromClient is to handle new comment from client
func commentFromClient(payload models.MsgPayload) {

	// valid := verifyTokenString(payload.Bearer)

	// if valid == false {
	// 	fmt.Println("invalid bearer ", payload.Bearer)
	// 	payload.Conn.Close()
	// 	return
	// }

	fmt.Println("commentFromClient")

	db := NewDBRepo(payload.Client)

	payload.Date = time.Now().Format("02-Jan-2006")

	_, err := db.InsertOneSubColByID("chirps", payload.PostID, "comment", map[string]interface{}{
		"ID":        payload.ID,
		"Username":  payload.Username,
		"Email":     payload.Email,
		"AvatarURL": payload.AvatarURL,
		"Text":      payload.Text,
		"Date":      payload.Date,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	//send to my self the comment
	go func() {
		fmt.Println("writing payload to user..")
		payload.Conn.WriteJSON(struct {
			Type      string
			ID        string
			Username  string
			Text      string
			Date      string
			AvatarURL string
		}{
			Type:      "commentFromServer",
			ID:        payload.ID,
			Username:  payload.Username,
			Text:      payload.Text,
			Date:      payload.Date,
			AvatarURL: payload.AvatarURL,
		})
	}()
	go broadcastComment(payload)
}

//broadcastComment is to broadcast comment to postID subscriber
func broadcastComment(payload models.MsgPayload) {
	for _, v := range postIDMap[payload.PostID] {
		if v != payload.Conn {
			v.WriteJSON(struct {
				Type      string
				ID        string
				Username  string
				Text      string
				Date      string
				AvatarURL string
			}{
				Type:      "commentFromServer",
				ID:        payload.ID,
				Username:  payload.Username,
				Text:      payload.Text,
				Date:      payload.Date,
				AvatarURL: payload.AvatarURL,
			})
		}
	}
}

//readComment is to read incoming comment from client, each ws.conn be assigned one readMsg goroutine
func readComment(ws *websocket.Conn, cancel context.CancelFunc, client *firestore.Client) {

	fmt.Println("readComment")

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
		case "initCommentFromClient":
			initCommentFromClientChan <- payload
		case "commentFromClient":
			commentFromClientChan <- payload
		}
	}
}

//routeComment is executed once to route every msg come to channel to corresponding function
func routeComment(ws *websocket.Conn, ctx context.Context) {

	fmt.Println("routeComment")

	defer func() {
		fmt.Println("after ctx.Done() ", runtime.NumGoroutine())
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-initCommentFromClientChan:
			initCommentFromClient(msg)
		case msg := <-commentFromClientChan:
			commentFromClient(msg)
		}
	}

}

func commentPingPonger(payload models.MsgPayload) {
	fmt.Println("pingPonger")

	payload.Conn.SetPongHandler(func(string) error {
		payload.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	timer := time.NewTicker(pingPeriod)
	defer func() {
		timer.Stop()
		if postIDMap[payload.PostID][payload.ID] == payload.Conn {
			delete(postIDMap[payload.PostID], payload.ID)
		}
		fmt.Println("postID map : ", postIDMap[payload.PostID])
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
