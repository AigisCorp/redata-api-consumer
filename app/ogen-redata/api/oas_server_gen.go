// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// Charge implements charge operation.
	//
	// Return True or False if the value is under mean for the day.
	//
	// GET /charge
	Charge(ctx context.Context) (ChargeRes, error)
	// GetCheap implements getCheap operation.
	//
	// Get X cheap hours.
	//
	// GET /cheap
	GetCheap(ctx context.Context, params GetCheapParams) (GetCheapRes, error)
	// GetCheapest implements getCheapest operation.
	//
	// Get cheapest hour and its price.
	//
	// GET /cheapest
	GetCheapest(ctx context.Context) (GetCheapestRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
