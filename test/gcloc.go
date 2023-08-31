package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JoaoDanielRufino/gcloc/internal/constants"
	"github.com/JoaoDanielRufino/gcloc/pkg/gcloc"
	"github.com/briandowns/spinner"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Repository struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

// Declare a struct for Config fields
type Configuration struct {
	AccessToken  string
	Organization string
	DevOps       string
	Repos        string
}

type Report struct {
	TotalFiles      int `json:",omitempty"`
	TotalLines      int
	TotalBlankLines int
	TotalComments   int
	TotalCodeLines  int
	Results         interface{}
}

// Read Config file : Config.json
func GetConfig(configjs Configuration) Configuration {

	fconfig, err := os.ReadFile("config.json")
	if err != nil {
		panic("Problem with the configuration file : config.json")
		os.Exit(1)
	}
	json.Unmarshal(fconfig, &configjs)
	return configjs
}

// Parse Result Files in JSON Format
func parseJSONFile(filePath, reponame string) int {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	var report Report
	err = json.Unmarshal(file, &report)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	//fmt.Printf("\nTotal Lines Of Code : %d\n\n", report.TotalCodeLines)

	return report.TotalCodeLines
}

func containsNextLink(linkHeader string) bool {
	return linkHeader != "" && (strings.Contains(linkHeader, `rel="next"`) || strings.Contains(linkHeader, `rel="last"`))
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

// Get Infos for 1 Repository in Organization for Main Branch
func getRepository(accessToken, organization, devops, repos string) (*Repository, error) {
	var repo Repository

	url := fmt.Sprintf("https://api.%s.com/repos/%s/%s", devops, organization, repos)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	//var repositories []Repository
	err = json.Unmarshal(body, &repo)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

// Get Infos for all Repositories in Organization for Main Branch
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

	/*	req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)


		err = json.Unmarshal(body, &repositories)
		if err != nil {
			return nil, err
		}*/

	return repositories, nil
}

func main() {

	var config1 Configuration
	var AppConfig = GetConfig(config1)
	var largestLineCounter int
	var nameRepos2 string
	var cpt = 0

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}
	DestinationResult := pwd + "/Results"
	if err := os.MkdirAll(DestinationResult, os.ModePerm); err != nil {
		panic(err)
	}

	GlobalReport := DestinationResult + "/GobalReport.txt"
	// Create Global Report File
	file, err := os.Create(GlobalReport)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if len(AppConfig.Repos) != 0 {
		repositories, err := getRepository(AppConfig.AccessToken, AppConfig.Organization, AppConfig.DevOps, AppConfig.Repos)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("\nAnalyse Repository %s in Organization: %s\n", repositories.Name, AppConfig.Organization)

		pathToScan := fmt.Sprintf("git::https://%s@%s.com/%s/%s", AppConfig.AccessToken, AppConfig.DevOps, AppConfig.Organization, AppConfig.Repos)
		outputFileName := fmt.Sprintf("Result_%s", repositories.Name)

		params := gcloc.Params{
			Path:              pathToScan,
			ByFile:            false,
			ExcludePaths:      []string{},
			ExcludeExtensions: []string{},
			IncludeExtensions: []string{},
			OrderByLang:       false,
			OrderByFile:       false,
			OrderByCode:       false,
			OrderByLine:       false,
			OrderByBlank:      false,
			OrderByComment:    false,
			Order:             "DESC",
			OutputName:        outputFileName,
			OutputPath:        DestinationResult,
			ReportFormats:     []string{"json"},
		}

		gc, err := gcloc.NewGCloc(params, constants.Languages)
		if err != nil {
			fmt.Println("Error:", err)
		}
		gc.Run()
		cpt++

	} else {
		repositories, err := getRepositoryList(AppConfig.AccessToken, AppConfig.Organization, AppConfig.DevOps)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("\nAnalyse Repositories in Organization: %s\n", AppConfig.Organization)
		for _, repo := range repositories {
			fmt.Printf("\nAnalyse Repository : %s\n", repo.Name)

			pathToScan := fmt.Sprintf("git::https://%s@%s.com/%s/%s", AppConfig.AccessToken, AppConfig.DevOps, AppConfig.Organization, repo.Name)
			outputFileName := fmt.Sprintf("Result_%s", repo.Name)

			params := gcloc.Params{
				Path:              pathToScan,
				ByFile:            false,
				ExcludePaths:      []string{},
				ExcludeExtensions: []string{},
				IncludeExtensions: []string{},
				OrderByLang:       false,
				OrderByFile:       false,
				OrderByCode:       false,
				OrderByLine:       false,
				OrderByBlank:      false,
				OrderByComment:    false,
				Order:             "DESC",
				OutputName:        outputFileName,
				OutputPath:        DestinationResult,
				ReportFormats:     []string{"json"},
			}

			gc, err := gcloc.NewGCloc(params, constants.Languages)
			if err != nil {
				fmt.Println("Error:", err)
			}

			gc.Run()
			cpt++

		}
	}

	spin := spinner.New(spinner.CharSets[35], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	spin.Suffix = " \nAnalyse Report..."
	spin.Start()

	// List files in the directory
	fileInfos, err := os.ReadDir(DestinationResult)
	if err != nil {
		fmt.Println("Error listing files:", err)
		return
	}

	// Loop through each file
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".json" {
			filePath := filepath.Join(DestinationResult, fileInfo.Name())
			nameRepos := strings.Split(fileInfo.Name(), "_")
			nameRepos1 := strings.Split(nameRepos[1], ".")
			TotalCodeLines := parseJSONFile(filePath, nameRepos1[0])
			if TotalCodeLines > largestLineCounter {
				largestLineCounter = TotalCodeLines
				nameRepos2 = nameRepos1[0]
			}
		}
	}
	spin.Stop()

	p := message.NewPrinter(language.English)
	s := strings.Replace(p.Sprintf("%d", largestLineCounter), ",", " ", -1)

	message0 := fmt.Sprintf("Number of Repository in Organization %s is %d \n", AppConfig.Organization, cpt)
	message1 := fmt.Sprintf("In Organization %s the largest number of line of code is %s and the repository is %s\n\nReports are located in the Results directory", AppConfig.Organization, s, nameRepos2)
	message2 := message0 + message1
	fmt.Println(message0)
	fmt.Println(message1)

	// Write message in Gobal Report File
	_, err = file.WriteString(message2)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
