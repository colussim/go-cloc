package utils

import (
	"bufio"
	"os"
	"strings"
)

func fileExists(fileexclusion string) bool {
	_, err := os.Stat(fileexclusion)
	return !os.IsNotExist(err)
}

func isFileEmpty(fileexclusion string) (bool, error) {
	fileInfo, err := os.Stat(fileexclusion)
	if err != nil {
		return false, err
	}

	return fileInfo.Size() == 0, nil
}

func searchStringInFile(fileexclusion string, target string) (bool, error) {
	file, err := os.Open(fileexclusion)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, target) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
