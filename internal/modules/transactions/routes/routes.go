package transactionsroutes

import (
	transactionscontroller "my_finance/internal/modules/transactions/controller"
	transactionsrepository "my_finance/internal/modules/transactions/repository"
	transactionsservices "my_finance/internal/modules/transactions/services"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTransactionsRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := transactionsrepository.NewTransactionsRepository(db)
	service := transactionsservices.NewTransactionsService(repo)
	controller := transactionscontroller.NewTransactionsController(service)

	r.HandleFunc("GET  /transactions", controller.GetAll)
	r.HandleFunc("POST /transactions", controller.Create)
	r.HandleFunc("GET  /transactions/{id}", controller.FindById)

}
