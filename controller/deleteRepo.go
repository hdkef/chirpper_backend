package controller

import (
	"context"

	"cloud.google.com/go/firestore"
)

type deleteRepo interface {
}

type DeleteRepoStruct struct {
	client *firestore.Client
}

//NewDeleteRepo return new memory of DeleteRepoStruct struct
func NewDeleteRepo(client *firestore.Client) *DeleteRepoStruct {
	return &DeleteRepoStruct{client}
}

//DeleteByID is to delete one document according to document's ref / ID
func (x *DeleteRepoStruct) DeleteByID(collection string, id string) error {
	ctx := context.Background()
	_, err := x.client.Collection(collection).Doc(id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
