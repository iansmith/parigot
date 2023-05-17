package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
)

func main() {
	links := readCSV()
	if links == nil {
		fmt.Println("Failed to find any links.")
		return
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
		fmt.Println("Failed to open CSV file: ", err)
		return nil
	}
	defer file.Close()

	// Create a new CSV reader reading from the opened file
	reader := csv.NewReader(file)

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Failed to read CSV records: ", err)
		return nil
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
		fmt.Println("Target column not found: ", targetColumn)
		return nil
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
		fmt.Printf("Failed to fetch %s: %s\n", webURL, err)
		return
	}

	defer resp.Body.Close()

	fmt.Printf("Successfully fetching data from %s\n", webURL)
}
