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
	HTMLURL     string `json:"html_url"`
}

func main() {
	var orgName string
	var repoName string

	flag.StringVar(&orgName, "org", "", "Nom de l'organisation GitHub")
	flag.StringVar(&repoName, "repo", "", "Nom du dépôt GitHub (facultatif)")

	flag.Parse()

	if orgName == "" {
		fmt.Println("Le nom de l'organisation est requis")
		return
	}

	baseURL := fmt.Sprintf("https://api.github.com/orgs/%s/repos", orgName)

	if repoName != "" {
		// Récupérer les informations sur un dépôt spécifique
		repoURL := fmt.Sprintf("%s/%s", baseURL, repoName)
		repo, err := fetchRepository(repoURL)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Nom du dépôt: %s\n", repo.Name)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("URL: %s\n", repo.HTMLURL)
	} else {
		// Liste de tous les dépôts de l'organisation
		repos, err := fetchRepositories(baseURL)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Dépôts dans l'organisation %s:\n", orgName)
		for _, repo := range repos {
			fmt.Printf("- %s\n", repo.Name)
		}
	}
}

func fetchRepository(url string) (*Repository, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repo Repository
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func fetchRepositories(url string) ([]Repository, error) {
	resp, err := http.Get(url)
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
