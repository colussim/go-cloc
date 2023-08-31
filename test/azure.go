package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {
	var organizationName string

	flag.StringVar(&organizationName, "org", "", "Nom de l'organisation Azure DevOps")
	flag.Parse()

	if organizationName == "" {
		fmt.Println("Le nom de l'organisation est requis")
		return
	}

	baseURL := fmt.Sprintf("https://dev.azure.com/%s/_apis/projects", organizationName)

	var allProjects []Project
	page := 1
	for {
		projects, err := fetchProjects(baseURL, page)
		if err != nil {
			log.Fatal(err)
		}
		allProjects = append(allProjects, projects...)
		if len(projects) == 0 {
			break
		}
		page++
	}

	fmt.Printf("Projets dans l'organisation %s:\n", organizationName)
	for _, project := range allProjects {
		fmt.Printf("- %s\n", project.Name)
	}
}

func fetchProjects(url string, page int) ([]Project, error) {
	resp, err := http.Get(fmt.Sprintf("%s?api-version=6.0&page=%d", url, page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects struct {
		Value []Project `json:"value"`
	}
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}

	return projects.Value, nil
}
