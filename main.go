package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const hereAPIKey = "g-gI_EzetmmNF8WDJNLXRu-hYIqtoRj8OiGtGZADXeM"

var laLong [][]string

func main() {
	laLong1, err := readCSV()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

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
			fmt.Printf("Name: %s, Latitude: %f, Longitude: %f\n", station.Title, station.Position.Lat, station.Position.Lng)
		}
	}
}

type EVChargingStation struct {
	Items []struct {
		Title      string `json:"title"`
		ID         string `json:"id"`
		Language   string `json:"language"`
		OntologyID string `json:"ontologyId"`
		ResultType string `json:"resultType"`
		Address    struct {
			Label       string `json:"label"`
			CountryCode string `json:"countryCode"`
			CountryName string `json:"countryName"`
			State       string `json:"state"`
			CountyCode  string `json:"countyCode"`
			County      string `json:"county"`
			City        string `json:"city"`
			District    string `json:"district"`
			Street      string `json:"street"`
			PostalCode  string `json:"postalCode"`
		} `json:"address"`
		Position struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"position"`
		Access []struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"access"`
		Distance   float64 `json:"distance"`
		Categories []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Primary bool   `json:"primary"`
		} `json:"categories"`
		Contacts []struct {
			Phone []struct {
				Value string `json:"value"`
			} `json:"phone"`
			WWW []struct {
				Value string `json:"value"`
			} `json:"www"`
		} `json:"contacts"`
		OpeningHours []struct {
			Text       []string `json:"text"`
			IsOpen     bool     `json:"isOpen"`
			Structured []struct {
				Start      string `json:"start"`
				Duration   string `json:"duration"`
				Recurrence string `json:"recurrence"`
			} `json:"structured"`
		} `json:"openingHours"`
	} `json:"items"`
}

func getEVChargePoints(lat, lng, limit string) (*EVChargingStation, error) {
	baseURL := "https://discover.search.hereapi.com/v1/discover"
	apiURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %w", err)
	}

	params := url.Values{}
	params.Set("apiKey", hereAPIKey)
	params.Set("q", "EV Charging Station")
	params.Set("at", fmt.Sprintf("%s,%s", lat, lng))
	params.Set("limit", limit)

	apiURL.RawQuery = params.Encode()

	resp, err := http.Get(apiURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response body: %w", err)
	}

	//fmt.Println("Raw JSON Response:", string(body))

	var evChargePoints EVChargingStation
	err = json.Unmarshal(body, &evChargePoints)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return &evChargePoints, nil
}

func readCSV() ([][]string, error) {
	file, err := os.Open("Wales.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		item := strings.Split(line, ",")
		listItem := item[9:11]
		laLong = append(laLong, listItem)
	}
	return laLong, nil
}
