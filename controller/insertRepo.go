package controller

import (
	"context"

	"cloud.google.com/go/firestore"
)

type insertRepo interface {
}

type InsertRepoStruct struct {
	client *firestore.Client
}

//NewInsertRepo return new memory of insertrepostruct struct
func NewInsertRepo(client *firestore.Client) *InsertRepoStruct {
	return &InsertRepoStruct{client}
}

//InsertOne insert one document to collection
func (x *InsertRepoStruct) InsertOne(collection string, payload map[string]interface{}) error {
	ctx := context.Background()
	_, _, err := x.client.Collection(collection).Add(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

func (x *InsertRepoStruct) InsertAll() {

}
