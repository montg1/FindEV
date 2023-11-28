package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const hereAPIKey = "g-gI_EzetmmNF8WDJNLXRu-hYIqtoRj8OiGtGZADXeM"

var laLong [][]string

func main() {
	laLong1 := readCSV()
	for i := 0; i < len(laLong1); i++ {
		lat := laLong1[i][0]
		lng := laLong1[i][1]
		limit := "1" // Replace with the desired limit

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

func readCSV() [][]string {
	file, err := os.Open("Wales.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		item := strings.Split(line, ",")
		listItem := item[9:11]
		laLong = append(laLong, listItem)
	}
	return laLong
}
