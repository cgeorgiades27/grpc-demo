package models

import (
	"github.com/cgeorgiades27/grpc-demo/pkg/constants"
)

type Xref struct {
	Value string
}

type XrefRequest struct {
	LastFour string
}

type XrefResponse struct {
	XREF   Xref
	Status constants.XrefStatus
}
