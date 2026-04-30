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
	TotalChunks        int           `json:"total_chunks"`
	MatchedChunks      int           `json:"matched_chunks"`
	PlagiarismPercent  float32       `json:"plagiarism_percent"`
	OriginalityPercent float32       `json:"originality_percent"`
	Matches            []model.Match `json:"matches"`
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

func FromCheckResult(r *model.CheckResult) PlagiarismResponse {
	matches := r.Matches
	if matches == nil {
		matches = []model.Match{}
	}
	return PlagiarismResponse{
		TotalChunks:        r.TotalChunks,
		MatchedChunks:      r.MatchedChunks,
		PlagiarismPercent:  r.PlagiarismPercent,
		OriginalityPercent: r.OriginalityPercent,
		Matches:            matches,
	}
}
