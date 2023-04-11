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

	fileCount = 20
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

	accessToken = token.AccessToken
	refreshToken = token.RefreshToken
	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Refresh Token:", token.RefreshToken)
}

func main() {
	loadEnvVariables()
	//Start the go routine to update the token and refresh token every X minutes

	go func() {
		count := 0
		for {
			count++
			time.Sleep(20 * time.Second)
			//Post request to get the new access,refresh tokens

			updateToken()

			fmt.Println("New Acces Token: " + accessToken)
			fmt.Println("New Refresh Token: " + refreshToken)
		}
	}()
	http.HandleFunc("/pdf-upload", handlePDFUpload)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
