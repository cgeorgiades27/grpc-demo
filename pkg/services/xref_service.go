package services

import (
	"context"
	"errors"
	"log"

	pb "github.com/cgeorgiades27/grpc-demo/pkg/services/proto"
	services "github.com/cgeorgiades27/grpc-demo/pkg/services/proto"
)

// type xrefService struct {
// 	db *gorm.DB
// }

type XrefServer struct {
	pb.UnimplementedXrefServiceServer
}

func (x *XrefServer) GetXref(ctx context.Context, in *pb.XrefRequest) (*pb.XrefResponse, error) {

	if len(in.GetLastfour()) < 4 {
		return nil, errors.New("invalid request")
	}

	log.Printf("request received: %s", in.GetLastfour())
	return &pb.XrefResponse{
		Token: &services.XREF{
			Value:  in.GetLastfour(),
			Source: "T",
			Status: 0,
		},
	}, nil
}
