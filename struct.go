package shoebox

import (
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// Entries represents a list of entries
type Entries struct {
	Entries []Entry
}

// Entry represents a single entry
type Entry struct {
	ID        string        `json:"id"`
	Content   []interface{} `json:"content" firestore:"content"`
	Published time.Time     `json:"published" firestore:"published"`
	Updated   time.Time     `json:"updated" firestore:"updated"`
}

// TextItem represents text
type TextItem struct {
	Type    string `json:"type" firestore:"type"`
	Content string `json:"text" firestore:"text"`
}

// GeoPointItem represents a latitude/longitude coordinate
type GeoPointItem struct {
	Type     string         `json:"type" firestore:"type"`
	GeoPoint *latlng.LatLng `json:"geopoint" firestore:"geopoint"`
}
