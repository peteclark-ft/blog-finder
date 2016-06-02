package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

var apiURL = "https://api.ft.com/lists/"
var apiKey *string
var documentStoreURL *string

var jsonPrint *bool

func main() {
	app := cli.App("blog-finder", "Given a list of uuids for Methode Editorial Lists, the Blog Finder will locate Wordpress blogs!")

	app.Version("v version", "1.0.0")

	apiKey = app.StringArg("APIKEY", "", "A valid FT.com API key.")
	documentStoreURL = app.StringArg("DSURL", "", "URL for a running DS API Instance.")

	jsonPrint = app.BoolOpt("j json-print", false, "Print the results in JSON.")

	app.Action = func() {
		if *jsonPrint {
			log.SetLevel(log.ErrorLevel)
		}

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			checkList(scanner.Text())
		}

		log.Info("Finished checking uuids.")
	}

	app.Run(os.Args)
}

func checkList(listUUID string) {
	log.Infof("Checking list [%v]", listUUID)
	listResponse, err := http.Get(apiURL + listUUID + "?apiKey=" + *apiKey)
	handleErr(err)

	list := parseListJSON(listResponse)

	for _, item := range list.Items {
		uuid := stripUUID(item.Id)
		response, err := http.Get(*documentStoreURL + uuid)
		handleErr(err)

		content := parseContentJSON(response)

		for _, identifier := range content.Identifiers {
			if !strings.Contains(identifier.Authority, "METHODE") {
				if *jsonPrint {
					println(toJSON(Result{listUUID, uuid}))
				}
				log.Infof("List [" + listUUID + "] may contain blog @ [" + uuid + "]")
			}
		}
	}
	log.Infof("Finished checking [%v] - it had [%v] items.", listUUID, len(list.Items))
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

func toJSON(result Result) string {
	b, _ := json.Marshal(result)
	return string(b)
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
