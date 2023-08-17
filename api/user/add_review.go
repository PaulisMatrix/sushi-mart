package user

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

func (v *Validator) AddReview(ctx context.Context, req *AddReviewReq, custId int) *common.ErrorResponse {
	return v.UsersService.AddReview(ctx, req, custId)
}

func (u *UsersServiceImpl) AddReview(ctx context.Context, req *AddReviewReq, custId int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "AddReview", "request": req})

	// first check if the user had placed an order for the product they are adding review for
	_, err := u.Queries.ValidateProductOrderReview(ctx, pgtype.Int4{Int32: int32(custId), Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			// user hasn't purchased that product yet
			logger.WithError(err).Info("user hasn't purchased this product yet")
			return &common.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "need to purchase this product before reviewing",
			}
		}
		logger.WithError(err).Error("failed to get order for the corresponding product_id")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	dbParams := database.AddReviewParams{
		Rating:     int32(req.Rating),
		ReviewText: req.ReviewText,
		ReviewDate: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		CustomerID: pgtype.Int4{Int32: int32(custId), Valid: true},
		ProductID:  pgtype.Int4{Int32: int32(req.ProductId), Valid: true},
	}

	reviewErr := u.Queries.AddReview(ctx, dbParams)
	if reviewErr != nil {
		logger.WithError(reviewErr).Error("error in adding a new customer review")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
