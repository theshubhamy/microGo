package catalog

// protoc --go_out=./pb --go-grpc_out=./pb ./catalog.proto

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/theshubhamy/microGo/services/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGrpcServer(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterCatalogServiceServer(server, &grpcServer{UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{}, service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (server *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := server.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (server *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := server.service.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (server *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []Product
	var err error

	if r.Query != "" {
		res, err = server.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = server.service.GetProductsbyIds(ctx, r.Ids)
	} else {
		res, err = server.service.GetProducts(ctx, r.Skip, r.Take)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*pb.Product{}
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
