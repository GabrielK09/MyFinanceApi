package routes

import (
	"log"
	"net/http"
	"os"

	transactionsroutes "my_finance/internal/modules/transactions/routes"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartServer(db *pgxpool.Pool) {
	r := http.NewServeMux()
	api := http.NewServeMux()

	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	transactionsroutes.RegisterTransactionsRoutes(api, db)
	r.Handle("/api/", http.StripPrefix("/api", api))

	log.Printf("Servidor rodando em http://localhost:%s/api", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
