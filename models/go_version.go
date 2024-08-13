package models

const (
	GoBaseURL     = "https://golang.org/dl/"
	GoVersionsURL = "https://golang.org/dl/?mode=json&include=all"
)

type GoVersion struct {
	Version string `json:"version"`
	Files   []File `json:"files"`
}

type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Sha256   string `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}
