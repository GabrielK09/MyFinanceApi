package transactionscontroller

import (
	"context"
	"encoding/json"
	"errors"
	"my_finance/internal/apperrors"
	getparamid "my_finance/internal/helpers/get_param_id"
	loggerHelper "my_finance/internal/logger"
	"my_finance/internal/response"
	transactions_model "my_finance/models/transactions"
	"net/http"
	"time"
)

type TransactionsService interface {
	GetAll(context.Context) ([]transactions_model.TransactionsModel, error)
	Create(context.Context, transactions_model.TransactionsModel) (int, error)
	FindById(context.Context, int) (transactions_model.TransactionsModel, error)
	Pay(context.Context, int, time.Time) error
}

type TransactionsController struct {
	service TransactionsService
}

func NewTransactionsController(service TransactionsService) *TransactionsController {
	return &TransactionsController{
		service: service,
	}
}

func (t *TransactionsController) GetAll(w http.ResponseWriter, r *http.Request) {
	transactions, err := t.service.GetAll(r.Context())

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao retornar todas as transações",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
		"Todas as transações",
		map[string]any{"transactions": transactions},
	))
}

func (t *TransactionsController) Create(w http.ResponseWriter, r *http.Request) {
	var transaction transactions_model.TransactionsModel

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao ler os dados da request:", err)

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao ler os dados da transação.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	transactionId, err := t.service.Create(r.Context(), transaction)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Erro ao criar a transação.",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao criar a transação.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.SuccessResponse(
		"Transação cadastrada com sucesso!",
		map[string]any{"transaction_id": transactionId},
	))
}

func (t *TransactionsController) FindById(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	transaction, err := t.service.FindById(r.Context(), id)

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Transação localizada com sucesso!",
		map[string]any{"transaction": transaction},
	))
}

func (t *TransactionsController) Pay(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := t.service.Pay(r.Context(), id, time.Now()); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao pagar a transação.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Transação paga com sucesso!",
		map[string]any{},
	))
}
