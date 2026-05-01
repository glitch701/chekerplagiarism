package dto

import "checker/go-backend/plagiarism/model"

type UploadReferenceResponse struct {
	Success    bool   `json:"success"`
	DocID      string `json:"doc_id"`
	Filename   string `json:"filename"`
	Category   string `json:"category"`
	ChunkCount int    `json:"chunk_count"`
	Message    string `json:"message"`
}

type SimilarityResponse struct {
	Page        int           `json:"page"`
	Limit       int           `json:"limit"`
	TotalChunks int           `json:"total_chunks"`
	TotalPages  int           `json:"total_pages"`
	Matches     []model.Match `json:"matches"`
}

type PlagiarismResponse struct {
	TotalChunks   int           `json:"total_chunks"`
	MatchedChunks int           `json:"matched_chunks"`
	Matches       []model.Match `json:"matches"`
}

type DocumentInfoResponse struct {
	DocID      string `json:"doc_id"`
	DocName    string `json:"doc_name"`
	Category   string `json:"category"`
	ChunkCount int    `json:"chunk_count"`
}

type ListReferencesResponse struct {
	Total     int                    `json:"total"`
	Documents []DocumentInfoResponse `json:"documents"`
}

type TextSearchHitResponse struct {
	MatchedChunk    string  `json:"matched_chunk"`
	MatchedDocument string  `json:"matched_document"`
	MatchedDocID    string  `json:"matched_doc_id"`
	Category        string  `json:"category"`
	Similarity      float32 `json:"similarity"`
}

type TextSearchResponse struct {
	Query string                  `json:"query"`
	Total int                     `json:"total"`
	Hits  []TextSearchHitResponse `json:"hits"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func FromUploadResult(r *model.UploadResult) UploadReferenceResponse {
	return UploadReferenceResponse{
		Success:    true,
		DocID:      r.DocID,
		Filename:   r.Filename,
		Category:   r.Category,
		ChunkCount: r.ChunkCount,
		Message:    "Reference document uploaded successfully",
	}
}

func FromSimilarityResult(r *model.SimilarityResult) SimilarityResponse {
	matches := r.Matches
	if matches == nil {
		matches = []model.Match{}
	}
	return SimilarityResponse{
		Page:        r.Page,
		Limit:       r.Limit,
		TotalChunks: r.TotalChunks,
		TotalPages:  r.TotalPages,
		Matches:     matches,
	}
}

func FromDocumentInfoList(docs []model.DocumentInfo) ListReferencesResponse {
	items := make([]DocumentInfoResponse, len(docs))
	for i, d := range docs {
		items[i] = DocumentInfoResponse{
			DocID:      d.DocID,
			DocName:    d.DocName,
			Category:   d.Category,
			ChunkCount: d.ChunkCount,
		}
	}
	return ListReferencesResponse{Total: len(docs), Documents: items}
}

func FromTextSearchResult(r *model.TextSearchResult) TextSearchResponse {
	hits := make([]TextSearchHitResponse, len(r.Hits))
	for i, h := range r.Hits {
		hits[i] = TextSearchHitResponse{
			MatchedChunk:    h.MatchedChunk,
			MatchedDocument: h.MatchedDocument,
			MatchedDocID:    h.MatchedDocID,
			Category:        h.Category,
			Similarity:      h.Similarity,
		}
	}
	return TextSearchResponse{Query: r.Query, Total: r.Total, Hits: hits}
}

func FromCheckResult(r *model.CheckResult) PlagiarismResponse {
	matches := r.Matches
	if matches == nil {
		matches = []model.Match{}
	}
	return PlagiarismResponse{
		TotalChunks:   r.TotalChunks,
		MatchedChunks: r.MatchedChunks,
		Matches:       matches,
	}
}
