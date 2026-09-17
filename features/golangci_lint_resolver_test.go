package features

import (
	"testing"

	"github.com/lemon956/go-switch/config"
	"github.com/stretchr/testify/require"
)

func TestResolveGolangCILintSelectionChoosesLatestStableMajor(t *testing.T) {
	releases := []GolangCILintRelease{
		{
			TagName: "v2.1.0",
			Assets: []GolangCILintAsset{{
				Name:               "golangci-lint-2.1.0-linux-amd64.tar.gz",
				BrowserDownloadURL: "https://example.test/v2.1.0.tar.gz",
			}},
		},
		{
			TagName: "v2.2.0",
			Assets: []GolangCILintAsset{{
				Name:               "golangci-lint-2.2.0-linux-amd64.tar.gz",
				BrowserDownloadURL: "https://example.test/v2.2.0.tar.gz",
			}},
		},
		{
			TagName: "v2.3.0",
			Draft:   true,
			Assets: []GolangCILintAsset{{
				Name:               "golangci-lint-2.3.0-linux-amd64.tar.gz",
				BrowserDownloadURL: "https://example.test/v2.3.0.tar.gz",
			}},
		},
		{
			TagName:    "v2.4.0",
			Prerelease: true,
			Assets: []GolangCILintAsset{{
				Name:               "golangci-lint-2.4.0-linux-amd64.tar.gz",
				BrowserDownloadURL: "https://example.test/v2.4.0.tar.gz",
			}},
		},
	}

	selection, err := ResolveGolangCILintSelection(releases, "v2", "linux", "amd64")
	require.NoError(t, err)
	require.Equal(t, "v2", selection.Major)
	require.Equal(t, "v2.2.0", selection.Version)
	require.Equal(t, "golangci-lint-2.2.0-linux-amd64.tar.gz", selection.Asset.Name)
}

func TestResolveGolangCILintSelectionMatchesExactVersion(t *testing.T) {
	releases := []GolangCILintRelease{
		{
			TagName: "v1.64.8",
			Assets: []GolangCILintAsset{{
				Name:               "golangci-lint-1.64.8-darwin-arm64.tar.gz",
				BrowserDownloadURL: "https://example.test/v1.64.8.tar.gz",
			}},
		},
		{
			TagName: "v1.63.4",
			Assets: []GolangCILintAsset{{
				Name:               "golangci-lint-1.63.4-darwin-arm64.tar.gz",
				BrowserDownloadURL: "https://example.test/v1.63.4.tar.gz",
			}},
		},
	}

	selection, err := ResolveGolangCILintSelection(releases, "v1.63.4", "darwin", "arm64")
	require.NoError(t, err)
	require.Equal(t, "v1", selection.Major)
	require.Equal(t, "v1.63.4", selection.Version)
	require.Equal(t, "https://example.test/v1.63.4.tar.gz", selection.Asset.BrowserDownloadURL)
}

func TestNormalizeGolangCILintPlatformUsesDarwinForMacTypo(t *testing.T) {
	osName, arch := NormalizeGolangCILintPlatform(config.Mac, "x86_64")
	require.Equal(t, "darwin", osName)
	require.Equal(t, "amd64", arch)
}
