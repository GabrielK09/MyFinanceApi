package incomereceiptsmodel

import (
	incomereceiptssconstants "my_finance/internal/constants/income_receipts"
	"time"
)

type ValidationErrors map[string]string

type IncomeReceiptsModel struct {
	Id             int                             `json:"id"`
	IncomeSourceId int                             `json:"income_source_id"`
	Description    string                          `json:"description"`
	Amount         float64                         `json:"amount"`
	ReceivedAt     time.Time                       `json:"received_at"`
	ReferenceMonth *int                            `json:"reference_month"`
	ReferenceYear  *int                            `json:"reference_year"`
	Notes          *string                         `json:"notes"`
	Status         incomereceiptssconstants.Status `json:"status"`
	CreatedAt      time.Time                       `json:"created_at"`
	UpdatedAt      time.Time                       `json:"updated_at"`
	DeletedAt      *time.Time                      `json:"deleted_at"`
}

func (e ValidationErrors) Error() string {
	return "campos obrigatórios ausentes ou inválidos"
}

func (i IncomeReceiptsModel) Validate() error {
	errors := ValidationErrors{}

	if i.IncomeSourceId <= 0 {
		errors["income_source_id"] = "o identificador da fonte de renda não pode ser menor que zero"
	}

	if i.Description == "" {
		errors["description"] = "a descrição do recebimento de renda é obrigatória"
	}

	if len(i.Description) > 100 {
		errors["description"] = "a descrição do recebimento deve ter no máximo 100 caracteres"
	}

	if i.Amount <= 0 {
		errors["amount"] = "o valor do recebimento não pode ser menor que zero"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
