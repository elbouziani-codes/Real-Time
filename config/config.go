package config

import "time"


type Config struct {
	Port            string
	DBPath          string
	SessionDuration time.Duration
}
func Load() Config {
	return Config{
		Port:            ":8080",
		DBPath:          "./forum.db",
		SessionDuration: 24 * time.Hour,
	}
}