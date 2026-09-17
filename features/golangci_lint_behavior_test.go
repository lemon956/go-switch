package features

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/lemon956/go-switch/config"
	"github.com/stretchr/testify/require"
)

func TestUpsertGolangCILintInstallReplacesExistingVersion(t *testing.T) {
	config.Conf = &config.Config{
		LocalGolangCILints: []config.GolangCILintVersion{
			{Major: "v1", Version: "v1.64.7", Path: "/old", BinaryPath: "/old/golangci-lint"},
			{Major: "v1", Version: "v1.64.8", Path: "/stale", BinaryPath: "/stale/golangci-lint"},
		},
	}

	upsertGolangCILintInstall(config.GolangCILintVersion{
		Major:      "v1",
		Version:    "v1.64.8",
		Path:       "/new",
		BinaryPath: "/new/golangci-lint",
	})

	require.Len(t, config.Conf.LocalGolangCILints, 2)
	require.Equal(t, "/new", config.Conf.LocalGolangCILints[1].Path)
	require.Equal(t, "/new/golangci-lint", config.Conf.LocalGolangCILints[1].BinaryPath)
}

func TestFindInstalledGolangCILintChoosesLatestMajorVersion(t *testing.T) {
	config.Conf = &config.Config{
		LocalGolangCILints: []config.GolangCILintVersion{
			{Major: "v1", Version: "v1.63.4", BinaryPath: "/v1.63.4/golangci-lint"},
			{Major: "v2", Version: "v2.1.0", BinaryPath: "/v2.1.0/golangci-lint"},
			{Major: "v1", Version: "v1.64.8", BinaryPath: "/v1.64.8/golangci-lint"},
		},
	}

	installed, ok := findInstalledGolangCILint("v1")

	require.True(t, ok)
	require.Equal(t, "v1.64.8", installed.Version)
	require.Equal(t, "/v1.64.8/golangci-lint", installed.BinaryPath)
}

func TestCreateGolangCILintStableEntryPointsToBinary(t *testing.T) {
	root := t.TempDir()
	goPath := filepath.Join(root, "go")
	config.GoPathDirPath = filepath.Join(root, "default-go")
	config.Conf = &config.Config{GoPath: goPath}
	sourceDir := filepath.Join(root, "tools", "golangci-lint", "v2.2.0")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	sourceBinary := filepath.Join(sourceDir, "golangci-lint")
	require.NoError(t, os.WriteFile(sourceBinary, []byte("#!/bin/sh\n"), 0755))

	linkPath, err := createGolangCILintStableEntry(sourceBinary)

	require.NoError(t, err)
	require.Equal(t, filepath.Join(goPath, "bin", "golangci-lint"), linkPath)
	target, err := os.Readlink(linkPath)
	require.NoError(t, err)
	require.Equal(t, sourceBinary, target)
}

func TestCreateGolangCILintStableEntryFallsBackToDefaultGoPath(t *testing.T) {
	root := t.TempDir()
	config.GoPathDirPath = filepath.Join(root, "default-go")
	config.Conf = &config.Config{}
	sourceDir := filepath.Join(root, "tools", "golangci-lint", "v1.64.8")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	sourceBinary := filepath.Join(sourceDir, "golangci-lint")
	require.NoError(t, os.WriteFile(sourceBinary, []byte("#!/bin/sh\n"), 0755))

	linkPath, err := createGolangCILintStableEntry(sourceBinary)

	require.NoError(t, err)
	require.Equal(t, filepath.Join(config.GoPathDirPath, "bin", "golangci-lint"), linkPath)
}

func TestRenderGolangCILintListMarksCurrentVersion(t *testing.T) {
	config.Conf = &config.Config{
		LocalGolangCILints: []config.GolangCILintVersion{
			{Major: "v1", Version: "v1.64.8", Path: "/v1.64.8"},
			{Major: "v2", Version: "v2.2.0", Path: "/v2.2.0"},
		},
		CurrentGolangCILint: config.GolangCILintCurrent{
			Major:      "v2",
			Version:    "v2.2.0",
			BinaryPath: "/v2.2.0/golangci-lint",
		},
	}
	var out bytes.Buffer

	renderGolangCILintList(&out)

	require.Contains(t, out.String(), "v1.64.8")
	require.Contains(t, out.String(), "* v2.2.0")
}

func TestRenderGolangCILintEnvShowsCurrentBinaryAndGoPathBin(t *testing.T) {
	config.Conf = &config.Config{
		GoPath: "/tmp/go-switch/go",
		CurrentGolangCILint: config.GolangCILintCurrent{
			Major:      "v2",
			Version:    "v2.2.0",
			BinaryPath: "/tmp/go-switch/tools/golangci-lint/v2.2.0/golangci-lint",
		},
	}
	var out bytes.Buffer

	renderGolangCILintEnv(&out)

	require.Contains(t, out.String(), "Current golangci-lint version: v2.2.0")
	require.Contains(t, out.String(), "GOPATH bin path: /tmp/go-switch/go/bin")
	require.Contains(t, out.String(), "Stable entry path: /tmp/go-switch/go/bin/golangci-lint")
	require.Contains(t, out.String(), "Binary path: /tmp/go-switch/tools/golangci-lint/v2.2.0/golangci-lint")
}
