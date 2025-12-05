package helper

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnTarGz(t *testing.T) {
	t.Parallel()

	// Create a temporary directory
	tmpDir := t.TempDir()
	tarGzPath := filepath.Join(tmpDir, "test.tar.gz")
	extractDir := filepath.Join(tmpDir, "extract")

	// Create a simple tar.gz with one file
	err := createTestTarGz(tarGzPath, "testfile.txt", []byte("hello goswitch"))
	require.NoError(t, err)

	// Extract using UntarGz
	err = UntarGz(tarGzPath, extractDir)
	require.NoError(t, err)

	// Verify extracted file exists and content matches
	extractedFile := filepath.Join(extractDir, "testfile.txt")
	data, err := os.ReadFile(extractedFile)
	require.NoError(t, err)
	require.Equal(t, []byte("hello goswitch"), data)
}

func TestRenameDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "go")
	err := os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	err = RenameDir(srcDir, "go1.22.6")
	require.NoError(t, err)

	// New path should exist
	newPath := filepath.Join(tmpDir, "go1.22.6")
	info, err := os.Stat(newPath)
	require.NoError(t, err)
	require.True(t, info.IsDir())
}

func TestSetPermissions(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "dir")
	err := os.MkdirAll(testPath, 0755)
	require.NoError(t, err)

	if GlobalSetPermissions != nil {
		err := GlobalSetPermissions.SetPermissions(testPath)
		require.NoError(t, err)
	}
}

func TestTruncateFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Create file with some content
	err := os.WriteFile(filePath, []byte("some content"), 0644)
	require.NoError(t, err)

	err = TruncateFile(filePath)
	require.NoError(t, err)

	info, err := os.Stat(filePath)
	require.NoError(t, err)
	require.Equal(t, int64(0), info.Size())
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "test", "config.toml")

	exists, created := FileExists(target)
	require.False(t, exists)
	require.True(t, created)

	// Subsequent call should report exists
	exists, created = FileExists(target)
	require.True(t, exists)
	require.False(t, created)
}

// createTestTarGz creates a simple tar.gz file containing a single file with given content.
func createTestTarGz(path, filename string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	hdr := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := tw.Write(data); err != nil {
		return err
	}
	return nil
}
