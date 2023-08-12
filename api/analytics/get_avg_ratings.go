package analytics

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"

	"github.com/sirupsen/logrus"
)

func (v *Validator) GetAvgCustomerRatings(ctx context.Context) (*AvgCustomerRatingsResp, *common.ErrorResponse) {
	return v.AnalyticsService.GetAvgCustomerRatings(ctx)
}

func (a *AnalyticsServiceImpl) GetAvgCustomerRatings(ctx context.Context) (*AvgCustomerRatingsResp, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetAvgCustomerRatings"})
	var avgRatings []AvgCustomerRatings

	resp, err := a.Queries.GetAvgCustomerRatings(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithError(err).Error("records not found, add ratings and products first")
			return nil, &common.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "records not found. add ratings and products first",
			}
		}
		logger.WithError(err).Error("error in fetching customer ratings")
		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	for _, r := range resp {
		avgRating, _ := strconv.ParseFloat(r.AverageRating, 64)
		avgRatings = append(avgRatings, AvgCustomerRatings{
			ProductName:     r.Name,
			ProductCategory: r.Category,
			AvgRating:       avgRating,
		})
	}

	return &AvgCustomerRatingsResp{
		AvgRatings: avgRatings,
	}, nil
}
