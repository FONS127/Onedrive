package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

func main() {
	LoadEnvVariables()
	//Start the go routine to update the token and refresh token every X minutes
	go func() {
		count := 0
		for {
			count++
			UpdateToken()
			fmt.Println("Token updated", count, "times")
			time.Sleep(5 * time.Minute)
		}
	}()

	http.HandleFunc("/pdf-upload", handlePDFUpload)
	http.HandleFunc("/get-token", getAccessToken)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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

func getAccessToken(w http.ResponseWriter, r *http.Request) {
	// Get the access token from the refresh token
	fmt.Println()
	authURL := fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%v&response_type=code&redirect_uri=%v&scope=%v&state=qwerty", clientId, redirectUri, scopes)
	fmt.Println(authURL)

	// // Create HTTP client
	// client := &http.Client{}

	// // Create request
	// req, err := http.NewRequest("GET", authURL, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// // Make request
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// // Get redirected URL
	// redirectedUrl := resp.Request.URL.String()

	// fmt.Printf("Redirected URL: %s\n", redirectedUrl)
	// // Parse URL string
	// u, err := url.Parse(redirectedUrl)
	// if err != nil {
	// 	panic(err)
	// }

	// // Get query parameters
	// code := u.Query().Get("code")
	// state := u.Query().Get("state")

	// fmt.Printf("Code: %s\n", code)
	// fmt.Printf("State: %s\n", state)

}
