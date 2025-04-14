package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	es "github.com/elastic/go-elasticsearch/v8"
)

type Book struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating"`
	Genre  string  `json:"genre"`
}

func (b *Book) Index(es *es.Client) {
	data, _ := json.Marshal(b)
	res, err := es.Index(
		"books",
		strings.NewReader(string(data)),
		es.Index.WithDocumentID("1"), // Optional ID
	)
	if err != nil {
		log.Fatalf("Error indexing book: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("Book indexed successfully!")
}

func SearchBooks(es *es.Client, query string) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": query,
			},
		},
	}
	json.NewEncoder(&buf).Encode(searchQuery)

	res, err := es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error searching: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("Search results:", res.String())
}

func AggregateRatingsByGenre(es *es.Client) {
	var buf bytes.Buffer
	aggQuery := map[string]interface{}{
		"aggs": map[string]interface{}{
			"genres": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "genre",
				},
				"aggs": map[string]interface{}{
					"avg_rating": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": "rating",
						},
					},
				},
			},
		},
	}
	json.NewEncoder(&buf).Encode(aggQuery)

	res, err := es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		log.Fatalf("Error aggregating: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("Aggregation results:", res.String())
}
