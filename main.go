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
	"time"
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
			if len(station.Connectors.Connector) > 0 {
				for _, connector := range station.Connectors.Connector {
					fmt.Printf("  Supplier: %s, Charge Capacity: %s\n", connector.SupplierName, connector.ChargeCapacity)
				}
			} else {
				fmt.Println("  No connectors available.")
			}
			fmt.Println("-------------------------------")
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
		PoolID                  string `json:"poolId,omitempty"`
		CpoID                   string `json:"cpoId,omitempty"`
		TotalNumberOfConnectors int    `json:"totalNumberOfConnectors,omitempty"`
		Connectors              struct {
			Connector []struct {
				SupplierName  string `json:"supplierName"`
				ConnectorType []struct {
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"connectorType"`
				ChargeCapacity        string  `json:"chargeCapacity"`
				MaxPowerLevel         float64 `json:"maxPowerLevel"`
				CustomerChargeLevel   string  `json:"customerChargeLevel"`
				CustomerConnectorName string  `json:"customerConnectorName"`
				FixedCable            bool    `json:"fixedCable"`
				ConnectorDetails      struct {
					PrivateAccess bool   `json:"privateAccess"`
					Open24x7      bool   `json:"open24x7"`
					OpeningTime   string `json:"openingTime"`
					OpeningHours  struct {
						RegularOpeningHours []struct {
							Daymask []int `json:"daymask"`
							Period  []struct {
								From string `json:"from"`
								To   string `json:"to"`
							} `json:"period"`
						} `json:"regularOpeningHours"`
						AnnualOpenings []interface{} `json:"annualOpenings"`
					} `json:"openingHours"`
					Pay          bool   `json:"pay"`
					Manufacturer string `json:"manufacturer"`
				} `json:"connectorDetails"`
				ChargingPoint []struct {
					ChargeMode          string    `json:"chargeMode"`
					VoltsRange          string    `json:"voltsRange"`
					Phases              int       `json:"phases"`
					AmpsRange           string    `json:"ampsRange"`
					NumberOfConnectors  int       `json:"numberOfConnectors"`
					NumberOfAvailable   int       `json:"numberOfAvailable"`
					LastUpdateTimestamp time.Time `json:"lastUpdateTimestamp"`
				} `json:"chargingPoint"`
				ConnectorStatuses struct {
					ConnectorStatus []struct {
						CpoEvseId         string `json:"cpoEvseId"`
						CpoEvseEMI3Id     string `json:"cpoEvseEMI3Id"`
						PhysicalReference string `json:"physicalReference"`
						CpoConnectorId    string `json:"cpoConnectorId"`
						State             string `json:"state"`
					} `json:"connectorStatus"`
				} `json:"connectorStatuses"`
			} `json:"connector,omitempty"`
		} `json:"connectors,omitempty"`
		EVStationDetails []struct {
			PrivateAccess     bool   `json:"privateAccess"`
			RestrictedAccess  bool   `json:"restrictedAccess"`
			AccessibilityType string `json:"accessibilityType"`
			PaymentMethods    struct {
				Subscription struct {
					Provider string `json:"provider"`
					Accept   bool   `json:"accept"`
				} `json:"subscription"`
				Note       string `json:"note"`
				CreditCard struct {
					Types struct {
						Type []string `json:"type"`
					} `json:"types"`
					Accept bool `json:"accept"`
				} `json:"creditCard"`
				Cash struct {
					Currencies struct {
						Currency []string `json:"currency"`
					} `json:"currencies"`
					Accept bool `json:"accept"`
				} `json:"cash"`
				DebitCard struct {
					Accept bool `json:"accept"`
				} `json:"debitCard"`
				Check struct {
					Types struct {
						Type []string `json:"type"`
					} `json:"types"`
					Accept bool `json:"accept"`
				} `json:"check"`
				EPayment struct {
					Types struct {
						Type []string `json:"type"`
					} `json:"types"`
					Accept bool `json:"accept"`
				} `json:"ePayment"`
				Other struct {
					Types struct {
						Type []string `json:"type"`
					} `json:"types"`
					Accept bool `json:"accept"`
				} `json:"other"`
			} `json:"paymentMethods"`
			Notes string `json:"notes"`
		} `json:"evStationDetails,omitempty"`
		LastUpdateTimestamp time.Time `json:"lastUpdateTimestamp,omitempty"`
		TimeZone            string    `json:"timeZone,omitempty"`
		Name                string    `json:"name,omitempty"`
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
