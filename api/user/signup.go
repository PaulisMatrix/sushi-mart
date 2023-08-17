package user

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Validator acts like a wrapper around UsersServiceImpl.AddUser
// It does validations for incoming requests
func (v *Validator) CreateUser(ctx context.Context, req *SignUpReq) *common.ErrorResponse {
	return v.UsersService.CreateUser(ctx, req)
}

func (u *UsersServiceImpl) CreateUser(ctx context.Context, req *SignUpReq) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "CreateUser", "request": req})

	//generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("error in generating hashed user password")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	dbParams := database.CreateCustomerParams{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    pgtype.Text{String: req.Phone, Valid: true},
		Address:  pgtype.Text{String: req.Address, Valid: true},
	}

	respErr := u.Queries.CreateCustomer(ctx, dbParams)
	if respErr != nil {
		logger.WithError(err).Error("failed to create a new user")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
