package controller

import (
	"context"
)

//DeleteByID is to delete one document according to document's ref / ID
func (x *DBRepoStruct) DeleteByID(collection string, id string) error {
	ctx := context.Background()
	_, err := x.client.Collection(collection).Doc(id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
