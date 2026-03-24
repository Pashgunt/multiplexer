package config

import (
	"errors"
	"os"
	"testing"
	"transport/pkg/logging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()

	assert.NotNil(t, validator)
	assert.IsType(t, &Validator{}, validator)
}

func TestValidateFileExists_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte("config: test")
	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	validator := NewValidator()
	err = validator.ValidateFileExists(tmpFile.Name())

	assert.NoError(t, err)
}

func TestValidateFileExists_FileNotFound(t *testing.T) {
	validator := NewValidator()
	err := validator.ValidateFileExists("/non/existent/path/config.yaml")

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestValidateFileExists_IsDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_config_dir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	validator := NewValidator()
	err = validator.ValidateFileExists(tmpDir)

	assert.Error(t, err)

	var appErr *logging.AppError

	if errors.As(err, &appErr) {
		assert.Equal(t, "Config path is a directory, not a file.", appErr.Error())
	} else {
		t.Errorf("Expected AppError, got %T", err)
	}
}

func TestValidateFileExists_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_empty_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	err = tmpFile.Close()
	require.NoError(t, err)

	validator := NewValidator()
	err = validator.ValidateFileExists(tmpFile.Name())

	assert.Error(t, err)

	var appErr *logging.AppError

	if errors.As(err, &appErr) {
		assert.Equal(t, "Config file is empty.", appErr.Error())
	} else {
		t.Errorf("Expected AppError, got %T", err)
	}
}

func TestValidateFileExists_FileWithSpaces(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test config with spaces*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte("config: test")
	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	validator := NewValidator()
	err = validator.ValidateFileExists(tmpFile.Name())

	assert.NoError(t, err)
}

func TestValidateFileExists_FileWithDifferentExtensions(t *testing.T) {
	extensions := []string{".json", ".toml", ".yml", ".conf", ".cfg"}

	for _, ext := range extensions {
		t.Run("extension_"+ext, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_config*"+ext)
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			content := []byte("config: test")
			_, err = tmpFile.Write(content)
			require.NoError(t, err)
			err = tmpFile.Close()
			require.NoError(t, err)

			validator := NewValidator()
			err = validator.ValidateFileExists(tmpFile.Name())

			assert.NoError(t, err)
		})
	}
}

func TestValidateFileExists_RelativePath(t *testing.T) {
	tmpFile, err := os.CreateTemp(".", "test_relative_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte("config: test")
	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	baseName := tmpFile.Name()

	if baseName[0:2] == "./" {
		baseName = baseName[2:]
	}

	validator := NewValidator()
	err = validator.ValidateFileExists(baseName)

	assert.NoError(t, err)
}

func TestValidateFileExists_Symlink(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_real_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte("config: test")
	_, err = tmpFile.Write(content)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)

	symlinkPath := tmpFile.Name() + ".symlink"
	err = os.Symlink(tmpFile.Name(), symlinkPath)
	require.NoError(t, err)
	defer os.Remove(symlinkPath)

	validator := NewValidator()
	err = validator.ValidateFileExists(symlinkPath)

	assert.NoError(t, err)
}

func TestValidateFileExists_BrokenSymlink(t *testing.T) {
	symlinkPath := "/tmp/broken_symlink_test.yaml"
	err := os.Symlink("/non/existent/file.yaml", symlinkPath)

	if err != nil {
		t.Skip("Cannot create symlink, skipping test")
	}

	defer os.Remove(symlinkPath)

	validator := NewValidator()
	err = validator.ValidateFileExists(symlinkPath)

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestValidateFileExists_EmptyStringPath(t *testing.T) {
	validator := NewValidator()
	err := validator.ValidateFileExists("")

	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestValidateFileExists_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		cleanup       func(t *testing.T, path string)
		expectedError bool
		errorContains string
	}{
		{
			name: "valid file with content",
			setupFile: func(t *testing.T) string {
				f, err := os.CreateTemp("", "valid_*.yaml")
				require.NoError(t, err)
				_, err = f.Write([]byte("content"))
				require.NoError(t, err)
				f.Close()
				return f.Name()
			},
			cleanup:       func(t *testing.T, path string) { os.Remove(path) },
			expectedError: false,
		},
		{
			name: "empty file",
			setupFile: func(t *testing.T) string {
				f, err := os.CreateTemp("", "empty_*.yaml")
				require.NoError(t, err)
				f.Close()
				return f.Name()
			},
			cleanup:       func(t *testing.T, path string) { os.Remove(path) },
			expectedError: true,
			errorContains: "Config file is empty",
		},
		{
			name: "directory",
			setupFile: func(t *testing.T) string {
				dir, err := os.MkdirTemp("", "test_dir")
				require.NoError(t, err)
				return dir
			},
			cleanup:       func(t *testing.T, path string) { os.RemoveAll(path) },
			expectedError: true,
			errorContains: "Config path is a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFile(t)
			if tt.cleanup != nil {
				defer tt.cleanup(t, path)
			}

			validator := NewValidator()
			err := validator.ValidateFileExists(path)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func BenchmarkValidateFileExists(b *testing.B) {
	tmpFile, err := os.CreateTemp("", "bench_*.yaml")

	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("config: test"))

	if err != nil {
		b.Fatal(err)
	}

	tmpFile.Close()

	validator := NewValidator()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = validator.ValidateFileExists(tmpFile.Name())
	}
}
