package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"strings"

	"saga-app/gen/paymentpb" 

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type PaymentServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	mu       sync.Mutex
	payments map[string]string // paymentID -> status
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{
		payments: make(map[string]string),
	}
}

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	paymentID := fmt.Sprintf("PAY-%s", req.OrderId)
	s.payments[paymentID] = "SUCCESS" // bisa juga simulasikan gagal kalau mau

	log.Printf("Payment Processed: %s", paymentID)

	if strings.Contains(req.OrderId, "fail-payment") {
		log.Println("‼️ Simulasi gagal payment")
		return nil, status.Error(codes.Internal, "Simulasi payment error")
	}

	return &paymentpb.ProcessPaymentResponse{
		PaymentId: paymentID,
		Status:    "SUCCESS",
	}, nil
}

func (s *PaymentServer) RefundPayment(ctx context.Context, req *paymentpb.RefundPaymentRequest) (*paymentpb.RefundPaymentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.payments[req.PaymentId]; exists {
		s.payments[req.PaymentId] = "REFUNDED"
		log.Printf("Payment Refunded: %s", req.PaymentId)
		return &paymentpb.RefundPaymentResponse{Status: "REFUNDED"}, nil
	}

	return &paymentpb.RefundPaymentResponse{Status: "NOT_FOUND"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, NewPaymentServer())

	log.Println("Payment Service gRPC server running on port :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
