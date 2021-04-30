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

//NewFindRepo return new memory for findrepostruct struct
func NewFindRepo(client *firestore.Client) *FindRepoStruct {
	return &FindRepoStruct{client}
}

//FindOneByField will get document in collection filtered by field and return one document
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

//FindOneByID return one full document which query is filtered by doc ref / ID
func (x *FindRepoStruct) FindOneByID(collection string, ID string) (map[string]interface{}, error) {
	ctx := context.Background()
	result, err := x.client.Collection(collection).Doc(ID).Get(ctx)
	if err != nil {
		return map[string]interface{}{}, err
	}

	return result.Data(), nil
}

//FindAllByField return all result of filtered by field query
func (x *FindRepoStruct) FindAllByField(collection string, field string, value string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	iter := x.client.Collection(collection).Where(field, "==", value).Documents(ctx)
	var result []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []map[string]interface{}{}, err
		}
		result = append(result, doc.Data())
	}
	return result, nil
}
