package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
)

func main() {
	f, err := os.Open("data/tx_law.txt")
	if err != nil {
		log.Fatal("unable to access .txt file ", err)
	}
	df := dataframe.ReadJSON(f)

	fmt.Println(df, df.Describe())
}
