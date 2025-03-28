package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"strings"

	"saga-app/gen/shippingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type ShippingServer struct {
	shippingpb.UnimplementedShippingServiceServer
	mu         sync.Mutex
	shipments  map[string]string // shippingID -> status
}

func NewShippingServer() *ShippingServer {
	return &ShippingServer{
		shipments: make(map[string]string),
	}
}

func (s *ShippingServer) StartShipping(ctx context.Context, req *shippingpb.StartShippingRequest) (*shippingpb.StartShippingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	shippingID := fmt.Sprintf("SHIP-%s", req.OrderId)
	s.shipments[shippingID] = "SHIPPED" // atau "PENDING" tergantung logikamu

	log.Printf("Shipping Started: %s", shippingID)

	if strings.Contains(req.OrderId, "fail-shipping") {
		log.Println("‼️ Simulasi gagal shipping")
		return nil, status.Error(codes.Internal, "Simulasi shipping error")
	}

	return &shippingpb.StartShippingResponse{
		ShippingId: shippingID,
		Status:     "SHIPPED",
	}, nil
}

func (s *ShippingServer) CancelShipping(ctx context.Context, req *shippingpb.CancelShippingRequest) (*shippingpb.CancelShippingResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.shipments[req.ShippingId]; exists {
		s.shipments[req.ShippingId] = "CANCELLED"
		log.Printf("Shipping Cancelled: %s", req.ShippingId)
		return &shippingpb.CancelShippingResponse{Status: "CANCELLED"}, nil
	}

	return &shippingpb.CancelShippingResponse{Status: "NOT_FOUND"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	shippingpb.RegisterShippingServiceServer(grpcServer, NewShippingServer())

	log.Println("Shipping Service gRPC server running on port :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
