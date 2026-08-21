package config


type Config struct {
	Port string
	DBPath          string
}

func Load() Config {
	return Config{
		Port: ":8081",
		DBPath:          "./realTime.db",
	}
}
