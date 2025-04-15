package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	es "github.com/elastic/go-elasticsearch/v9"
)

func CreateClient() (*es.Client, error) {
	client, err := es.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("error creating ElasticSearch client: %w", err)
	}
	return client, nil
}

type Book struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Year   int     `json:"year"`
	Rating float64 `json:"rating"`
	Genre  string  `json:"genre"`
}

func (b *Book) Index(es *es.Client) error {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("error marshaling book: %w", err)
	}
	res, err := es.Index(
		"books",
		strings.NewReader(string(data)),
		es.Index.WithDocumentID(fmt.Sprintf("%d", b.Year)), // Using year as ID for uniqueness
	)
	if err != nil {
		return fmt.Errorf("error indexing book: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	fmt.Println("Book indexed successfully!")
	return nil
}

func SearchBooks(es *es.Client, query string) error {
	var buf bytes.Buffer
	searchQuery := map[string]any{
		"query": map[string]any{
			"match": map[string]any{
				"title": query,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return fmt.Errorf("error encoding query: %w", err)
	}

	res, err := es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		return fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	fmt.Println("Search results:")
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		fmt.Printf("- %+v\n", source)
	}

	return nil
}

func AggregateRatingsByGenre(es *es.Client) error {
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
	if err := json.NewEncoder(&buf).Encode(aggQuery); err != nil {
		return fmt.Errorf("error encoding aggregation query: %w", err)
	}

	res, err := es.Search(
		es.Search.WithIndex("books"),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		return fmt.Errorf("error executing aggregation: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("error parsing aggregation response: %w", err)
	}

	fmt.Println("Aggregation results:")
	genres := result["aggregations"].(map[string]interface{})["genres"].(map[string]interface{})["buckets"].([]interface{})
	for _, genre := range genres {
		g := genre.(map[string]interface{})
		fmt.Printf("- Genre: %s, Avg Rating: %f\n", g["key"], g["avg_rating"].(map[string]interface{})["value"])
	}

	return nil
}
