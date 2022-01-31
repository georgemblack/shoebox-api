package shoebox

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"google.golang.org/genproto/googleapis/type/latlng"
)

func decodeTextItem(raw map[string]interface{}) (TextItem, error) {
	expectedKeys := []string{"type", "text"}

	var metadata mapstructure.Metadata
	var output TextItem

	config := &mapstructure.DecoderConfig{
		Metadata: &metadata,
		Result:   &output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return output, fmt.Errorf("failed to create decode config for text item; %w", err)
	}

	err = decoder.Decode(raw)
	if err != nil {
		return output, fmt.Errorf("failed to decode text item; %w", err)
	}

	if !stringSlicesEqual(metadata.Keys, expectedKeys) {
		return output, fmt.Errorf("missing keys %v in %v", expectedKeys, metadata.Keys)
	}

	return output, nil
}

func decodeGeoPointItem(raw map[string]interface{}) (GeoPointItem, error) {
	expectedKeys := []string{"type", "geopoint", "geopoint.latitude", "geopoint.longitude"}

	var metadata mapstructure.Metadata
	var output GeoPointItem

	config := &mapstructure.DecoderConfig{
		Metadata: &metadata,
		Result:   &output,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return output, fmt.Errorf("failed to create decode config for geopoint item; %w", err)
	}

	err = decoder.Decode(raw)
	if err != nil {
		return output, fmt.Errorf("failed to decode geopoint item; %w", err)
	}

	if !stringSlicesEqual(metadata.Keys, expectedKeys) {
		return output, fmt.Errorf("missing keys %v in %v", expectedKeys, metadata.Keys)
	}

	return output, nil
}

func decodeGeoPointItemFromFirestore(raw map[string]interface{}) (GeoPointItem, error) {
	output := GeoPointItem{}

	geoPoint, ok := raw["geopoint"]
	if !ok {
		return GeoPointItem{}, errors.New("missing 'geopoint' field in firestore document")
	}

	parsedGeoPoint, ok := geoPoint.(*latlng.LatLng)
	if !ok {
		return GeoPointItem{}, errors.New("failed to parse 'geopoint' field as latlng.LatLng in firestore document")
	}

	output.Type = "geopoint"
	output.GeoPoint = parsedGeoPoint
	return output, nil
}
