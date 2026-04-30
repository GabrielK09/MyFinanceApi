package categoriescontroller

import (
	"context"
	"encoding/json"
	"my_finance/internal/response"
	categories_model "my_finance/models/categories"
	"net/http"
	"strconv"
)

type CategoryService interface {
	GetAll(context.Context) ([]categories_model.CategoryModel, error)
	Create(context.Context, categories_model.CategoryModel) (int, error)
	Update(context.Context, categories_model.CategoryModel) (categories_model.CategoryModel, error)
	FindById(context.Context, int) (categories_model.CategoryModel, error)
	Delete(context.Context, int) error
	Active(context.Context, int) error
}

type CategoryController struct {
	service CategoryService
}

func NewCategoryController(service CategoryService) *CategoryController {
	return &CategoryController{
		service: service,
	}
}

func (c *CategoryController) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := c.service.GetAll(r.Context())

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao retornar todas as categorias",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Todas as categorias",
		map[string]any{"categories": categories},
	))
}

func (c *CategoryController) Create(w http.ResponseWriter, r *http.Request) {
	var categoryPayLoad categories_model.CategoryModel

	if err := json.NewDecoder(r.Body).Decode(&categoryPayLoad); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao ler os dados",
			map[string]any{"error": err.Error()},
		))
		return
	}

	categoryId, err := c.service.Create(r.Context(), categoryPayLoad)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao cadastrar a categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.SuccessResponse(
		"Categoria cadastrada com sucesso!",
		map[string]any{"category_id": categoryId},
	))
}

func (c *CategoryController) Update(w http.ResponseWriter, r *http.Request) {
	var categoryPayLoad categories_model.CategoryModel

	param := r.PathValue("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&categoryPayLoad); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao ler os dados",
			map[string]any{"error": err.Error()},
		))
		return
	}

	categoryPayLoad.Id = id

	category, err := c.service.Update(r.Context(), categoryPayLoad)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao alterar a categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.SuccessResponse(
		"Categoria alterada com sucesso!",
		map[string]any{"categories": category},
	))

}

func (c *CategoryController) FindById(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	category, err := c.service.FindById(r.Context(), id)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao desativar a categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Categoria desativada com sucesso!",
		map[string]any{"category": category},
	))
}

func (c *CategoryController) Delete(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao desativar a categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusNoContent, response.SuccessResponse(
		"Categoria desativada com sucesso!",
		map[string]any{},
	))
}

func (c *CategoryController) Active(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Active(r.Context(), id); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao ativar a categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusNoContent, response.SuccessResponse(
		"Categoria ativada com sucesso!",
		map[string]any{},
	))
}
