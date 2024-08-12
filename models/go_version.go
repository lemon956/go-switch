package models

const GoVersionsURL = "https://golang.org/dl/?mode=json&include=all"

type GoVersion struct {
	Version string `json:"version"`
	Files   []struct {
		Filename string `json:"filename"`
		OS       string `json:"os"`
		Arch     string `json:"arch"`
		Download string `json:"download"`
	} `json:"files"`
}
