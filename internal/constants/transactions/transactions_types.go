package transactionsconstants

type ConstantType string

const (
	TypeIncome  ConstantType = "Entrada"
	TypeExpense ConstantType = "Saída"
)

func (t ConstantType) IsValidType() bool {
	switch t {
	case TypeIncome, TypeExpense:
		return true
	default:
		return false
	}
}
