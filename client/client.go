package client

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"kalshi-starter-code-go/config"
)

type HttpClient struct {
	keyID      string
	privateKey *rsa.PrivateKey
	env        config.Environment
	baseURL    string
	client     *http.Client
	lastCall   time.Time
}

type Balance struct {
	AvailableBalance int64 `json:"available_balance"`
	TotalBalance     int64 `json:"total_balance"`
}

func NewHttpClient(keyID string, privateKey *rsa.PrivateKey, env config.Environment) *HttpClient {
	var baseURL string
	if env == config.Demo {
		baseURL = "https://demo-api.kalshi.co"
	} else {
		baseURL = "https://api.elections.kalshi.com"
	}

	return &HttpClient{
		keyID:      keyID,
		privateKey: privateKey,
		env:        env,
		baseURL:    baseURL,
		client:     &http.Client{Timeout: 10 * time.Second},
		lastCall:   time.Now(),
	}
}

func (c *HttpClient) rateLimit() {
	const threshold = 100 * time.Millisecond
	now := time.Now()
	if diff := now.Sub(c.lastCall); diff < threshold {
		time.Sleep(threshold - diff)
	}
	c.lastCall = now
}

func (c *HttpClient) createSignature(method, path string) (string, string, error) {
	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	timestampString := strconv.FormatInt(timestamp, 10)

	parsedPath, err := url.Parse(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse path: %w", err)
	}

	msgString := timestampString + method + parsedPath.Path
	fmt.Println("MESSAGE STRING", msgString)
	hashed := sha256.Sum256([]byte(msgString))
	signature, err := rsa.SignPSS(rand.Reader, c.privateKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto})
	if err != nil {
		return "", "", fmt.Errorf("failed to create signature: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), timestampString, nil
}

func (c *HttpClient) createRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	signature, timestamp, err := c.createSignature(method, path)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader = http.NoBody
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KALSHI-ACCESS-KEY", c.keyID)
	req.Header.Set("KALSHI-ACCESS-SIGNATURE", signature)
	req.Header.Set("KALSHI-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *HttpClient) doRequest(req *http.Request, v interface{}) error {
	c.rateLimit()

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode repsonse: %w", err)
		}
	}

	return nil
}

func (c *HttpClient) GetBalance(ctx context.Context) (*Balance, error) {
	req, err := c.createRequest(ctx, http.MethodGet, "/trade-api/v2/portfolio/balance", nil)
	if err != nil {
		return nil, err
	}

	var response Balance
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
