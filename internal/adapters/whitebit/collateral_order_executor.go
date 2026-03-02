package whitebit

import (
	"context"
	"encoding/json"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
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

// PlaceCollateralLimitOrder maps app request to WhiteBIT payload and executes signed request.
func (adapter *CollateralOrderExecutorAdapter) PlaceCollateralLimitOrder(
	ctx context.Context,
	credential domainauth.Credential,
	request ports.CollateralLimitOrderRequest,
) (json.RawMessage, error) {
	postOnly := request.PostOnly

	return adapter.client.PlaceCollateralLimitOrder(ctx, credential, CollateralLimitOrderRequest{
		Market:        request.Market,
		Side:          OrderSide(request.Side),
		Amount:        request.Amount,
		Price:         request.Price,
		ClientOrderID: request.ClientOrderID,
		PostOnly:      &postOnly,
	})
}
