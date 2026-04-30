package categoriesroutes

import (
	categoriescontroller "my_finance/internal/modules/categories/controller"
	categoriesrepository "my_finance/internal/modules/categories/repository"
	categoriesservice "my_finance/internal/modules/categories/services"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterCategoriesRoutes(r *http.ServeMux, db *pgxpool.Pool) {
	repo := categoriesrepository.NewCategoryRepository(db)
	service := categoriesservice.NewCategoryService(repo)
	controller := categoriescontroller.NewCategoryController(service)

	r.HandleFunc("GET /categories", controller.GetAll)
	r.HandleFunc("GET /category/{id}", controller.FindById)
	r.HandleFunc("POST /category", controller.Create)
	r.HandleFunc("PUT /category/{id}", controller.Update)
	r.HandleFunc("DELETE /category/delete/{id}", controller.Delete)
	r.HandleFunc("PATCH /category/active/{id}", controller.Active)
}
