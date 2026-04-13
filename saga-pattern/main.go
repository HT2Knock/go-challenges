package main

import (
	"fmt"

	"github.com/google/uuid"
)

type PaymentStatus int

const (
	PaymentPending PaymentStatus = iota
	PaymentCompleted
	PaymentFailed
)

type OrderCreatedEvent struct {
	OrderID uuid.UUID
	Amount  float64
}

type PaymentProcessedEvent struct {
	PaymentID uuid.UUID
	OrderID   uuid.UUID
	Status    PaymentStatus
}

func createOrder(orderID uuid.UUID, amount float64) OrderCreatedEvent {
	return OrderCreatedEvent{OrderID: orderID, Amount: amount}
}

func processedPayment(orderID uuid.UUID) PaymentProcessedEvent {
	return PaymentProcessedEvent{PaymentID: uuid.New(), OrderID: orderID, Status: PaymentFailed}
}

func compensatedOrder(sagaID uuid.UUID, orderID uuid.UUID, version int) {
	fmt.Printf("Compensated sagaID: %v for order entityID: %v version: %d", sagaID, orderID, version)
}

func handleOrderSaga(sagaID uuid.UUID, version int) {
	orderCreated := createOrder(uuid.New(), 1500)
	processedPayment := processedPayment(orderCreated.OrderID)
	if processedPayment.Status == PaymentFailed {
		compensatedOrder(sagaID, orderCreated.OrderID, version)
	}
}

func main() {
	sagaID := uuid.New()
	version := 1

	handleOrderSaga(sagaID, version)
}
