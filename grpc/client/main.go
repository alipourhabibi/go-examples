package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	proto "github.com/alipourhabibi/grpc/client/proto/currency"
)

func main() {
	r := gin.Default()
	r.POST("datas", datas)
	r.POST("stream", stream)

	r.Run(":9090")

}

func datas(c *gin.Context) {
	con, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer con.Close()

	base := c.Query("base")
	destination := c.Query("destination")

	rr := &proto.RateRequest{
		Base:        proto.Currencies(proto.Currencies_value[base]),
		Destination: proto.Currencies(proto.Currencies_value[destination]),
	}

	cc := proto.NewCurrencyClient(con)
	resp, err := cc.GetRate(context.Background(), rr)
	if err != nil {
		if theStatus, ok := status.FromError(err); ok {
			md := theStatus.Details()[0].(*proto.RateRequest)
			fmt.Println(md.Base.String(), md.Destination.String())
			fmt.Println("[ERROR]")
		}
	}
	fmt.Println("resp", resp)
}

func stream(c *gin.Context) {
	con, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer con.Close()

	cc := proto.NewCurrencyClient(con)
	sub, err := cc.Subscribe(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	for {
		srr, err := sub.Recv()
		if err != nil {
			return
		}
		ge := srr.GetError()
		// handle error
		fmt.Println(ge)
	}
}
