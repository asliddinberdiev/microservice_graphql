package account

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/asliddinberdiev/microservice_graphql/account/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	proto.UnimplementedAccountServiceServer
}

func NewGRPCServer(s Service, port uint16) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterAccountServiceServer(srv, &grpcServer{service: s, UnimplementedAccountServiceServer: proto.UnimplementedAccountServiceServer{}})
	reflection.Register(srv)

	return srv.Serve(lis)
}

func (s *grpcServer) CreateAccount(ctx context.Context, r *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	account, err := s.service.CreateAccount(ctx, r.Name)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	log.Println(account)
	return &proto.CreateAccountResponse{Account: &proto.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *grpcServer) GetAccountByID(ctx context.Context, r *proto.GetAccountByIDRequest) (*proto.GetAccountByIDResponse, error) {
	account, err := s.service.GetAccountByID(ctx, r.Id)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get account: %v", err)
	}

	if account == nil {
		return &proto.GetAccountByIDResponse{}, nil
	}

	log.Println(account)
	return &proto.GetAccountByIDResponse{Account: &proto.Account{Id: account.ID, Name: account.Name}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *proto.GetAccountsRequest) (*proto.GetAccountsResponse, error) {
	resp, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to get accounts: %v", err)
	}

	accounts := make([]*proto.Account, len(resp))
	for i, account := range resp {
		accounts[i] = &proto.Account{Id: account.ID, Name: account.Name}
	}

	log.Println(accounts)
	return &proto.GetAccountsResponse{Accounts: accounts}, nil
}
