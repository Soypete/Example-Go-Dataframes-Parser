package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kniren/gota/dataframe"
)

// GetData makes an http call with the given url to get a dataframe.
func GetData(client http.Client, apiKey string, maps []map[string]interface{}, url, headerKey, header string) (df dataframe.DataFrame, err error) {
	// create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return df, fmt.Errorf("unable to create request on url: %s, error: %s", url, err)
	}

	// add api key header
	req.Header.Add("X-Api-Key", apiKey)
	if headerKey != "" {
		req.Header.Add("Accept", "application/vnd.api+json")
	}

	// perform request
	resp, err := client.Do(req)
	if err != nil {
		return df, fmt.Errorf("unable to perform request on url: %s, error: %s", url, err)
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusOK {
		return df, fmt.Errorf("failure getting data from url: %s, status code: %s", url, resp.Status)
	}

	// json decode the response body
	m := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return df, fmt.Errorf("unable to perform request on url: %s, error: %s", url, err)
	}

	// add result to maps
	maps = append(maps, m)

	// return dataframe from the maps
	return dataframe.LoadMaps(maps), nil
}

func main() {
	// get api key from environment
	apiKey := os.Getenv("API_KEY")

	// setup map to hold dataframe results
	maps := make([]map[string]interface{}, 0)

	// open file with list of urls
	f, err := os.Open("url.csv")
	if err != nil {
		log.Fatalf("unable to access csv: %s", err)
		return
	}

	// build dataframe from csv
	df := dataframe.ReadCSV(f)
	if df.Nrow() == 0 || df.Ncol() == 0 {
		log.Fatal("your csv is empty")
		return
	}

	// create client to get data
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// process all the dataframes and get the data
	for i := 0; i < df.Nrow(); i++ {
		url := df.Elem(i, 0).String()
		headerKey := df.Elem(i, 1).String()
		header := df.Elem(i, 2).String()
		dataFrame, err := GetData(*client, apiKey, maps, url, headerKey, header)
		if err != nil {
			fmt.Printf("unable to get data for url: %s, error: %s", url, err)
			continue
		}
		fmt.Println(dataFrame)
	}
}
