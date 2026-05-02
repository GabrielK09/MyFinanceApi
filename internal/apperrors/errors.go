package apperrors

import "errors"

type AppError struct {
	message string
	err     error
}

var (
	ErrNotFound            = errors.New("registro não encontrado")
	ErrUnprocessableEntity = errors.New("campos obrigatórios ausentes ou inválidos")
	// DB errors
	ErrRelationViolation   = errors.New("referência incorreta")
	ErrForeignKeyViolation = errors.New("referência incorreta")
	ErrUniqueConstraint    = errors.New("valores duplicados foram informados")
	ErrCheckViolation      = errors.New("valores inseridos fora do padrão")
)

func (e AppError) Error() string {
	return e.message
}

func (e AppError) Unwrap() error {
	return e.err
}

func NewUnprocessableEntity(message string) error {
	return AppError{
		message: message,
		err:     ErrUnprocessableEntity,
	}
}
