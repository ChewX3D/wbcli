package collateral

import (
	"context"
	"fmt"
	"strings"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

const (
	singleOrderMode       = "single"
	collateralOrderPrefix = "order"
)

// PlaceOrderRequest is input for collateral single order placement use-case.
type PlaceOrderRequest struct {
	Market        string
	Side          string
	Amount        string
	Price         string
	ClientOrderID string
}

// PlaceOrderResult is normalized output for collateral single order placement use-case.
type PlaceOrderResult struct {
	RequestID       string   `json:"request_id"`
	Mode            string   `json:"mode"`
	OrdersPlanned   int      `json:"orders_planned"`
	OrdersSubmitted int      `json:"orders_submitted"`
	OrdersFailed    int      `json:"orders_failed"`
	Errors          []string `json:"errors"`
}

// PlaceOrderService orchestrates collateral single order placement.
type PlaceOrderService struct {
	credentialStore ports.CredentialStore
	orderExecutor   ports.CollateralOrderExecutor
	clock           ports.Clock
}

// NewPlaceOrderService constructs PlaceOrderService.
func NewPlaceOrderService(
	credentialStore ports.CredentialStore,
	orderExecutor ports.CollateralOrderExecutor,
	clock ports.Clock,
) *PlaceOrderService {
	return &PlaceOrderService{
		credentialStore: credentialStore,
		orderExecutor:   orderExecutor,
		clock:           clock,
	}
}

// Execute places one collateral post-only limit order.
func (service *PlaceOrderService) Execute(ctx context.Context, request PlaceOrderRequest) (PlaceOrderResult, error) {
	credential, err := service.credentialStore.Load(ctx)
	if err != nil {
		return PlaceOrderResult{}, fmt.Errorf("load credential: %w", err)
	}

	_, err = service.orderExecutor.PlaceCollateralLimitOrder(ctx, credential, ports.CollateralLimitOrderRequest{
		Market:        strings.TrimSpace(request.Market),
		Side:          strings.TrimSpace(request.Side),
		Amount:        strings.TrimSpace(request.Amount),
		Price:         strings.TrimSpace(request.Price),
		ClientOrderID: request.ClientOrderID,
		PostOnly:      true,
	})
	if err != nil {
		return PlaceOrderResult{}, fmt.Errorf("place collateral limit order: %w", err)
	}

	return PlaceOrderResult{
		RequestID:       fmt.Sprintf("%s-%d", collateralOrderPrefix, service.clock.Now().UTC().UnixNano()),
		Mode:            singleOrderMode,
		OrdersPlanned:   1,
		OrdersSubmitted: 1,
		OrdersFailed:    0,
		Errors:          []string{},
	}, nil
}
