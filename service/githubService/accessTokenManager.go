package githubService

import (
	"GitOperator/logger"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// STRUCT THAT REPRESENTS GITHUB ACCESS TOKEN RESPONSE
type githubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// SERVICE TO GET GITHUB ACCESS TOKEN
func GetGithubAccessToken(code string) (string, error) {
	githubClientID := GetGithubClientID()
	githubClientSecret := GetGithubClientSecret()

	reqMap := map[string]string{"client_id": githubClientID, "client_secret": githubClientSecret, "code": code}

	reqJSON, err := json.Marshal(reqMap)

	if err != nil {
		logger.ErrorLogger.Printf("ERROR IN MARSHALLING OF REQ MAP FOR ACCESS TOKEN: %v\n", err)
		return "", err
	}

	req, reqErr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(reqJSON))

	if reqErr != nil {
		logger.ErrorLogger.Printf("ERROR IN CREATING REQUEST FOR ACCESS TOKEN: %v\n", err)
		return "", reqErr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, respErr := http.DefaultClient.Do(req)

	if respErr != nil {
		logger.ErrorLogger.Printf("ERROR IN REQUEST PROCESSING FOR ACCESS TOKEN: %v\n", err)
		return "", respErr
	}

	respBytes, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		logger.ErrorLogger.Printf("ERROR IN READING FROM ACCESS TOKEN RESPONSE: %v\n", err)
		return "", readErr
	}
	gitAccTokenResp := githubAccessTokenResponse{}

	unmarshalErr := json.Unmarshal(respBytes, &gitAccTokenResp)

	if unmarshalErr != nil {
		logger.ErrorLogger.Printf("ERROR IN UNMARSHALLING ACCESS TOKEN RESPONSE: %v", err)
		return "", unmarshalErr
	}
	return gitAccTokenResp.AccessToken, nil
}
