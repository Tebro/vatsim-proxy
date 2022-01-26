package main

import (
	"fmt"
	"io"
	"time"

	"net/http"
)

const URL = "https://data.vatsim.net/v3/vatsim-data.json"

func getData() ([]byte, error) {

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func main() {
	var dataCache []byte

	go func() {
		for {
			fmt.Println("Attempting to update data")
			result, err := getData()
			if err != nil {
				fmt.Printf("Failed to get data: %v\n", err)
			} else {
				dataCache = result
				fmt.Println("Data updated")
			}
			time.Sleep(time.Second * 60)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", dataCache)
	})
	http.ListenAndServe(":3000", nil)

}
