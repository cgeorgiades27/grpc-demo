package xref

import (
	"context"
	"errors"
	"log"
)

// type xrefService struct {
// 	db *gorm.DB
// }

func NewXrefService() *xrefServer {
	return &xrefServer{
		UnimplementedXrefServiceServer: services.UnimplementedXrefServiceServer{},
	}
}

type xrefServer struct {
	services.UnimplementedXrefServiceServer
}

func (x *xrefServer) GetXref(ctx context.Context, in *services.XrefRequest) (*services.XrefResponse, error) {

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
