package main

import (
	"checker/go-backend/pkg/chunker"
	"checker/go-backend/pkg/embedder"
	"checker/go-backend/pkg/extractor"
	"checker/go-backend/pkg/qdrantclient"
	"checker/go-backend/pkg/storage"
	"checker/go-backend/plagiarism/service"
	transport "checker/go-backend/transport/service"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cast"
)

type Env struct {
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

func main() {
	var env Env
	loadEnv(&env)

	stor, err := storage.New(env.UploadDir)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}

	qdrant, err := qdrantclient.New(env.QdrantHost, env.QdrantGRPCPort, env.QdrantCollection, env.QdrantVectorSize)
	if err != nil {
		log.Fatalf("qdrant: %v", err)
	}

	svc := service.New(service.Deps{
		Extractor: extractor.New(),
		Chunker:   chunker.New(env.ChunkSize, env.ChunkOverlap, env.ChunkMinLen),
		Embedder:  embedder.New(env.EmbedderURL, env.EmbedderBatch, env.EmbedderTimeout),
		Qdrant:    qdrant,
		Storage:   stor,
		Threshold: env.SimilarityThreshold,
	})

	var router = echo.New()
	setMiddlewares(router)
	transport.NewService(router, svc)
	runHTTPServerOnAddr(router, env.ServerPort)
}

func runHTTPServerOnAddr(router *echo.Echo, port string) {
	log.Printf("Server starting on :%s", port)
	if err := router.Start(":" + port); err != nil {
		panic(err)
	}
}

func setMiddlewares(router *echo.Echo) {
	router.HideBanner = true
	router.Use(middleware.RemoveTrailingSlash())
	router.Use(middleware.RequestID())
	router.Use(middleware.Recover())
}

func loadEnv(env *Env) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env not found, using environment variables")
	}
	env.ServerPort = getOrDefault("SERVER_PORT", "8000")
	env.UploadDir = getOrDefault("UPLOAD_DIR", "./uploads")
	env.QdrantHost = getOrDefault("QDRANT_HOST", "localhost")
	env.QdrantGRPCPort = cast.ToInt(getOrDefault("QDRANT_GRPC_PORT", "6334"))
	env.QdrantCollection = getOrDefault("QDRANT_COLLECTION", "documents")
	env.QdrantVectorSize = cast.ToUint64(getOrDefault("QDRANT_VECTOR_SIZE", "384"))
	env.EmbedderURL = getOrDefault("EMBEDDER_URL", "http://localhost:5000")
	env.EmbedderBatch = cast.ToInt(getOrDefault("EMBEDDER_BATCH_SIZE", "64"))
	env.EmbedderTimeout = cast.ToInt(getOrDefault("EMBEDDER_TIMEOUT_SEC", "300"))
	env.ChunkSize = cast.ToInt(getOrDefault("CHUNK_SIZE", "500"))
	env.ChunkOverlap = cast.ToInt(getOrDefault("CHUNK_OVERLAP", "100"))
	env.ChunkMinLen = cast.ToInt(getOrDefault("CHUNK_MIN_LEN", "50"))
	env.SimilarityThreshold = cast.ToFloat32(getOrDefault("SIMILARITY_THRESHOLD", "0.88"))
}

func getOrDefault(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// package main

// import (
// 	"checker/go-backend/pkg/chunker"
// 	"checker/go-backend/pkg/embedder"
// 	"checker/go-backend/pkg/extractor"
// 	"checker/go-backend/pkg/qdrantclient"
// 	"checker/go-backend/pkg/storage"
// 	"checker/go-backend/plagiarism/service"
// 	transport "checker/go-backend/transport/service"
// 	"log"

// 	"git.sriss.uz/shared/shared_service/sharedutil"
// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// type Env struct {
// 	ServerPort          string
// 	UploadDir           string
// 	QdrantHost          string
// 	QdrantGRPCPort      int
// 	QdrantCollection    string
// 	QdrantVectorSize    uint64
// 	EmbedderURL         string
// 	EmbedderBatch       int
// 	EmbedderTimeout     int
// 	ChunkSize           int
// 	ChunkOverlap        int
// 	ChunkMinLen         int
// 	SimilarityThreshold float32
// }

// func main() {
// 	var env Env
// 	loadEnv(&env)

// 	stor, err := storage.New(env.UploadDir)
// 	if err != nil {
// 		log.Fatalf("storage: %v", err)
// 	}

// 	qdrant, err := qdrantclient.New(env.QdrantHost, env.QdrantGRPCPort, env.QdrantCollection, env.QdrantVectorSize)
// 	if err != nil {
// 		log.Fatalf("qdrant: %v", err)
// 	}

// 	svc := service.New(service.Deps{
// 		Extractor: extractor.New(),
// 		Chunker:   chunker.New(env.ChunkSize, env.ChunkOverlap, env.ChunkMinLen),
// 		Embedder:  embedder.New(env.EmbedderURL, env.EmbedderBatch, env.EmbedderTimeout),
// 		Qdrant:    qdrant,
// 		Storage:   stor,
// 		Threshold: env.SimilarityThreshold,
// 	})

// 	var router = echo.New()
// 	setMiddlewares(router)
// 	transport.NewService(router, svc)
// 	runHTTPServerOnAddr(router, env.ServerPort)
// }

// func runHTTPServerOnAddr(router *echo.Echo, port string) {
// 	log.Printf("Server starting on :%s", port)
// 	if err := router.Start(":" + port); err != nil {
// 		panic(err)
// 	}
// }

// func setMiddlewares(router *echo.Echo) {
// 	router.HideBanner = true
// 	router.Use(middleware.RemoveTrailingSlash())
// 	router.Use(middleware.RequestID())
// 	router.Use(middleware.Recover())
// }

// func loadEnv(e *Env) {
// 	if err := sharedutil.Load(e, ".env"); err != nil {
// 		log.Println("Warning: .env not found, using environment variables")
// 		sharedutil.MustLoad(e)
// 	}
// }
