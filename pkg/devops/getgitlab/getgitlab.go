package getgitlab

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const baseURL = "https://gitlab.com/api/v4"

type Repository struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

/*func FetchRepositories2(url string, page int) ([]Repository, error) {
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
}*/

// Browsing number of pages
func FetchRepositories2(url string, page int) ([]Repository, string, error) {
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

// Get Infos for all Repositories in Organization for Main Branch
func GetRepoGitlabList(accessToken, organization string) ([]Repository, error) {
	var url = ""
	var repositories []Repository

	url = fmt.Sprintf("%s/groups/%s/projects?include_subgroups=true", baseURL, organization)
	fmt.Print(url)

	page := 1
	for {
		repos, nextPageURL, err := FetchRepositories2(url, page)
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
