package controller

import (
	"chirpper_backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

//Setting is an endpoints for setting
func (x *EndPoints) Setting(client *firestore.Client) http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Setting")

		valid := verifyToken(req)
		if valid != true {
			utils.ResClearSite(&res)
			utils.ResError(res, http.StatusUnauthorized, errors.New("INVALID TOKEN"))
			return
		}

		if err := req.ParseMultipartForm(1024); err != nil {
			fmt.Println(err)
			utils.ResError(res, http.StatusInternalServerError, err)
			return
		}

		if _, _, err := req.FormFile("Avatar"); err == nil && req.FormValue("Desc") != "" {
			//Both
			setAvaAndDesc(res, req, client)
		} else if _, _, err := req.FormFile("Avatar"); err == nil && req.FormValue("Desc") == "" {
			//Avatar only
			setAvaOnly(res, req, client)
		} else if _, _, err := req.FormFile("Avatar"); err != nil && req.FormValue("Desc") != "" {
			//Desc only
			setDescOnly(res, req, client)
		} else {
			utils.ResError(res, http.StatusInternalServerError, errors.New("NO PAYLOAD"))
		}
	}
}

//setAvaOnly will store ava and response with AvatarURL only
func setAvaOnly(res http.ResponseWriter, req *http.Request, client *firestore.Client) {

	avaLocation, err := storeImage(res, req, "Avatar", "avatar")
	if err != nil {
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}

	//TOBEIMPLEMENT update AvatarURL
	err = updateAva(req.FormValue("ID"), avaLocation, client)
	if err != nil {
		//Implement error handling
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}

	payloadToBeSent := struct {
		AvatarURL string
	}{
		AvatarURL: avaLocation,
	}

	err = json.NewEncoder(res).Encode(&payloadToBeSent)
	if err != nil {
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}
}

//setDescOnly will set desc and response with Desc only
func setDescOnly(res http.ResponseWriter, req *http.Request, client *firestore.Client) {

	//TOBEIMPLEMENT update Desc
	desc := req.FormValue("Desc")

	err := updateDesc(req.FormValue("ID"), desc, client)
	if err != nil {
		//Implement error handling
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}

	payloadToBeSent := struct {
		Desc string
	}{
		Desc: desc,
	}

	err = json.NewEncoder(res).Encode(&payloadToBeSent)
	if err != nil {
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}
}

//setAvaAndDesc will set desc and store avara and response with Desc and AvatarURL
func setAvaAndDesc(res http.ResponseWriter, req *http.Request, client *firestore.Client) {

	desc := req.FormValue("Desc")

	avaLocation, err := storeImage(res, req, "Avatar", "avatar")
	if err != nil {
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}

	err = updateAva(req.FormValue("ID"), avaLocation, client)
	if err != nil {
		//Implement error handling
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}
	err = updateDesc(req.FormValue("ID"), desc, client)
	if err != nil {
		//Implement error handling
		fmt.Println(err)
		desc = ""
	}
	//TOBEIMPLEMENT update Desc and AvatarURL

	payloadToBeSent := struct {
		Desc      string
		AvatarURL string
	}{
		Desc:      desc,
		AvatarURL: avaLocation,
	}

	err = json.NewEncoder(res).Encode(&payloadToBeSent)
	if err != nil {
		fmt.Println(err)
		utils.ResError(res, http.StatusInternalServerError, err)
		return
	}
}

//updateDesc will update Desc firestore
func updateDesc(ID string, desc string, client *firestore.Client) error {
	value := []firestore.Update{
		{
			Path:  "Desc",
			Value: desc,
		},
	}
	err := update("users", ID, value, client)
	if err != nil {
		return err
	}

	return nil
}

//updateAva will update ava firestore
func updateAva(ID string, avatarurl string, client *firestore.Client) error {
	value := []firestore.Update{
		{
			Path:  "AvatarURL",
			Value: avatarurl,
		},
	}
	err := update("users", ID, value, client)
	if err != nil {
		return err
	}

	return nil
}

func update(collection string, ID string, value []firestore.Update, client *firestore.Client) error {
	db := NewDBRepo(client)
	err := db.UpdateOneByID(collection, ID, value)
	if err != nil {
		return err
	}
	return nil
}

//handlingAvaErr will delete avatar in storage and reset ava firestore
// func handlingAvaErr(ID string, client *firestore.Client) {

// }
