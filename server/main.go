package main

import (
	"log"
	"net/http"

	"go-grpc-chat/config"
	"go-grpc-chat/handlers"
)

func main() {
	// Connect to MongoDB
	err := config.ConnectMongo()
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// REST API routes with CORS middleware
	http.HandleFunc("/register", withCORS(handlers.Register))
	http.HandleFunc("/login", withCORS(handlers.Login))

	// WebSocket with JWT Auth
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Println("✅ Server is running at http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil)) // ✅ changed from ":8080"
}

// CORS middleware
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (not for production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
}
