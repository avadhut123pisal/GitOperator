package githubService

import (
	"GitOperator/logger"
	"os"
)

// RETURNS GITHUB CLIENT ID
func GetGithubClientID() string {
	githubClientID, isExists := os.LookupEnv("CLIENT_ID")
	if !isExists {
		logger.FatalLogger.Fatalf("GITHUB CLIENT ID NOT FOUND")
	}
	return githubClientID
}

// RETURNS GITHUB SECRET
func GetGithubClientSecret() string {
	githubClientSecret, isExists := os.LookupEnv("CLIENT_SECRET")
	if !isExists {
		logger.FatalLogger.Fatalf("GITHUB CLIENT SECRET NOT FOUND")
	}
	return githubClientSecret
}
