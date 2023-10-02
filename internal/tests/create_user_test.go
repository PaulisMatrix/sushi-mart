package tests

import (
	"context"
	"sushi-mart/internal/database"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (db *DatabaseSuite) TestCreateUser() {
	//generate hashed password
	testPassword := "testpassword"
	testUserName := "testuser"
	testUserEmail := "testing@gmail.com"
	testUserPhone := "98789876554"
	testUserAddres := "Pune,Maharashtra"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	db.Nil(err)

	dbParams := database.CreateCustomerParams{
		Username: testUserName,
		Password: string(hashedPassword),
		Email:    testUserEmail,
		Phone:    pgtype.Text{String: testUserPhone, Valid: true},
		Address:  pgtype.Text{String: testUserAddres, Valid: true},
	}

	respErr := db.queries.CreateCustomer(context.Background(), dbParams)
	db.Nil(respErr)

	//verify customer creation.
	resp, err := db.queries.GetCustomer(context.Background(), testUserEmail)
	db.Nil(err)
	db.NotNil(resp)
	db.EqualValues(testUserName, resp.Username)
}
