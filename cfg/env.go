package cfg

import (
	"github.com/IvanProdaiko94/ssh-test/persistence"
	"github.com/spf13/viper"
)

func ReadEnv() *Config {
	viper.AutomaticEnv()

	viper.SetEnvPrefix("APP")

	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("HEALTH_CHECK_PORT", 8888)
	viper.SetDefault("POLICY_FILE_PATH", "./logic/policy.json")

	viper.SetDefault("SQLDB_HOST", "localhost")
	viper.SetDefault("SQLDB_PORT", 5432)
	viper.SetDefault("SQLDB_USER", "postgres")
	viper.SetDefault("SQLDB_PASS", "")
	viper.SetDefault("SQLDB_DB_NAME", "db")
	viper.SetDefault("SQLDB_MAX_OPEN_CONNS", 10)

	return &Config{
		LogLevel:        viper.GetString("LOG_LEVEL"),
		HealthCheckPort: viper.GetInt("HEALTH_CHECK_PORT"),
		GameMode:        viper.GetString("GAME_MODE"),
		PolicyFilePath:  viper.GetString("POLICY_FILE_PATH"),
		PostgresConfig: persistence.SQLDBConfig{
			Host:         viper.GetString("SQLDB_HOST"),
			Port:         viper.GetInt("SQLDB_PORT"),
			User:         viper.GetString("SQLDB_USER"),
			Pass:         viper.GetString("SQLDB_PASS"),
			DBName:       viper.GetString("SQLDB_DB_NAME"),
			MaxOpenConns: viper.GetInt("SQLDB_MAX_OPEN_CONNS"),
		},
	}
}
