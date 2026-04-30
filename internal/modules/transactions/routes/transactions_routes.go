package transactionsroutes

import (
	categoriesrepository "my_finance/internal/modules/categories/repository"
	transactionscontroller "my_finance/internal/modules/transactions/controller"
	transactionsrepository "my_finance/internal/modules/transactions/repository"
	transactionsservices "my_finance/internal/modules/transactions/services"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTransactionsRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := transactionsrepository.NewTransactionsRepository(db)
	categorRepo := categoriesrepository.NewCategoryRepository(db)

	service := transactionsservices.NewTransactionsService(repo, categorRepo)
	controller := transactionscontroller.NewTransactionsController(service)

	r.HandleFunc("GET  /transactions", controller.GetAll)
	r.HandleFunc("POST /transactions", controller.Create)
	r.HandleFunc("GET  /transaction/{id}", controller.FindById)

	r.HandleFunc("PATCH /transaction/{id}/pay", controller.Pay)
	//r.HandleFunc("PATCH /transaction/{id}/cancel")
	//r.HandleFunc("PATCH /transaction/{id}/peding")
}
