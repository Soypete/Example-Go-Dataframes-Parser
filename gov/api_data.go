package gov

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/gota/dataframe"
)

var (
	// Key the ai key needed to make requests to government supported APIs.
	Key  = os.Getenv("API_KEY")
	fURL = "data/data_%d.txt"
)

// GetData querries all gov apis for data and writes that data to a file.
func GetData(client http.Client, url, headerKey, header string, i int) (err error) {
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
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("unable to reach ", url, resp.Status)
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("unable to read body from ", url, err)
		return err
	}
	//	fmt.Println(string(b))
	file := fmt.Sprintf(fURL, i+1)
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
	df, err := MakeDataframe(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(df)
	time.Sleep(15 * time.Second)
	return nil
}

// MakeDataframe takes all text files and reads them into dataframes.
func MakeDataframe(file string) (df dataframe.DataFrame, err error) {
	f, err := os.Open(file)
	if err != nil {
		return df, err
	}
	defer f.Close()
	spew.Dump(f)
	df = dataframe.ReadJSON(f)
	return df, err
}

func main() {
	f, err := os.Open("data_2/url.csv")
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
		err := GetData(*client, df.Elem(i, 0).String(), df.Elem(i, 1).String(), df.Elem(i, 2).String(), i)
		if err != nil {
			log.Fatal(err)
		}
	}
	//	df1 := dataframe.LoadMaps(Maps)
	//	fmt.Println(df1)

}
