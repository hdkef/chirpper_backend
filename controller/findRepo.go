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

func (x *FindRepoStruct) FindOneByField(collection string, field string, value string) (map[string]interface{}, error) {
	ctx := context.Background()
	iter := x.client.Collection(collection).Where(field, "==", value).Documents(ctx)
	var dataFound map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return map[string]interface{}{}, err
		}
		dataFound = doc.Data()
		dataFound["ID"] = doc.Ref.ID
		break
	}

	return dataFound, nil
}

func (x *FindRepoStruct) FindAll() {

}
