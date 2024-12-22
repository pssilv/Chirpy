package main

import (
	"database/sql"
	"os"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/pssilv/Chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {
  godotenv.Load()

  dbURL := os.Getenv("DB_URL")
  if dbURL == "" {
    log.Fatal("DB_URL must be set")
  }

  platform := os.Getenv("PLATFORM")
  if platform == "" {
    log.Fatal("PLATFORM must be set")
  }

  jwtSecret := os.Getenv("JWT_SECRET")
  if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is not set")
  }

  polkaKey := os.Getenv("POLKA_KEY")
  if polkaKey == "" {
    log.Fatal("POLKA_KEY environment variable is not set")
  }

  dbConn, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Fatalf("Failed to open database: %v", err)
  }

  dbQueries := database.New(dbConn)


  const filePathRoot = "."
  const port = "8080"

  apiCfg := apiConfig{
    fileserverHits: atomic.Int32{},
    db:             dbQueries,
    platform:       platform,   
    jwtSecret: jwtSecret,
    polkaKey: polkaKey,
  }

  if dbURL == "" {
    log.Fatal("DB_URL must be set")
  }


  mux := http.NewServeMux()
  handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
  mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

  mux.HandleFunc("GET /api/healthz", handlerReadiness)

  mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
  mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
  mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

  mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
  mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

  mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
  mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

  mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
  mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
  mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
  mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp) 

  mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

  server := &http.Server{
    Handler: mux,
    Addr: ":" + port,
  }

  log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)

  log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Add(1)
    next.ServeHTTP(w, r)
  })
}

type apiConfig struct {
  fileserverHits atomic.Int32
  db             *database.Queries
  platform       string
  jwtSecret      string
  polkaKey         string
}

