/*
ElasticSearch distributed system from scratch using Go.
An ES cluster hosts indices.
Documents are stored in shards (physical partitions), and each index stores the shards.
Shards also store inverted indices (for search), doc values (for aggregations), and metadata.
An index represents a "category" of (related) documents.
ES is disk-based with filesystem caching. Only hot data is kept in memory (OS page cache).
A collection of documents can be retrieved by clients using an index (user-friendly category name, such as "books") and a query string.
ES then fans out the query string to all shards to find the requested document(s).
The results are fanned/merged into a single collection of results. If the result is too large, ES supports pagination.
Every time a document is retrieved by a client, what is actually returned is a copy of the document, not the original.
In Kubernetes:
Since ES is disk-based, data must be persisted using PersistentVolumes, with each pod (shard) using a PersistentVolumeClaim.
ES shard allocation handles recovery from PVs. Hot data is kept in memory.
Rather than a Deployment (which creates a ReplicaSet), 1 StatefulSet is created per node role (e.g., data, master).
Each StatefulSet Pod can host multiple shards (from different indices). This custom resource is managed by ECK Operator (not manual StatefulSets).
Each shard has N replicas (configured per index). Kubernetes ensures replicas land on different nodes via anti-affinity.
For reliability, availability, and scalability, the pods in each StatefulSet can have anti-affinity specified to give each pod a chance to be scheduled on a different node.
This strategy is useful when using EKS because nodes are distributed among availability zones, reducing the risk of AZ outages affecting the entire application.
Shard rebalancing is expensive if using HPA. Therefore, VPA for StatefulSets is preferred; use Index Lifecycle Management (ILM) for auto-scaling indices.
*/
package elastic

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Core types =================================
type Document struct {
	ID    string         `json:"id"`
	Index string         `json:"index"`
	Body  map[string]any `json:"body"`
}

// In Kafka, this would be a partition
type Shard struct {
	ID       int
	Data     map[string]Document // Storage
	Inverted map[string][]string // Index: term -> docIDs
	sync.RWMutex
}

// In Kafka, this would be a topic
type Index struct {
	Name   string
	Shards []*Shard
}

// Cluster simulation =========================
var (
	indices    = make(map[string]*Index)
	shardCount = 3
)

// NewIndex creates an index with shards
func NewIndex(name string) *Index {
	idx := &Index{
		Name:   name,
		Shards: make([]*Shard, shardCount),
	}

	for i := range shardCount {
		idx.Shards[i] = &Shard{
			ID:       i,
			Data:     make(map[string]Document),
			Inverted: make(map[string][]string),
		}
	}
	indices[name] = idx
	return idx
}

// Hash routing (like our Kafka implementation) ========
func getShard(id string) int {
	h := fnv.New32a()
	h.Write([]byte(id))
	return int(h.Sum32()) % shardCount
}

// CRUD Operations ===========================
func IndexDocument(doc Document) error {
	idx, exists := indices[doc.Index]
	if !exists {
		idx = NewIndex(doc.Index)
	}

	shard := idx.Shards[getShard(doc.ID)]

	shard.Lock()
	defer shard.Unlock()

	// Store document
	shard.Data[doc.ID] = doc

	// Update inverted index
	for field, value := range doc.Body {
		term := fmt.Sprintf("%s:%v", field, value)
		shard.Inverted[term] = append(shard.Inverted[term], doc.ID)
	}

	return nil
}

func Search(index, query string) ([]Document, error) {
	idx, exists := indices[index]
	if !exists {
		return nil, fmt.Errorf("index not found")
	}

	var results []Document
	term := parseQuery(query) // Simplified

	// Search across all shards
	for _, shard := range idx.Shards {
		shard.RLock()
		if docIDs, exists := shard.Inverted[term]; exists {
			for _, id := range docIDs {
				results = append(results, shard.Data[id])
			}
		}
		shard.RUnlock()
	}

	return results, nil
}

// Helper functions ==========================
func parseQuery(q string) string {
	// Simplified - real ES would parse properly
	return q
}
