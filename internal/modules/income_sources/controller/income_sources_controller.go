package incomesourcescontroller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"my_finance/internal/apperrors"
	getparamid "my_finance/internal/helpers/get_param_id"
	loggerHelper "my_finance/internal/logger"
	"my_finance/internal/response"
	incomesourcesmodel "my_finance/models/income_sources"
	"net/http"
)

type IncomeSourcesService interface {
	GetAll(context.Context) ([]incomesourcesmodel.IncomeSourcesModel, error)
	Create(context.Context, incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error)
	Update(context.Context, incomesourcesmodel.IncomeSourcesModel) (incomesourcesmodel.IncomeSourcesModel, error)
	Delete(context.Context, int) error
	Active(context.Context, int) error
}

type IncomeSourcesController struct {
	service IncomeSourcesService
}

func NewIncomeSourcesController(service IncomeSourcesService) *IncomeSourcesController {
	return &IncomeSourcesController{
		service: service,
	}
}

func (c *IncomeSourcesController) GetAll(w http.ResponseWriter, r *http.Request) {
	incomeSources, err := c.service.GetAll(r.Context())

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao retornar todas as fontes de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Todas as fontes de renda",
		map[string]any{"income_sources": incomeSources},
	))
}

func (c *IncomeSourcesController) Create(w http.ResponseWriter, r *http.Request) {
	var payLoad incomesourcesmodel.IncomeSourcesModel

	err := json.NewDecoder(r.Body).Decode(&payLoad)

	if err != nil && !errors.Is(err, io.EOF) {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"JSON inválido",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := payLoad.Validate(); err != nil {
		var fields incomesourcesmodel.ValidationErrors

		if errors.As(err, &fields) {
			loggerHelper.ErrorLogger.Println("Pelo validate do model:", err)

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

	incomeSource, err := c.service.Create(r.Context(), payLoad)

	if err != nil {
		if errors.Is(err, apperrors.ErrUnprocessableEntity) {
			response.WriteJSON(w, http.StatusUnprocessableEntity, response.ErrorResponse(
				"Campos ausentes",
				map[string]any{"error": err.Error()},
			))
			return
		}

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao cadastrar a fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.SuccessResponse(
		"Fonte de renda cadastrada com sucesso!",
		map[string]any{"income_source": incomeSource},
	))
}

func (c *IncomeSourcesController) Update(w http.ResponseWriter, r *http.Request) {
	var payLoad incomesourcesmodel.IncomeSourcesModel

	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		loggerHelper.ErrorLogger.Println("Erro ao identificar o parametro da fonte de renda", err)

		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	err = json.NewDecoder(r.Body).Decode(&payLoad)

	if err != nil && !errors.Is(err, io.EOF) {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"JSON inválido",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := payLoad.Validate(); err != nil {
		var fields incomesourcesmodel.ValidationErrors

		if errors.As(err, &fields) {
			loggerHelper.ErrorLogger.Println("Pelo validate do model:", err)

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

	loggerHelper.InfoLogger.Println("Dados:", payLoad)

	updated, err := c.service.Update(r.Context(), payLoad)

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao alterar a fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.SuccessResponse(
		"Fonte de renda alterada com sucesso!",
		map[string]any{"income_source": updated},
	))
}

func (c *IncomeSourcesController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Delete(r.Context(), id); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao desativar a fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *IncomeSourcesController) Active(w http.ResponseWriter, r *http.Request) {
	id, err := getparamid.HandleParamIdUrl(r.PathValue("id"))

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao identificar o parametro da fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	if err := c.service.Active(r.Context(), id); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.ErrorResponse(
			"Erro ao ativar a fonte de renda",
			map[string]any{"error": err.Error()},
		))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
