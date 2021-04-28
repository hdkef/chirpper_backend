package driver

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

//ConnectDB initate connection to firebase and return connection (the firestore.Client)
func ConnectDB() *firestore.Client {
	ctx := context.Background()
	opt := option.WithCredentialsFile("fskey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
