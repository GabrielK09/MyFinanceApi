package transactionscontroller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"my_finance/internal/apperrors"
	getparamid "my_finance/internal/helpers/get_param_id"
	loggerHelper "my_finance/internal/logger"
	"my_finance/internal/response"
	transactions_model "my_finance/models/transactions"
	"net/http"
)

type TransactionsService interface {
	GetAll(context.Context) ([]transactions_model.TransactionsModel, error)
	Create(context.Context, transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error)
	Update(context.Context, transactions_model.TransactionsModel) (transactions_model.TransactionsModel, error)
	FindById(context.Context, int) (transactions_model.TransactionsModel, error)
	Pay(context.Context, int) error
	Cancel(context.Context, int) error
	Delete(context.Context, int) error
}

type TransactionsController struct {
	service TransactionsService
}

func NewTransactionsController(service TransactionsService) *TransactionsController {
	return &TransactionsController{
		service: service,
	}
}

func (c *TransactionsController) GetAll(w http.ResponseWriter, r *http.Request) {
	transactions, err := c.service.GetAll(r.Context())

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao retornar todas as transações",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Todas as transações",
		map[string]any{"transactions": transactions},
	))
}

func (c *TransactionsController) Create(w http.ResponseWriter, r *http.Request) {
	var payLoad transactions_model.TransactionsModel

	err := json.NewDecoder(r.Body).Decode(&payLoad)

	if err != nil && !errors.Is(err, io.EOF) {
		loggerHelper.ErrorLogger.Println("JSON inválido", err)

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"JSON inválido",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := payLoad.Validate(); err != nil {
		var fields transactions_model.ValidationErrors

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

	transactionData, err := c.service.Create(r.Context(), payLoad)

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
		map[string]any{"transaction": transactionData},
	))
}

func (c *TransactionsController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	var payLoad transactions_model.TransactionsModel
	err = json.NewDecoder(r.Body).Decode(&payLoad)

	if err != nil && !errors.Is(err, io.EOF) {
		loggerHelper.ErrorLogger.Println("JSON inválido", err)

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"JSON inválido",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := payLoad.Validate(); err != nil {
		var fields transactions_model.ValidationErrors

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

	payLoad.Id = id
	transactionData, err := c.service.Update(r.Context(), payLoad)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Erro ao alterar a transação.",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao alterar a transação.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.SuccessResponse(
		"Transação alterada com sucesso!",
		map[string]any{"transaction": transactionData},
	))
}

func (c *TransactionsController) FindById(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	transaction, err := c.service.FindById(r.Context(), id)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Transação não localizada",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao localizar a transação",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Transação localizada com sucesso!",
		map[string]any{"transaction": transaction},
	))
}

func (c *TransactionsController) Pay(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da categoria",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Pay(r.Context(), id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.WriteJSON(w, http.StatusNotFound, response.ErrorResponse(
				"Transação não localizada",
				map[string]any{"error": err.Error()},
			))
			return
		}

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

func (c *TransactionsController) Cancel(w http.ResponseWriter, r *http.Request) {
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
				"Transação não localizada",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao cancelar a transação.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Transação cancelada com sucesso!",
		map[string]any{},
	))
}

func (c *TransactionsController) Delete(w http.ResponseWriter, r *http.Request) {
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
				"Transação não localizada",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao deletar a transação.",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Transação deletada com sucesso!",
		map[string]any{},
	))
}
