package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielfs/paredao/backend/handlers"
	"github.com/danielfs/paredao/backend/repositories"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	repositories.InitDB()
	defer repositories.CloseDB()

	// Initialize Redis client
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	repositories.InitRedis(redisHost, redisPort)
	defer repositories.CloseRedis()

	r := mux.NewRouter()

	// Participante routes
	r.HandleFunc("/participantes", handlers.GetParticipantes).Methods("GET")
	r.HandleFunc("/participantes/{id}", handlers.GetParticipante).Methods("GET")
	r.HandleFunc("/participantes", handlers.CreateParticipante).Methods("POST")
	r.HandleFunc("/participantes/{id}", handlers.UpdateParticipante).Methods("PUT")
	r.HandleFunc("/participantes/{id}", handlers.DeleteParticipante).Methods("DELETE")

	// Votacao routes
	r.HandleFunc("/votacoes", handlers.GetVotacoes).Methods("GET")
	r.HandleFunc("/votacoes/{id}", handlers.GetVotacao).Methods("GET")
	r.HandleFunc("/votacoes", handlers.CreateVotacao).Methods("POST")
	r.HandleFunc("/votacoes/{id}", handlers.UpdateVotacao).Methods("PUT")
	r.HandleFunc("/votacoes/{id}", handlers.DeleteVotacao).Methods("DELETE")
	r.HandleFunc("/votacoes/{id}/participantes", handlers.GetVotacaoParticipantes).Methods("GET")
	r.HandleFunc("/votacoes/{id}/participantes", handlers.AddParticipanteToVotacao).Methods("POST")

	// Voto routes
	r.HandleFunc("/votos", handlers.GetVotos).Methods("GET")
	r.HandleFunc("/votos/{participanteId}/{votacaoId}", handlers.GetVoto).Methods("GET")
	r.HandleFunc("/votos", handlers.CreateVoto).Methods("POST")

	// Estatisticas routes
	r.HandleFunc("/estatisticas/votacoes/{id}/total", handlers.GetVotacaoTotal).Methods("GET")
	r.HandleFunc("/estatisticas/votacoes/{id}/participantes", handlers.GetVotacaoTotalByParticipante).Methods("GET")
	r.HandleFunc("/estatisticas/votacoes/{id}/hourly", handlers.GetVotacaoTotalByHour).Methods("GET")

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Add CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}

	// Apply CORS middleware
	handler := corsMiddleware(r)

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on port 8080...")
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")
}
