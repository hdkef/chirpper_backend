package controller

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type findRepo interface {
}

type FindRepoStruct struct {
	client *firestore.Client
}

func NewFindRepo(client *firestore.Client) *FindRepoStruct {
	return &FindRepoStruct{client}
}

func (x *FindRepoStruct) FindOneByUsername(collection string, username string) (map[string]interface{}, error) {
	ctx := context.Background()
	iter := x.client.Collection("users").Where("Username", "==", username).Documents(ctx)
	var usernameFound map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return map[string]interface{}{}, err
		}
		usernameFound = doc.Data()
		usernameFound["ID"] = doc.Ref.ID
		break
	}

	return usernameFound, nil
}

func (x *FindRepoStruct) FindAll() {

}
