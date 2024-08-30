package config

var dbConfig DbConfig

type DbConfig struct {
	URI    string `yaml:"uri"`
	DBName string `yaml:"db-name"`
}

func DBConfig() DbConfig {
	return dbConfig
}
