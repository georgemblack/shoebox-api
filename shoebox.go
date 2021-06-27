package shoebox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

var client *firestore.Client

func Init() error {
	if client != nil {
		return nil
	}
	firestoreClient, err := firestore.NewClient(context.Background(), "shoeboxweb")
	if err != nil {
		return fmt.Errorf("failed to initialize firestore client; %w", err)
	}
	client = firestoreClient
	return nil
}

func ParseEntry(raw map[string]interface{}) (Entry, error) {
	entry := Entry{}

	var content []interface{}
	rawContent, ok := raw["content"]
	if !ok {
		return entry, errors.New("missing 'content' in entry")
	}
	parsedContent, ok := rawContent.([]interface{})
	if !ok {
		return entry, errors.New("'content' must be an array")
	}

	for _, rawItem := range parsedContent {
		parsedItem, ok := rawItem.(map[string]interface{})
		if !ok {
			return entry, errors.New("item in 'content' must be a dictionary")
		}

		itemType, ok := parsedItem["type"]
		if !ok {
			return entry, errors.New("item in 'content' missing a 'type'")
		}
		itemTypeStr, ok := itemType.(string)
		if !ok {
			return entry, errors.New("item 'type' must be a string")
		}

		switch itemTypeStr {
		case "text":
			output, err := decodeTextItem(parsedItem)
			if err != nil {
				return entry, fmt.Errorf("failed to decode text item; %w", err)
			}
			content = append(content, output)
		case "geopoint":
			output, err := decodeGeoPointItem(parsedItem)
			if err != nil {
				return entry, fmt.Errorf("failed to decode geopoint item; %w", err)
			}
			content = append(content, output)
		default:
			return entry, errors.New("failed to identify item type")
		}
	}

	entry.Published = time.Now()
	entry.Updated = time.Now()
	entry.Content = content
	return entry, nil
}

func CreateEntry(entry Entry) error {
	reference := client.Doc("shoebox-entries/" + uuid.NewString())
	_, err := reference.Create(context.Background(), entry)
	if err != nil {
		return fmt.Errorf("failed to initialize create firestore document; %w", err)
	}
	return nil
}
