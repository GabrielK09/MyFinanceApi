package incomereceiptscontroller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"my_finance/internal/apperrors"
	getparamid "my_finance/internal/helpers/get_param_id"
	"my_finance/internal/response"
	incomereceiptsmodel "my_finance/models/income_receipts"
	"net/http"
)

type IncomeReceiptsService interface {
	GetAll(context.Context) ([]incomereceiptsmodel.IncomeReceiptsModel, error)
	Create(context.Context, incomereceiptsmodel.IncomeReceiptsModel) (incomereceiptsmodel.IncomeReceiptsModel, error)
	FindById(context.Context, int) (*incomereceiptsmodel.IncomeReceiptsModel, error)
	Cancel(context.Context, int) error
	Delete(context.Context, int) error
}

type IncomeReceiptsController struct {
	service IncomeReceiptsService
}

func NewIncomeReceiptsController(service IncomeReceiptsService) *IncomeReceiptsController {
	return &IncomeReceiptsController{
		service: service,
	}
}

func (c *IncomeReceiptsController) GetAll(w http.ResponseWriter, r *http.Request) {
	incomeReceipts, err := c.service.GetAll(r.Context())

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao retornar todos os recebimento.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Todos os recebimentos.",
		map[string]any{"income_receipts": incomeReceipts},
	))
}

func (c *IncomeReceiptsController) Create(w http.ResponseWriter, r *http.Request) {
	var payLoad incomereceiptsmodel.IncomeReceiptsModel

	err := json.NewDecoder(r.Body).Decode(&payLoad)

	if err != nil && !errors.Is(err, io.EOF) {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"JSON inválido",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := payLoad.Validate(); err != nil {
		var fields incomereceiptsmodel.ValidationErrors

		if errors.As(err, &fields) {
			response.WriteJSON(w, http.StatusUnprocessableEntity, response.ErrorResponse(
				"Campos ausentes",
				map[string]any{"error": fields},
			))
			return
		}

		response.WriteJSON(w, http.StatusUnprocessableEntity, response.ErrorResponse(
			"Campos ausentes",
			map[string]any{"error": err.Error()},
		))
		return
	}

	incomeReceipt, err := c.service.Create(r.Context(), payLoad)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Fonte de renda não localizada",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao cadastrar o recebimento",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.SuccessResponse(
		"Recebimento cadastrado com sucesso",
		map[string]any{"income_receipt": incomeReceipt},
	))
}

func (c *IncomeReceiptsController) Cancel(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Cancel(r.Context(), id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Recebimento não localizado.",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao cancelar o recebimento.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Recebimento cancelado com sucesso!.",
		map[string]any{},
	))
}

func (c *IncomeReceiptsController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Recebimento não localizado.",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao deletar o recebimento.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *IncomeReceiptsController) FindById(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	incomeReceipt, err := c.service.FindById(r.Context(), id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Recebimento não localizada",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao localizar o recebimento",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Transação localizada com sucesso!",
		map[string]any{"income_receipt": incomeReceipt},
	))
}
