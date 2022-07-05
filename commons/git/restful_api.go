package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v41/github"
)

func GetAccessTokenByClient(clientID, clientSecret, code string) (string, error) {
	requestUrl := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)

	req, err := http.NewRequest("POST", requestUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &result)

	if result["error_description"] != nil {
		return "", errors.New(result["error_description"].(string))
	}

	return result["access_token"].(string), nil
}

// Get Git user by user access token using GitHub API
// The reason why not using Git client: to avoid the authentication token confliction.
func GetUserByToken(token string) (*github.User, error) {
	requestUrl := "https://api.github.com/user"

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	githubUser := &github.User{}
	err = json.Unmarshal(body, githubUser)
	if err != nil {
		return nil, err
	}

	return githubUser, nil
}
