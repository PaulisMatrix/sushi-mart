package user

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"

	"github.com/sirupsen/logrus"
)

// Validator acts like a wrapper around UsersServiceImpl.AddUser
// It does validations for incoming requests
func (v *Validator) CreateUser(ctx context.Context, req database.CreateCustomerParams) (*database.Customer, error) {
	return v.UsersService.CreateUser(ctx, req)
}

func (u *UsersServiceImpl) CreateUser(ctx context.Context, req database.CreateCustomerParams) (*database.Customer, error) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "CreateUser"})

	insertedUser, err := u.Queries.CreateCustomer(ctx, req)
	if err != nil {
		logger.WithError(err).Error("error in creating a user")
		return nil, err
	}
	return &insertedUser, err
}
