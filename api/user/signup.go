package user

import (
	"context"
	"sushi-mart/internal/database"
)

// Validator acts like a wrapper around UsersServiceImpl.AddUser
// It does validations for incoming requests
func (v *Validator) CreateUser(ctx context.Context, req database.CreateCustomerParams) (*database.Customer, error) {
	return v.UsersService.CreateUser(ctx, req)
}

func (u *UsersServiceImpl) CreateUser(ctx context.Context, req database.CreateCustomerParams) (*database.Customer, error) {
	insertedUser, err := u.Queries.CreateCustomer(ctx, req)
	if err != nil {
		return nil, err
	}
	return &insertedUser, err
}
