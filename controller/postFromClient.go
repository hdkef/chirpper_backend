package controller

import (
	"chirpper_backend/models"
	"context"
	"fmt"
)

//postFromClient is a function to create a new post
func postFromClient(payload models.MsgPayload) {
	fmt.Println("postFromClient")

	go postFeed(payload)
	db := NewFindRepo(payload.Client)
	result, err := db.FindOneByID("users", payload.ID)
	fmt.Println(result)
	if err != nil {
		return
	}
	followers := result["followers"].([]interface{})
	go broadcastPostFeed(payload, followers)
	go afterPostFeed(payload, followers)
	//post chirp and ref it to user's feed
	//post ref to user's followers
	//send update to user's followers that online
}

//postFeed is post feed to user's document
func postFeed(payload models.MsgPayload) {
	fmt.Println("postFeed")
	postFeedByID(payload.ID, payload)
}

//postFeedByID is intended for afterPostFeed alias insert new feed to user's followers's feed
func postFeedByID(ID string, payload models.MsgPayload) {
	fmt.Println("postFeedByID")
}

//afterPostFeed is update every user's followers feed
func afterPostFeed(payload models.MsgPayload, followers []interface{}) {

	fmt.Println("afterPostFeed")

	var folIDChan chan string = make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	go folIDSender(followers, cancel, folIDChan)
	go folIDAfterPostHandler(ctx, folIDChan, payload)
}

//brodcastNewFeed broadcasts new post to online followers
func broadcastPostFeed(payload models.MsgPayload, followers []interface{}) {

	fmt.Println("broadcastPostFeed")

	var folIDChan chan string = make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	go folIDSender(followers, cancel, folIDChan)
	go folIDBroadcastPostHandler(ctx, folIDChan, payload)

}

//folIDSender will send follower ID to channel
func folIDSender(followers []interface{}, cancel context.CancelFunc, folIDChan chan string) {
	defer cancel()
	for _, v := range followers {
		fmt.Println("send fol ID to channel", v.(string))
		folIDChan <- v.(string)
	}
}

//receiveUpdate receive follower ID and send payload to correspondingID
func folIDBroadcastPostHandler(ctx context.Context, folIDChan chan string, payload models.MsgPayload) {
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

func folIDAfterPostHandler(ctx context.Context, folIDChan chan string, payload models.MsgPayload) {
	defer func() {
		fmt.Println("channel closed")
		close(folIDChan)
	}()
	for {
		select {
		case ID := <-folIDChan:
			postFeedByID(ID, payload)
		case <-ctx.Done():
			return
		}
	}
}

//sendPayload will send payload to online user
func sendPayload(folID string, payload models.MsgPayload) {
	ws, res := onlineMap[folID]
	if res != true || ws == payload.Conn {
		return
	}
	fmt.Println("found online user, writing json..")
	ws.WriteJSON(payload)
}
