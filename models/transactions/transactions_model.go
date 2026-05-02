package transactions_model

import (
	"fmt"
	"my_finance/internal/constants"
	"strings"
	"time"
)

// time.Time layout RFC3339

type TransactionsModel struct {
	Id          int                    `json:"id"`
	CategoryId  int                    `json:"category_id"`
	Origin      *string                `json:"origin"`
	OriginId    *int                   `json:"origin_id"`
	Type        constants.ConstantType `json:"type"`
	Description string                 `json:"description"`
	Amount      float64                `json:"amount"`
	DueDate     *time.Time             `json:"due_date"`
	PaidAt      *time.Time             `json:"paid_at"`
	Status      constants.Status       `json:"status"`
	Notes       string                 `json:"notes"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   *time.Time             `json:"deleted_at"`
}

func (t TransactionsModel) Validate() error {
	if t.CategoryId <= 0 {
		return fmt.Errorf("A categoria é obrigatória.")
	}

	if !t.Type.IsValidType() {
		return fmt.Errorf("O tipo deve ser entre Entrada ou Saída.")
	}

	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("A descrição é obrigatória.")
	}

	if len(t.Description) > 150 {
		return fmt.Errorf("A descrição deve ter no máximo 150 caracteres.")
	}

	if t.Amount <= 0 {
		return fmt.Errorf("O valor deve ser maior que zero.")
	}

	if t.PaidAt != nil && t.Status == constants.StatusPendente {
		return fmt.Errorf("Uma transação não pode ser possuir o status pendente junto a uma data de pagamento.")
	}

	return nil
}
