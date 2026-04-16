package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"pb"
)

type AccountClient struct {
	conn   *grpc.ClientConn
	client pb.AccountGrpcClient
}

func NewAccountClient(addr string) (*AccountClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &AccountClient{
		conn:   conn,
		client: pb.NewAccountGrpcClient(conn),
	}, nil
}

func (c *AccountClient) GetBalance(ctx context.Context, accountNo string) (*pb.GetBalanceResponse, error) {
	return c.client.GetBalance(ctx, &pb.GetBalanceRequest{AccountNo: accountNo})
}

func (c *AccountClient) UpdateBalance(ctx context.Context, accountNo string, newBalance float64) (*pb.UpdateBalanceResponse, error) {
	return c.client.UpdateBalance(ctx, &pb.UpdateBalanceRequest{
		AccountNo:  accountNo,
		NewBalance: newBalance,
	})
}

func (c *AccountClient) Close() error {
	return c.conn.Close()
}
