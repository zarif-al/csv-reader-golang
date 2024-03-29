package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Customer struct {
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	Enabled           bool      `json:"enabled"`
	EmailScheduleTime time.Time `json:"emailScheduleTime"`
}

func main() {

	records := make(chan []string)
	customersArray := make(chan string)

	go reader(records)

	go createArray(records, customersArray)

	response := <-customersArray

	fmt.Println(response)

}

func createArray(records chan []string, customersArray chan string) {

	defer close(customersArray)
	for record := range records {

		customer := Customer{}

		customer.Name = record[0]
		customer.Email = record[1]
		customer.Address = record[2]

		enabled, err := strconv.ParseBool(record[3])
		if err != nil {
			log.Fatal(err)
		}

		customer.Enabled = enabled

		t, err := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
		if err != nil {
			log.Fatal(err)
		}

		customer.EmailScheduleTime = t

	}
	customersArray <- "done"
}

func reader(records chan []string) {
	defer close(records)

	file, err := os.Open("MOCK_DATA.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	parser := csv.NewReader(file)

	// Skipping the first line
	if _, err := parser.Read(); err != nil {
		panic(err)
	}

	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		records <- record
	}

}
