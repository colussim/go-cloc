package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/wcharczuk/go-chart/v2"
)

type Record struct {
	Language  string `json:"Language"`
	CodeLines int    `json:"CodeLines"`
}

type JSONData struct {
	Results []Record `json:"Results"`
}

func main() {
	// Specify the directory where JSON files are located
	directory := "/Users/manu/Documents/App/Dev/Results"

	// Initialize a map to store the total code lines by language
	totalCodeLinesByLanguage := make(map[string]int)

	// List all JSON files in the directory
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Read JSON file
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Parse JSON data
			var jsonData JSONData
			if err := json.Unmarshal(data, &jsonData); err != nil {
				return err
			}

			// Update the total code lines by language
			for _, record := range jsonData.Results {
				totalCodeLinesByLanguage[record.Language] += record.CodeLines
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create a new doughnut chart

	pie := chart.DonutChart{
		Width:  612,
		Height: 612,
	}

	// Populate the chart with data
	var chartData []chart.Value

	// Print the total code lines by language
	fmt.Println("Total Code Lines by Language:")
	for language, codeLines := range totalCodeLinesByLanguage {
		fmt.Printf("%s: %s\n", language, humanize.Commaf(float64(codeLines)))
		chartData = append(chartData, chart.Value{
			Label: language,
			Value: float64(codeLines),
		})
	}
	pie.Values = chartData

	// Save the chart to a file or display it
	file, _ := os.Create("doughnut_chart.png")
	defer file.Close()
	pie.Render(chart.PNG, file)

	fmt.Println("Doughnut chart saved as doughnut_chart.png")
}
