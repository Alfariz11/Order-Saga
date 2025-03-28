package client

import (
	"context"
	"log"
	"time"

	"saga-app/gen/shippingpb"
	"google.golang.org/grpc"
)

type ShippingClient struct {
	client shippingpb.ShippingServiceClient
}

func NewShippingClient(cc grpc.ClientConnInterface) *ShippingClient {
	return &ShippingClient{
		client: shippingpb.NewShippingServiceClient(cc),
	}
}

func (c *ShippingClient) StartShipping(orderID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &shippingpb.StartShippingRequest{
		OrderId: orderID,
	}

	resp, err := c.client.StartShipping(ctx, req)
	if err != nil {
		return "", err
	}

	log.Printf("[ShippingClient] Shipping started: %s → %s", resp.ShippingId, resp.Status)
	return resp.ShippingId, nil
}

func (c *ShippingClient) CancelShipping(shippingID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &shippingpb.CancelShippingRequest{
		ShippingId: shippingID,
	}

	resp, err := c.client.CancelShipping(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("[ShippingClient] Shipping cancelled: %s → %s", shippingID, resp.Status)
	return nil
}
