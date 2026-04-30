package service

import (
	"checker/go-backend/pkg/chunker"
	"checker/go-backend/pkg/embedder"
	"checker/go-backend/pkg/extractor"
	"checker/go-backend/pkg/qdrantclient"
	"checker/go-backend/pkg/storage"
	"checker/go-backend/plagiarism/model"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type PlagiarismService interface {
	UploadReference(ctx context.Context, filename, category string, data []byte) (*model.UploadResult, error)
	CheckSimilarity(ctx context.Context, filename string, data []byte, page, limit int) (*model.SimilarityResult, error)
	CheckPlagiarism(ctx context.Context, filename string, data []byte) (*model.CheckResult, error)
}

type plagiarismService struct {
	extractor *extractor.Extractor
	chunker   *chunker.Chunker
	embedder  *embedder.Client
	qdrant    *qdrantclient.Client
	storage   *storage.Storage
	threshold float32
}

type Deps struct {
	Extractor *extractor.Extractor
	Chunker   *chunker.Chunker
	Embedder  *embedder.Client
	Qdrant    *qdrantclient.Client
	Storage   *storage.Storage
	Threshold float32
}

func New(d Deps) PlagiarismService {
	if d.Threshold == 0 {
		d.Threshold = 0.75
	}
	return &plagiarismService{
		extractor: d.Extractor,
		chunker:   d.Chunker,
		embedder:  d.Embedder,
		qdrant:    d.Qdrant,
		storage:   d.Storage,
		threshold: d.Threshold,
	}
}

func (s *plagiarismService) UploadReference(ctx context.Context, filename, category string, data []byte) (*model.UploadResult, error) {
	text, err := s.extractor.ExtractFromBytes(filename, data)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	chunks := s.chunker.Chunk(text)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no text content found in document")
	}

	embeddings, err := s.embedder.EmbedBatch(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("embed: %w", err)
	}

	docID := uuid.New().String()

	if err := s.storage.Save(docID, filename, data); err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}

	if err := s.qdrant.Upsert(ctx, docID, filename, category, chunks, embeddings); err != nil {
		return nil, fmt.Errorf("qdrant upsert: %w", err)
	}

	return &model.UploadResult{
		DocID:      docID,
		Filename:   filename,
		Category:   category,
		ChunkCount: len(chunks),
	}, nil
}

func (s *plagiarismService) CheckSimilarity(ctx context.Context, filename string, data []byte, page, limit int) (*model.SimilarityResult, error) {
	text, err := s.extractor.ExtractFromBytes(filename, data)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	chunks := s.chunker.Chunk(text)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no text content found in document")
	}

	total := len(chunks)
	totalPages := (total + limit - 1) / limit

	start := (page - 1) * limit
	if start >= total {
		return &model.SimilarityResult{
			Page:        page,
			Limit:       limit,
			TotalChunks: total,
			TotalPages:  totalPages,
			Matches:     []model.Match{},
		}, nil
	}
	end := start + limit
	if end > total {
		end = total
	}

	pageChunks := chunks[start:end]
	embeddings, err := s.embedder.EmbedBatch(ctx, pageChunks)
	if err != nil {
		return nil, fmt.Errorf("embed: %w", err)
	}

	seen := make(map[string]struct{})
	var matches []model.Match
	for i, emb := range embeddings {
		if _, dup := seen[pageChunks[i]]; dup {
			continue
		}
		results, err := s.qdrant.Search(ctx, emb, s.threshold, 1)
		if err != nil {
			continue
		}
		if len(results) == 0 {
			continue
		}
		seen[pageChunks[i]] = struct{}{}
		r := results[0]
		matches = append(matches, model.Match{
			UploadedChunk:   pageChunks[i],
			MatchedChunk:    r.Text,
			MatchedDocument: r.DocName,
			MatchedDocID:    r.DocID,
			Category:        r.Category,
			Similarity:      r.Score * 100,
		})
	}

	return &model.SimilarityResult{
		Page:        page,
		Limit:       limit,
		TotalChunks: total,
		TotalPages:  totalPages,
		Matches:     matches,
	}, nil
}

func (s *plagiarismService) CheckPlagiarism(ctx context.Context, filename string, data []byte) (*model.CheckResult, error) {
	text, err := s.extractor.ExtractFromBytes(filename, data)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	chunks := s.chunker.Chunk(text)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no text content found in document")
	}

	embeddings, err := s.embedder.EmbedBatch(ctx, chunks)
	if err != nil {
		return nil, fmt.Errorf("embed: %w", err)
	}

	seen := make(map[string]struct{})
	var matches []model.Match
	matchedSet := make(map[int]struct{})

	for i, emb := range embeddings {
		if _, dup := seen[chunks[i]]; dup {
			matchedSet[i] = struct{}{}
			continue
		}
		results, err := s.qdrant.Search(ctx, emb, s.threshold, 1)
		if err != nil {
			continue
		}
		if len(results) == 0 {
			continue
		}
		seen[chunks[i]] = struct{}{}
		r := results[0]
		matches = append(matches, model.Match{
			UploadedChunk:   chunks[i],
			MatchedChunk:    r.Text,
			MatchedDocument: r.DocName,
			MatchedDocID:    r.DocID,
			Category:        r.Category,
			Similarity:      r.Score * 100,
		})
		matchedSet[i] = struct{}{}
	}

	total := len(chunks)
	matched := len(matchedSet)
	plagPercent := float32(0)
	if total > 0 {
		plagPercent = float32(matched) / float32(total) * 100
	}

	return &model.CheckResult{
		TotalChunks:        total,
		MatchedChunks:      matched,
		PlagiarismPercent:  plagPercent,
		OriginalityPercent: 100 - plagPercent,
		Matches:            matches,
	}, nil
}
