package model

import "crypto/ecdsa"

type Environment struct {
	ConnectionString     string `env:"CONNECTION_STRING"`
	RedisAdress          string `env:"REDIS_ADRESS"`
	RedisPassword        string `env:"REDIS_PASSWORD"`
	APIPort              string `env:"API_PORT"`
	FrontURL             string `env:"FRONT_URL"`
	CloudFlareAccountAPI string `env:"CLOUD_FLARE_ACCOUNT_API"`
	RedisDB              int    `env:"REDIS_DB"`
	SessionExp           int    `env:"SESSION_EXP"`
	PrivateKey           *ecdsa.PrivateKey
	PublicKey            *ecdsa.PublicKey
}
