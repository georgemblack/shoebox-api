package types

import (
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

type Entry struct {
	ID       string         `json:"id" firestore:"-"`
	Text     string         `json:"text" firestore:"text,omitempty"`
	GeoPoint *latlng.LatLng `json:"geopoint" firestore:"geopoint,omitempty"`
	Created  time.Time      `json:"created" firestore:"created"`
	Updated  time.Time      `json:"updated" firestore:"updated"`
}
