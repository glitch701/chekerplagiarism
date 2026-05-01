package qdrantclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type Client struct {
	inner      *qdrant.Client
	collection string
	vectorSize uint64
}

type SearchResult struct {
	Score    float32
	Text     string
	DocName  string
	DocID    string
	Category string
}

func New(host string, grpcPort int, collection string, vectorSize uint64) (*Client, error) {
	inner, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: grpcPort,
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant connect: %w", err)
	}

	c := &Client{inner: inner, collection: collection, vectorSize: vectorSize}
	if err := c.ensureCollection(context.Background()); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) ensureCollection(ctx context.Context) error {
	exists, err := c.inner.CollectionExists(ctx, c.collection)
	if err != nil {
		return fmt.Errorf("collection exists check: %w", err)
	}
	if exists {
		return nil
	}
	return c.inner.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: c.collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     c.vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (c *Client) Upsert(ctx context.Context, docID, docName, category string, chunks []string, embeddings [][]float32) error {
	points := make([]*qdrant.PointStruct, 0, len(chunks))
	for i, chunk := range chunks {
		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDUUID(uuid.New().String()),
			Vectors: qdrant.NewVectors(embeddings[i]...),
			Payload: qdrant.NewValueMap(map[string]any{
				"doc_id":      docID,
				"doc_name":    docName,
				"category":    category,
				"chunk_index": i,
				"text":        chunk,
			}),
		})
	}
	_, err := c.inner.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: c.collection,
		Points:         points,
	})
	return err
}

func (c *Client) Search(ctx context.Context, embedding []float32, threshold float32, limit uint64) ([]SearchResult, error) {
	q := &qdrant.QueryPoints{
		CollectionName: c.collection,
		Query:          qdrant.NewQuery(embedding...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
	}
	if threshold > 0 {
		q.ScoreThreshold = &threshold
	}
	scored, err := c.inner.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("qdrant query: %w", err)
	}

	results := make([]SearchResult, 0, len(scored))
	for _, p := range scored {
		results = append(results, SearchResult{
			Score:    p.Score,
			Text:     p.Payload["text"].GetStringValue(),
			DocName:  p.Payload["doc_name"].GetStringValue(),
			DocID:    p.Payload["doc_id"].GetStringValue(),
			Category: p.Payload["category"].GetStringValue(),
		})
	}
	return results, nil
}

type DocumentInfo struct {
	DocID      string
	DocName    string
	Category   string
	ChunkCount int
}

func (c *Client) ListDocuments(ctx context.Context) ([]DocumentInfo, error) {
	seen := make(map[string]*DocumentInfo)
	var offset *qdrant.PointId

	for {
		points, next, err := c.inner.ScrollAndOffset(ctx, &qdrant.ScrollPoints{
			CollectionName: c.collection,
			WithPayload:    qdrant.NewWithPayload(true),
			WithVectors:    qdrant.NewWithVectors(false),
			Limit:          qdrant.PtrOf(uint32(100)),
			Offset:         offset,
		})
		if err != nil {
			return nil, fmt.Errorf("qdrant scroll: %w", err)
		}

		for _, p := range points {
			docID := p.Payload["doc_id"].GetStringValue()
			if info, ok := seen[docID]; ok {
				info.ChunkCount++
			} else {
				seen[docID] = &DocumentInfo{
					DocID:      docID,
					DocName:    p.Payload["doc_name"].GetStringValue(),
					Category:   p.Payload["category"].GetStringValue(),
					ChunkCount: 1,
				}
			}
		}

		if next == nil {
			break
		}
		offset = next
	}

	docs := make([]DocumentInfo, 0, len(seen))
	for _, info := range seen {
		docs = append(docs, *info)
	}
	return docs, nil
}

func (c *Client) DeleteByDocID(ctx context.Context, docID string) error {
	_, err := c.inner.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: c.collection,
		Points: qdrant.NewPointsSelectorFilter(&qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatch("doc_id", docID),
			},
		}),
	})
	return err
}
