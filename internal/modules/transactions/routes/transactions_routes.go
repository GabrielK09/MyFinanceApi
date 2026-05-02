package transactionsroutes

import (
	categoriesrepository "my_finance/internal/modules/categories/repository"
	incomereceiptsrepository "my_finance/internal/modules/income_receipts/repository"
	transactionscontroller "my_finance/internal/modules/transactions/controller"
	transactionsrepository "my_finance/internal/modules/transactions/repository"
	transactionsservices "my_finance/internal/modules/transactions/services"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTransactionsRoutes(r *http.ServeMux, db *pgxpool.Pool) {

	repo := transactionsrepository.NewTransactionsRepository(db)
	categoryRepo := categoriesrepository.NewCategoryRepository(db)
	incomeReceiptsRepo := incomereceiptsrepository.NewIncomeReceiptsRepository(db)

	service := transactionsservices.NewTransactionsService(repo, categoryRepo, incomeReceiptsRepo)
	controller := transactionscontroller.NewTransactionsController(service)

	r.HandleFunc("GET  /transactions", controller.GetAll)
	r.HandleFunc("POST /transactions", controller.Create)
	r.HandleFunc("PUT /transactions/{id}", controller.Update)
	r.HandleFunc("GET  /transactions/{id}", controller.FindById)

	r.HandleFunc("PATCH /transaction/{id}/pay", controller.Pay)
	r.HandleFunc("PATCH /transaction/{id}/cancel", controller.Cancel)
	r.HandleFunc("DELETE /transaction/{id}", controller.Delete)
}
