package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"os"

	"github.com/GSVillas/movie-pass-api/config/model"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

var Env model.Environment

func LoadEnvironments() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&Env)
	if err != nil {
		panic(err)
	}

	Env.PrivateKey, err = loadPrivateKey()
	if err != nil {
		panic(err)
	}

	Env.PublicKey, err = loadPublicKey()
	if err != nil {
		panic(err)
	}

}

func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	keyData, err := os.ReadFile("ec_private_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func loadPublicKey() (*ecdsa.PublicKey, error) {
	keyData, err := os.ReadFile("ec_public_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not ECDSA public key")
	}

	return ecdsaPubKey, nil
}

func ConfigureLogger() {
	handler := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))
	slog.SetDefault(handler)
}
