// ROUTES PACKAGE HANDLES THE REQUESTS FOR GITHUB OPERATIONS
package routes

import (
	"GitOperator/logger"
	"GitOperator/service/githubService"
	"fmt"
	"io"
	"net/http"
)

// REGISTER THE HANDLERS
func InitialiseRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/loginWithGithub", githubLoginHandler)
	mux.HandleFunc("/login/github/callback", githubCallbackHandler)
}

// ROOT MAPPING HANDLER
func rootHandler(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, `<a href="/loginWithGithub">Login With Github</a>`)
}

// MAPPING FOR LOGIN
func githubLoginHandler(rw http.ResponseWriter, req *http.Request) {
	githubClientID := githubService.GetGithubClientID()
	// CREATE GITHUB REDIRECTION URL
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s", githubClientID, "user%20public_repo", "http://localhost:3030/login/github/callback")
	http.Redirect(rw, req, redirectURL, 301)
}

// IT WILL BE CALLED BY GIHUB AFTER AUTHENTICATION
func githubCallbackHandler(rw http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	accessToken, err := githubService.GetGithubAccessToken(code)
	if err != nil {
		logger.ErrorLogger.Printf("COULD NOT GET ACCESS TOKEN: %v\n", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	gitOpErr := githubService.HandleGithubOperations(accessToken)
	if gitOpErr != nil {
		logger.ErrorLogger.Printf("ERROR OCCURED IN PROCESSING GIT OPERATIONS: %v\n", err)
		http.Error(rw, gitOpErr.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, `<b>Git operator has performed some changes in your repo...</b>`)
}
