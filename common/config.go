package common

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PgDbName  string
	PgUser    string
	PgPass    string
	JwtSktKey string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil
	}
	return &Config{
		PgDbName:  os.Getenv("POSTGRES_DB"),
		PgUser:    os.Getenv("POSTGRES_USER"),
		PgPass:    os.Getenv("sushimartpass"),
		JwtSktKey: os.Getenv("MyLittleSecret"),
	}
}
