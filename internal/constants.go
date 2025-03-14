package internal

import (
	"github.com/joho/godotenv"
	"os"
)

var PostgresUser string
var PostgresPassword string
var PostgresDB string
var PostgresHost string
var PostgresPort string

func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPassword = os.Getenv("POSTGRES_PW")
	PostgresDB = os.Getenv("POSTGRES_DB")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	return nil
}

type StatusResponse struct {
	Status string `json:"status"`
}

type InfoResponse struct {
	*StatusResponse
	Description string `json:"description"`
}
