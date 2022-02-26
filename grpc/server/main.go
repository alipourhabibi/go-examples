package main

import (
	"net"
	"os"

	protos "github.com/alipourhabibi/grpc/server/protos/currency"
	"github.com/alipourhabibi/grpc/server/server"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, cs)

	// Enable reflection for server
	// Disable in production
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen to", err)
		os.Exit(1)
	}
	gs.Serve(l)
}
