package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	ServerPort          string
	UploadDir           string
	QdrantHost          string
	QdrantGRPCPort      int
	QdrantCollection    string
	QdrantVectorSize    uint64
	EmbedderURL         string
	EmbedderBatch       int
	EmbedderTimeout     int
	ChunkSize           int
	ChunkOverlap        int
	ChunkMinLen         int
	SimilarityThreshold float32
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env not found, using environment variables")
	}
	return Config{
		ServerPort:       getOrDefault("SERVER_PORT", "8000"),
		UploadDir:        getOrDefault("UPLOAD_DIR", "./uploads"),
		QdrantHost:       getOrDefault("QDRANT_HOST", "localhost"),
		QdrantGRPCPort:   cast.ToInt(getOrDefault("QDRANT_GRPC_PORT", "6334")),
		QdrantCollection: getOrDefault("QDRANT_COLLECTION", "documents"),
		QdrantVectorSize: cast.ToUint64(getOrDefault("QDRANT_VECTOR_SIZE", "384")),
		EmbedderURL:      getOrDefault("EMBEDDER_URL", "http://localhost:5000"),
		EmbedderBatch:    cast.ToInt(getOrDefault("EMBEDDER_BATCH_SIZE", "64")),
		EmbedderTimeout:  cast.ToInt(getOrDefault("EMBEDDER_TIMEOUT_SEC", "300")),
		ChunkSize:           cast.ToInt(getOrDefault("CHUNK_SIZE", "500")),
		ChunkOverlap:        cast.ToInt(getOrDefault("CHUNK_OVERLAP", "100")),
		ChunkMinLen:         cast.ToInt(getOrDefault("CHUNK_MIN_LEN", "50")),
		SimilarityThreshold: cast.ToFloat32(getOrDefault("SIMILARITY_THRESHOLD", "0.88")),
	}
}

func getOrDefault(key string, defaultValue interface{}) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return cast.ToString(defaultValue)
}
