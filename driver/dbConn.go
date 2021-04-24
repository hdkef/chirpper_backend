package driver

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func ConnectDB() {
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
	defer client.Close()

	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}
}
