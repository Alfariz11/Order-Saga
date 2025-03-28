package client

import (
	"context"
	"log"
	"time"

	"saga-app/gen/orderpb"
	"google.golang.org/grpc"
)

type OrderClient struct {
	client orderpb.OrderServiceClient
}

func NewOrderClient(cc grpc.ClientConnInterface) *OrderClient {
	return &OrderClient{
		client: orderpb.NewOrderServiceClient(cc),
	}
}

func (c *OrderClient) CreateOrder(userID, itemID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &orderpb.CreateOrderRequest{
		UserId: userID,
		ItemId: itemID,
	}

	resp, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		return "", err
	}

	log.Printf("[OrderClient] Created Order ID: %s, Status: %s", resp.OrderId, resp.Status)
	return resp.OrderId, nil
}

func (c *OrderClient) CancelOrder(orderID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &orderpb.CancelOrderRequest{
		OrderId: orderID,
	}

	resp, err := c.client.CancelOrder(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("[OrderClient] Cancel Order: %s → %s", orderID, resp.Status)
	return nil
}
