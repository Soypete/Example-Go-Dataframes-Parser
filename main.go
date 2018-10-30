package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gota/dataframe"
)

// JSONMap is a map string interface for querried data
type JSONMap map[string]interface{}

var (
	// Key the api key needed to make requests to government supported APIs.
	Key = os.Getenv("API_KEY")
	//Map is a singular json map
	Map JSONMap
	// Maps are a slice of all maps
	Maps []JSONMap
)

// GetData querries all gov apis for data.
func GetData(client http.Client, url, headerKey, header string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("unable to create request on url: ", url, err)
	}

	req.Header.Add("X-Api-Key", Key)
	req.Header.Add("Accept", "application/vnd.api+json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("unable to perform request on url: ", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("unable to reach ", url)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to read body from ", url, err)
	}
	err = json.Unmarshal(body, &Map)
	if err != nil {
		log.Fatal("cannot Unmarshal, ", err)
	}
	fmt.Println(len(Map))

	//	df := dataframe.ReadJSON(resp.Body)
	//	fmt.Println(df)
}

func main() {
	f, err := os.Open("url.csv")
	if err != nil {
		log.Fatal("unable to access csv ", err)
	}

	//	loop := func(s series.Series) {
	//		client := &http.Client{}
	//		GetData(*client)
	//	}

	df := dataframe.ReadCSV(f)
	fmt.Println(df)
	//	df.Rapply(loop)
}
