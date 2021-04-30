package controller

import (
	"chirpper_backend/models"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

//initFromClient handle if ws conn has been established
func initFromClient(payload models.MsgPayload) {
	fmt.Println("initFromClient")
	onlineMap[payload.ID] = payload.Conn
	go pingPonger(payload)
	//implement send feed
}

//ping the client, if offline then delete id in onlineMap
func pingPonger(payload models.MsgPayload) {

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
