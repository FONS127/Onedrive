package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	clientId = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectUri = os.Getenv("REDIRECT_URI")
	scopes = os.Getenv("SCOPES")

	accessToken = os.Getenv("ACCESS_TOKEN")
	refreshToken = os.Getenv("REFRESH_TOKEN")

	authUrl = os.Getenv("AUTH_URL")
	tokenUrl = os.Getenv("TOKEN_URL")
	tenantId = os.Getenv("TENANT_ID")

	fileCount = 20
}

func UpdateEnvFileTokens(newAccessToken string, newRefreshToken string) {

	envMap := make(map[string]string)
	envMap["CLIENT_ID"] = clientId
	envMap["CLIENT_SECRET"] = clientSecret
	envMap["REDIRECT_URI"] = redirectUri
	envMap["SCOPES"] = scopes
	envMap["TENANT_ID"] = tenantId
	envMap["ACCESS_TOKEN"] = newAccessToken
	envMap["REFRESH_TOKEN"] = newRefreshToken
	envMap["AUTH_URL"] = authUrl
	envMap["TOKEN_URL"] = tokenUrl

	err := godotenv.Write(envMap, ".env")
	if err != nil {
		log.Fatal("Error writing .env file")
	}
	os.Setenv("ACCESS_TOKEN", newAccessToken)
	os.Setenv("REFRESH_TOKEN", newRefreshToken)

}

func UpdateToken() {
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	data := []byte(fmt.Sprintf("grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s&scope=%s", clientId, clientSecret, refreshToken, scopes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var token TokenResponse
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

	errAccesToken := os.Setenv("ACCESS_TOKEN", token.AccessToken)
	if err != nil {
		panic(errAccesToken)
	}
	errRefreshToken := os.Setenv("REFRESH_TOKEN", token.RefreshToken)
	if err != nil {
		panic(errRefreshToken)
	}

	UpdateEnvFileTokens(token.AccessToken, token.RefreshToken)

	accessToken = token.AccessToken
	refreshToken = token.RefreshToken
	if accessToken != "" && refreshToken != "" {
		fmt.Println("Token updated")
		return
	}
	fmt.Println("Token not updated")
}


