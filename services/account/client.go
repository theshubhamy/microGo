package account

import (
	"context"

	"github.com/theshubhamy/microGo/services/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name, email, phone, password string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name, Email: email, Phone: phone, Password: password})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:    r.Id,
		Name:  r.Name,
		Email: r.Email,
		Phone: r.Phone,
	}, nil
}

func (c *Client) LoginAccount(ctx context.Context, emailorphone, password, ip, userAgent string) (*Account, string, string, error) {
	res, err := c.service.LoginAccount(ctx, &pb.LoginRequest{Emailorphone: emailorphone, Password: password, Ip: ip, UserAgent: userAgent})
	if err != nil {
		return nil, "", "", err
	}
	account := &Account{
		ID:    res.Id,
		Name:  res.Name,
		Email: res.Email,
		Phone: res.Phone,
	}

	return account, res.AccessToken, res.RefreshToken, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:    r.Id,
		Name:  r.Name,
		Email: r.Email,
		Phone: r.Phone,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) (*[]Account, error) {
	res, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, acc := range res.Accounts {
		accounts = append(accounts, Account{
			ID:   acc.Id,
			Name: acc.Name,
		})
	}
	return &accounts, nil
}
