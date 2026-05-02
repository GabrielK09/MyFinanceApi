package transactions_model

import (
	transactionsconstants "my_finance/internal/constants/transactions"
	"strings"
	"time"
)

type ValidationErrors map[string]string

// time.Time layout RFC3339

type TransactionsModel struct {
	Id          int                                      `json:"id"`
	CategoryId  int                                      `json:"category_id"`
	Origin      transactionsconstants.TransactionsOrigin `json:"origin"`
	OriginId    *int                                     `json:"origin_id"`
	Type        transactionsconstants.ConstantType       `json:"type"`
	Description string                                   `json:"description"`
	Amount      float64                                  `json:"amount"`
	DueDate     *time.Time                               `json:"due_date"`
	PaidAt      *time.Time                               `json:"paid_at"`
	Status      transactionsconstants.Status             `json:"status"`
	Notes       string                                   `json:"notes"`
	CreatedAt   time.Time                                `json:"created_at"`
	UpdatedAt   time.Time                                `json:"updated_at"`
	DeletedAt   *time.Time                               `json:"deleted_at"`
}

func (e ValidationErrors) Error() string {
	return "campos obrigatórios ausentes ou inválidos"
}

func (t TransactionsModel) Validate() error {
	errors := ValidationErrors{}

	if t.CategoryId <= 0 {
		errors["category_id"] = "a categoria é obrigatória"
	}

	if !t.Type.IsValidType() {
		errors["type"] = "o tipo deve ser entre Entrada ou Saída"
	}

	if strings.TrimSpace(t.Description) == "" {
		errors["description"] = "a descrição é obrigatória"
	}

	if len(t.Description) > 150 {
		errors["description"] = "a descrição deve ter no máximo 150 caracteres"
	}

	if t.Amount <= 0 {
		errors["amount"] = "o valor deve ser maior que zero"
	}

	if t.PaidAt != nil && t.Status == transactionsconstants.StatusPendente {
		errors["status"] = "uma transação não pode ser possuir o status pendente junto a uma data de pagamento"
	}

	if t.Origin == "" {
		errors["origin"] = "a origem de uma transação é obrigatória"
	}

	if !t.Origin.IsValidTransactionsOrigin() {
		errors["origin"] = "a origem de uma transação deve ser manual ou de um recebimento"
	}

	if t.Origin == transactionsconstants.ManualOrigin && t.OriginId != nil {
		errors["origin_id"] = "origem manual não deve possuir um identificador de um recebimento"
	}

	if t.Origin != "" && t.Origin != transactionsconstants.ManualOrigin {
		if t.OriginId == nil || *t.OriginId <= 0 {
			errors["origin_id"] = "origem não manual deve possuir um identificador de um recebimento"
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
