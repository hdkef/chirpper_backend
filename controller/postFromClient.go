package controller

import (
	"context"
	"fmt"
)

//postFromClient is a function to create a new post
func postFromClient(payload MsgPayload) {
	fmt.Println("postFromClient")
	go broadcastNewFeed(payload)
	//post chirp and ref it to user's feed
	//post ref to user's followers
	//send update to user's followers that online
}

//brodcastNewFeed broadcasts new post to online followers
func broadcastNewFeed(payload MsgPayload) {

	fmt.Println("broadcastNewFeed")

	db := NewFindRepo(payload.Client)
	result, err := db.FindOneByID("users", payload.ID)
	fmt.Println(result)
	if err != nil {
		return
	}
	followers := result["followers"].([]interface{})

	var folIDChan chan string = make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	go sendUpdate(followers, cancel, folIDChan)
	go receiveUpdate(ctx, folIDChan, payload)

}

//sendUpdate will send follower ID to channel
func sendUpdate(followers []interface{}, cancel context.CancelFunc, folIDChan chan string) {
	defer cancel()
	for _, v := range followers {
		fmt.Println("send fol ID to channel", v.(string))
		folIDChan <- v.(string)
	}
}

//receiveUpdate receive follower ID and send payload to correspondingID
func receiveUpdate(ctx context.Context, folIDChan chan string, payload MsgPayload) {
	defer func() {
		fmt.Println("channel closed")
		close(folIDChan)
	}()
	for {
		select {
		case ID := <-folIDChan:
			sendPayload(ID, payload)
		case <-ctx.Done():
			return
		}
	}
}

//sendPayload will send payload to online user
func sendPayload(folID string, payload MsgPayload) {
	ws, res := onlineMap[folID]
	if res != true || ws == payload.Conn {
		return
	}
	fmt.Println("found online user, writing json..")
	ws.WriteJSON(payload)
}
