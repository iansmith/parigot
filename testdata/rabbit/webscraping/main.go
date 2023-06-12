package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	links := readCSV()
	if links == nil {
		log.Fatalf("Failed to find any links.")
	}

	for _, link := range links {
		getDataFromWeb(link)
	}

}

func readCSV() []string {
	filePath := "./top500Domains.csv"
	targetColumn := "Root Domain"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CSV file: ", err)
	}

	defer file.Close()

	// Create a new CSV reader reading from the opened file
	reader := csv.NewReader(file)

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV records: ", err)
	}

	header := records[0]
	// Find the index of the target column
	targetColumnIndex := -1
	for i, column := range header {
		if column == targetColumn {
			targetColumnIndex = i
			break
		}
	}

	if targetColumnIndex == -1 {
		log.Fatalf("Target column not found: ", targetColumn)
	}

	links := make([]string, len(records)-1)
	for i, record := range records[1:] {
		// add a valid protocol scheme "https://"
		links[i] = "https://" + record[targetColumnIndex]
	}
	return links
}

func getDataFromWeb(webURL string) {
	resp, err := http.Get(webURL)
	if err != nil {
		log.Printf("Failed to fetch %s: %s\n", webURL, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Server returned non-200 status code: ", resp.StatusCode)
		return
	}

	limitBody = io.LimitReader(resp.Body, 1)
	if body, err := io.Copy(os.Stdout, limitBody); err != nil {
		log.Printf("Error reading response body: %s\n", err)
		return
	}

	if len(body) == 0 {
		log.Println("Response body is empty.")
		return
	}

	fmt.Printf("Successfully fetched data from %s\n", webURL)
}
