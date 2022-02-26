package shoebox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
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

func parseFirestoreDoc(ID string, raw map[string]interface{}) (Entry, error) {
	entry := Entry{}
	var content []interface{}

	// Parse published timestamp
	_, ok := raw["published"]
	if !ok {
		return entry, errors.New("missing 'published' field in firestore document")
	}
	published, ok := raw["published"].(time.Time)
	if !ok {
		return entry, errors.New("failed to parse 'published' field as time.Time in firestore document")
	}

	// Parse updated timestamp
	_, ok = raw["updated"]
	if !ok {
		return entry, errors.New("missing 'updated' field in firestore document")
	}
	updated, ok := raw["updated"].(time.Time)
	if !ok {
		return entry, errors.New("failed to parse 'updated' field as time.Time in firestore document")
	}

	// Parse content blocks
	_, ok = raw["content"]
	if !ok {
		return entry, errors.New("missing 'content' field in firestore document")
	}
	parsedContent, ok := raw["content"].([]interface{})
	if !ok {
		return entry, errors.New("'content' field must be an array")
	}

	for _, rawItem := range parsedContent {
		parsedItem, ok := rawItem.(map[string]interface{})
		if !ok {
			return entry, errors.New("item in 'content' must be a dictionary")
		}

		_, ok = parsedItem["type"]
		if !ok {
			return entry, errors.New("item in 'content' missing a 'type'")
		}

		itemType, ok := parsedItem["type"].(string)
		if !ok {
			return entry, errors.New("item 'type' must be a string")
		}

		switch itemType {
		case "text":
			output, err := decodeTextItem(parsedItem)
			if err != nil {
				return entry, fmt.Errorf("failed to decode text item; %w", err)
			}
			content = append(content, output)
		case "geopoint":
			output, err := decodeGeoPointItemFromFirestore(parsedItem)
			if err != nil {
				return entry, fmt.Errorf("failed to decode geopoint item; %w", err)
			}
			content = append(content, output)
		default:
			return entry, errors.New("failed to identify item type")
		}
	}

	entry.ID = ID
	entry.Published = published
	entry.Updated = updated
	entry.Content = content
	return entry, nil
}

func GetEntries() (Entries, error) {
	entries := Entries{}
	entries.Entries = make([]Entry, 0)

	collectionRef := client.Collection("shoebox-entries")
	iter := collectionRef.Documents(context.Background())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Bad document: ", err)
			return Entries{}, fmt.Errorf("failed to iterate over firestore documents; %w", err)
		}

		entry, err := parseFirestoreDoc(doc.Ref.ID, doc.Data())
		if err != nil {
			fmt.Println(err)
		}
		entries.Entries = append(entries.Entries, entry)
	}

	return entries, nil
}
