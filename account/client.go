package account

import (
	"context"
	"log"

	"github.com/asliddinberdiev/microservice_graphql/account/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service proto.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := proto.NewAccountServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.CreateAccount(ctx, &proto.CreateAccountRequest{Name: name})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(r)
	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (c *Client) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccountByID(ctx, &proto.GetAccountByIDRequest{Id: id})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(r)
	return &Account{ID: r.Account.Id, Name: r.Account.Name}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(ctx, &proto.GetAccountsRequest{Skip: skip, Take: take})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	accounts := make([]Account, len(r.Accounts))
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{ID: a.Id, Name: a.Name})
	}

	log.Println(accounts)
	return accounts, nil
}
