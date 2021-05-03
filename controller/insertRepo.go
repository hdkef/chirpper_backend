package controller

import (
	"context"
)

//InsertOne insert one document to collection
func (x *DBRepoStruct) InsertOne(collection string, payload map[string]interface{}) (string, error) {
	ctx := context.Background()
	_, _, err := x.client.Collection(collection).Add(ctx, payload)
	if err != nil {
		return "", err
	}
	return x.client.Collection(collection).NewDoc().ID, nil
}

//InsertOneSubColByID will insert doc into subcollection filtered by doc ID
func (x *DBRepoStruct) InsertOneSubColByID(collectionOne string, ID string, collectionTwo string, payload map[string]interface{}) (string, error) {
	ctx := context.Background()
	_, _, err := x.client.Collection(collectionOne).Doc(ID).Collection(collectionTwo).Add(ctx, payload)
	if err != nil {
		return "", err
	}
	return x.client.Collection(collectionOne).Doc(ID).Collection(collectionTwo).NewDoc().ID, nil
}

func (x *DBRepoStruct) InsertAll() {

}
