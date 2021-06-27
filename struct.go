package shoebox

import (
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// Entry represents a single entry
type Entry struct {
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
