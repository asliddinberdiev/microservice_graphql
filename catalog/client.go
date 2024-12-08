package catalog

import (
	"context"
	"github.com/asliddinberdiev/microservice_graphql/catalog/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service proto.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := proto.NewCatalogServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateProduct(ctx context.Context, name, description string, price, quantity uint64) (*Product, error) {
	r, err := c.service.CreateProduct(ctx, &proto.CreateProductRequest{Name: name, Description: description, Price: price, Quantity: quantity})
	if err != nil {
		return nil, err
	}

	return &Product{ID: r.Product.Id, Name: r.Product.Name, Description: r.Product.Description, Price: r.Product.Price, Quantity: r.Product.Quantity}, nil
}

func (c *Client) GetProductByID(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProductByID(ctx, &proto.GetProductByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &Product{ID: r.Product.Id, Name: r.Product.Name, Description: r.Product.Description, Price: r.Product.Price, Quantity: r.Product.Quantity}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, ids []string, query string) ([]*Product, error) {
	r, err := c.service.GetProducts(ctx, &proto.GetProductsRequest{Skip: skip, Take: take, Ids: ids, Query: query})
	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0, len(r.GetProducts()))
	for _, p := range r.GetProducts() {
		products = append(products, &Product{ID: p.Id, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity})
	}

	return products, nil
}
