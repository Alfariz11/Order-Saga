package client

import (
	"context"
	"log"
	"time"

	"saga-app/gen/paymentpb" 
	"google.golang.org/grpc"
)

type PaymentClient struct {
	client paymentpb.PaymentServiceClient
}

func NewPaymentClient(cc grpc.ClientConnInterface) *PaymentClient {
	return &PaymentClient{
		client: paymentpb.NewPaymentServiceClient(cc),
	}
}

func (c *PaymentClient) ProcessPayment(orderID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &paymentpb.ProcessPaymentRequest{
		OrderId: orderID,
	}

	resp, err := c.client.ProcessPayment(ctx, req)
	if err != nil {
		return "", err
	}

	log.Printf("[PaymentClient] Payment processed: %s → %s", resp.PaymentId, resp.Status)
	return resp.PaymentId, nil
}

func (c *PaymentClient) RefundPayment(paymentID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &paymentpb.RefundPaymentRequest{
		PaymentId: paymentID,
	}

	resp, err := c.client.RefundPayment(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("[PaymentClient] Payment refund: %s → %s", paymentID, resp.Status)
	return nil
}
