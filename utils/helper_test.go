package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnTarGz(t *testing.T) {
	targetPath := "/home/hellotalk/.go-switch/gos/go1.22.6.linux-amd64.tar.gz"
	destPath := "/home/hellotalk/.go-switch/gos/"
	err := UnTarGz(targetPath, destPath)
	require.Nil(t, err)
}

func TestRenameDir(t *testing.T) {
	srcPath := "/home/hellotalk/.go-switch/gos/go"
	err := RenameDir(srcPath, "go1.22.6")
	require.Nil(t, err)
}

func Test(t *testing.T) {
	SetPermissions("/home/hellotalk/.go-switch/gos/go1.22.6")
}

func TestTruncateFile(t *testing.T) {
	err := TruncateFile("/home/hellotalk/.go-switch/config/config.toml.bck")
	t.Log(err)
	require.Nil(t, err)
}
