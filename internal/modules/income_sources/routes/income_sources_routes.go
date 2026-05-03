package incomesourcesroutes

import (
	incomesourcescontroller "my_finance/internal/modules/income_sources/controller"
	incomesourcesrepository "my_finance/internal/modules/income_sources/repository"
	incomesourcesservices "my_finance/internal/modules/income_sources/services"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterIncomeSourcesRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := incomesourcesrepository.NewIncomeSourcesRepository(db)
	service := incomesourcesservices.NewIncomeSourcesService(repo)
	controller := incomesourcescontroller.NewIncomeSourcesController(service)

	r.HandleFunc("GET /income-sources", controller.GetAll)
	r.HandleFunc("GET /income-sources/{id}", controller.FindById)
	r.HandleFunc("POST /income-sources", controller.Create)
	r.HandleFunc("PUT /income-source/{id}", controller.Update)
	r.HandleFunc("DELETE /income-source/delete/{id}", controller.Delete)
	r.HandleFunc("PATCH /income-source/active/{id}", controller.Active)
}
