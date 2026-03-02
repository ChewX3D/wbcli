package ports

// ErrorCode classifies APIError for programmatic handling such as retry decisions or UI routing.
type ErrorCode string

const (
	// CodeUnauthorized means the credentials were rejected by the exchange.
	CodeUnauthorized ErrorCode = "unauthorized"
	// CodeForbidden means credentials are valid but the API token lacks the required endpoint permission.
	CodeForbidden ErrorCode = "forbidden"
	// CodeInvalidRequest means the exchange rejected the payload as malformed or invalid.
	CodeInvalidRequest ErrorCode = "invalid_request"
	// CodeBusinessRule means the exchange rejected the request due to a trading constraint.
	CodeBusinessRule ErrorCode = "business_rule"
	// CodeUnavailable means the exchange or transport is temporarily unreachable.
	CodeUnavailable ErrorCode = "unavailable"
)

// APIError is the unified error type for all exchange-facing operation failures.
// Message is the short user-facing explanation, always displayed first.
// Details carries full context: endpoint, exchange response, and actionable hints.
type APIError struct {
	Code    ErrorCode
	Message string
	Details string
}

// Error returns Message alone when Details is empty; otherwise returns Message followed by
// Details on a new line so callers can display both independently.
func (e *APIError) Error() string {
	if e.Details != "" {
		return e.Message + "\n" + e.Details
	}

	return e.Message
}
