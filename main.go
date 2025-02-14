package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"kalshi-starter-code-go/client"
	"kalshi-starter-code-go/config"

	"github.com/joho/godotenv"
)

func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	env := config.Demo

	var (
		keyID   string
		keyFile string
	)

	if env == config.Demo {
		keyID = os.Getenv("DEMO_KEYID")
		keyFile = os.Getenv("DEMO_KEYFILE")
	} else {
		keyID = os.Getenv("PROD_KEYID")
		keyFile = os.Getenv("PROD_KEYFILE")
	}

	if keyID == "" || keyFile == "" {
		log.Fatal("Required environment variables are not set")
	}

	privateKey, err := loadPrivateKey(keyFile)
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	ctx := context.Background()
	httpClient := client.NewHttpClient(keyID, privateKey, env)

	balance, err := httpClient.GetBalance(ctx)
	if err != nil {
		log.Fatal("Failed to get balance: ", err)
	}

	fmt.Println("Balance:", balance)
}
