package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getAccessToken(clientId string, redirectUri string, scopes string) {
	// Get the access token from the refresh token
	fmt.Println()
	fmt.Printf("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=%v&response_type=code&redirect_uri=%v&scope=%v&state=qwerty", clientId, redirectUri, scopes)
	fmt.Println()

}

func UploadPDF(fileName string) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	if accessToken == "" && refreshToken == "" {
		getAccessToken(clientId, redirectUri, scopes)
		accessToken = "INITIAL_TOKEN"
		refreshToken = "INITIAL_REFRESH"
		return
	}

	// Set the URL for uploading the file
	uploadURL := "https://graph.microsoft.com/v1.0/me/drive/root:/FILE" + strconv.Itoa(fileCount) + ".pdf:/content"
	fileCount++
	// Open the PDF file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file content
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new file part
	part, err := writer.CreateFormFile("file", "FILE_NAME.pdf")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Write the file content to the file part
	_, err = part.Write(fileBytes)
	if err != nil {
		fmt.Println("Error writing file part:", err)
		return
	}

	// Close the multipart form
	writer.Close()

	// Create a new HTTP request with the multipart form data
	req, err := http.NewRequest("PUT", uploadURL, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the access token in the request header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Error uploading file: HTTP status", resp.StatusCode)
		//print the error
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)

		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return
	}

	// Print the upload success message
	fmt.Println("File uploaded successfully!")
}
