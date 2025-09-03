package GithubHttpMethods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/models"
	"net/http"
)

func CreateUserRepo(token string, name string) (interface{}, error) {
	repo := models.Repositary{
		Name:        name,
		Description: "This Project is Created by the Imagine.bo",
		Private:     false,
	}

	body, err := json.Marshal(repo)
	if err != nil {
		return "", err
	}
	reqPath := fmt.Sprintf("https://api.github.com/user/repos")
	req, err := http.NewRequest("POST", reqPath, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("GitHub API error: %s", resp.Status)
	}
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	// cloneURI := result["clone_url"].(string)
	return result, nil
}
