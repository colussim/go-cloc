package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Links       struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

func main() {
	var username string

	flag.StringVar(&username, "user", "", "Nom d'utilisateur Bitbucket")
	flag.Parse()

	if username == "" {
		fmt.Println("Le nom d'utilisateur est requis")
		return
	}

	baseURL := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s", username)

	var allRepos []Repository
	page := 1
	for {
		repos, err := fetchRepositories(baseURL, page)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if len(repos) == 0 {
			break
		}
		page++
	}

	fmt.Printf("Dépôts pour l'utilisateur %s:\n", username)
	for _, repo := range allRepos {
		fmt.Printf("- %s\n", repo.Name)
	}
}

func fetchRepositories(url string, page int) ([]Repository, error) {
	resp, err := http.Get(fmt.Sprintf("%s?page=%d", url, page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos struct {
		Values []Repository `json:"values"`
	}
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos.Values, nil
}
