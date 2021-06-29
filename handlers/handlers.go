package handlers

import (
	"context"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	DStore "github.com/tonymj76/cuvva-sample-code/datastore"
	Models "github.com/tonymj76/cuvva-sample-code/models"
)

type Metcher interface {
	CreateMerchant(context.Context, *Models.CreateRequest) (*Models.CreateResponse, error)
}

// Service serve the rpc
type Service struct {
	MRepository DStore.MRepository
}

// NewService connects all services to postgresDB
func NewService(connect *DStore.Connection) *Service {
	return &Service{
		MRepository: connect,
	}
}

// CreateMerchant create merchant sevice
func (s *Service) CreateMerchant(ctx context.Context, req *Models.CreateRequest) (*Models.CreateResponse, error) {
	clog.Get(ctx).Info("creating a merchant...")
	return s.MRepository.Create(ctx, req)
}
