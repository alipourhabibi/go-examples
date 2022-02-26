package server

import (
	"context"
	"io"

	protos "github.com/alipourhabibi/grpc/server/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct {
	protos.UnimplementedCurrencyServer
	log  hclog.Logger
	subs map[protos.Currency_SubscribeServer][]*protos.RateRequest
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	// Error handling
	if rr.GetBase().String() == rr.GetDestination().String() {
		theStatus := status.Newf(
			codes.InvalidArgument,
			"something",
		)
		// if serializing rr has error it will return that
		x, err := theStatus.WithDetails(rr)
		if err != nil {
			return nil, err
		}
		return nil, x.Err()
	}

	return &protos.RateResponse{Rate: "0.5"}, nil
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}

func (c *Currency) Subscribe(src protos.Currency_SubscribeServer) error {
	for {
		rr, err := src.Recv()

		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}
		if err != nil {
			c.log.Error("Unable to read data form client", "error", err)
			return err
		}

		rrs, ok := c.subs[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		for _, r := range rrs {
			if r.Base == rr.Base && r.Destination == rr.Destination {
				grpcError := status.New(codes.AlreadyExists, "already exists")
				// metadata
				grpcError, err := grpcError.WithDetails(rr)
				if err != nil {
					c.log.Info("Internal Error in metadata")
					continue
				}

				rrs := &protos.StreamingRateRequest_Error{Error: grpcError.Proto()}
				src.Send(&protos.StreamingRateRequest{Message: rrs})

			}
		}

		rrs = append(rrs, rr)
		c.subs[src] = rrs
	}

	return nil
}
