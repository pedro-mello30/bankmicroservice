package usecase

import (
	"encoding/json"
	"os"
	"time"

	"github.com/pedro-mello30/bankmicroservice/domain"
	"github.com/pedro-mello30/bankmicroservice/dto"
	"github.com/pedro-mello30/bankmicroservice/infrastructure/kafka"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
	KafkaProducer         kafka.KafkaProducer
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	creditCard := u.hydrateCreditCard(transactionDto)
	ccDB, err := u.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}

	creditCard.ID = ccDB.ID
	creditCard.Limit = ccDB.Limit
	creditCard.Balance = ccDB.Balance

	transaction := u.newTransaction(transactionDto, *creditCard)
	transaction.ProcessAndValidade(creditCard)

	err = u.TransactionRepository.SaveTransaction(*transaction, *creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}

	transactionDto.ID = transaction.ID
	transactionDto.CreatedAt = transaction.CreatedAt

	transactionJson, err := json.Marshal(transactionDto)
	if err != nil {
		return domain.Transaction{}, err
	}

	err = u.KafkaProducer.Publish(string(transactionJson), os.Getenv("KafkaTransactionsTopic"))
	if err != nil {
		return domain.Transaction{}, err
	}

	return *transaction, nil

}

func (u UseCaseTransaction) hydrateCreditCard(transactionDto dto.Transaction) *domain.CreditCard {
	cc := domain.NewCreditCard()
	cc.Name = transactionDto.Name
	cc.Number = transactionDto.Number
	cc.ExpirationMonth = transactionDto.ExpirationMonth
	cc.ExpirationYear = transactionDto.ExpirationYear
	cc.CVV = transactionDto.CVV
	return cc
}

func (u UseCaseTransaction) newTransaction(transactionDto dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	transaction := domain.NewTransaction()
	transaction.CreditCardID = cc.ID
	transaction.Amount = transactionDto.Amount
	transaction.Description = transactionDto.Description
	transaction.Store = transactionDto.Store
	transaction.CreditCardID = cc.ID
	transaction.CreatedAt = time.Now()
	return transaction
}
