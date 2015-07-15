package parse_test

import (
	"fmt"
	"log"
	"os"

	"github.com/tmc/parse"
)

func ExampleNewClient() {
	appID := os.Getenv("APPLICATION_ID")
	apiKey := os.Getenv("REST_API_KEY")
	_, err := parse.NewClient(appID, apiKey)
	fmt.Println(err)
	// output: <nil>
}

// Shows saving, querying, and deletion.
func ExampleClient_Query() {
	appID := os.Getenv("APPLICATION_ID")
	apiKey := os.Getenv("REST_API_KEY")
	client, err := parse.NewClient(appID, apiKey)
	if err != nil {
		log.Fatalln("parse: error creating client:", err)
	}

	type Widget struct {
		parse.ParseObject
		Name     string `json:"name"`
		Quantity int    `json:"qty"`
	}
	widgets := []*Widget{
		{Name: "widget a", Quantity: 42},
		{Name: "widget b", Quantity: 41},
		{Name: "widget c", Quantity: 40},
	}
	for _, w := range widgets {
		if _, err := client.Create(w); err != nil {
			log.Fatalln("parse: error saving widget:", err)
		}
	}

	results := []*Widget{}
	if err := client.Query(nil, &results); err != nil {
		log.Fatalln("parse: error querying widgets:", err)
	}

	for _, w := range results {
		client.Delete(w)
	}
	fmt.Println(len(results))
	// output:
	// 3
}
