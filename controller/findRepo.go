package controller

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type findRepo interface {
}

type DBRepoStruct struct {
	client *firestore.Client
}

//NewFindRepo return new memory for findrepostruct struct
func NewDBRepo(client *firestore.Client) *DBRepoStruct {
	return &DBRepoStruct{client}
}

//FindOneByField will get document in collection filtered by field and return one document
func (x *DBRepoStruct) FindOneByField(collection string, field string, value string) (map[string]interface{}, error) {
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

	if dataFound == nil {
		return nil, errors.New("NO RESULT")
	}

	return dataFound, nil
}

//FindOneByID return one full document which query is filtered by doc ref / ID
func (x *DBRepoStruct) FindOneByID(collection string, ID string) (map[string]interface{}, error) {
	ctx := context.Background()
	result, err := x.client.Collection(collection).Doc(ID).Get(ctx)
	if err != nil {
		return map[string]interface{}{}, err
	}

	if result == nil {
		return nil, errors.New("NO RESULT")
	}

	return result.Data(), nil
}

//FindAllByField return all result of filtered by field query
func (x *DBRepoStruct) FindAllByField(collection string, field string, value string) ([]map[string]interface{}, error) {
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
	if result == nil {
		return nil, errors.New("NO RESULT")
	}
	return result, nil
}

//FindAllSubColByID will find all doc in subcollection filtered by id
func (x *DBRepoStruct) FindAllSubColByID(collectionOne string, ID string, collectionTwo string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	iter := x.client.Collection(collectionOne).Doc(ID).Collection(collectionTwo).Documents(ctx)
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
	if result == nil {
		return nil, errors.New("NO RESULT")
	}
	return result, nil
}

//FindAllSubColByIDField will find all doc in subcollection filtered by field
func (x *DBRepoStruct) FindAllSubColByIDField(collectionOne string, ID string, field string, value string, collectionTwo string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	iter := x.client.Collection(collectionOne).Doc(ID).Collection(collectionTwo).Where(field, "==", value).Documents(ctx)
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
	if result == nil {
		return nil, errors.New("NO RESULT")
	}
	return result, nil
}
