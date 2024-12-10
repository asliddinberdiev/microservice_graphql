package order

import (
	"context"
	"time"

	"github.com/asliddinberdiev/microservice_graphql/order/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service proto.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := proto.NewOrderServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	protoProducts := make([]*proto.CreateOrderRequest_OrderedProduct, len(products))

	for _, p := range products {
		protoProducts = append(protoProducts, &proto.CreateOrderRequest_OrderedProduct{ProductId: p.ID, Quantity: p.Quantity})
	}

	r, err := c.service.CreateOrder(ctx, &proto.CreateOrderRequest{AccountId: accountID, Products: protoProducts})
	if err != nil {
		return nil, err
	}

	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)

	return &Order{ID: newOrder.Id, AccountID: newOrder.AccountId, TotalPrice: newOrder.TotalPrice, CreatedAt: newOrderCreatedAt, Products: products}, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	r, err := c.service.GetOrdersForAccount(ctx, &proto.GetOrdersForAccountRequest{AccountId: accountID})
	if err != nil {
		return nil, err
	}

	orders := make([]*Order, len(r.Orders))
	for _, orderProto := range r.Orders {
		newOrder := &Order{ID: orderProto.Id, AccountID: orderProto.AccountId, TotalPrice: orderProto.TotalPrice}

		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)

		products := make([]OrderedProduct, len(orderProto.Products))
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{ID: p.Id, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity})
		}

		newOrder.Products = products
		orders = append(orders, newOrder)
	}

	return orders, nil
}
