package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Content-Type", "text/html")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf(`
  <html>

  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>

  </html>
      `, cfg.fileserverHits.Load())))
}
