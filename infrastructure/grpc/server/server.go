package server

import (
	"log"
	"net"

	"github.com/pedro-mello30/bankmicroservice/infrastructure/grpc/service"
	"github.com/pedro-mello30/bankmicroservice/usecase"
)

type GRPCServer struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
}

func NewGRPCServer() GRPCServer {
	return GRPCServer{}
}

func (g GRPCServer) Serve() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("could not listen tcp port")
	}

	transactionService := service.NewTransactionService()
	transactionService.ProcessTransactionUseCase = g.ProcessTransactionUseCase
	GRPCServer := grpc.NewServer()
	
}