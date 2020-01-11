package cfg

import (
	"github.com/IvanProdaiko94/ssh-test/logic"
	"github.com/IvanProdaiko94/ssh-test/persistence"
)

type Config struct {
	LogLevel        string
	HealthCheckPort int
	GameMode        logic.Mode
	PostgresConfig  persistence.SQLDBConfig
}
