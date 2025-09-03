package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/GithubHttpMethods"
	"main/utils"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var accessToken = make(map[string]string)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<a href = "/login">Click for Github Login</a>`))
}

func LoginHandler(w http.ResponseWriter, r *http.Request, clientID, clientSecret, redirectURI string) {
	githubUrl := "https://github.com/login/oauth/authorize"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "repo")
	params.Add("allow_signup", "true")

	http.Redirect(w, r, githubUrl+"?"+params.Encode(), http.StatusFound)
}

func callBackHandler(w http.ResponseWriter, r *http.Request, clientID, clientSecret, redirectURI string) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	token, err := GithubHttpMethods.GetAccesstoken(code, clientID, clientSecret, redirectURI)
	accessToken["token"] = token
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get Access Token", http.StatusInternalServerError)
		return
	}

	filename := utils.GetWordHandler() + utils.GetWordHandler()
	response, err := GithubHttpMethods.CreateUserRepo(token, filename)
	if err != nil {
		fmt.Printf("Create User Repo failed: %s\n", err.Error())
		http.Error(w, "Error creating User Repositary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"clone_url": response,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URL")

	http.HandleFunc("/", HomePageHandler)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, clientID, clientSecret, redirectURI)
	})
	http.HandleFunc("/github/callback", func(w http.ResponseWriter, r *http.Request) {
		callBackHandler(w, r, clientID, clientSecret, redirectURI)
	})

	fmt.Println("Listening on port 8000: http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
