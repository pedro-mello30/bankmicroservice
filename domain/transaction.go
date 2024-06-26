package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransactionRepository interface {
	SaveTransaction(transaction Transaction, creditCard CreditCard) error
	CreateCreditCard(creditCard CreditCard)  error
	GetCreditCard(creditCard CreditCard) (CreditCard, error)
}

type Transaction struct {
	ID string
	Amount float64
	Status string
	Description string
	Store string
	CreditCardID string
	CreatedAt time.Time
}

func NewTransaction() *Transaction {
	t := &Transaction{}
	t.ID = uuid.NewV4().String()
	t.CreatedAt = time.Now()
	return t
}

func (t *Transaction) ProcessAndValidade(creditCard *CreditCard) {
	if (t.Amount + creditCard.Balance > creditCard.Limit){
		t.Status = "rejected"
	} else {
		t.Status = "approved"
		creditCard.Balance = creditCard.Balance + t.Amount
	}
}