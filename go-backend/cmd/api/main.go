package main

import (
	"checker/go-backend/config"
	"checker/go-backend/pkg/chunker"
	"checker/go-backend/pkg/embedder"
	"checker/go-backend/pkg/extractor"
	"checker/go-backend/pkg/qdrantclient"
	"checker/go-backend/pkg/storage"
	"checker/go-backend/plagiarism/handler"
	"checker/go-backend/plagiarism/service"
	transport "checker/go-backend/transport/http"
	"log"
)

func main() {
	cfg := config.Load()

	stor, err := storage.New(cfg.UploadDir)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	qdrant, err := qdrantclient.New(cfg.QdrantHost, cfg.QdrantGRPCPort, cfg.QdrantCollection, cfg.QdrantVectorSize)
	if err != nil {
		log.Fatalf("qdrant: %v", err)
	}

	svc := service.New(service.Deps{
		Extractor: extractor.New(),
		Chunker:   chunker.New(cfg.ChunkSize, cfg.ChunkOverlap, cfg.ChunkMinLen),
		Embedder:  embedder.New(cfg.EmbedderURL, cfg.EmbedderBatch, cfg.EmbedderTimeout),
		Qdrant:    qdrant,
		Storage:   stor,
		Threshold: cfg.SimilarityThreshold,
	})

	h := handler.New(svc)
	srv := transport.NewServer(h, cfg.ServerPort)

	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("server: %v", err)
	}
}
