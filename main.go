package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const hereAPIKey = "g-gI_EzetmmNF8WDJNLXRu-hYIqtoRj8OiGtGZADXeM"

func main() {
	lat := "37.7749"   // Replace with the desired latitude
	lng := "-122.4194" // Replace with the desired longitude
	limit := "5"       // Replace with the desired limit

	evChargePoints, err := getEVChargePoints(lat, lng, limit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("EV Charging Stations:")
	for _, station := range evChargePoints.Items {
		fmt.Printf("Name: %s, Latitude: %f, Longitude: %f\n", station.Title, station.Address.Position.Lat, station.Address.Position.Lng)
	}
}

type EVChargePointResponse struct {
	Items []struct {
		Title   string `json:"title"`
		Address struct {
			Label    string `json:"label"`
			Country  string `json:"country"`
			State    string `json:"state"`
			County   string `json:"county"`
			City     string `json:"city"`
			District string `json:"district"`
			Street   string `json:"street"`
			Postal   string `json:"postalCode"`
			Position struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"position"`
		} `json:"address"`
	} `json:"items"`
}

func getEVChargePoints(lat, lng, limit string) (*EVChargePointResponse, error) {
	baseURL := "https://discover.search.hereapi.com/v1/discover"
	apiURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("apiKey", hereAPIKey)
	params.Set("q", "EV Charging Station")
	params.Set("at", fmt.Sprintf("%s,%s", lat, lng))
	params.Set("limit", limit)

	apiURL.RawQuery = params.Encode()

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var evChargePoints EVChargePointResponse
	err = json.Unmarshal(body, &evChargePoints)
	if err != nil {
		return nil, err
	}

	return &evChargePoints, nil
}
