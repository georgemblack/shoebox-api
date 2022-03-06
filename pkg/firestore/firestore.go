package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/georgemblack/shoebox/pkg/config"
	"github.com/georgemblack/shoebox/pkg/types"
	"google.golang.org/api/iterator"
)

type Datastore interface {
	GetEntries() ([]types.Entry, error)
	CreateEntry(entry types.Entry) error
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

func (client datastoreClient) GetEntries() ([]types.Entry, error) {
	ctx := context.Background()

	collection := client.realClient.Collection(client.collectionName)
	iter := collection.Documents(ctx)
	defer iter.Stop()

	entries := make([]types.Entry, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over firestore documents; %w", err)
		}
		var entry types.Entry

		err = doc.DataTo(&entry)
		if err != nil {
			return nil, fmt.Errorf("failed to parse firestore document; %w", err)
		}
		entry.ID = doc.Ref.ID
		entries = append(entries, entry)
	}

	return entries, nil
}

func (client datastoreClient) CreateEntry(entry types.Entry) error {
	ctx := context.Background()

	_, err := client.realClient.Doc(client.collectionName+"/"+entry.ID).Set(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to write document; %w", err)
	}

	return nil
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
