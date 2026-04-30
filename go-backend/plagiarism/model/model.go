package model

type Document struct {
	ID         string
	Name       string
	Category   string
	Chunks     []string
	ChunkCount int
}

type DocumentInfo struct {
	DocID      string
	DocName    string
	Category   string
	ChunkCount int
}

type Match struct {
	UploadedChunk   string  `json:"uploaded_chunk"`
	MatchedChunk    string  `json:"matched_chunk"`
	MatchedDocument string  `json:"matched_document"`
	MatchedDocID    string  `json:"matched_doc_id"`
	Category        string  `json:"category"`
	Similarity      float32 `json:"similarity"`
}

type UploadResult struct {
	DocID      string `json:"doc_id"`
	Filename   string `json:"filename"`
	Category   string `json:"category"`
	ChunkCount int    `json:"chunk_count"`
}

type SimilarityResult struct {
	Page        int     `json:"page"`
	Limit       int     `json:"limit"`
	TotalChunks int     `json:"total_chunks"`
	TotalPages  int     `json:"total_pages"`
	Matches     []Match `json:"matches"`
}

type CheckResult struct {
	TotalChunks        int     `json:"total_chunks"`
	MatchedChunks      int     `json:"matched_chunks"`
	PlagiarismPercent  float32 `json:"plagiarism_percent"`
	OriginalityPercent float32 `json:"originality_percent"`
	Matches            []Match `json:"matches"`
}

type TextSearchHit struct {
	MatchedChunk    string  `json:"matched_chunk"`
	MatchedDocument string  `json:"matched_document"`
	MatchedDocID    string  `json:"matched_doc_id"`
	Category        string  `json:"category"`
	Similarity      float32 `json:"similarity"`
}

type TextSearchResult struct {
	Query   string          `json:"query"`
	Total   int             `json:"total"`
	Hits    []TextSearchHit `json:"hits"`
}
