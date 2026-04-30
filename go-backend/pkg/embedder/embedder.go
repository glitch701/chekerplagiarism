package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL   string
	batchSize int
	http      *http.Client
}

func New(baseURL string, batchSize, timeoutSec int) *Client {
	return &Client{
		baseURL:   baseURL,
		batchSize: batchSize,
		http:      &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

func (c *Client) Health(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("embedder unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("embedder health status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return c.embedBatchWithPrefix(ctx, texts, "passage")
}

func (c *Client) EmbedQuery(ctx context.Context, texts []string) ([][]float32, error) {
	return c.embedBatchWithPrefix(ctx, texts, "query")
}

func (c *Client) embedBatchWithPrefix(ctx context.Context, texts []string, prefix string) ([][]float32, error) {
	var all [][]float32
	for i := 0; i < len(texts); i += c.batchSize {
		end := i + c.batchSize
		if end > len(texts) {
			end = len(texts)
		}
		batch, err := c.embedRequest(ctx, texts[i:end], prefix)
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
	}
	return all, nil
}

type batchReq struct {
	Texts  []string `json:"texts"`
	Prefix string   `json:"prefix"`
}

type batchResp struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func (c *Client) embedRequest(ctx context.Context, texts []string, prefix string) ([][]float32, error) {
	body, _ := json.Marshal(batchReq{Texts: texts, Prefix: prefix})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/embed-batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedder request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedder returned status %d", resp.StatusCode)
	}

	var result batchResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("embedder decode: %w", err)
	}
	return result.Embeddings, nil
}
