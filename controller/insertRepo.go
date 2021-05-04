package controller

import (
	"context"
)

//InsertOne insert one document to collection
func (x *DBRepoStruct) InsertOne(collection string, payload map[string]interface{}) (string, error) {
	ctx := context.Background()
	ref, _, err := x.client.Collection(collection).Add(ctx, payload)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

//InsertOneSubColByID will insert doc into subcollection filtered by col doc ID
func (x *DBRepoStruct) InsertOneSubColByID(collectionOne string, ID string, collectionTwo string, payload map[string]interface{}) (string, error) {
	ctx := context.Background()
	ref, _, err := x.client.Collection(collectionOne).Doc(ID).Collection(collectionTwo).Add(ctx, payload)
	if err != nil {
		return "", err
	}
	return ref.ID, nil
}

//InsertOneSubColByIDID will insert doc into subcollection filtered by col doc ID and subcol doc ID
func (x *DBRepoStruct) InsertOneSubColByIDID(collectionOne string, ID string, collectionTwo string, DocID string, payload map[string]interface{}) error {
	ctx := context.Background()
	_, err := x.client.Collection(collectionOne).Doc(ID).Collection(collectionTwo).Doc(DocID).Set(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

func (x *DBRepoStruct) InsertAll() {

}
