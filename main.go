package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pedro-mello30/bankmicroservice/domain"
	"github.com/pedro-mello30/bankmicroservice/infrastructure/repository"
	"github.com/pedro-mello30/bankmicroservice/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := setupDb()
	defer db.Close()

	cc := domain.NewCreditCard()
	cc.Number = "1234"
	cc.Name = "Pedro Mello"
	cc.ExpirationMonth = 6
	cc.ExpirationYear = 2022
	cc.CVV = 123
	cc.Limit = 10000
	cc.Balance = 0

	repo := repository.NewTransactionRepositoryDb(db)
	err := repo.CreateCreditCard(*cc)
	if err != nil {
		fmt.Println(err)
	}

}

func setupTransactionUseCase(db *sql.DB) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db",
		"5432",
		"postgres",
		"root",
		"codebank",
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%s", err)
		log.Fatal("error connect to database")
	}
	return db
}
