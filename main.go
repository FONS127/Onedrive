package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	// Global variables
	appId        string
	clientSecret string
	redirectUri  string
	scopes       string
	accessToken  string
	refreshToken string
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

	fmt.Fprint(w, "PDF uploaded successfully")

}

func main() {
	//Load env variables as global variables
	//Start the go routine to update the token and refresh token every X minutes
	//Call the handlePDFUpload function
	//Call the UploadPDF function with the env variables and the pdf file name
	accessToken = os.Getenv("TOKEN")
	http.HandleFunc("/pdf-upload", handlePDFUpload)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
