package incomereceiptsroutes

import (
	incomereceiptscontroller "my_finance/internal/modules/income_receipts/controller"
	incomereceiptsrepository "my_finance/internal/modules/income_receipts/repository"
	incomereceiptsservice "my_finance/internal/modules/income_receipts/services"
	incomesourcesrepository "my_finance/internal/modules/income_sources/repository"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterIncomeReceiptsRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	incomeSourcesRepo := incomesourcesrepository.NewIncomeSourcesRepository(db)
	repo := incomereceiptsrepository.NewIncomeReceiptsRepository(db)

	service := incomereceiptsservice.NewIncomeReceiptsService(repo, incomeSourcesRepo)
	controller := incomereceiptscontroller.NewIncomeReceiptsController(service)

	r.HandleFunc("GET /income-receipts", controller.GetAll)
	r.HandleFunc("GET /income-receipts/{id}", controller.FindById)
	r.HandleFunc("POST /income-receipts", controller.Create)
	r.HandleFunc("PATCH /income-receipts/cancel/{id}", controller.Cancel)
	r.HandleFunc("DELETE /income-receipts/delete/{id}", controller.Delete)
}
