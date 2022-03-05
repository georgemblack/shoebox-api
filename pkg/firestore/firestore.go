package firestore

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

var client *firestore.Client

func GetFirestoreClient() *firestore.Client {
	if client != nil {
		return client
	}
	firestoreClient, err := firestore.NewClient(context.Background(), "shoeboxweb")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize firestore client; %w", err))
	}
	client = firestoreClient
	return client
}

func DeleteEntry(ID string) error {
	ctx := context.Background()
	client := GetFirestoreClient()
	doc, err := client.Doc("shoebox-entries/" + ID).Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve document; %w", err)
	}
	_, err = doc.Ref.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document; %w", err)
	}
	return nil
}
