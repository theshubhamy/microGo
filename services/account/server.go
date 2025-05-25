package account

//go:generate protoc ./account.proto --go_out=plugins=grpc:./pb

import (
	"context"
	"fmt"
	"net"

	"github.com/theshubhamy/microGo/services/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGrpcServer(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &grpcServer{UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{}, service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (server *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	account, err := server.service.PostAccount(ctx, r.Name, r.Email, r.Phone, r.Password)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Id:    account.ID,
		Name:  account.Name,
		Email: account.Email,
		Phone: account.Phone,
	}, nil
}

func (server *grpcServer) LoginAccount(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	account, accessToken, refreshToken, err := server.service.LoginAccount(ctx, r.Emailorphone, r.Password, r.Ip, r.UserAgent)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{
		Id:           account.ID,
		Name:         account.Name,
		Email:        account.Email,
		Phone:        account.Phone,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (server *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := server.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	accounts := []*pb.Account{}
	for _, account := range res {
		accounts = append(accounts, &pb.Account{
			Id:    account.ID,
			Name:  account.Name,
			Email: account.Email,
			Phone: account.Phone,
		})
	}
	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
