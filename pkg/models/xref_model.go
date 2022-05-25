package models

import (
	"github.com/cgeorgiades27/grpc-demo/pkg/constants"
	"github.com/google/uuid"
)

type Xref struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Value  string
	Source string
	Status constants.XrefStatus
}
