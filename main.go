package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gota/dataframe"
)

// JSONMap is a map string interface for querried data
type JSONMap map[string]interface{}

var (
	// Key the api key needed to make requests to government supported APIs.
	Key = os.Getenv("API_KEY")
	// Maps are a slice of all maps
	Maps []map[string]interface{}
)

// GetData querries all gov apis for data.
func GetData(client http.Client, url, headerKey, header string) (df dataframe.DataFrame, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("unable to create request on url: ", url, err)
	}

	req.Header.Add("X-Api-Key", Key)
	if headerKey != "" {
		req.Header.Add("Accept", "application/vnd.api+json")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("unable to perform request on url:", url, err)
		return df, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("unable to reach ", url, resp.Status)
		return df, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to read body from ", url, err)
		return df, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Println("cannot Unmarshal, ", err)
		return df, err
	}
	Maps := append(Maps, m)
	fmt.Println(len(m))

	df = dataframe.LoadMaps(Maps)
	return df, nil
}

func main() {
	f, err := os.Open("url.csv")
	if err != nil {
		log.Fatal("unable to access csv ", err)
	}
	df := dataframe.ReadCSV(f)
	fmt.Println(df)
	if df.Nrow() == 0 || df.Ncol() == 0 {
		log.Fatal("your cvs is empty")
	}
	for i := 0; i < df.Nrow(); i++ {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		multidf, err := GetData(*client, df.Elem(i, 0).String(), df.Elem(i, 1).String(), df.Elem(i, 2).String())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(multidf)
	}
}
