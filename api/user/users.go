package user

import (
	"context"
	"sushi-mart/internal/database"
)

type UsersService interface {
	CreateUser(context.Context, database.CreateCustomerParams) (*database.Customer, error)
	GetUser(context.Context, string) (*database.Customer, error)
}

type UsersServiceImpl struct {
	Queries *database.Queries
}

type Cache struct {
	UsersService
}
type Validator struct {
	UsersService
}

type RoutesWrapper struct {
	UsersService
}

func New(Queries *database.Queries) *RoutesWrapper {
	return &RoutesWrapper{
		UsersService: &Validator{
			UsersService: &Cache{
				UsersService: &UsersServiceImpl{
					Queries: Queries,
				},
			},
		},
	}
}
