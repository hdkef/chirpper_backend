package controller

import (
	"chirpper_backend/models"
	"time"

	"github.com/gorilla/websocket"
)

//initFromClient handle if ws conn has been established
func initFromClient(payload models.MsgPayload) {

	valid := verifyTokenString(payload.Bearer)

	if valid == false {
		payload.Conn.Close()
		return
	}

	onlineMap[payload.ID] = payload.Conn
	go pingPonger(payload)
	go initFromServer(payload)
	//implement sinitFromServer
}

func initFromServer(payload models.MsgPayload) {
	db := NewDBRepo(payload.Client)
	result, err := db.FindAllSubColByID("users", payload.ID, "feed")
	if err != nil {
		return
	}

	payloadToBeSent := struct {
		Type string
		Data []map[string]interface{}
	}{
		Type: "initFromServer",
		Data: result,
	}
	payload.Conn.WriteJSON(payloadToBeSent)
}

//ping the client, if offline then delete id in onlineMap
func pingPonger(payload models.MsgPayload) {

	payload.Conn.SetPongHandler(func(string) error {
		payload.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	timer := time.NewTicker(pingPeriod)
	defer func() {
		timer.Stop()
		if onlineMap[payload.ID] == payload.Conn {
			delete(onlineMap, payload.ID)
		}
	}()
	for {
		select {
		case <-timer.C:
			if err := payload.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
