package config

import (
	"log"
	"os"
)

type Config struct {
	Postgres `json:"database"`
	Server   `json:"server"`
}

// Contains default values. Overwritten by environment variables
type Postgres struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

type Server struct {
	Port string `json:"port"`
}

func GetEnv() string {
	return getEnv("ENV")
}

func GetDatabaseName() string {
	return getEnv("POSTGRES_DB")
}

func GetDatabaseUser() string {
	return getEnv("POSTGRES_USER")
}

func GetDatabasePassword() string {
	return getEnv("POSTGRES_PASSWORD")
}

func GetDatabaseHost() string {
	return getEnv("POSTGRES_HOST")
}

func GetDatabasePort() string {
	return getEnv("POSTGRES_PORT")
}

func getEnv(param string) string {
	val := os.Getenv(param)
	if val == "" {
		log.Printf("env not set: %s", param)
	}
	return val
}
