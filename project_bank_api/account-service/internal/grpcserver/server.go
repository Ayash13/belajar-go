package grpcserver

import (
	"account-service/internal/service"
	"context"

	"pb"

	"google.golang.org/grpc"
)

type accountGrpcServer struct {
	pb.UnimplementedAccountGrpcServer
	svc service.AccountService
}

func RegisterAccountServer(s *grpc.Server, svc service.AccountService) {
	pb.RegisterAccountGrpcServer(s, &accountGrpcServer{svc: svc})
}

func (s *accountGrpcServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	account, err := s.svc.GetBalance(ctx, req.AccountNo)
	if err != nil {
		return &pb.GetBalanceResponse{Found: false}, nil
	}

	return &pb.GetBalanceResponse{
		AccountNo: account.AccountNo,
		Name:      account.Name,
		Balance:   account.Balance,
		Currency:  account.Currency,
		Status:    account.Status,
		Found:     true,
	}, nil
}

func (s *accountGrpcServer) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*pb.UpdateBalanceResponse, error) {
	err := s.svc.UpdateBalance(ctx, req.AccountNo, req.NewBalance)
	if err != nil {
		return &pb.UpdateBalanceResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpdateBalanceResponse{
		Success: true,
		Message: "balance updated",
	}, nil
}
