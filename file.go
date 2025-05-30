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

func CleanOpen(path string) (*os.File, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return nil, err
	}

	path = filepath.Clean(path)

	return os.Open(path)
}

func CleanOpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(expandedPath)

	return os.OpenFile(cleanPath, flag, perm) // #nosec
}

// CreateDirPath creates a directory path if it doesn't exist.
func CreateDirPath(path string, defaultPath string) (string, error) {
	if path == "" {
		path = defaultPath
	}

	path, err := ExpandPath(path)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path, 0750)
	if err != nil {
		return "", err
	}

	return path, nil
}

// ExpandPath expands a path to an absolute path.
// It also expands ~ and environment variables.
func ExpandPath(path string) (string, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}

	path = os.ExpandEnv(path)

	path = filepath.Clean(path)

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

// WaitForFiles waits for a set of files to exist, it will check every interval seconds up until max seconds.
func WaitForFiles(ctx context.Context, interval, max uint, files ...string) error {
	i := time.Duration(interval) * time.Second
	return WaitFor(ctx, i, max, func() bool {
		return FilesExist(files...)
	})
}

func fileExists(filename string) error {
	_, err := os.Stat(filename)
	return err
}

// FilesExist checks if all file names exist.
func FilesExist(files ...string) bool {
	return generics.Apply(fileExists, files) == nil
}

type decoder interface {
	Decode(v any) error
}

type encoder interface {
	Encode(v any) error
}

type decoderFunc func(r io.Reader) decoder
type encoderFunc func(w io.Writer) encoder

func yamlDecoderFunc(r io.Reader) decoder {
	return yaml.NewDecoder(r)
}

func yamlEncoderFunc(w io.Writer) encoder {
	return yaml.NewEncoder(w)
}

func jsonDecoderFunc(r io.Reader) decoder {
	return json.NewDecoder(r)
}

func jsonEncoderFunc(w io.Writer) encoder {
	return json.NewEncoder(w)
}

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

func saveStructToWriterWithEncoder[T any](v *T, w io.Writer, eFunc encoderFunc) error {
	enc := eFunc(w)
	return enc.Encode(v)
}

func loadStructFromReaderWithDecoder[T any](r io.Reader, dFunc decoderFunc) (*T, error) {
	var data T

	dec := dFunc(r)
	err := dec.Decode(&data)
	if err != nil {
		return nil, err
	}

	if generics.IsZeroValue(data) {
		return nil, fmt.Errorf("failed to load data from file")
	}

	return &data, nil
}

func LoadStructFromFile[T any](filePath string) (*T, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}

	decFunc := decoderFuncFromFilePath(filePath)

	if decFunc == nil {
		return nil, fmt.Errorf("unrecognised file type. expected yaml/yml or json")
	}

	structFile, err := CleanOpen(filePath)
	if err != nil {
		return nil, err
	}

	data, err := loadStructFromReaderWithDecoder[T](structFile, decFunc)

	if err != nil {
		closeErr := structFile.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("%w: %v", err, closeErr)
		}
		return nil, err
	}

	return data, structFile.Close()
}

func SaveStructToFile[T any](v *T, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is empty")
	}

	encFunc := encoderFuncFromFilePath(filePath)

	if encFunc == nil {
		return fmt.Errorf("unrecognised file type. expected yaml/yml or json")
	}

	filePathDir := filepath.Dir(filePath)
	_, err := CreateDirPath(filePathDir, "")
	if err != nil {
		return fmt.Errorf("failed to create directory path: %w", err)
	}

	structFile, err := CleanOpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	err = saveStructToWriterWithEncoder[T](v, structFile, encFunc)

	if err != nil {
		closeErr := structFile.Close()
		if closeErr != nil {
			return fmt.Errorf("%w: %v", err, closeErr)
		}
		return err
	}

	return structFile.Close()
}
