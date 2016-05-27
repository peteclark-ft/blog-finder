package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var apiKey string
var apiURL = "https://api.ft.com/lists/"
var documentStoreURL string

func main() {
	if len(os.Args) != 3 {
		println("No API Key or DS API URL specified!")
		println("Usage: echo 'uuids' | blog-finder $apiKey $dsApiUrl ")
		os.Exit(1)
	}

	apiKey = "?apiKey=" + os.Args[1]
	documentStoreURL = os.Args[2]

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		checkList(scanner.Text())
	}

	println("Finished checking uuids.")
}

func checkList(listUUID string) {
	println("Checking list [" + listUUID + "]")
	listResponse, err := http.Get(apiURL + listUUID + apiKey)
	handleErr(err)

	list := parseListJSON(listResponse)

	for _, item := range list.Items {
		uuid := stripUUID(item.Id)
		response, err := http.Get(documentStoreURL + uuid)
		handleErr(err)

		content := parseContentJSON(response)

		for _, identifier := range content.Identifiers {
			if !strings.Contains(identifier.Authority, "METHODE") {
				println("List [" + listUUID + "] may contain blog @ [" + uuid + "]")
			}
		}
	}
	fmt.Println("Finished checking [", listUUID, "] - it had [", len(list.Items), "] items.")
}

func stripUUID(id string) string {
	return id[strings.LastIndex(id, "/")+1 : len(id)]
}

func handleErr(err error) {
	if err != nil {
		println(err)
		return
	}
}

func parseListJSON(response *http.Response) *List {
	defer response.Body.Close()

	list := new(List)
	json.NewDecoder(response.Body).Decode(list)
	return list
}

func parseContentJSON(response *http.Response) *Content {
	defer response.Body.Close()

	content := new(Content)
	json.NewDecoder(response.Body).Decode(content)
	return content
}
