package controller

import (
	"chirpper_backend/models"
	"chirpper_backend/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/firestore"
)

//postWithImage handle post feed request when storing an image file is needed
func (x *EndPoints) PostWithImage(client *firestore.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			utils.ResOK(res, "OK")
			return
		}()

		valid := verifyToken(req)
		if valid != true {
			utils.ResClearSite(&res)
			utils.ResError(res, http.StatusUnauthorized, errors.New("INVALID TOKEN"))
			return
		}

		if err := req.ParseMultipartForm(1024); err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		var payload models.MsgPayload = models.MsgPayload{
			Conn:      onlineMap[req.FormValue("ID")],
			Type:      req.FormValue("Type"),
			Client:    client,
			ID:        req.FormValue("ID"),
			Username:  req.FormValue("Username"),
			Email:     req.FormValue("Email"),
			AvatarURL: req.FormValue("AvatarURL"),
			Text:      req.FormValue("Text"),
			Bearer:    req.FormValue("Bearer"),
		}

		//implement store image in server and retrieve img directory location

		imgLocation, err := storeImage(res, req, "Image", "post")

		if err != nil {
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		payload.ImageURL = imgLocation

		go postFromClient(payload)
	}
}

//storeImage will store image file and return image file path/dir
func storeImage(res http.ResponseWriter, req *http.Request, formfilename string, foldername string) (string, error) {
	uploadedFile, handler, err := req.FormFile(formfilename)
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filename := handler.Filename
	folderpath := filepath.Join(dir, os.Getenv("STATICPATH"), "assets", foldername)
	fileLocation := filepath.Join(folderpath, filename)

	err = createNewFolderIfNotExist(folderpath)
	if err != nil {
		return "", err
	}

	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return "", err
	}

	return fmt.Sprintf("assets/%s/%s", foldername, filename), nil
}

//createNewFolderIfNotExist will check folder directory and create new folder if not exist
func createNewFolderIfNotExist(path string) error {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

//postFromClient is a function to create a new post
func postFromClient(payload models.MsgPayload) {

	valid := verifyTokenString(payload.Bearer)

	if valid == false {
		payload.Conn.Close()
		return
	}

	payload.Date = time.Now().Format("02-Jan-2006")

	db := NewDBRepo(payload.Client)

	insertedID, err := db.InsertOne("chirps", map[string]interface{}{
		"ID":        payload.ID,
		"Username":  payload.Username,
		"ImageURL":  payload.ImageURL,
		"Text":      payload.Text,
		"Date":      payload.Date,
		"AvatarURL": payload.AvatarURL,
	})

	payload.PostID = insertedID

	go postFeedByIDID(payload.ID, payload)

	result, err := db.FindOneByID("users", payload.ID)
	if err != nil {
		return
	}
	followers := result["followers"].([]interface{})

	// payload.PostID = insertedID
	go sendSelfPayload(payload)

	if followers == nil {
		return
	}

	go broadcastPostFeed(payload, followers)
	go afterPostFeed(payload, followers)
	//post chirp and ref it to user's feed
	//post ref to user's followers
	//send update to user's followers that online
}

//postFeedByID is intended for afterPostFeed alias insert new feed to user's followers's feed
func postFeedByIDID(ID string, payload models.MsgPayload) {

	db := NewDBRepo(payload.Client)
	err := db.InsertOneSubColByIDID("users", ID, "feed", payload.PostID, map[string]interface{}{
		"ID":        payload.ID,
		"PostID":    payload.PostID,
		"Username":  payload.Username,
		"ImageURL":  payload.ImageURL,
		"Text":      payload.Text,
		"Date":      payload.Date,
		"AvatarURL": payload.AvatarURL,
	})
	if err != nil {
		return
	}
}

//afterPostFeed is update every user's followers feed
func afterPostFeed(payload models.MsgPayload, followers []interface{}) {

	var folIDChan chan string = make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	go folIDSender(followers, cancel, folIDChan)
	go folIDAfterPostHandler(ctx, folIDChan, payload)
}

//brodcastNewFeed broadcasts new post to online followers
func broadcastPostFeed(payload models.MsgPayload, followers []interface{}) {

	var folIDChan chan string = make(chan string)

	ctx, cancel := context.WithCancel(context.Background())

	go folIDSender(followers, cancel, folIDChan)
	go folIDBroadcastPostHandler(ctx, folIDChan, payload)

}

//folIDSender will send follower ID to channel
func folIDSender(followers []interface{}, cancel context.CancelFunc, folIDChan chan string) {
	defer cancel()
	for _, v := range followers {
		folIDChan <- v.(string)
	}
}

//receiveUpdate receive follower ID and send payload to correspondingID
func folIDBroadcastPostHandler(ctx context.Context, folIDChan chan string, payload models.MsgPayload) {
	defer func() {
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
		close(folIDChan)
	}()
	for {
		select {
		case ID := <-folIDChan:
			postFeedByIDID(ID, payload)
		case <-ctx.Done():
			return
		}
	}
}

//sendPayload will send payload to online user
func sendPayload(folID string, payload models.MsgPayload) {
	ws, res := onlineMap[folID]
	if res != true {
		return
	}
	ws.WriteJSON(struct {
		Type      string
		ID        string
		PostID    string
		Username  string
		ImageURL  string
		Text      string
		Date      string
		AvatarURL string
	}{
		Type:      "postFromServer",
		ID:        payload.ID,
		PostID:    payload.PostID,
		Username:  payload.Username,
		ImageURL:  payload.ImageURL,
		Text:      payload.Text,
		Date:      payload.Date,
		AvatarURL: payload.AvatarURL,
	})
}

func sendSelfPayload(payload models.MsgPayload) {
	payload.Conn.WriteJSON(struct {
		Type      string
		ID        string
		PostID    string
		Username  string
		ImageURL  string
		Text      string
		Date      string
		AvatarURL string
	}{
		Type:      "postFromServer",
		ID:        payload.ID,
		PostID:    payload.PostID,
		Username:  payload.Username,
		ImageURL:  payload.ImageURL,
		Text:      payload.Text,
		Date:      payload.Date,
		AvatarURL: payload.AvatarURL,
	})
}
