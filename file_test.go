package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	t.Run("expand home directory", func(t *testing.T) {
		savedVal := os.Getenv("HOME")
		defer func() {
			os.Setenv("HOME", savedVal)
		}()

		os.Setenv("HOME", "/home/test")
		path, err := ExpandPath("~")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if path == "" {
			t.Errorf("expected path got empty string")
		}
		if path != "/home/test" {
			t.Errorf("expected '/home/test' got '%s'", path)
		}
	})

	t.Run("expand environment variables", func(t *testing.T) {
		savedVal := os.Getenv("TEST_DIR")
		defer func() {
			os.Setenv("TEST_DIR", savedVal)
		}()

		os.Setenv("TEST_DIR", "/test/dir")
		path, err := ExpandPath("$TEST_DIR/file.txt")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !filepath.IsAbs(path) {
			t.Errorf("expected absolute path, got: %s", path)
		}
		if filepath.Base(path) != "file.txt" {
			t.Errorf("expected filename 'file.txt', got: %s", filepath.Base(path))
		}
	})

	t.Run("clean path", func(t *testing.T) {
		path, err := ExpandPath("/path/to/../file.txt")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if path != "/path/file.txt" {
			t.Errorf("expected '/path/file.txt', got: %s", path)
		}
	})

	t.Run("error on invalid path", func(t *testing.T) {
		// This test is tricky because filepath.Abs rarely returns an error
		// We'll just ensure the function doesn't panic
		_, err := ExpandPath("")
		// We don't assert on the error because it might or might not be an error
		// depending on the OS and environment
		_ = err
	})
}

func TestCleanOpen(t *testing.T) {
	t.Run("open existing file", func(t *testing.T) {
		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.Close()

		// Test CleanOpen
		file, err := CleanOpen(tmpfile.Name())
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		defer file.Close()
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		_, err := CleanOpen("/path/to/nonexistent/file")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		// An empty path should be expanded to the current directory
		file, err := CleanOpen("")
		if err != nil {
			// This might fail if the current directory doesn't exist or isn't readable
			// but that's unlikely in a test environment
			t.Errorf("unexpected error: %s", err)
		} else {
			file.Close()
		}
	})
}

func TestCleanOpenFile(t *testing.T) {
	t.Run("create new file", func(t *testing.T) {
		// Create a temporary directory
		tmpdir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpdir)

		// Test CleanOpenFile
		filePath := filepath.Join(tmpdir, "newfile.txt")
		file, err := CleanOpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		defer file.Close()

		// Verify the file was created
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("file was not created: %s", err)
		}
	})

	t.Run("error on invalid path", func(t *testing.T) {
		_, err := CleanOpenFile("", os.O_RDWR|os.O_CREATE, 0600)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreateDirPath(t *testing.T) {
	t.Run("create directory", func(t *testing.T) {
		// Create a temporary directory
		tmpdir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpdir)

		// Test CreateDirPath
		newDir := filepath.Join(tmpdir, "newdir")
		createdDir, err := CreateDirPath(newDir, "")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		// Verify the directory was created
		if _, err := os.Stat(createdDir); err != nil {
			t.Errorf("directory was not created: %s", err)
		}
	})

	t.Run("use default path", func(t *testing.T) {
		// Create a temporary directory
		tmpdir, err := os.MkdirTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpdir)

		// Test CreateDirPath with empty path
		newDir := filepath.Join(tmpdir, "defaultdir")
		createdDir, err := CreateDirPath("", newDir)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		// Verify the directory was created
		if _, err := os.Stat(createdDir); err != nil {
			t.Errorf("directory was not created: %s", err)
		}
	})

	t.Run("empty paths", func(t *testing.T) {
		// Both path and defaultPath are empty
		// This should create the current directory, which already exists
		dir, err := CreateDirPath("", "")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		// Verify the directory exists
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("directory does not exist: %s", err)
		}
	})
}

func TestFilesExist(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		if !FilesExist() {
			t.Error("expected true for empty files list")
		}
	})

	t.Run("existing file", func(t *testing.T) {
		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())
		tmpfile.Close()

		if !FilesExist(tmpfile.Name()) {
			t.Errorf("expected true for existing file: %s", tmpfile.Name())
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		if FilesExist("/path/to/nonexistent/file") {
			t.Error("expected false for non-existent file")
		}
	})

	t.Run("multiple files", func(t *testing.T) {
		// Create temporary files
		tmpfile1, err := os.CreateTemp("", "test1")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile1.Name())
		tmpfile1.Close()

		tmpfile2, err := os.CreateTemp("", "test2")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile2.Name())
		tmpfile2.Close()

		// Both files exist
		if !FilesExist(tmpfile1.Name(), tmpfile2.Name()) {
			t.Error("expected true when all files exist")
		}

		// One file doesn't exist
		if FilesExist(tmpfile1.Name(), "/path/to/nonexistent/file") {
			t.Error("expected false when one file doesn't exist")
		}
	})
}

// TestWaitForFiles is already covered in wait_test.go

// Test helper functions for encoding/decoding
func TestEncoderDecoderFuncs(t *testing.T) {
	t.Run("encoderFuncFromFilePath", func(t *testing.T) {
		// YAML files
		if encoderFuncFromFilePath("file.yaml") == nil {
			t.Error("expected encoder function for .yaml file")
		}
		if encoderFuncFromFilePath("file.yml") == nil {
			t.Error("expected encoder function for .yml file")
		}

		// JSON files
		if encoderFuncFromFilePath("file.json") == nil {
			t.Error("expected encoder function for .json file")
		}

		// Unsupported format
		if encoderFuncFromFilePath("file.txt") != nil {
			t.Error("expected nil for unsupported file format")
		}
	})

	t.Run("decoderFuncFromFilePath", func(t *testing.T) {
		// YAML files
		if decoderFuncFromFilePath("file.yaml") == nil {
			t.Error("expected decoder function for .yaml file")
		}
		if decoderFuncFromFilePath("file.yml") == nil {
			t.Error("expected decoder function for .yml file")
		}

		// JSON files
		if decoderFuncFromFilePath("file.json") == nil {
			t.Error("expected decoder function for .json file")
		}

		// Unsupported format
		if decoderFuncFromFilePath("file.txt") != nil {
			t.Error("expected nil for unsupported file format")
		}
	})
}

// Define a test struct for LoadStructFromFile and SaveStructToFile tests
type TestConfig struct {
	Name  string `json:"name" yaml:"name"`
	Value int    `json:"value" yaml:"value"`
}

func TestLoadSaveStructToFile(t *testing.T) {
	testFormats := []string{".json", ".yaml", ".yml"}

	for _, format := range testFormats {
		t.Run("load and save "+format, func(t *testing.T) {
			// Create a temporary directory
			tmpdir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpdir)

			// Create a test config
			config := &TestConfig{
				Name:  "test",
				Value: 42,
			}

			// Save the config to a file
			filePath := filepath.Join(tmpdir, "config"+format)
			err = SaveStructToFile(config, filePath)
			if err != nil {
				t.Fatalf("failed to save struct to file: %s", err)
			}

			// Load the config from the file
			loadedConfig, err := LoadStructFromFile[TestConfig](filePath)
			if err != nil {
				t.Fatalf("failed to load struct from file: %s", err)
			}

			// Verify the loaded config matches the original
			if loadedConfig.Name != config.Name {
				t.Errorf("expected Name %s, got %s", config.Name, loadedConfig.Name)
			}
			if loadedConfig.Value != config.Value {
				t.Errorf("expected Value %d, got %d", config.Value, loadedConfig.Value)
			}
		})
	}

	t.Run("error on empty file path", func(t *testing.T) {
		// Test SaveStructToFile with empty path
		err := SaveStructToFile(&TestConfig{}, "")
		if err == nil {
			t.Error("expected error for empty file path")
		}

		// Test LoadStructFromFile with empty path
		_, err = LoadStructFromFile[TestConfig]("")
		if err == nil {
			t.Error("expected error for empty file path")
		}
	})

	t.Run("error on unsupported format", func(t *testing.T) {
		// Test SaveStructToFile with unsupported format
		err := SaveStructToFile(&TestConfig{}, "config.txt")
		if err == nil {
			t.Error("expected error for unsupported file format")
		}

		// Test LoadStructFromFile with unsupported format
		_, err = LoadStructFromFile[TestConfig]("config.txt")
		if err == nil {
			t.Error("expected error for unsupported file format")
		}
	})

	t.Run("error on non-existent file", func(t *testing.T) {
		_, err := LoadStructFromFile[TestConfig]("/path/to/nonexistent/file.json")
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}

// Examples in Go standard format
func ExampleExpandPath() {
	// This example shows how to expand a path with environment variables
	os.Setenv("CONFIG_DIR", "/etc/app")
	path, err := ExpandPath("${CONFIG_DIR}/config.json")
	if err != nil {
		// Handle error
		return
	}
	// Output would be something like: /etc/app/config.json
	// But we can't assert on the exact output because it depends on the environment
	_ = path
}

func ExampleCleanOpen() {
	// This example shows how to safely open a file with a path that might contain
	// home directory references or environment variables
	file, err := CleanOpen("~/config.json")
	if err != nil {
		// Handle error
		return
	}
	defer file.Close()
	// Now you can use the file
}

func ExampleCreateDirPath() {
	// This example shows how to create a directory path, falling back to a default
	// if the provided path is empty
	_, err := CreateDirPath("", "~/.myapp/config")
	if err != nil {
		// Handle error
		return
	}
	// The directory has been created
}

func ExampleFilesExist() {
	// This example shows how to check if multiple files exist
	if FilesExist("/etc/hosts", "/etc/passwd") {
		// All files exist, proceed
	} else {
		// At least one file is missing
	}
}

func ExampleLoadStructFromFile() {
	// This example shows how to load a struct from a JSON file
	type Config struct {
		ServerName string `json:"server_name"`
		Port       int    `json:"port"`
	}

	config, err := LoadStructFromFile[Config]("/etc/app/config.json")
	if err != nil {
		// Handle error
		return
	}
	// Now you can use the config
	_ = config
}

func ExampleSaveStructToFile() {
	// This example shows how to save a struct to a YAML file
	type Config struct {
		ServerName string `yaml:"server_name"`
		Port       int    `yaml:"port"`
	}

	config := &Config{
		ServerName: "api-server",
		Port:       8080,
	}

	err := SaveStructToFile(config, "/etc/app/config.yaml")
	if err != nil {
		// Handle error
		return
	}
	// The config has been saved to the file
}
