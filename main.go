package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	// Global variables
	clientId     string
	clientSecret string
	redirectUri  string
	scopes       string
	accessToken  string
	refreshToken string
	tenantId     string
	authUrl      string
	tokenUrl     string
	fileCount    int
)

func handlePDFUpload(w http.ResponseWriter, r *http.Request) {

	// Check that the request method is POST
	pdfSavedName := "uploaded.pdf"
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	pdfBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Save the PDF bytes to a file
	err = ioutil.WriteFile(pdfSavedName, pdfBytes, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to save PDF file", http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "PDF saved successfully")

	UploadPDF(pdfSavedName)

	fmt.Fprint(w, "\nPDF uploaded successfully")

}
func loadEnvVariables() {
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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func updateEnvFileTokens(newAccessToken string, newRefreshToken string) {

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
func updateToken() {
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

	updateEnvFileTokens(token.AccessToken, token.RefreshToken)

	accessToken = token.AccessToken
	refreshToken = token.RefreshToken
	if accessToken != "" && refreshToken != "" {
		fmt.Println("Token updated")
		return
	}
	fmt.Println("Token not updated")
}

func main() {
	loadEnvVariables()
	//Start the go routine to update the token and refresh token every X minutes

	go func() {
		count := 0
		for {
			count++
			updateToken()
			time.Sleep(5 * time.Minute)
		}
	}()
	if accessToken == "" && refreshToken == "" {
		getAccessToken(clientId, redirectUri, scopes)
		accessToken = "INITIAL_TOKEN"
		refreshToken = "INITIAL_REFRESH"
		return
	}
	http.HandleFunc("/pdf-upload", handlePDFUpload)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
