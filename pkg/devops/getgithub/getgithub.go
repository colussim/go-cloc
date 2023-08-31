package getgithub

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const baseURL = "https://api.github.com"

type Repository struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

// Browsing number of pages
func FetchRepositories(url string, page int) ([]Repository, string, error) {
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

// Get Infos for 1 Repository in Organization for Main Branch
func GetRepoGithub(accessToken, organization, repos string) (*Repository, error) {
	var repo Repository

	url := fmt.Sprintf("%s/repos/%s/%s", baseURL, organization, repos)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &repo)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

// Get Infos for all Repositories in Organization for Main Branch
func GetRepoGithubList(accessToken, organization string) ([]Repository, error) {
	var url = ""
	var repositories []Repository

	url = fmt.Sprintf("%s/orgs/%s/repos", baseURL, organization)

	page := 1
	for {
		repos, nextPageURL, err := FetchRepositories(url, page)
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
