package constants

type Status string

const (
	StatusPendente  Status = "Pendente"
	StatusPago      Status = "Pago"
	StatusCancelado Status = "Cancelado"
)

func (s Status) IsValidStatus() bool {
	switch s {
	case StatusPendente, StatusPago, StatusCancelado:
		return true
	default:
		return false
	}
}
