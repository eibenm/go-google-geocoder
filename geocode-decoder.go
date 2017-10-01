package main

import (
	"bytes"
	"encoding/json"
)

type geocodeResponse struct {
	Results []geocodeResult `json:"results"`
	Status  string          `json:"status"`
}

type geocodeResult struct {
	AddressComponents interface{} `json:"address_components"`
	FormattedAddress  string      `json:"formatted_address"`
	Geometry          geometry    `json:"geometry"`
	PlaceID           string      `json:"place_id"`
	Types             interface{} `json:"types"`
}

type geometry struct {
	Location     map[string]float64 `json:"location"`
	LocationType string             `json:"location_type"`
	ViewPort     interface{}        `json:"viewport"`
}

func decodeResponse(b []byte) (*geocodeResponse, error) {
	var err error
	data := &geocodeResponse{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	err = decoder.Decode(&data)
	return data, err
}
