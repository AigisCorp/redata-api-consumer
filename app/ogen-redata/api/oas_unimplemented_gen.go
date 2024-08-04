// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// Charge implements charge operation.
//
// Return True or False if the value is under mean for the day.
//
// GET /charge
func (UnimplementedHandler) Charge(ctx context.Context) (r ChargeRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCheap implements getCheap operation.
//
// Get X cheap hours.
//
// GET /cheap
func (UnimplementedHandler) GetCheap(ctx context.Context, params GetCheapParams) (r GetCheapRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCheapest implements getCheapest operation.
//
// Get cheapest hour and its price.
//
// GET /cheapest
func (UnimplementedHandler) GetCheapest(ctx context.Context) (r GetCheapestRes, _ error) {
	return r, ht.ErrNotImplemented
}
