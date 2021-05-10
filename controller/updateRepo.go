package controller

import (
	"context"

	"cloud.google.com/go/firestore"
)

func (x *DBRepoStruct) UpdateOneByID(collection string, ID string, value []firestore.Update) error {
	ctx := context.Background()
	_, err := x.client.Collection(collection).Doc(ID).Update(ctx, value)
	if err != nil {
		return err
	}
	return nil
}
