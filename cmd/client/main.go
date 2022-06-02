package main

import (
	"context"
	"flag"
	"log"
	"strconv"

	"github.com/cgeorgiades27/grpc-demo/pkg/xref"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("grpcsvr", "localhost:50051", "grpc server address")
)

func main() {

	flag.Parse()

	// set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to connect to grpc server: %v", err)
	}
	defer conn.Close()

	xsvc := xref.NewXrefServiceClient(conn)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("grpcsvr", conn)
		c.Set("xrefsvc", xsvc)
	})

	// grpc route group
	rg := r.Group("/grpc")
	{
		rg.GET("/xref/:num", getXref)
		rg.GET("/xrefs/", getXrefs)
	}
	r.Run()
}

func getXref(c *gin.Context) {
	xsvc := c.MustGet("xrefsvc").(xref.XrefServiceClient)
	lf := c.Param("num")
	res, err := xsvc.GetXref(context.Background(), &xref.XrefRequest{
		Lastfour: lf,
	})
	if err != nil {
		log.Printf("err: %v", err)
	} else {
		log.Println("received", res.Token.GetValue())
	}
}

func getXrefs(c *gin.Context) {
	xsvc := c.MustGet("xrefsvc").(xref.XrefServiceClient)
	stream, err := xsvc.AddXrefs(context.Background())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	for i := 1300; i < 1400; i++ {
		if err := stream.Send(&xref.XrefRequest{Lastfour: strconv.Itoa(i)}); err != nil {
			log.Printf("error: %v", err)
			return
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Summary: %v", res)
}
