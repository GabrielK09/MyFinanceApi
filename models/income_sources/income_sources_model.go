package incomesourcesmodel

import "time"

type ValidationErrors map[string]string

type IncomeSourcesModel struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (e ValidationErrors) Error() string {
	return "campos obrigatórios ausentes ou inválidos"
}

func (i IncomeSourcesModel) Validate() error {
	errors := ValidationErrors{}

	if i.Name == "" {
		errors["name"] = "o nome da fonte de renda é obrigatória"
	}

	if len(i.Name) > 100 {
		errors["name"] = "o nome da fonte de renda deve ter no máximo 100 caracteres"
	}

	if i.Description == "" {
		errors["description"] = "a descrição da fonte de renda é obrigatória"
	}

	if len(i.Description) > 500 {
		errors["description"] = "a descrição da fonte de renda deve ter no máximo 500 caracteres"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
