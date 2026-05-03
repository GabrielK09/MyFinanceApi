package incomereceiptssconstants

type Status string

const (
	StatusAtivo     Status = "Ativo"
	StatusCancelado Status = "Cancelado"
)

func (s Status) IsValidStatus() bool {
	switch s {
	case StatusAtivo, StatusCancelado:
		return true
	default:
		return false
	}
}
