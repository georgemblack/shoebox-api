package shoebox

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
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
