package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"saga-app/gen/orderpb" 

	"google.golang.org/grpc"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	mu     sync.Mutex
	orders map[string]string // orderID -> status
}

func NewOrderServer() *OrderServer {
	return &OrderServer{
		orders: make(map[string]string),
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	orderID := fmt.Sprintf("ORD-%s-%s", req.UserId, req.ItemId)
	s.orders[orderID] = "PENDING"

	log.Printf("Order Created: %s\n", orderID)

	return &orderpb.CreateOrderResponse{
		OrderId: orderID,
		Status:  "PENDING",
	}, nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[req.OrderId]; exists {
		s.orders[req.OrderId] = "CANCELLED"
		log.Printf("Order Cancelled: %s\n", req.OrderId)
		return &orderpb.CancelOrderResponse{Status: "CANCELLED"}, nil
	}

	return &orderpb.CancelOrderResponse{Status: "NOT_FOUND"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, NewOrderServer())

	log.Println("Order Service gRPC server running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
