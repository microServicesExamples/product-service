package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/microServicesExamples/gRPC/product/productpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	productpb.UnimplementedProductServiceServer
}

func (*server) GetProductDetails(ctx context.Context, req *productpb.GetProductDetailsRequest) (resp *productpb.GetProductDetailsResponse, err error) {
	fmt.Println("Getting the product details via GRPC")

	// get request params
	productId := req.GetId()

	// execute the business logic
	result, ok := products[productId]

	if !ok {
		fmt.Println("product with id:", productId, "does not exist")
		return &productpb.GetProductDetailsResponse{}, fmt.Errorf("product with id: %v does not exist", productId)
	}
	fmt.Printf("found product with details: %+v\n", result)

	// prepare the response
	resp = &productpb.GetProductDetailsResponse{
		Id:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Category:    string(result.Category),
		Price:       result.Price,
		Quantity:    result.Quantity,
	}

	// return the response
	return resp, nil
}

func (*server) ListProductDetails(ctx context.Context, req *productpb.ListProductDetailsRequest) (resp *productpb.ListProductDetailsResponse, err error) {
	fmt.Println("Getting the list of product details via GRPC")

	// get request params
	productIds := req.GetIds()

	// execute the business logic
	var details []*productpb.GetProductDetailsResponse
	for _, productId := range productIds {
		result, ok := products[productId.Id]

		if !ok {
			fmt.Println("product with id:", productId, "does not exist")
			return &productpb.ListProductDetailsResponse{}, fmt.Errorf("product with id: %v does not exist", productId)
		}

		details = append(details, &productpb.GetProductDetailsResponse{
			Id:          result.ID,
			Name:        result.Name,
			Description: result.Description,
			Category:    string(result.Category),
			Price:       result.Price,
			Quantity:    result.Quantity,
		})
	}
	fmt.Println("found products with details:", details)

	// prepare the response
	resp = &productpb.ListProductDetailsResponse{
		Details: details,
	}

	// return the response
	return resp, nil
}

func (*server) UpdateProductQuantity(ctx context.Context, req *productpb.UpdateProductQuantityRequest) (resp *productpb.UpdateProductQuantityResponse, err error) {
	fmt.Println("Getting the product details via GRPC")

	// get request params
	productId := req.GetId()
	quantity := req.GetQuantity()

	// execute the business logic
	result, ok := products[productId]

	if !ok {
		fmt.Println("product with id:", productId, "does not exist")
		return &productpb.UpdateProductQuantityResponse{}, fmt.Errorf("product with id: %v does not exist", productId)
	}
	fmt.Println("found product with details:", result)

	// update the product quantity
	result.Quantity = quantity
	products[productId] = result

	// prepare the response
	resp = &productpb.UpdateProductQuantityResponse{}

	// return the response
	return resp, nil
}

func startGRPCServer() {
	fmt.Println("Starting the gRPC server")

	// Create a connection
	listen, err := net.Listen("tcp", "0.0.0.0:5051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a server
	s := grpc.NewServer()

	// Register the service on server
	productpb.RegisterProductServiceServer(s, &server{})

	// Register reflection service on gRPC service
	reflection.Register(s)

	// Start listening on the created connection
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
