package features

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lemon956/go-switch/config"
	"github.com/lemon956/go-switch/helper"
)

const GolangCILintReleasesURL = "https://api.github.com/repos/golangci/golangci-lint/releases?per_page=100"

type GolangCILintRelease struct {
	TagName    string              `json:"tag_name"`
	Draft      bool                `json:"draft"`
	Prerelease bool                `json:"prerelease"`
	Assets     []GolangCILintAsset `json:"assets"`
}

type GolangCILintAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GolangCILintSelection struct {
	Major   string
	Version string
	Asset   GolangCILintAsset
}

type golangCILintSemver struct {
	major int
	minor int
	patch int
}

func Lint(args []string) error {
	if len(args) == 0 {
		return errors.New(lintUsage())
	}

	switch args[0] {
	case "help", "-h", "--help":
		fmt.Print(lintUsage())
		return nil
	case "install":
		requested, err := parseGolangCILintInstallArgs(args[1:])
		if err != nil {
			return err
		}
		return InstallGolangCILint(requested)
	case "switch":
		major, err := parseGolangCILintMajorArgs(args[1:], "switch")
		if err != nil {
			return err
		}
		return SwitchGolangCILint(major)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("lint list does not accept extra arguments")
		}
		ListGolangCILints()
		return nil
	case "env":
		if len(args) != 1 {
			return fmt.Errorf("lint env does not accept extra arguments")
		}
		EnvGolangCILint()
		return nil
	default:
		return fmt.Errorf("unknown lint command %q\n%s", args[0], lintUsage())
	}
}

func parseGolangCILintInstallArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("lint install requires v1 or v2")
	}
	if len(args) > 2 {
		return "", fmt.Errorf("Usage: goswitch lint install v1|v2 [version]")
	}

	_, major, err := normalizeGolangCILintRequest(args[0])
	if err != nil {
		return "", err
	}
	if args[0] != "v1" && args[0] != "1" && args[0] != "v2" && args[0] != "2" {
		return "", fmt.Errorf("lint install first argument must be v1 or v2")
	}
	if len(args) == 1 {
		return major, nil
	}

	exactVersion, exactMajor, err := normalizeGolangCILintRequest(args[1])
	if err != nil {
		return "", err
	}
	if exactVersion == "" {
		return "", fmt.Errorf("lint install version override must be a full version like v1.64.8")
	}
	if exactMajor != major {
		return "", fmt.Errorf("lint install major %s does not match version %s", major, exactVersion)
	}
	return exactVersion, nil
}

func parseGolangCILintMajorArgs(args []string, command string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Usage: goswitch lint %s v1|v2", command)
	}
	exactVersion, major, err := normalizeGolangCILintRequest(args[0])
	if err != nil {
		return "", err
	}
	if exactVersion != "" {
		return "", fmt.Errorf("lint %s accepts only v1 or v2", command)
	}
	return major, nil
}

func lintUsage() string {
	return `Usage: goswitch lint <command>

Command:
	help	- Show golangci-lint help message
	install	- Install latest golangci-lint v1/v2, optionally with a full version
	switch	- Choose current golangci-lint major version
	list	- List installed golangci-lint versions
	env	- Show golangci-lint environment
`
}

func NormalizeGolangCILintPlatform(system config.Env, arch string) (string, string) {
	osName := string(system)
	if system == config.Mac || osName == "drawin" {
		osName = "darwin"
	}

	switch arch {
	case "x86_64":
		arch = "amd64"
	case "aarch64":
		arch = "arm64"
	}

	return osName, arch
}

func ResolveGolangCILintSelection(releases []GolangCILintRelease, requested, system, arch string) (GolangCILintSelection, error) {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return GolangCILintSelection{}, fmt.Errorf("golangci-lint version is required")
	}

	systemEnv := config.Env(system)
	osName, normalizedArch := NormalizeGolangCILintPlatform(systemEnv, arch)
	exactVersion, requestedMajor, err := normalizeGolangCILintRequest(requested)
	if err != nil {
		return GolangCILintSelection{}, err
	}

	candidates := make([]GolangCILintRelease, 0, len(releases))
	for _, release := range releases {
		if release.Draft || release.Prerelease {
			continue
		}
		version, ok := parseGolangCILintSemver(release.TagName)
		if !ok {
			continue
		}
		major := fmt.Sprintf("v%d", version.major)
		if exactVersion != "" {
			if release.TagName != exactVersion {
				continue
			}
		} else if major != requestedMajor {
			continue
		}
		if _, ok := findGolangCILintAsset(release, osName, normalizedArch); !ok {
			continue
		}
		candidates = append(candidates, release)
	}

	if len(candidates) == 0 {
		return GolangCILintSelection{}, fmt.Errorf("no golangci-lint release found for %s on %s/%s", requested, osName, normalizedArch)
	}

	sort.Slice(candidates, func(i, j int) bool {
		left, _ := parseGolangCILintSemver(candidates[i].TagName)
		right, _ := parseGolangCILintSemver(candidates[j].TagName)
		return compareGolangCILintSemver(left, right) > 0
	})

	selected := candidates[0]
	asset, _ := findGolangCILintAsset(selected, osName, normalizedArch)
	selectedVersion, _ := parseGolangCILintSemver(selected.TagName)

	return GolangCILintSelection{
		Major:   fmt.Sprintf("v%d", selectedVersion.major),
		Version: selected.TagName,
		Asset:   asset,
	}, nil
}

func normalizeGolangCILintRequest(requested string) (exactVersion, major string, err error) {
	if !strings.HasPrefix(requested, "v") {
		requested = "v" + requested
	}

	if requested == "v1" || requested == "v2" {
		return "", requested, nil
	}

	version, ok := parseGolangCILintSemver(requested)
	if !ok {
		return "", "", fmt.Errorf("invalid golangci-lint version %q, use v1, v2, or a full version like v1.64.8", requested)
	}
	return requested, fmt.Sprintf("v%d", version.major), nil
}

func parseGolangCILintSemver(version string) (golangCILintSemver, bool) {
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return golangCILintSemver{}, false
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return golangCILintSemver{}, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return golangCILintSemver{}, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return golangCILintSemver{}, false
	}

	return golangCILintSemver{major: major, minor: minor, patch: patch}, true
}

func compareGolangCILintSemver(left, right golangCILintSemver) int {
	if left.major != right.major {
		return left.major - right.major
	}
	if left.minor != right.minor {
		return left.minor - right.minor
	}
	return left.patch - right.patch
}

func findGolangCILintAsset(release GolangCILintRelease, osName, arch string) (GolangCILintAsset, bool) {
	version := strings.TrimPrefix(release.TagName, "v")
	requiredSegment := fmt.Sprintf("-%s-%s", osName, arch)
	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, "golangci-lint-"+version+"-") &&
			strings.Contains(asset.Name, requiredSegment) &&
			(strings.HasSuffix(asset.Name, ".tar.gz") || strings.HasSuffix(asset.Name, ".zip")) {
			return asset, true
		}
	}
	return GolangCILintAsset{}, false
}

func InstallGolangCILint(requested string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	releases, err := FetchGolangCILintReleases(client)
	if err != nil {
		return err
	}

	selection, err := ResolveGolangCILintSelection(releases, requested, string(config.SystemEnv), config.SystemArch)
	if err != nil {
		return err
	}

	for _, installed := range config.Conf.LocalGolangCILints {
		if installed.Version == selection.Version {
			upsertGolangCILintInstall(installed)
			config.Conf.SaveConfig()
			fmt.Printf("golangci-lint %s is already installed\n", selection.Version)
			return nil
		}
	}

	installed, err := installGolangCILintSelection(selection)
	if err != nil {
		return err
	}
	upsertGolangCILintInstall(installed)
	config.Conf.SaveConfig()

	fmt.Printf("golangci-lint %s has been installed successfully\n", installed.Version)
	fmt.Println("Use 'goswitch lint switch' to switch to this major version")
	return nil
}

func FetchGolangCILintReleases(client *http.Client) ([]GolangCILintRelease, error) {
	req, err := http.NewRequest(http.MethodGet, GolangCILintReleasesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "go-switch")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch golangci-lint releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch golangci-lint releases: %s", resp.Status)
	}

	var releases []GolangCILintRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode golangci-lint releases: %w", err)
	}
	return releases, nil
}

func installGolangCILintSelection(selection GolangCILintSelection) (config.GolangCILintVersion, error) {
	if selection.Asset.BrowserDownloadURL == "" {
		return config.GolangCILintVersion{}, fmt.Errorf("golangci-lint %s has no download URL", selection.Version)
	}
	if err := os.MkdirAll(config.GolangCILintsPath, 0755); err != nil {
		return config.GolangCILintVersion{}, err
	}

	archivePath := filepath.Join(config.GolangCILintsPath, selection.Asset.Name)
	extractPath := filepath.Join(config.GolangCILintsPath, ".extract-"+selection.Version)
	installPath := filepath.Join(config.GolangCILintsPath, selection.Version)

	if err := os.RemoveAll(extractPath); err != nil {
		return config.GolangCILintVersion{}, err
	}
	if err := os.RemoveAll(installPath); err != nil {
		return config.GolangCILintVersion{}, err
	}
	if err := os.MkdirAll(extractPath, 0755); err != nil {
		return config.GolangCILintVersion{}, err
	}
	defer os.RemoveAll(extractPath)
	defer os.Remove(archivePath)

	if err := helper.DownloadFile(selection.Asset.BrowserDownloadURL, archivePath); err != nil {
		return config.GolangCILintVersion{}, fmt.Errorf("failed to download golangci-lint %s: %w", selection.Version, err)
	}
	if err := helper.Decompress(archivePath, extractPath); err != nil {
		return config.GolangCILintVersion{}, fmt.Errorf("failed to decompress golangci-lint %s: %w", selection.Version, err)
	}

	extractedBinary, err := findGolangCILintBinary(extractPath)
	if err != nil {
		return config.GolangCILintVersion{}, err
	}
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return config.GolangCILintVersion{}, err
	}

	binaryPath := filepath.Join(installPath, filepath.Base(extractedBinary))
	if err := copyFile(extractedBinary, binaryPath); err != nil {
		return config.GolangCILintVersion{}, err
	}
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return config.GolangCILintVersion{}, err
	}
	if helper.GlobalSetPermissions != nil {
		if err := helper.GlobalSetPermissions.SetPermissions(installPath); err != nil {
			return config.GolangCILintVersion{}, err
		}
	}

	return config.GolangCILintVersion{
		Major:      selection.Major,
		Version:    selection.Version,
		Path:       installPath,
		BinaryPath: binaryPath,
	}, nil
}

func SwitchGolangCILint(requestedMajor string) error {
	_, major, err := normalizeGolangCILintRequest(requestedMajor)
	if err != nil {
		return err
	}
	installed, ok := findInstalledGolangCILint(major)
	if !ok {
		return fmt.Errorf("golangci-lint %s is not installed, please run 'goswitch lint install %s' first", major, major)
	}
	if _, err := os.Stat(installed.BinaryPath); err != nil {
		return fmt.Errorf("golangci-lint binary does not exist: %s", installed.BinaryPath)
	}

	linkPath, err := createGolangCILintStableEntry(installed.BinaryPath)
	if err != nil {
		return err
	}

	config.Conf.CurrentGolangCILint = config.GolangCILintCurrent{
		Major:      installed.Major,
		Version:    installed.Version,
		BinaryPath: installed.BinaryPath,
	}
	config.Conf.SaveConfig()

	fmt.Printf("Switched to golangci-lint %s (%s) successfully\n", installed.Version, installed.Major)
	fmt.Printf("Current golangci-lint entry: %s\n", linkPath)
	return nil
}

func ListGolangCILints() {
	renderGolangCILintList(os.Stdout)
}

func EnvGolangCILint() {
	renderGolangCILintEnv(os.Stdout)
}

func upsertGolangCILintInstall(installed config.GolangCILintVersion) {
	if config.Conf == nil {
		config.Conf = &config.Config{}
	}
	for idx, existing := range config.Conf.LocalGolangCILints {
		if existing.Version == installed.Version {
			config.Conf.LocalGolangCILints[idx] = installed
			return
		}
	}
	config.Conf.LocalGolangCILints = append(config.Conf.LocalGolangCILints, installed)
}

func findInstalledGolangCILint(major string) (config.GolangCILintVersion, bool) {
	if config.Conf == nil {
		return config.GolangCILintVersion{}, false
	}

	var selected config.GolangCILintVersion
	selectedVersion := golangCILintSemver{}
	found := false
	for _, installed := range config.Conf.LocalGolangCILints {
		if installed.Major != major {
			continue
		}
		parsed, ok := parseGolangCILintSemver(installed.Version)
		if !ok {
			continue
		}
		if !found || compareGolangCILintSemver(parsed, selectedVersion) > 0 {
			selected = installed
			selectedVersion = parsed
			found = true
		}
	}
	return selected, found
}

func createGolangCILintStableEntry(binaryPath string) (string, error) {
	goPathBin := golangCILintGoPathBin()
	if goPathBin == "" {
		return "", fmt.Errorf("GOPATH is not configured, please run 'goswitch init' first")
	}
	if err := os.MkdirAll(goPathBin, 0755); err != nil {
		return "", err
	}

	entryName := config.GolangCILintBinary
	if strings.HasSuffix(binaryPath, ".exe") {
		entryName += ".exe"
	}
	linkPath := filepath.Join(goPathBin, entryName)
	if _, err := os.Lstat(linkPath); err == nil {
		if err := os.Remove(linkPath); err != nil {
			return "", fmt.Errorf("failed to remove existing golangci-lint entry: %w", err)
		}
	}
	if err := os.Symlink(binaryPath, linkPath); err != nil {
		return "", fmt.Errorf("failed to create golangci-lint entry: %w", err)
	}
	return linkPath, nil
}

func golangCILintGoPathBin() string {
	goPath := ""
	if config.Conf != nil {
		goPath = config.Conf.GoPath
	}
	if goPath == "" {
		goPath = config.GoPathDirPath
	}
	if goPath == "" {
		return ""
	}
	return filepath.Join(goPath, "bin")
}

func golangCILintStableEntryPath() string {
	goPathBin := golangCILintGoPathBin()
	if goPathBin == "" {
		return ""
	}
	entryName := config.GolangCILintBinary
	if config.SystemEnv == config.Windows {
		entryName += ".exe"
	}
	return filepath.Join(goPathBin, entryName)
}

func renderGolangCILintList(w io.Writer) {
	if config.Conf == nil || len(config.Conf.LocalGolangCILints) == 0 {
		fmt.Fprintln(w, "No golangci-lint versions installed")
		return
	}

	installed := append([]config.GolangCILintVersion(nil), config.Conf.LocalGolangCILints...)
	sort.Slice(installed, func(i, j int) bool {
		if installed[i].Major != installed[j].Major {
			return installed[i].Major < installed[j].Major
		}
		left, leftOK := parseGolangCILintSemver(installed[i].Version)
		right, rightOK := parseGolangCILintSemver(installed[j].Version)
		if leftOK && rightOK {
			return compareGolangCILintSemver(left, right) < 0
		}
		return installed[i].Version < installed[j].Version
	})

	for _, item := range installed {
		prefix := "  "
		if config.Conf.CurrentGolangCILint.Version == item.Version {
			prefix = "* "
		}
		fmt.Fprintf(w, "%s%s (%s) %s\n", prefix, item.Version, item.Major, item.Path)
	}
}

func renderGolangCILintEnv(w io.Writer) {
	fmt.Fprintln(w, "=== golangci-lint Environment ===")
	fmt.Fprintf(w, "GOPATH bin path: %s\n", golangCILintGoPathBin())
	fmt.Fprintf(w, "Stable entry path: %s\n", golangCILintStableEntryPath())

	if config.Conf == nil || config.Conf.CurrentGolangCILint.Version == "" {
		fmt.Fprintln(w, "Current golangci-lint version: not set")
		fmt.Fprintln(w, "Hint: use 'goswitch lint switch v1' or 'goswitch lint switch v2'")
		return
	}

	fmt.Fprintf(w, "Current golangci-lint version: %s\n", config.Conf.CurrentGolangCILint.Version)
	fmt.Fprintf(w, "Current golangci-lint major: %s\n", config.Conf.CurrentGolangCILint.Major)
	fmt.Fprintf(w, "Binary path: %s\n", config.Conf.CurrentGolangCILint.BinaryPath)
}

func findGolangCILintBinary(root string) (string, error) {
	var result string
	errFound := errors.New("golangci-lint binary found")
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if name == config.GolangCILintBinary || name == config.GolangCILintBinary+".exe" {
			result = path
			return errFound
		}
		return nil
	})
	if errors.Is(err, errFound) {
		return result, nil
	}
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", fmt.Errorf("golangci-lint binary not found in archive")
	}
	return result, nil
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
