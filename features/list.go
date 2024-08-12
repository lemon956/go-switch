package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const goVersionsURL = "https://golang.org/dl/?mode=json&include=all"

type GoVersion struct {
	Version string `json:"version"`
}

func ListAll() {
	resp, err := http.Get(goVersionsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching Go versions: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: received non-200 response code: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var versions []GoVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %v\n", err)
		os.Exit(1)
	}

	for _, version := range versions {
		fmt.Println(version.Version)
	}
}
