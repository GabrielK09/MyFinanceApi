package routes

import (
	"log"
	"net/http"
	"os"

	categoriesroutes "my_finance/internal/modules/categories/routes"
	incomereceiptsroutes "my_finance/internal/modules/income_receipts/routes"
	incomesourcesroutes "my_finance/internal/modules/income_sources/routes"
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
	categoriesroutes.RegisterCategoriesRoutes(api, db)
	incomesourcesroutes.RegisterIncomeSourcesRoutes(api, db)
	incomereceiptsroutes.RegisterIncomeReceiptsRoutes(api, db)

	r.Handle("/api/", http.StripPrefix("/api", api))

	log.Printf("Servidor rodando em http://localhost:%s/api", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
