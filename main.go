package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var apiBaseURL = "https://maps.googleapis.com/maps/api/geocode/json"
var apiKey = "AIzaSyD2VRI5OGSgNUSsNzm_EEyr1O22IigLI2E"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("This program requires an argument of <file.csv>")
	}

	f, err := os.Open(args[0])

	if err != nil {
		log.Fatalf("Could not open file with error: %v", err)
	}

	r := csv.NewReader(bufio.NewReader(f))

	records, _ := r.ReadAll()

	header := records[0]
	data := records[1:]

	addressIndex := indexOf("address", &header)
	cityIndex := indexOf("city", &header)
	stateIndex := indexOf("state", &header)
	zipIndex := indexOf("zip", &header)

	c := make(chan string)

	// Iterate through csv data
	for i := range data {
		row := data[i]
		address := fmt.Sprintf("%s, %s, %s %s", row[addressIndex], row[cityIndex], row[stateIndex], row[zipIndex])

		// Geocode addresses concurrently
		// Be sure to scope "address" for concurrent use
		go func(address string) {
			geocodeAddress(address, c)
		}(address)
	}

	// Print all the responses
	for i := 0; i < len(data); i++ {
		fmt.Println(<-c)
	}
}

func geocodeAddress(address string, c chan string) {

	url := apiBaseURL + "?address=" + url.PathEscape(address) + "&key=" + apiKey

	resp, err := http.Get(url)

	if err != nil {
		c <- fmt.Sprintf("Google geocoding api error: %v\n", err)
		return
	}

	if resp.StatusCode == http.StatusOK {

		var errWebResp error
		var errDecode error
		var b []byte
		var data *geocodeResponse

		b, errWebResp = ioutil.ReadAll(resp.Body)

		if errWebResp != nil {

			c <- fmt.Sprintf("Error reading web resposne: %v\n", errWebResp)
			return
		}

		data, errDecode = decodeResponse(b)

		if errDecode != nil {
			c <- fmt.Sprintf("There was an error decoding json: %v", err)
			return
		}

		if len(data.Results) == 0 {
			c <- fmt.Sprintln("No Data")
		} else {
			lat := data.Results[0].Geometry.Location["lat"]
			lon := data.Results[0].Geometry.Location["lng"]
			locationType := data.Results[0].Geometry.LocationType

			c <- fmt.Sprintf("\n") +
				fmt.Sprintf("Address: %s\n", address) +
				fmt.Sprintf("Lat: %0.2f\n", lat) +
				fmt.Sprintf("Lon: %0.2f\n", lon) +
				fmt.Sprintf("Location Type: %s\n", locationType) +
				fmt.Sprintf("\n")
		}
	}
}

func indexOf(word string, data *[]string) int {
	for k, v := range *data {
		if word == v {
			return k
		}
	}
	return -1
}
