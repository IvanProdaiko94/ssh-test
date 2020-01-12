package cfg

import (
	"github.com/IvanProdaiko94/ssh-test/persistence"
)

type Config struct {
	LogLevel        string
	HealthCheckPort int
	GameMode        string
	PolicyFilePath  string
	PostgresConfig  persistence.SQLDBConfig
}
