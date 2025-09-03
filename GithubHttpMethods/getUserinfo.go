package GithubHttpMethods

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func getUserRepoInfo(token string) ([]map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {

		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var repo []map[string]interface{}
	json.Unmarshal(body, &repo)
	return repo, nil
}
