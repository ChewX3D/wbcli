package whitebit

import (
	"context"
	"encoding/json"

	"github.com/ChewX3D/crypto/internal/app/ports"
	domainauth "github.com/ChewX3D/crypto/internal/domain/auth"
)

// CollateralOrderExecutorAdapter adapts app order port to WhiteBIT transport client.
type CollateralOrderExecutorAdapter struct {
	client *Client
}

// NewCollateralOrderExecutorAdapter constructs order executor adapter.
func NewCollateralOrderExecutorAdapter(client *Client) *CollateralOrderExecutorAdapter {
	return &CollateralOrderExecutorAdapter{client: client}
}

// NewDefaultCollateralOrderExecutorAdapter constructs order executor adapter with default client.
func NewDefaultCollateralOrderExecutorAdapter() *CollateralOrderExecutorAdapter {
	return NewCollateralOrderExecutorAdapter(NewDefaultClient())
}

// GetCollateralAccountHedgeMode reads current account hedge mode from WhiteBIT.
func (adapter *CollateralOrderExecutorAdapter) GetCollateralAccountHedgeMode(
	ctx context.Context,
	credential domainauth.Credential,
) (bool, error) {
	response, err := adapter.client.GetCollateralAccountHedgeMode(ctx, credential)
	if err != nil {
		return false, buildAPIError(err, collateralAccountHedgeModePath, "hedge mode query")
	}

	return response.HedgeMode, nil
}

// PlaceCollateralLimitOrder maps app request to WhiteBIT payload and executes signed request.
func (adapter *CollateralOrderExecutorAdapter) PlaceCollateralLimitOrder(
	ctx context.Context,
	credential domainauth.Credential,
	request ports.CollateralLimitOrderRequest,
) (json.RawMessage, error) {
	postOnly := request.PostOnly

	result, err := adapter.client.PlaceCollateralLimitOrder(ctx, credential, CollateralLimitOrderRequest{
		Market:        request.Market,
		Side:          OrderSide(request.Side),
		PositionSide:  PositionSide(request.PositionSide),
		Amount:        request.Amount,
		Price:         request.Price,
		ClientOrderID: request.ClientOrderID,
		PostOnly:      &postOnly,
	})
	if err != nil {
		return nil, buildAPIError(err, collateralLimitOrderPath, "order placement")
	}

	return result, nil
}
