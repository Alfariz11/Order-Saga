package main

import (
	"log"

	orderclient "saga-app/order-service/client"
	paymentclient "saga-app/payment-service/client"
	shippingclient "saga-app/shipping-service/client"

	"google.golang.org/grpc"
)

func main() {
	// Connect to all services
	orderConn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer orderConn.Close()
	order := orderclient.NewOrderClient(orderConn)

	paymentConn, _ := grpc.Dial("localhost:50052", grpc.WithInsecure())
	defer paymentConn.Close()
	payment := paymentclient.NewPaymentClient(paymentConn)

	shippingConn, _ := grpc.Dial("localhost:50053", grpc.WithInsecure())
	defer shippingConn.Close()
	shipping := shippingclient.NewShippingClient(shippingConn)

	// ================================
	// 💥 TEST: Shipping Failure → Semua Kompensasi
	// ================================
	log.Println("=== TEST SHIPPING GAGAL ===")
	testShippingFail(order, payment, shipping)

	// ================================
	// 💥 TEST: Payment Failure → Cancel Order
	// ================================
	log.Println("\n=== TEST PAYMENT GAGAL ===")
	testPaymentFail(order, payment)

	// ================================
	// ✅ TEST: Semua Berhasil
	// ================================
	log.Println("\n=== TEST SUKSES ===")
	testSuccess(order, payment, shipping)
}

// ========== TEST CASE IMPLEMENTATION ==========

func testShippingFail(order *orderclient.OrderClient, payment *paymentclient.PaymentClient, shipping *shippingclient.ShippingClient) {
	orderID, err := order.CreateOrder("solon", "fail-shipping")
	if err != nil {
		log.Println("Gagal buat order:", err)
		return
	}
	paymentID, err := payment.ProcessPayment(orderID)
	if err != nil {
		log.Println("Gagal bayar:", err)
		order.CancelOrder(orderID)
		return
	}
	shippingID, err := shipping.StartShipping(orderID)
	if err != nil {
		log.Println("GAGAL shipping:", err)
		log.Println("🔁 Kompensasi: Cancel Shipping → Refund → Cancel Order")
		shipping.CancelShipping(shippingID)
		payment.RefundPayment(paymentID)
		order.CancelOrder(orderID)
		return
	}
	log.Println("SUKSES shipping:", shippingID)
}

func testPaymentFail(order *orderclient.OrderClient, payment *paymentclient.PaymentClient) {
	orderID, err := order.CreateOrder("solon", "fail-payment")
	if err != nil {
		log.Println("Gagal buat order:", err)
		return
	}
	_, err = payment.ProcessPayment(orderID)
	if err != nil {
		log.Println("GAGAL bayar:", err)
		log.Println("🔁 Kompensasi: Cancel Order")
		order.CancelOrder(orderID)
		return
	}
	log.Println("SUKSES bayar")
}

func testSuccess(order *orderclient.OrderClient, payment *paymentclient.PaymentClient, shipping *shippingclient.ShippingClient) {
	orderID, err := order.CreateOrder("solon", "sepatu")
	if err != nil {
		log.Println("Gagal buat order:", err)
		return
	}
	paymentID, err := payment.ProcessPayment(orderID)
	if err != nil {
		log.Println("Gagal bayar:", err)
		order.CancelOrder(orderID)
		return
	}
	shippingID, err := shipping.StartShipping(orderID)
	if err != nil {
		log.Println("Gagal shipping:", err)
		payment.RefundPayment(paymentID)
		order.CancelOrder(orderID)
		return
	}
	log.Println("ORDER SUKSES sampai dikirim:", shippingID)
}
