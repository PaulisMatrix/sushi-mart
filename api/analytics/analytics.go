package analytics

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"
)

type AnalyticsService interface {
	GetAvgCustomerRatings(context.Context) (*AvgCustomerRatingsResp, *common.ErrorResponse)
	GetMostOrdersPlaced(context.Context, int) (*MostOrdersPlacedResp, *common.ErrorResponse)
}

type AnalyticsServiceImpl struct {
	Queries *database.Queries
}

type Validator struct {
	AnalyticsService
}

type Cache struct {
	AnalyticsService
}

type RoutesWrapper struct {
	AnalyticsService
}

func New(Queries *database.Queries) *RoutesWrapper {
	return &RoutesWrapper{
		AnalyticsService: &Cache{
			AnalyticsService: &Validator{
				AnalyticsService: &AnalyticsServiceImpl{
					Queries: Queries,
				},
			},
		},
	}
}
