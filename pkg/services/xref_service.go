package services

import (
	"context"
	"log"

	services "github.com/cgeorgiades27/grpc-demo/pkg/services/proto"
	"gorm.io/gorm"
)

type xrefService struct {
	db *gorm.DB
}

type xrefServer struct {
	services.UnimplementedXrefServiceServer
}

func (x *xrefServer) GetXref(ctx context.Context, in *services.XrefRequest) (*services.XrefResponse, error) {
	log.Printf("request received: %s", in.GetLastfour())
	return &services.XrefResponse{
		Token: &services.XREF{
			Value:  "",
			Source: "",
			Status: 0,
		},
	}, nil
}
