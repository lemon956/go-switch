package features

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGolangCILintInstallArgsDefaultsToMajor(t *testing.T) {
	requested, err := parseGolangCILintInstallArgs([]string{"v2"})

	require.NoError(t, err)
	require.Equal(t, "v2", requested)
}

func TestParseGolangCILintInstallArgsUsesExactVersionOverride(t *testing.T) {
	requested, err := parseGolangCILintInstallArgs([]string{"v1", "v1.64.8"})

	require.NoError(t, err)
	require.Equal(t, "v1.64.8", requested)
}

func TestLintRejectsUnknownSubcommand(t *testing.T) {
	err := Lint([]string{"unknown"})

	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown lint command")
}

func TestLintRequiresSubcommand(t *testing.T) {
	err := Lint(nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Usage: goswitch lint")
}
