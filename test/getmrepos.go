package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Repository struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

func main() {

	orgName := "SonarSource"

	baseURL := fmt.Sprintf("https://api.github.com/orgs/%s/repos", orgName)

	var allRepos []Repository
	page := 1
	for {
		repos, nextPageURL, err := fetchRepositories(baseURL, page)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if nextPageURL == "" {
			break
		}
		page++
	}

	fmt.Printf("Dépôts dans l'organisation %s:\n", orgName)
	for _, repo := range allRepos {
		fmt.Printf("- %s\n", repo.Name)
	}
}

func fetchRepositories(url string, page int) ([]Repository, string, error) {
	resp, err := http.Get(fmt.Sprintf("%s?page=%d", url, page))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, "", err
	}

	nextPageURL := ""
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		links := strings.Split(linkHeader, ",")
		for _, link := range links {
			parts := strings.Split(strings.TrimSpace(link), ";")
			if len(parts) == 2 && strings.TrimSpace(parts[1]) == `rel="next"` {
				nextPageURL = strings.Trim(parts[0], "<>")
			}
		}
	}

	return repos, nextPageURL, nil
}

func fetchRepositoriesGitlab(url string, page int) ([]Repository, error) {
	resp, err := http.Get(fmt.Sprintf("%s?page=%d", url, page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
