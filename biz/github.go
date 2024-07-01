package biz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IgorEulalio/notificationservice/model"
)

func CreateRepository(repository model.Repository) error {
	if repository.Visibility == "PUBLIC" {
		return fmt.Errorf("user doesn't have access to create public repositories")
	}
	log.Println("Creating repository on GitHub...")

	// Load GitHub token from environment variable
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Create the repository on GitHub
	url := "https://api.github.com/user/repos"
	jsonRepoData, err := json.Marshal(repository)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonRepoData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create repository: %s", resp.Status)
	}

	log.Println("Repository created on GitHub successfully!")
	return nil
}
