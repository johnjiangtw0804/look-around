package main

import (
	"fmt"
	entity "look-around/external/entity"
	"net/http"
	"net/url"

	"github.com/pierrre/geohash"
)

const (
	discoveryURL             = "https://app.ticketmaster.com/discovery/v2/events.json"
	discoveryDEFAULT_KEYWORD = ""
	discoveryDEFAULT_RADIUS  = 50
)

type events interface {
	listEvents(latitude, longitude float64, keyword string) (entity.Item, error)
}

type ticketMasterEventsSearcher struct {
	apiKey string
}

func NewTicketMasterEventsSearcher(discoveryURL, apiKey string) events {
	return &ticketMasterEventsSearcher{
		apiKey: apiKey,
	}
}

type ticketMasterDiscoveryResponse struct {
	Embedded embedded `json:"_embedded"`
}

type embedded struct {
	Events []event `json:"events"`
}
type event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (t *ticketMasterEventsSearcher) listEvents(latitude, longitude float64, keyword string) (entity.Item, error) {
	if len(keyword) == 0 {
		keyword = discoveryDEFAULT_KEYWORD
	}
	client := &http.Client{}
	// function performs URL encoding in UTF-8, so we don't need to specify the encoding type explicitly.
	encodedKeyword := url.QueryEscape(keyword)
	geoHash := geohash.Encode(latitude, longitude, 8)

	payload := url.Values{}
	payload.Add("apikey", t.apiKey)
	payload.Add("geoPoint", geoHash)
	payload.Add("keyword", encodedKeyword)
	payload.Add("radius", "50")
	req, err := http.NewRequest(http.MethodGet, discoveryURL+"?"+payload.Encode(), nil)
	if err != nil {
		return entity.Item{}, err
	}
	res, _ := client.Do(req)
	// build item from response
	fmt.Println("test api", res)

	return entity.Item{}, nil
}
func main() {
	e := NewTicketMasterEventsSearcher(discoveryURL, "XbJGpDikv93zCALbmNKU6l5NWK26BP1T")
	e.listEvents(51.503364, -95.295410, "")
}
