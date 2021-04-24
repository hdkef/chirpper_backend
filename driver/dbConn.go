package driver

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func ConnectDB() *firestore.Client {
	ctx := context.Background()
	opt := option.WithCredentialsFile("fskey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Println("error initializing firebase firestore")
		panic(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
