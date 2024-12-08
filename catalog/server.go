package catalog

import (
	"context"
	"fmt"
	"github.com/asliddinberdiev/microservice_graphql/catalog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	service Service
	proto.UnimplementedCatalogServiceServer
}

func NewGRPCServer(s Service, port uint16) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterCatalogServiceServer(srv, &grpcServer{service: s, UnimplementedCatalogServiceServer: proto.UnimplementedCatalogServiceServer{}})
	reflection.Register(srv)

	return srv.Serve(lis)
}

func (s *grpcServer) CreateProduct(ctx context.Context, r *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	p, err := s.service.CreateProduct(ctx, r.Name, r.Description, r.Price, r.Quantity)
	if err != nil {
		return nil, err
	}

	return &proto.CreateProductResponse{Product: &proto.Product{Id: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity}}, nil
}

func (s *grpcServer) GetProductByID(ctx context.Context, r *proto.GetProductByIDRequest) (*proto.GetProductByIDResponse, error) {
	p, err := s.service.GetProductByID(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &proto.GetProductByIDResponse{Product: &proto.Product{Id: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity}}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *proto.GetProductsRequest) (*proto.GetProductsResponse, error) {
	res := make([]*Product, 0)
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		return nil, err
	}

	products := make([]*proto.Product, len(res))
	for _, p := range res {
		products = append(products, &proto.Product{Id: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Quantity: p.Quantity})
	}

	return &proto.GetProductsResponse{Products: products}, nil
}
