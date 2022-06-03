package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/cgeorgiades27/grpc-demo/pkg/xref"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AVAILABLE   string = "available"
	UNAVAILABLE string = "unavailable"
)

var (
	serverAddr = flag.String("grpcsvr", "localhost:50051", "grpc server address")
	ctx        context.Context
)

func main() {

	flag.Parse()

	// set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to connect to grpc server: %v", err)
	}
	defer conn.Close()

	ctx = context.Background()
	xsvc := xref.NewXrefServiceClient(conn)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("grpcsvr", conn)
		c.Set("xsvc", xsvc)
	})

	// grpc route group
	rg := r.Group("/grpc")
	{
		rg.GET("/getxref/:num", getXref)                                // simple rpc
		rg.GET("/addxrefs/:min/:max", addXrefs)                         // client streaming rpc
		rg.GET("/getxrefs/:min/:max", getXrefs)                         // bidirectional streaming rpc
		rg.GET("/getmagicnumbers/:status", getMagicNumbers)             // server streaming rpc
		rg.GET("/getmagicnumbersummary/:status", getMagicNumberSummary) // simple rpc
	}
	r.Run()
}

func getXref(c *gin.Context) {
	xsvc := c.MustGet("xsvc").(xref.XrefServiceClient)
	lf := c.Param("num")

	/*******
	* gRPC *
	********/

	res, err := xsvc.GetXref(ctx, &xref.XrefRequest{
		Lastfour: lf,
	})
	if err != nil {
		log.Printf("err: %v", err)
	} else {
		log.Println("received", res.Token.GetValue())
	}
}

func addXrefs(c *gin.Context) {

	// get params and validate request
	minP := c.Param("min")
	maxP := c.Param("max")
	if len(minP) < 4 || len(maxP) < 4 {
		c.AbortWithError(500, errors.New("min 4 digit num"))
		return
	}
	min, err := strconv.Atoi(c.Param("min"))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	max, err := strconv.Atoi(c.Param("max"))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	if max < min {
		c.AbortWithError(500, errors.New("max must be > min"))
		return
	}

	/*******
	* gRPC *
	********/

	xsvc := c.MustGet("xsvc").(xref.XrefServiceClient)
	stream, err := xsvc.AddXrefs(ctx)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	for i := min; i < max; i++ {
		if err := stream.Send(&xref.XrefRequest{Lastfour: strconv.Itoa(i)}); err != nil {
			return
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return
	}

	// event summary
	log.Printf("\n\n******************************\n%-10s%10d\n%-10s%10d\n%-10s%15d\n******************************\n", "Total New", res.TotalNew, "Total Updated", res.TotalUpdated, "Elapsed Time", res.ElapsedTime)
}

func getMagicNumbers(c *gin.Context) {

	var s xref.Status_STATUS
	statP := c.Param("status")
	switch statP {
	case AVAILABLE:
		s = xref.Status_AVAILABLE
	case UNAVAILABLE:
		s = xref.Status_UNAVAILABLE
	default:
		c.AbortWithError(400, errors.New("unknown status requested"))
		return
	}

	/*******
	* gRPC *
	********/

	xsvc := c.MustGet("xsvc").(xref.XrefServiceClient)
	stream, err := xsvc.GetMagicNumbers(ctx, &xref.Status{Status: s})
	if err != nil {
		return
	}

	for {
		magicNum, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		log.Println(magicNum.Value)
	}
}

func getXrefs(c *gin.Context) {

	// get params and validate request
	minP := c.Param("min")
	maxP := c.Param("max")
	if len(minP) < 4 || len(maxP) < 4 {
		c.AbortWithError(500, errors.New("min 4 digit num"))
		return
	}
	min, err := strconv.Atoi(c.Param("min"))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	max, err := strconv.Atoi(c.Param("max"))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	if max < min {
		c.AbortWithError(500, errors.New("max must be > min"))
		return
	}

	/*******
	* gRPC *
	********/

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	xsvc := c.MustGet("xsvc").(xref.XrefServiceClient)

	stream, err := xsvc.GetXrefs(ctx)
	if err != nil {
		return
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				return
			}
			log.Println("received", in.Token.GetValue())

		}
	}()

	for i := min; i < max; i++ {
		if err := stream.Send(&xref.XrefRequest{Lastfour: strconv.Itoa(i)}); err != nil {
			return
		}
	}

	stream.CloseSend()
	<-waitc
}

func getMagicNumberSummary(c *gin.Context) {

	var s xref.Status_STATUS
	statP := c.Param("status")
	switch statP {
	case AVAILABLE:
		s = xref.Status_AVAILABLE
	case UNAVAILABLE:
		s = xref.Status_UNAVAILABLE
	default:
		c.AbortWithError(400, errors.New("unknown status requested"))
		return
	}

	/*******
	* gRPC *
	********/

	xsvc := c.MustGet("xsvc").(xref.XrefServiceClient)
	summary, err := xsvc.GetMagicNumberSummary(ctx, &xref.Status{Status: s})
	if err != nil {
		return
	}
	log.Printf("Total %s: %d", statP, summary.Total)
}
