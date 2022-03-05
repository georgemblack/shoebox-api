package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/georgemblack/shoebox/pkg/config"
)

type Datastore interface {
	DeleteEntry(ID string) error
}

type datastoreClient struct {
	realClient     *firestore.Client
	collectionName string
}

func GetDatastoreClient(config config.Config) (Datastore, error) {
	var client datastoreClient

	// Initialize Firestore client
	firestoreClient, err := firestore.NewClient(context.Background(), config.FirestoreProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client; %w", err)
	}
	client.realClient = firestoreClient

	// Set config values
	client.collectionName = config.FirestoreCollectionName

	return client, nil
}

func (client datastoreClient) DeleteEntry(ID string) error {
	ctx := context.Background()
	doc, err := client.realClient.Doc(client.collectionName + "/" + ID).Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve document; %w", err)
	}
	_, err = doc.Ref.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document; %w", err)
	}
	return nil
}
