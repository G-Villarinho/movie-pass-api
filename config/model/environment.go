package model

import "crypto/ecdsa"

type Environment struct {
	ConnectionString string `env:"CONNECTION_STRING"`
	RedisAdress      string `env:"REDIS_ADRESS"`
	RedisPassword    string `env:"REDIS_PASSWORD"`
	RedisDB          int    `env:"REDIS_DB"`
	APIPort          string `env:"API_PORT"`
	SessionExp       int    `env:"SESSION_EXP"`
	PrivateKey       *ecdsa.PrivateKey
	PublicKey        *ecdsa.PublicKey
}
