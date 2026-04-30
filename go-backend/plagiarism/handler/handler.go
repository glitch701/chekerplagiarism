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
