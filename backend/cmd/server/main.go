package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	httpserver "github.com/alexmented/ai-for-developers-project-386/backend/internal/http"
)

func main() {
	addr := ":4010"

	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.HasPrefix(port, ":") {
			addr = port
		} else {
			addr = ":" + port
		}
	} else if backendAddr := strings.TrimSpace(os.Getenv("BACKEND_ADDR")); backendAddr != "" {
		addr = backendAddr
	}

	handler := httpserver.NewDefaultHandler()
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	log.Printf("backend listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
