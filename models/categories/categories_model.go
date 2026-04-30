package categories_model

import (
	"fmt"
	"strings"
	"time"
)

type CategoryModel struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	MonthlyLimit float64    `json:"monthly_limit"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

func (c CategoryModel) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("O nome da categoria é obrigatória.")
	}

	if len(c.Name) > 100 {
		return fmt.Errorf("O nome da categoria deve ter no máximo 100 caracteres.")
	}

	if len(c.Type) > 100 {
		return fmt.Errorf("O tipo da categoria deve ter no máximo 100 caracteres.")
	}

	if strings.TrimSpace(c.Type) == "" {
		return fmt.Errorf("O tipo da categoria é obrigatória.")
	}

	if c.MonthlyLimit <= 0 {
		return fmt.Errorf("O limite da categoria precisa ser maior que zero.")
	}

	return nil
}
