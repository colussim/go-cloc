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
	Description string `json:"description"`
	WebURL      string `json:"web_url"`
}

func main() {
	var groupName string

	flag.StringVar(&groupName, "group", "", "Nom du groupe GitLab")
	flag.Parse()

	if groupName == "" {
		fmt.Println("Le nom du groupe est requis")
		return
	}

	baseURL := fmt.Sprintf("https://gitlab.com/api/v4/groups/%s/projects", groupName)

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

	fmt.Printf("Projets dans le groupe %s:\n", groupName)
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

	var repos []Repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func getRepositoryList(accessToken, organization, devops string) ([]Repository, error) {
	var url = ""
	var repositories []Repository

	url = fmt.Sprintf("https://api.%s.com/orgs/%s/repos", devops, organization)

	page := 1
	for {
		repos, nextPageURL, err := fetchRepositories(url, page)
		if err != nil {
			log.Fatal(err)
		}
		repositories = append(repositories, repos...)
		if nextPageURL == "" {
			break
		}
		page++
	}

	return repositories, nil
}
