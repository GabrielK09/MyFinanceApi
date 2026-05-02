package transactionsconstants

type TransactionsOrigin string

const (
	ManualOrigin        TransactionsOrigin = "manual"
	IncomeReceiptOrigin TransactionsOrigin = "income_receipt"
)

func (t TransactionsOrigin) IsValidTransactionsOrigin() bool {
	switch t {
	case ManualOrigin, IncomeReceiptOrigin:
		return true
	default:
		return false
	}
}
