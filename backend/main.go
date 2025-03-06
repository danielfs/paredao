package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/danielfs/paredao/backend/handlers"
	"github.com/danielfs/paredao/backend/repositories"
)

func main() {
	// Inicializa conexão com o banco de dados
	repositories.InitDB()
	defer repositories.CloseDB()

	// Inicializa cliente Redis
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	repositories.InitRedis(redisHost, redisPort)
	defer repositories.CloseRedis()

	r := mux.NewRouter()

	// Rotas de Participante
	r.HandleFunc("/participantes", handlers.GetParticipantes).Methods("GET")
	r.HandleFunc("/participantes/{id}", handlers.GetParticipante).Methods("GET")
	r.HandleFunc("/participantes", handlers.CreateParticipante).Methods("POST")
	r.HandleFunc("/participantes/{id}", handlers.UpdateParticipante).Methods("PUT")
	r.HandleFunc("/participantes/{id}", handlers.DeleteParticipante).Methods("DELETE")

	// Rotas de Votação
	r.HandleFunc("/votacoes", handlers.GetVotacoes).Methods("GET")
	r.HandleFunc("/votacoes/{id}", handlers.GetVotacao).Methods("GET")
	r.HandleFunc("/votacoes", handlers.CreateVotacao).Methods("POST")
	r.HandleFunc("/votacoes/{id}", handlers.UpdateVotacao).Methods("PUT")
	r.HandleFunc("/votacoes/{id}", handlers.DeleteVotacao).Methods("DELETE")
	r.HandleFunc("/votacoes/{id}/participantes", handlers.GetVotacaoParticipantes).Methods("GET")
	r.HandleFunc("/votacoes/{id}/participantes", handlers.AddParticipanteToVotacao).Methods("POST")

	// Rotas de Voto
	r.HandleFunc("/votos", handlers.GetVotos).Methods("GET")
	r.HandleFunc("/votos/{participanteId}/{votacaoId}", handlers.GetVoto).Methods("GET")
	r.HandleFunc("/votos", handlers.CreateVoto).Methods("POST")

	// Rotas de Estatísticas
	r.HandleFunc("/estatisticas/votacoes/{id}/total", handlers.GetVotacaoTotal).Methods("GET")
	r.HandleFunc("/estatisticas/votacoes/{id}/participantes", handlers.GetVotacaoTotalByParticipante).Methods("GET")
	r.HandleFunc("/estatisticas/votacoes/{id}/hourly", handlers.GetVotacaoTotalByHour).Methods("GET")

	// Configura encerramento gracioso
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Adiciona middleware CORS
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Define cabeçalhos CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Trata requisições preflight
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Chama o próximo handler
			next.ServeHTTP(w, r)
		})
	}

	// Aplica middleware CORS
	handler := corsMiddleware(r)

	// Cria um servidor com timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Inicia servidor em uma goroutine
	go func() {
		log.Println("Server starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Aguarda sinal de interrupção
	<-stop
	log.Println("Shutting down server...")

	// Cria um prazo para o encerramento do servidor
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Encerra o servidor graciosamente
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		// Não use os.Exit aqui para garantir que as declarações defer sejam executadas
	}

	log.Println("Server exited properly")
}
