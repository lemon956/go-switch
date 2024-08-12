package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/xulimeng/go-switch/models"
)

func ListAll() {
	resp, err := http.Get(models.GoVersionsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching Go versions: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: received non-200 response code: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var versions []models.GoVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %v\n", err)
		os.Exit(1)
	}
	for i := len(versions) - 1; i >= 0; i-- {
		fmt.Println(versions[i].Version)
	}
}
