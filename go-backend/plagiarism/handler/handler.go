package handler

import (
	"checker/go-backend/plagiarism/dto"
	"checker/go-backend/plagiarism/service"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc service.PlagiarismService
}

func New(svc service.PlagiarismService) *Handler {
	return &Handler{svc: svc}
}

// POST /api/v1/reference
func (h *Handler) UploadReference(c echo.Context) error {
	category := c.FormValue("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "category is required"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "file is required"})
	}

	data, err := readFile(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	result, err := h.svc.UploadReference(c.Request().Context(), file.Filename, category, data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.FromUploadResult(result))
}

// POST /api/v1/similarity?page=1&limit=10
func (h *Handler) CheckSimilarity(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "file is required"})
	}

	data, err := readFile(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	result, err := h.svc.CheckSimilarity(c.Request().Context(), file.Filename, data, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.FromSimilarityResult(result))
}

// POST /api/v1/plagiarism
func (h *Handler) CheckPlagiarism(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "file is required"})
	}

	data, err := readFile(file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	result, err := h.svc.CheckPlagiarism(c.Request().Context(), file.Filename, data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.FromCheckResult(result))
}

// GET /api/v1/references
func (h *Handler) ListReferences(c echo.Context) error {
	docs, err := h.svc.ListReferences(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.FromDocumentInfoList(docs))
}

// DELETE /api/v1/reference/:id
func (h *Handler) DeleteReference(c echo.Context) error {
	docID := c.Param("id")
	if docID == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "doc_id is required"})
	}
	if err := h.svc.DeleteReference(c.Request().Context(), docID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"success": true, "doc_id": docID})
}

// POST /api/v1/search
func (h *Handler) SearchText(c echo.Context) error {
	var req struct {
		Query     string  `json:"query"`
		Limit     int     `json:"limit"`
		Threshold float32 `json:"threshold"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
	}
	if req.Query == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "query is required"})
	}
	if req.Limit < 1 {
		req.Limit = 5
	}
	if req.Threshold <= 0 {
		req.Threshold = 0.65
	}

	result, err := h.svc.SearchText(c.Request().Context(), req.Query, req.Limit, req.Threshold)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.FromTextSearchResult(result))
}

// GET /health
func (h *Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func readFile(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
