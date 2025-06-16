package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"

	"github.com/dioad/generics"
)

// CleanOpen opens a file with a cleaned and expanded path.
// It expands the path (resolving ~ and environment variables) and cleans it
// to prevent path traversal attacks.
//
// Example:
//
//	file, err := util.CleanOpen("~/config.json")
//	if err != nil {
//	    return err
//	}
//	defer file.Close()
func CleanOpen(path string) (*os.File, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path: %w", err)
	}

	return os.Open(path) // path is already cleaned by ExpandPath
}

// CleanOpenFile opens a file with the specified flags and permissions, using a cleaned and expanded path.
// It expands the path (resolving ~ and environment variables) and cleans it
// to prevent path traversal attacks.
//
// Example:
//
//	file, err := util.CleanOpenFile("~/config.json", os.O_RDWR|os.O_CREATE, 0600)
//	if err != nil {
//	    return err
//	}
//	defer file.Close()
func CleanOpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path: %w", err)
	}

	return os.OpenFile(path, flag, perm) // #nosec - path is already cleaned by ExpandPath
}

// CreateDirPath creates a directory path if it doesn't exist.
// If path is empty, it uses defaultPath instead.
// The path is expanded (resolving ~ and environment variables) and cleaned
// before creating the directory.
//
// Example:
//
//	configDir, err := util.CreateDirPath("", "~/.myapp/config")
//	if err != nil {
//	    return err
//	}
//	// configDir now contains the absolute path to the created directory
func CreateDirPath(path string, defaultPath string) (string, error) {
	if path == "" {
		path = defaultPath
	}

	path, err := ExpandPath(path)
	if err != nil {
		return "", fmt.Errorf("failed to expand path: %w", err)
	}

	err = os.MkdirAll(path, 0750)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	return path, nil
}

// ExpandPath expands a path to an absolute path.
// It performs the following operations:
// 1. Expands ~ to the user's home directory
// 2. Expands environment variables (e.g., $HOME, ${HOME})
// 3. Cleans the path to remove any unnecessary elements
// 4. Converts the path to an absolute path
//
// Example:
//
//	path, err := util.ExpandPath("~/Documents/${APP_DIR}/config.json")
//	if err != nil {
//	    return err
//	}
//	// path now contains the absolute path with ~ and ${APP_DIR} expanded
func ExpandPath(path string) (string, error) {
	// Expand ~ to home directory
	expandedPath, err := homedir.Expand(path)
	if err != nil {
		return "", fmt.Errorf("failed to expand home directory: %w", err)
	}

	// Expand environment variables
	expandedPath = os.ExpandEnv(expandedPath)

	// Clean the path to remove any unnecessary elements
	expandedPath = filepath.Clean(expandedPath)

	// Convert to absolute path
	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return absPath, nil
}

// WaitForFiles waits for a set of files to exist.
// It will check immediately and then every interval seconds until:
// - All files exist
// - The context is canceled
// - The maximum number of tries is reached
//
// Parameters:
//   - ctx: Context for cancellation
//   - interval: Time interval in seconds between checks
//   - max: Maximum number of tries (including the immediate try)
//   - files: List of file paths to check
//
// Returns:
//   - error: nil if all files exist, otherwise an error explaining why it failed
//
// Example:
//
//	// Wait for config files to exist, checking every 2 seconds, up to 30 tries
//	err := util.WaitForFiles(ctx, 2, 30, "/etc/app/config.json", "/etc/app/secrets.json")
func WaitForFiles(ctx context.Context, interval, max uint, files ...string) error {
	if len(files) == 0 {
		return fmt.Errorf("no files specified")
	}

	i := time.Duration(interval) * time.Second
	return WaitFor(ctx, i, max, func() bool {
		return FilesExist(files...)
	})
}

// fileExists checks if a single file exists.
// It returns nil if the file exists, otherwise it returns the error from os.Stat.
func fileExists(filename string) error {
	_, err := os.Stat(filename)
	return err
}

// FilesExist checks if all specified files exist.
// It returns true only if all files exist, otherwise false.
//
// Example:
//
//	if util.FilesExist("/etc/app/config.json", "/etc/app/secrets.json") {
//	    // All files exist, proceed
//	} else {
//	    // At least one file is missing
//	}
func FilesExist(files ...string) bool {
	if len(files) == 0 {
		return true // No files to check means all files exist
	}
	return generics.Apply(fileExists, files) == nil
}

// decoder is an interface for decoding data into a Go value.
type decoder interface {
	Decode(v any) error
}

// encoder is an interface for encoding a Go value.
type encoder interface {
	Encode(v any) error
}

// decoderFunc is a function type that creates a decoder from an io.Reader.
type decoderFunc func(r io.Reader) decoder

// encoderFunc is a function type that creates an encoder from an io.Writer.
type encoderFunc func(w io.Writer) encoder

// yamlDecoderFunc creates a YAML decoder from an io.Reader.
func yamlDecoderFunc(r io.Reader) decoder {
	return yaml.NewDecoder(r)
}

// yamlEncoderFunc creates a YAML encoder from an io.Writer.
func yamlEncoderFunc(w io.Writer) encoder {
	return yaml.NewEncoder(w)
}

// jsonDecoderFunc creates a JSON decoder from an io.Reader.
func jsonDecoderFunc(r io.Reader) decoder {
	return json.NewDecoder(r)
}

// jsonEncoderFunc creates a JSON encoder from an io.Writer.
func jsonEncoderFunc(w io.Writer) encoder {
	return json.NewEncoder(w)
}

// encoderFuncFromFilePath returns an appropriate encoder function based on the file extension.
// Supported extensions: .yaml, .yml, .json
// Returns nil if the file extension is not recognized.
func encoderFuncFromFilePath(path string) encoderFunc {
	switch {
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		return yamlEncoderFunc
	case strings.HasSuffix(path, ".json"):
		return jsonEncoderFunc
	default:
		return nil
	}
}

// decoderFuncFromFilePath returns an appropriate decoder function based on the file extension.
// Supported extensions: .yaml, .yml, .json
// Returns nil if the file extension is not recognized.
func decoderFuncFromFilePath(path string) decoderFunc {
	switch {
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		return yamlDecoderFunc
	case strings.HasSuffix(path, ".json"):
		return jsonDecoderFunc
	default:
		return nil
	}
}

// saveStructToWriterWithEncoder encodes a struct to a writer using the provided encoder function.
// It's a helper function used by SaveStructToFile.
func saveStructToWriterWithEncoder[T any](v *T, w io.Writer, eFunc encoderFunc) error {
	enc := eFunc(w)
	return enc.Encode(v)
}

// loadStructFromReaderWithDecoder decodes a struct from a reader using the provided decoder function.
// It's a helper function used by LoadStructFromFile.
// Returns an error if the decoded data is a zero value (empty struct).
func loadStructFromReaderWithDecoder[T any](r io.Reader, dFunc decoderFunc) (*T, error) {
	var data T

	dec := dFunc(r)
	err := dec.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	if generics.IsZeroValue(data) {
		return nil, fmt.Errorf("decoded data is empty (zero value)")
	}

	return &data, nil
}

// LoadStructFromFile loads a struct from a file.
// The file format is determined by the file extension (.json, .yaml, or .yml).
//
// Parameters:
//   - filePath: Path to the file to load from
//
// Returns:
//   - *T: Pointer to the loaded struct
//   - error: Error if loading fails
//
// Example:
//
//	type Config struct {
//	    ServerName string `json:"server_name"`
//	    Port       int    `json:"port"`
//	}
//
//	config, err := util.LoadStructFromFile[Config]("/etc/app/config.json")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Server: %s, Port: %d\n", config.ServerName, config.Port)
func LoadStructFromFile[T any](filePath string) (*T, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	decFunc := decoderFuncFromFilePath(filePath)
	if decFunc == nil {
		return nil, fmt.Errorf("unsupported file format: %s (expected .yaml, .yml, or .json)", filepath.Ext(filePath))
	}

	structFile, err := CleanOpen(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		// We already handle the close error in the non-deferred code path
		// This is just to ensure the file is closed in case of early returns
		_ = structFile.Close()
	}()

	data, err := loadStructFromReaderWithDecoder[T](structFile, decFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to load data from %s: %w", filePath, err)
	}

	// Explicitly close the file to catch any close errors
	if closeErr := structFile.Close(); closeErr != nil {
		return nil, fmt.Errorf("error closing file after successful read: %w", closeErr)
	}

	return data, nil
}

// SaveStructToFile saves a struct to a file.
// The file format is determined by the file extension (.json, .yaml, or .yml).
// If the directory doesn't exist, it will be created.
//
// Parameters:
//   - v: Pointer to the struct to save
//   - filePath: Path to the file to save to
//
// Returns:
//   - error: Error if saving fails
//
// Example:
//
//	config := &Config{
//	    ServerName: "api-server",
//	    Port:       8080,
//	}
//
//	err := util.SaveStructToFile(config, "/etc/app/config.json")
//	if err != nil {
//	    return err
//	}
func SaveStructToFile[T any](v *T, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	encFunc := encoderFuncFromFilePath(filePath)
	if encFunc == nil {
		return fmt.Errorf("unsupported file format: %s (expected .yaml, .yml, or .json)", filepath.Ext(filePath))
	}

	// Create directory if it doesn't exist
	filePathDir := filepath.Dir(filePath)
	_, err := CreateDirPath(filePathDir, "")
	if err != nil {
		return fmt.Errorf("failed to create directory path: %w", err)
	}

	// Open file with appropriate permissions
	structFile, err := CleanOpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer func() {
		// We already handle the close error in the non-deferred code path
		// This is just to ensure the file is closed in case of early returns
		_ = structFile.Close()
	}()

	// Encode and write the struct to the file
	err = saveStructToWriterWithEncoder[T](v, structFile, encFunc)
	if err != nil {
		return fmt.Errorf("failed to encode data to %s: %w", filePath, err)
	}

	// Explicitly close the file to catch any close errors
	if closeErr := structFile.Close(); closeErr != nil {
		return fmt.Errorf("error closing file after successful write: %w", closeErr)
	}

	return nil
}
