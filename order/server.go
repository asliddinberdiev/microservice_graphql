package order

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/asliddinberdiev/microservice_graphql/account"
	"github.com/asliddinberdiev/microservice_graphql/catalog"
	"github.com/asliddinberdiev/microservice_graphql/order/proto"
)

type grpcServer struct {
	service Service
	proto.UnimplementedOrderServiceServer
	account *account.Client
	catalog *catalog.Client
}

func NewGRPCServer(service Service, accountURL, catalogURL string, port uint16) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	srv := grpc.NewServer()
	proto.RegisterOrderServiceServer(srv, &grpcServer{service: service, account: accountClient, catalog: catalogClient, UnimplementedOrderServiceServer: proto.UnimplementedOrderServiceServer{}})
	reflection.Register(srv)

	return srv.Serve(lis)
}

func (s *grpcServer) CreateOrder(ctx context.Context, r *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	if _, err := s.account.GetAccountByID(ctx, r.AccountId); err != nil {
		return nil, errors.New("account not found")
	}

	productIDs := make([]string, len(r.Products))
	orderedProducts, err := s.catalog.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		return nil, errors.New("products not found")
	}

	products := make([]OrderedProduct, len(orderedProducts))
	for _, p := range orderedProducts {
		product := OrderedProduct{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: 0}

		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	order, err := s.service.Create(ctx, r.AccountId, products)
	if err != nil {
		return nil, errors.New("failed to create order")
	}

	orderProto := proto.Order{Id: order.ID, AccountId: order.AccountID, Products: make([]*proto.Order_OrderedProduct, len(order.Products))}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &proto.Order_OrderedProduct{Id: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity})
	}

	return &proto.CreateOrderResponse{Order: &orderProto}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *proto.GetOrdersForAccountRequest) (*proto.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		return nil, err
	}

	productIDMap := make(map[string]bool, len(accountOrders))

	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}

	productIDs := make([]string, 0, len(productIDMap))
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	products, err := s.catalog.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		return nil, err
	}

	orders := make([]*proto.Order, len(accountOrders))
	for _, o := range accountOrders {
		op := &proto.Order{Id: o.ID, AccountId: o.AccountID, Products: make([]*proto.Order_OrderedProduct, len(o.Products))}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, p := range o.Products {
			for _, pr := range products {
				if pr.ID == p.ID {
					op.Products = append(op.Products, &proto.Order_OrderedProduct{Id: pr.ID, Name: pr.Name, Description: pr.Description, Price: pr.Price, Quantity: pr.Quantity})
					break
				}
			}
			op.Products = append(op.Products, &proto.Order_OrderedProduct{Id: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity})
		}

		orders = append(orders, op)
	}

	return &proto.GetOrdersForAccountResponse{Orders: orders}, nil
}
