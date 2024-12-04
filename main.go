package main

import (
	"log"
	"net/http"
)

func main() {
  serverMux := http.NewServeMux()
  server := http.Server{
    Handler: serverMux,
    Addr: ":8080",
  }

  if err := server.ListenAndServe(); err != nil {
    log.Fatalf("Unable to start server: %v", err)
  }
}
