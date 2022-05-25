package main

import "os"

func main() {
	logger := NewLogrus("26ab3bf6-1cd3-4050-8650-aab5e08d39db-ls.logit.io", "30823")
	server := NewServer(logger)
	server.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := server.Start(port); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
