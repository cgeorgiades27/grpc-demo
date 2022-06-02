package xref

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cgeorgiades27/grpc-demo/pkg/constants"
	"github.com/cgeorgiades27/grpc-demo/pkg/models"
	"github.com/go-redis/redis/v8"
)

const (
	AVAILABLE   = "available"
	UNAVAILABLE = "unavailable"
)

func NewXrefService(ctx context.Context, rds *redis.Client) *xrefServer {
	return &xrefServer{
		UnimplementedXrefServiceServer: UnimplementedXrefServiceServer{},
		ctx:                            ctx,
		redis:                          rds,
	}
}

type xrefServer struct {
	UnimplementedXrefServiceServer
	ctx   context.Context
	redis *redis.Client
}

func (x *xrefServer) InitData(path string) error {

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	r := bufio.NewScanner(f)
	for r.Scan() {
		x.redis.RPush(x.ctx, AVAILABLE, r.Text())
	}
	return nil
}

// GetXref accepts an Xref Request (last 4) and returns a Xref Response with XREF num
func (x *xrefServer) GetXref(ctx context.Context, in *XrefRequest) (*XrefResponse, error) {
	xrefRes, err := x.getXref(&models.XrefRequest{LastFour: in.GetLastfour()})
	if err != nil {
		return nil, err
	}
	return &XrefResponse{
		Token: &XREF{
			Value: xrefRes.XREF.Value,
		},
	}, nil
}

func (x *xrefServer) AddXrefs(stream XrefService_AddXrefsServer) error {

	// counters
	totalNew, totalUpdated := 0, 0
	startTime := time.Now()

	// loop stream and build summary
	for {
		xrefReq, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&XrefSummary{
				TotalNew:     int32(totalNew),
				TotalUpdated: int32(totalUpdated),
				ElapsedTime:  int32(endTime.Sub(startTime)),
			})
		}
		if err != nil {
			return err
		}

		xrefRes, err := x.getXref(&models.XrefRequest{LastFour: xrefReq.GetLastfour()})
		if err != nil {
			return err
		}

		switch xrefRes.Status {
		case constants.NEW:
			totalNew++
		case constants.EXISTING:
			totalUpdated++
		}
	}
}

func (x *xrefServer) getXref(xrefReq *models.XrefRequest) (*models.XrefResponse, error) {

	// has to be len 4
	if len(xrefReq.LastFour) != 4 {
		return nil, errors.New("bad request")
	}

	// magic num status
	status := constants.EXISTING

	val, err := x.redis.Get(x.ctx, xrefReq.LastFour).Result()
	if err != nil {
		magicNum, err := x.redis.RPop(x.ctx, AVAILABLE).Result()
		if err != nil {
			return nil, err
		}

		status = constants.NEW
		val = magicNum + xrefReq.LastFour
		x.redis.Set(x.ctx, xrefReq.LastFour, val, 0)

		// write magic num to UNAVAIL list
		go x.redis.RPush(x.ctx, UNAVAILABLE, magicNum)
	}

	return &models.XrefResponse{
		XREF:   models.Xref{Value: val},
		Status: status,
	}, nil
}
