package config

type Config struct {
	GinMode         string `env:"GIN_MODE" env-default:"debug"`
	Port            string `env:"PORT" env-default:"8080"`
	PaginationLimit int    `env:"PAGINATION_LIMIT" env-default:"5"`
	DbSource        string `env:"DB_SOURCE" env-default:"host=db user=postgres password=example dbname=postgres port=5432 sslmode=disable"`
}
