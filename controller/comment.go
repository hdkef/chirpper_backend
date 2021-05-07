package controller

import (
	"chirpper_backend/models"
	"fmt"
	"time"
)

//commentFromClient is to handle new comment from client
func commentFromClient(payload models.MsgPayload) {

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
			ImageURL  string
			Text      string
			Date      string
			AvatarURL string
		}{
			Type:      "commentFromServer",
			ID:        payload.ID,
			Username:  payload.Username,
			ImageURL:  payload.ImageURL,
			Text:      payload.Text,
			Date:      payload.Date,
			AvatarURL: payload.AvatarURL,
		})
	}()
}
