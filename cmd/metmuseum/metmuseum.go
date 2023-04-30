package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiURL = "https://collectionapi.metmuseum.org/public/collection/v1/search"

type SearchResponse struct {
	ObjectIDs []int `json:"objectIDs"`
}

func main() {
	query := "books"
	queryParams := url.Values{}
	queryParams.Set("q", query)
	queryParams.Set("isHighlight", "true")

	resp, err := http.Get(fmt.Sprintf("%s?%s", apiURL, queryParams.Encode()))
	if err != nil {
		fmt.Printf("Error fetching data from API: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading API response: %v\n", err)
		return
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	objectIDs := searchResponse.ObjectIDs
	if len(objectIDs) == 0 {
		fmt.Println("No results found.")
		return
	}

	// Fetch the first 10 results, or fewer if less than 10 are available
	resultCount := 10
	if len(objectIDs) < 10 {
		resultCount = len(objectIDs)
	}

	fmt.Printf("Showing the first %d results for '%s':\n\n", resultCount, query)
	for i := 0; i < resultCount; i++ {
		objectID := objectIDs[i]
		objectURL := fmt.Sprintf("https://collectionapi.metmuseum.org/public/collection/v1/objects/%d", objectID)

		resp, err := http.Get(objectURL)
		if err != nil {
			fmt.Printf("Error fetching object data: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading object data: %v\n", err)
			continue
		}

		var prettyBuffer bytes.Buffer
		err = json.Indent(&prettyBuffer, body, "", "  ")
		if err != nil {
			fmt.Printf("Error pretty-printing JSON: %v\n", err)
			continue
		}

		prettyJSON := prettyBuffer.String()

		fmt.Printf("Object ID: %d\n", objectID)
		fmt.Println("JSON Response Body:")
		fmt.Println(prettyJSON)
		fmt.Println()
	}
}
