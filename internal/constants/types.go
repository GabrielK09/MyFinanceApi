package constants

type TransactionType string

const (
	TypeIncome  TransactionType = "Entrada"
	TypeExpense TransactionType = "Saída"
)

func (t TransactionType) IsValidType() bool {
	switch t {
	case TypeIncome, TypeExpense:
		return true
	default:
		return false
	}
}
