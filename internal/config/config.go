package config

type Config struct {
	Postgres `json:"database"`
	Server   `json:"server"`
}

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
