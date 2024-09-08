package main

// Package requirements:
// - go get github.com/joho/godotenv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type RequestData struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ResponseData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	accessToken, err := getAccessToken(clientId, clientSecret)
	if err != nil {
		fmt.Println("Error getting access token:", err)
		return
	}

	fmt.Println("Access Token:", accessToken)
}

func getAccessToken(clientId, clientSecret string) (string, error) {
	// API endpoint
	baseUrl := "https://accounts.spotify.com/api/token"

	// Create a map of query parameters
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", clientId)
	params.Add("client_secret", clientSecret)

	// Create a new request
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Not needed here, but retained for reference
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	// Access the desired value
	return responseData.AccessToken, nil
}
