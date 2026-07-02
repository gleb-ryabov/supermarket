package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net"

	"github.com/spf13/viper"
)

const (
	defaultReadTimeout  = 10
	defaultWriteTimeout = 10
)

// Init - initialize viper configuration.
func Init() error {
	setDefaultsConfig()

	viper.AutomaticEnv()

	configFile := viper.GetString(File)
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		var pe *fs.PathError
		if errors.As(err, &pe) {
			return nil
		}

		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	return nil
}

// setDefaultsConfig - configure default settings.
func setDefaultsConfig() {
	viper.SetDefault(File, ".env")
	viper.SetDefault(LogLevel, "ERROR")
	viper.SetDefault(DBHost, "localhost")
	viper.SetDefault(DBPort, "5432")
	viper.SetDefault(DBUser, "postgres")
	viper.SetDefault(DBPassword, "postgres")
	viper.SetDefault(DBName, "supermarket")
	viper.SetDefault(DBSSL, "disable")
	viper.SetDefault(HTTPAddress, "0.0.0.0")
	viper.SetDefault(HTTPPort, "8080")
	viper.SetDefault(ReadTimeout, defaultReadTimeout)
	viper.SetDefault(WriteTimeout, defaultWriteTimeout)
}

// GetPostgresDSN - return url for connect to PostgreSQL.
func GetPostgresDSN() string {
	hostPort := net.JoinHostPort(
		viper.GetString(DBHost),
		viper.GetString(DBPort),
	)

	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		viper.GetString(DBUser),
		viper.GetString(DBPassword),
		hostPort,
		viper.GetString(DBName),
		viper.GetString(DBSSL),
	)
}

// GetServerURL - return url for connect to server.
func GetServerURL() string {
	return net.JoinHostPort(
		viper.GetString(HTTPAddress),
		viper.GetString(HTTPPort),
	)
}
