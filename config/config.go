package config

import "os"

// RepositoryURL stores the repository URL of the project to be analyzed.
var RepositoryURL string

// HuskyAPI stores the address of Husky's API.
var HuskyAPI string

// SetConfigs sets all configuration needed to start the client.
func SetConfigs() {
	RepositoryURL = os.Getenv(`REPO_URL`)
	HuskyAPI = os.Getenv(`HUSKY_API`)
}
