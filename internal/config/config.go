package config

type Config struct {
	Postgres `json:"database"`
	Server   `json:"server"`
}

type Postgres struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Server struct {
	Port string `json:"port"`
}
