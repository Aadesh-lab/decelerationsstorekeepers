package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetWordHandler() string {
	req, _ := http.NewRequest("GET", "https://random-word-api.herokuapp.com/word", nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Failed to GET word"
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result []string
	json.Unmarshal(body, &result)

	return result[0]
}
